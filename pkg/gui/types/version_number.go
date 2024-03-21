package types

import (
	"errors"
	"regexp"
	"strconv"
)

type VersionNumber struct {
	Major, Minor, Patch int
}

func (v *VersionNumber) IsOlderThan(otherVersion *VersionNumber) bool {
	this := v.Major*1000*1000 + v.Minor*1000 + v.Patch
	other := otherVersion.Major*1000*1000 + otherVersion.Minor*1000 + otherVersion.Patch
	return this < other
}

func ParseVersionNumber(versionStr string) (*VersionNumber, error) {
	re := regexp.MustCompile(`^v?(\d+)\.(\d+)(?:\.(\d+))?$`)
	matches := re.FindStringSubmatch(versionStr)
	if matches == nil {
		return nil, errors.New("unexpected version format: " + versionStr)
	}

	v := &VersionNumber{}
	var err error

	if v.Major, err = strconv.Atoi(matches[1]); err != nil {
		return nil, err
	}
	if v.Minor, err = strconv.Atoi(matches[2]); err != nil {
		return nil, err
	}
	if len(matches[3]) > 0 {
		if v.Patch, err = strconv.Atoi(matches[3]); err != nil {
			return nil, err
		}
	}
	return v, nil
}
