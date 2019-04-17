// +build codegen

package api

type wafregionalExamplesBuilder struct {
	defaultExamplesBuilder
}

func (builder wafregionalExamplesBuilder) Imports(a *API) string {
	return `"fmt"
	"strings"
	"time"

	"` + SDKImportRoot + `/aws"
	"` + SDKImportRoot + `/aws/awserr"
	"` + SDKImportRoot + `/aws/session"
	"` + SDKImportRoot + `/service/waf"
	"` + SDKImportRoot + `/service/` + a.PackageName() + `"
	`
}
