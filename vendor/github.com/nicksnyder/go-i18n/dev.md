# Development notes

## How to upgrade CLDR data

1.  Go to http://cldr.unicode.org/index/downloads to find the latest version.
1.  Download the latest version of cldr-common (e.g. http://unicode.org/Public/cldr/33/cldr-common-33.0.zip)
1.  Unzip and copy `common/supplemental/plurals.xml` to `v2/i18n/internal/plural/codegen/plurals.xml`
1.  Run `generate.sh` in `v2/i18n/internal/plural/codegen/`
