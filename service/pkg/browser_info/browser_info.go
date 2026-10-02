package browser_info

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pkg/helper/context_helper"

	"github.com/Masterminds/semver"
	"github.com/pattfy/useragent"
)

type BrowserInfo struct {
	DeviceID        string
	PlatformID      string
	PlatformVersion string
	BrowserID       string
	BrowserVersion  string
}

func FromString(s string) (BrowserInfo, error) {
	var info BrowserInfo
	err := json.Unmarshal([]byte(s), &info)
	if err != nil {
		return BrowserInfo{}, err
	}
	return info, nil
}

func (b BrowserInfo) String() string {
	bytes, err := json.Marshal(b)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func New(ctx context.Context) BrowserInfo {
	agent := context_helper.Request(ctx).UserAgent()
	userAgent := useragent.New(agent)

	platformVer := func() string {
		v, err := semver.NewVersion(userAgent.Platform.Version())
		if err != nil {
			return ""
		}
		return fmt.Sprint(v.Major())
	}()

	browserVer := func() string {
		v, err := semver.NewVersion(userAgent.Browser.Version())
		if err != nil {
			return ""
		}
		return fmt.Sprint(v.Major())
	}()
	return BrowserInfo{
		DeviceID:        userAgent.Device.ID(),
		PlatformID:      userAgent.Platform.ID(),
		PlatformVersion: platformVer,
		BrowserID:       userAgent.Browser.ID(),
		BrowserVersion:  browserVer,
	}
}
