package processor

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/njkevlani/aur-autoupdater/internal/aur/pkgbuild"
	"github.com/njkevlani/aur-autoupdater/internal/aurversion"
	"github.com/njkevlani/aur-autoupdater/internal/internalerrors"
	"github.com/njkevlani/aur-autoupdater/internal/latestversion"
	"github.com/njkevlani/aur-autoupdater/internal/packageinfo"
	"github.com/njkevlani/aur-autoupdater/internal/version"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func Process(packageName string) error {
	var (
		packageInfo packageinfo.PackageInfo
		ok          bool
	)

	if packageInfo, ok = packageinfo.PackageInfos[packageName]; !ok {
		return internalerrors.ErrUnknownPackage
	}

	repo, fs, err := getRepo(packageName)
	if err != nil {
		return err
	}

	srcinfoFile, err := fs.OpenFile(".SRCINFO", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	aurVersion, err := aurversion.GetAURVersion(srcinfoFile)
	if err != nil {
		return err
	}

	if err = srcinfoFile.Close(); err != nil {
		return err
	}

	latestVersion, err := latestversion.GetLatestVersion(packageInfo.GitHubInfo.Owner, packageInfo.GitHubInfo.Repo)
	if err != nil {
		return err
	}

	if version.Equal(aurVersion, latestVersion) {
		logrus.Infof("versions are same for %s", packageName)
	} else {
		logrus.Infof("versions are not same for %s", packageName)

		changeVersion(repo, fs, packageInfo, latestVersion)
	}

	return nil
}

func changeVersion(repo *git.Repository, fs billy.Filesystem, packageInfo packageinfo.PackageInfo, latestVersion version.Version) error {
	workTree, err := repo.Worktree()
	if err != nil {
		return err
	}

	sha256Sum, err := getSHA256Sum(packageInfo.GitHubInfo.ReleaseAssetURL(version.StripV(latestVersion.Version())))
	if err != nil {
		return err
	}

	srcinfoFile, err := fs.OpenFile(".SRCINFO", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	srcinfoFile.Truncate(0)
	srcinfoFile.Seek(0, 0)

	if err = pkgbuild.RenderSrcinfo(
		packageInfo.Name,
		pkgbuild.Pkgbuild{
			Pkgver:    version.StripV(latestVersion.Version()),
			SHA256Sum: sha256Sum,
		},
		srcinfoFile,
	); err != nil {
		return err
	}

	if err = srcinfoFile.Close(); err != nil {
		return err
	}

	workTree.Add(".SRCINFO")

	pkgbuildFile, err := fs.OpenFile("PKGBUILD", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	pkgbuildFile.Truncate(0)
	pkgbuildFile.Seek(0, 0)

	if err = pkgbuild.RenderPkgbuild(
		packageInfo.Name,
		pkgbuild.Pkgbuild{
			Pkgver:    version.StripV(latestVersion.Version()),
			SHA256Sum: sha256Sum,
		},
		pkgbuildFile,
	); err != nil {
		return err
	}

	if err = pkgbuildFile.Close(); err != nil {
		return err
	}

	workTree.Add("PKGBUILD")

	status, err := workTree.Status()
	if err != nil {
		return err
	}

	logrus.WithField("status", status).Info("git status")

	commit, err := workTree.Commit(fmt.Sprintf("Updated to %s", latestVersion.Version()), &git.CommitOptions{
		Author: &object.Signature{Name: "Nilesh", Email: "njkevlani@gmail.com", When: time.Now()},
	})
	if err != nil {
		return err
	}

	logrus.WithField("commit", commit).Info("commit")

	commitObj, err := repo.CommitObject(commit)
	if err != nil {
		return err
	}

	logrus.WithField("commitObj", commitObj).Info("commitObj")

	return pushRepo(repo)
}

func getRepo(packageName string) (*git.Repository, billy.Filesystem, error) {
	fs := memfs.New()

	sshPrivateKeyPassword := os.Getenv("SSH_KEY_PASSWORD")
	sshPrivateKey := os.Getenv("SSH_KEY")
	publicKey, err := gitssh.NewPublicKeys("aur", []byte(sshPrivateKey), sshPrivateKeyPassword)
	publicKey.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	if err != nil {
		return nil, nil, err
	}

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:  fmt.Sprintf("ssh://aur@aur.archlinux.org/%s.git", packageName),
		Auth: publicKey,
	})

	if err != nil {
		return nil, nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, err
	}

	logrus.WithField("head", head.Hash()).Info("cloned repo")

	return repo, fs, nil
}

func getSHA256Sum(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	sha256Sum := fmt.Sprintf("%x", sha256.Sum256(respBodyBytes))

	logrus.WithField("url", url).
		WithField("sha256Sum", sha256Sum).
		Info("sha256Sum")

	return sha256Sum, nil
}

func pushRepo(repo *git.Repository) error {
	sshPrivateKeyPassword := os.Getenv("SSH_KEY_PASSWORD")
	sshPrivateKey := os.Getenv("SSH_KEY")
	publicKey, err := gitssh.NewPublicKeys("aur", []byte(sshPrivateKey), sshPrivateKeyPassword)
	publicKey.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	if err != nil {
		return err
	}

	return repo.Push(&git.PushOptions{Auth: publicKey})
}
