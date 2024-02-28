package env

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	OS      = runtime.GOOS
	Arch    = runtime.GOARCH
	Version = ""
)

func Env() OSArchVersion {
	return OSArchVersion{
		OS:      OS,
		Arch:    Arch,
		Version: Version,
	}
}

func Platform() string {
	return Env().OSArch()
}

func PlatformVersion() string {
	return Env().String()
}

type OSArchVersion struct {
	OS      string
	Arch    string
	Version string
}

func (oav OSArchVersion) OSArch() string {
	return strings.Join([]string{oav.OS, oav.Arch}, "/")
}

func (oav OSArchVersion) String() string {
	return strings.Join([]string{oav.OS, oav.Arch, oav.Version}, "/")
}

func GetOSArchVersions(osArchVersions ...string) []OSArchVersion {
	if len(osArchVersions) == 0 {
		return []OSArchVersion{Env()}
	}

	// get OS archs
	crossOSArch := make([]OSArchVersion, len(osArchVersions))
	for i, osArchVersion := range osArchVersions {
		oavSlice := strings.Split(osArchVersion, "/")
		var os, arch, version string
		if len(oavSlice) > 0 {
			os = oavSlice[0]
		}
		if len(oavSlice) > 1 {
			arch = oavSlice[1]
		}
		if len(oavSlice) > 2 {
			version = oavSlice[2]
		}
		if len(oavSlice) > 3 {
			panic(fmt.Errorf("osArchVersion %s has too many components", osArchVersion))
		}
		crossOSArch[i] = OSArchVersion{
			OS:      os,
			Arch:    arch,
			Version: version,
		}
	}
	return crossOSArch
}
