package version

import (
	"fmt"
	"runtime"
)

// nolint:gochecknoglobals
var (
	Version    = ""
	Commit     = ""
	VendorHash = ""
	BuildTags  = ""
)

type Info struct {
	Version    string `json:"version"`
	Commit     string `json:"commit"`
	VendorHash string `json:"vendor_hash"`
	Go         string `json:"go"`
}

func (i Info) String() string {
	return fmt.Sprintf(`%s
git commit: %s
vendor hash: %s
%s`, i.Version, i.Commit, i.VendorHash, i.Go)
}

func newInfo() Info {
	return Info{
		Version:    Version,
		Commit:     Commit,
		VendorHash: VendorHash,
		Go:         fmt.Sprintf("go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
}
