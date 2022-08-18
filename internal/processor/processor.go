package processor

import (
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/njkevlani/aur-autoupdater/internal/aurversion"
	"github.com/njkevlani/aur-autoupdater/internal/internalerrors"
	"github.com/njkevlani/aur-autoupdater/internal/latestversion"
	"github.com/njkevlani/aur-autoupdater/internal/packageinfo"
	"github.com/njkevlani/aur-autoupdater/internal/version"
	"github.com/sirupsen/logrus"
)

var (
	sshPrivateKeyPassword = os.Getenv("SSH_KEY_PASSWORD")
	sshPrivateKey         = os.Getenv("SSH_KEY")
)

func Process(packageName string) error {
	// git clone
	// get version from git repo
	// get version from github
	// if version different,
	//     - edit version in repo
	//     - edit sha in repo
	//     - git commit and git push

	var (
		packageInfo packageinfo.PackageInfo
		ok          bool
	)

	if packageInfo, ok = packageinfo.PackageInfos[packageName]; !ok {
		return internalerrors.ErrUnknownPackage
	}

	fs := memfs.New()

	// publicKey, err := ssh.NewPublicKeys("aur", []byte(sshPrivateKey), sshPrivateKeyPassword)

	// if err != nil {
	// 	return err
	// }

	// repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
	// 	URL:  fmt.Sprintf("ssh://aur@aur.archlinux.org/%s.git", packageName),
	// 	Auth: publicKey,
	// })

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: "file:///home/njkevlani/git/njkevlani/xremap-x11-bin",
	})

	if err != nil {
		return err
	}

	head, err := repo.Head()

	if err != nil {
		return err
	}

	logrus.WithField("head", head.Hash()).Info("cloned repo")

	srcinfoFile, err := fs.Open(".SRCINFO")

	if err != nil {
		return err
	}

	defer func() {
		if err = srcinfoFile.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	aurVersion, err := aurversion.GetAURVersion(srcinfoFile)

	if err != nil {
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
		// TODO:
		//   1. Get sha256sum for current zip file.
		//   2. Make new PKGBUILD
		//   3. Make new .SRCINFO
		//   4. Add files in repo
		//   5. Commit
		//   6. Push
	}

	return nil
}
