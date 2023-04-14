package container

import (
	"path"

	warewulfconf "github.com/hpcng/warewulf/internal/pkg/config"
)

func SourceParentDir() string {
	conf := warewulfconf.Get()
	return conf.Paths.WWChrootdir
}

func SourceDir(name string) string {
	return path.Join(SourceParentDir(), name)
}

func RootFsDir(name string) string {
	return path.Join(SourceDir(name), "rootfs")
}

func ImageParentDir() string {
	conf := warewulfconf.Get()
	return path.Join(conf.Paths.WWProvisiondir, "container/")
}

func ImageFile(name string) string {
	return path.Join(ImageParentDir(), name+".img")
}
