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
	"github.com/sirupsen/logrus"
	"github.com/xremap/aur-autoupdater/internal/aur/pkgbuild"
	"github.com/xremap/aur-autoupdater/internal/aurversion"
	"github.com/xremap/aur-autoupdater/internal/internalerrors"
	"github.com/xremap/aur-autoupdater/internal/latestversion"
	"github.com/xremap/aur-autoupdater/internal/packageinfo"
	"github.com/xremap/aur-autoupdater/internal/version"
	"golang.org/x/crypto/ssh"
)

func Process(packageName string) error {
	var (
		packageInfo packageinfo.PackageInfo
		ok          bool
	)

	if packageInfo, ok = packageinfo.PackageInfos[packageName]; !ok {
		logrus.WithField("packageName", packageName).Error("unknown packageName")
		return internalerrors.ErrUnknownPackage
	}

	repo, fs, err := getRepo(packageName)
	if err != nil {
		logrus.WithError(err).Error("failed to get repo")
		return err
	}

	srcinfoFile, err := fs.OpenFile(".SRCINFO", os.O_RDWR, 0644)
	if err != nil {
		logrus.WithError(err).Error("failed to open scrinfo file")
		return err
	}

	aurVersion, err := aurversion.GetAURVersion(srcinfoFile)
	if err != nil {
		logrus.WithError(err).Error("failed to get aur version")
		return err
	}

	if err = srcinfoFile.Close(); err != nil {
		logrus.WithError(err).Error("failed to close srcinfo file")
		return err
	}

	latestVersion, err := latestversion.GetLatestVersion(packageInfo.GitHubInfo.Owner, packageInfo.GitHubInfo.Repo)
	if err != nil {
		logrus.WithError(err).Error("failed to get latest version")
		return err
	}

	if version.Equal(aurVersion, latestVersion) {
		logrus.Infof("versions are same for %s", packageName)
	} else {
		logrus.Infof("versions are not same for %s", packageName)

		err := changeVersion(repo, fs, packageInfo, latestVersion)

		if err != nil {
			logrus.WithError(err).Error("failed to change version")
			return err
		}
	}

	return nil
}

func changeVersion(repo *git.Repository, fs billy.Filesystem, packageInfo packageinfo.PackageInfo, latestVersion version.Version) error {
	workTree, err := repo.Worktree()
	if err != nil {
		logrus.WithError(err).Error("failed to get worktree")
		return err
	}

	sha256Sum, err := getSHA256Sum(packageInfo.GitHubInfo.ReleaseAssetURL(version.StripV(latestVersion.Version())))
	if err != nil {
		logrus.WithError(err).Error("failed to get sha256sum")
		return err
	}

	var sha256SumAarch64 string
	if packageInfo.GitHubInfoAarch64.ReleaseAssetURL != nil {
		var err error
		sha256SumAarch64, err = getSHA256Sum(packageInfo.GitHubInfoAarch64.ReleaseAssetURL(version.StripV(latestVersion.Version())))
		if err != nil {
			logrus.WithError(err).Error("failed to get aarch64 sha256sum")
			return err
		}
	}

	srcinfoFile, err := fs.OpenFile(".SRCINFO", os.O_RDWR, 0644)
	if err != nil {
		logrus.WithError(err).Error("failed to open srcinfo file")
		return err
	}

	srcinfoFile.Truncate(0)
	srcinfoFile.Seek(0, 0)

	if err = pkgbuild.RenderSrcinfo(
		packageInfo.Name,
		pkgbuild.Pkgbuild{
			Pkgver:           version.StripV(latestVersion.Version()),
			SHA256Sum:        sha256Sum,
			SHA256SumAarch64: sha256SumAarch64,
		},
		srcinfoFile,
	); err != nil {
		logrus.WithError(err).Error("failed to render srcinfo file")
		return err
	}

	if err = srcinfoFile.Close(); err != nil {
		logrus.WithError(err).Error("failed to close srcinfo file")
		return err
	}

	workTree.Add(".SRCINFO")

	pkgbuildFile, err := fs.OpenFile("PKGBUILD", os.O_RDWR, 0644)
	if err != nil {
		logrus.WithError(err).Error("failed to open pkgbuild file")
		return err
	}

	pkgbuildFile.Truncate(0)
	pkgbuildFile.Seek(0, 0)

	if err = pkgbuild.RenderPkgbuild(
		packageInfo.Name,
		pkgbuild.Pkgbuild{
			Pkgver:           version.StripV(latestVersion.Version()),
			SHA256Sum:        sha256Sum,
			SHA256SumAarch64: sha256SumAarch64,
		},
		pkgbuildFile,
	); err != nil {
		logrus.WithError(err).Error("failed to render pkgbuild file")
		return err
	}

	if err = pkgbuildFile.Close(); err != nil {
		logrus.WithError(err).Error("failed to close pkgbuild file")
		return err
	}

	workTree.Add("PKGBUILD")

	status, err := workTree.Status()
	if err != nil {
		logrus.WithError(err).Error("failed to get worktree status")
		return err
	}

	logrus.WithField("status", status).Info("git status")

	commit, err := workTree.Commit(fmt.Sprintf("Updated to %s", latestVersion.Version()), &git.CommitOptions{
		Author: &object.Signature{Name: "k0kubun", Email: "takashikkbn@gmail.com", When: time.Now()},
	})
	if err != nil {
		logrus.WithError(err).Error("failed to commit")
		return err
	}

	logrus.WithField("commit", commit).Info("commit")

	commitObj, err := repo.CommitObject(commit)
	if err != nil {
		logrus.WithError(err).Error("failed to get commit object")
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
	if err != nil {
		logrus.WithError(err).Error("failed to get ssh key")
		return nil, nil, err
	}

	publicKey.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:  fmt.Sprintf("ssh://aur@aur.archlinux.org/%s.git", packageName),
		Auth: publicKey,
	})

	if err != nil {
		logrus.WithError(err).WithField("packageName", packageName).Error("failed to clone reo")
		return nil, nil, err
	}

	head, err := repo.Head()
	if err != nil {
		logrus.WithError(err).Error("failed to get head of repo")
		return nil, nil, err
	}

	logrus.WithField("head", head.Hash()).Info("cloned repo")

	return repo, fs, nil
}

func getSHA256Sum(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("failed to get response")
		return "", err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("failed to read response body")
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
