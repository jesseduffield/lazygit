package packp

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing"
	"net/url"
	"strings"
)

var ErrUnsupportedObjectFilterType = errors.New("unsupported object filter type")

// Filter values enable the partial clone capability which causes
// the server to omit objects that match the filter.
//
// See [Git's documentation] for more details.
//
// [Git's documentation]: https://github.com/git/git/blob/e02ecfcc534e2021aae29077a958dd11c3897e4c/Documentation/rev-list-options.txt#L948
type Filter string

type BlobLimitPrefix string

const (
	BlobLimitPrefixNone BlobLimitPrefix = ""
	BlobLimitPrefixKibi BlobLimitPrefix = "k"
	BlobLimitPrefixMebi BlobLimitPrefix = "m"
	BlobLimitPrefixGibi BlobLimitPrefix = "g"
)

// FilterBlobNone omits all blobs.
func FilterBlobNone() Filter {
	return "blob:none"
}

// FilterBlobLimit omits blobs of size at least n bytes (when prefix is
// BlobLimitPrefixNone), n kibibytes (when prefix is BlobLimitPrefixKibi),
// n mebibytes (when prefix is BlobLimitPrefixMebi) or n gibibytes (when
// prefix is BlobLimitPrefixGibi). n can be zero, in which case all blobs
// will be omitted.
func FilterBlobLimit(n uint64, prefix BlobLimitPrefix) Filter {
	return Filter(fmt.Sprintf("blob:limit=%d%s", n, prefix))
}

// FilterTreeDepth omits all blobs and trees whose depth from the root tree
// is larger or equal to depth.
func FilterTreeDepth(depth uint64) Filter {
	return Filter(fmt.Sprintf("tree:%d", depth))
}

// FilterObjectType omits all objects which are not of the requested type t.
// Supported types are TagObject, CommitObject, TreeObject and BlobObject.
func FilterObjectType(t plumbing.ObjectType) (Filter, error) {
	switch t {
	case plumbing.TagObject:
		fallthrough
	case plumbing.CommitObject:
		fallthrough
	case plumbing.TreeObject:
		fallthrough
	case plumbing.BlobObject:
		return Filter(fmt.Sprintf("object:type=%s", t.String())), nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedObjectFilterType, t.String())
	}
}

// FilterCombine combines multiple Filter values together.
func FilterCombine(filters ...Filter) Filter {
	var escapedFilters []string

	for _, filter := range filters {
		escapedFilters = append(escapedFilters, url.QueryEscape(string(filter)))
	}

	return Filter(fmt.Sprintf("combine:%s", strings.Join(escapedFilters, "+")))
}
