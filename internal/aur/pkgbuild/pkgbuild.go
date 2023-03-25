package pkgbuild

import (
	"io"
	"text/template"

	"github.com/xremap/aur-autoupdater/internal/internalerrors"
	"github.com/xremap/aur-autoupdater/internal/packageinfo"
)

type Pkgbuild struct {
	Pkgver    string
	SHA256Sum string
}

type renderType int

const (
	renderTypePkgbuid = renderType(iota)
	renderTypeSrcinfo
)

func render(renderType renderType, packageName string, pkgbuild Pkgbuild, dst io.Writer) error {
	var (
		packageInfo      packageinfo.PackageInfo
		ok               bool
		templateFilepath string
	)

	if packageInfo, ok = packageinfo.PackageInfos[packageName]; !ok {
		return internalerrors.ErrUnknownPackage
	}

	switch renderType {
	case renderTypePkgbuid:
		templateFilepath = packageInfo.PkgbuildInfo.PkgbuildTemplateFilepath

	case renderTypeSrcinfo:
		templateFilepath = packageInfo.PkgbuildInfo.SrcinfoTemplateFilepath

	default:
		panic("should not happen")
	}

	template, err := template.ParseFiles(templateFilepath)
	if err != nil {
		return err
	}

	err = template.Execute(dst, pkgbuild)
	if err != nil {
		return err
	}

	return nil
}

func RenderPkgbuild(packageName string, pkgbuild Pkgbuild, dst io.Writer) error {
	return render(renderTypePkgbuid, packageName, pkgbuild, dst)
}

func RenderSrcinfo(packageName string, pkgbuild Pkgbuild, dst io.Writer) error {
	return render(renderTypeSrcinfo, packageName, pkgbuild, dst)
}
