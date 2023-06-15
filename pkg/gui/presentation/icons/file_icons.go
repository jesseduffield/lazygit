package icons

import (
	"path/filepath"
)

// https://github.com/ogham/exa/blob/master/src/output/icons.rs
const (
	DEFAULT_FILE_ICON      = "\uf15b" // 
	DEFAULT_SUBMODULE_ICON = "\uf1d3" // 
	DEFAULT_DIRECTORY_ICON = "\uf114" // 
)

var nameIconMap = map[string]map[int]string{
	".Trash":             map[int]string{2: "\uf1f8"}, // 
	".atom":              map[int]string{2: "\ue764"}, // 
	".bashprofile":       map[int]string{2: "\ue615"}, // 
	".bashrc":            map[int]string{2: "\uf489"}, // 
	".idea":              map[int]string{2: "\ue7b5"}, // 
	".git":               map[int]string{2: "\uf1d3"}, // 
	".gitattributes":     map[int]string{2: "\uf1d3"}, // 
	".gitconfig":         map[int]string{2: "\uf1d3"}, // 
	".github":            map[int]string{2: "\uf408"}, // 
	".gitignore":         map[int]string{2: "\uf1d3"}, // 
	".gitmodules":        map[int]string{2: "\uf1d3"}, // 
	".rvm":               map[int]string{2: "\ue21e"}, // 
	".vimrc":             map[int]string{2: "\ue62b"}, // 
	".vscode":            map[int]string{2: "\ue70c"}, // 
	".zshrc":             map[int]string{2: "\uf489"}, // 
	"Cargo.lock":         map[int]string{2: "\ue7a8"}, // 
	"Cargo.toml":         map[int]string{2: "\ue7a8"}, // 
	"bin":                map[int]string{2: "\ue5fc"}, // 
	"config":             map[int]string{2: "\ue5fc"}, // 
	"docker-compose.yml": map[int]string{2: "\uf308"}, // 
	"Dockerfile":         map[int]string{2: "\uf308"}, // 
	"ds_store":           map[int]string{2: "\uf179"}, // 
	"gitignore_global":   map[int]string{2: "\uf1d3"}, // 
	"go.mod":             map[int]string{2: "\ue626"}, // 
	"go.sum":             map[int]string{2: "\ue626"}, // 
	"gradle":             map[int]string{2: "\ue256"}, // 
	"gruntfile.coffee":   map[int]string{2: "\ue611"}, // 
	"gruntfile.js":       map[int]string{2: "\ue611"}, // 
	"gruntfile.ls":       map[int]string{2: "\ue611"}, // 
	"gulpfile.coffee":    map[int]string{2: "\ue610"}, // 
	"gulpfile.js":        map[int]string{2: "\ue610"}, // 
	"gulpfile.ls":        map[int]string{2: "\ue610"}, // 
	"hidden":             map[int]string{2: "\uf023"}, // 
	"include":            map[int]string{2: "\ue5fc"}, // 
	"lib":                map[int]string{2: "\uf121"}, // 
	"localized":          map[int]string{2: "\uf179"}, // 
	"Makefile":           map[int]string{2: "\uf489"}, // 
	"node_modules":       map[int]string{2: "\ue718"}, // 
	"npmignore":          map[int]string{2: "\ue71e"}, // 
	"PKGBUILD":           map[int]string{2: "\uf303"}, // 
	"rubydoc":            map[int]string{2: "\ue73b"}, // 
	"yarn.lock":          map[int]string{2: "\ue718"}, // 
}

var extIconMap = map[string]map[int]string{
	".ai":             map[int]string{2: "\ue7b4"}, // 
	".android":        map[int]string{2: "\ue70e"}, // 
	".apk":            map[int]string{2: "\ue70e"}, // 
	".apple":          map[int]string{2: "\uf179"}, // 
	".avi":            map[int]string{2: "\uf03d"}, // 
	".avif":           map[int]string{2: "\uf1c5"}, // 
	".avro":           map[int]string{2: "\ue60b"}, // 
	".awk":            map[int]string{2: "\uf489"}, // 
	".bash":           map[int]string{2: "\uf489"}, // 
	".bash_history":   map[int]string{2: "\uf489"}, // 
	".bash_profile":   map[int]string{2: "\uf489"}, // 
	".bashrc":         map[int]string{2: "\uf489"}, // 
	".bat":            map[int]string{2: "\uf17a"}, // 
	".bats":           map[int]string{2: "\uf489"}, // 
	".bmp":            map[int]string{2: "\uf1c5"}, // 
	".bz":             map[int]string{2: "\uf410"}, // 
	".bz2":            map[int]string{2: "\uf410"}, // 
	".c":              map[int]string{2: "\ue61e"}, // 
	".c++":            map[int]string{2: "\ue61d"}, // 
	".cab":            map[int]string{2: "\ue70f"}, // 
	".cc":             map[int]string{2: "\ue61d"}, // 
	".cfg":            map[int]string{2: "\ue615"}, // 
	".class":          map[int]string{2: "\ue256"}, // 
	".clj":            map[int]string{2: "\ue768"}, // 
	".cljs":           map[int]string{2: "\ue76a"}, // 
	".cls":            map[int]string{2: "\uf034"}, // 
	".cmd":            map[int]string{2: "\ue70f"}, // 
	".coffee":         map[int]string{2: "\uf0f4"}, // 
	".conf":           map[int]string{2: "\ue615"}, // 
	".cp":             map[int]string{2: "\ue61d"}, // 
	".cpio":           map[int]string{2: "\uf410"}, // 
	".cpp":            map[int]string{2: "\ue61d"}, // 
	".cs":             map[int]string{2: "\uf81a"}, // 
	".csh":            map[int]string{2: "\uf489"}, // 
	".cshtml":         map[int]string{2: "\uf1fa"}, // 
	".csproj":         map[int]string{2: "\uf81a"}, // 
	".css":            map[int]string{2: "\ue749"}, // 
	".csv":            map[int]string{2: "\uf1c3"}, // 
	".csx":            map[int]string{2: "\uf81a"}, // 
	".cxx":            map[int]string{2: "\ue61d"}, // 
	".d":              map[int]string{2: "\ue7af"}, // 
	".dart":           map[int]string{2: "\ue798"}, // 
	".db":             map[int]string{2: "\uf1c0"}, // 
	".deb":            map[int]string{2: "\ue77d"}, // 
	".diff":           map[int]string{2: "\uf440"}, // 
	".djvu":           map[int]string{2: "\uf02d"}, // 
	".dll":            map[int]string{2: "\ue70f"}, // 
	".doc":            map[int]string{2: "\uf1c2"}, // 
	".docx":           map[int]string{2: "\uf1c2"}, // 
	".ds_store":       map[int]string{2: "\uf179"}, // 
	".DS_store":       map[int]string{2: "\uf179"}, // 
	".dump":           map[int]string{2: "\uf1c0"}, // 
	".ebook":          map[int]string{2: "\ue28b"}, // 
	".ebuild":         map[int]string{2: "\uf30d"}, // 
	".editorconfig":   map[int]string{2: "\ue615"}, // 
	".ejs":            map[int]string{2: "\ue618"}, // 
	".elm":            map[int]string{2: "\ue62c"}, // 
	".env":            map[int]string{2: "\uf462"}, // 
	".eot":            map[int]string{2: "\uf031"}, // 
	".epub":           map[int]string{2: "\ue28a"}, // 
	".erb":            map[int]string{2: "\ue73b"}, // 
	".erl":            map[int]string{2: "\ue7b1"}, // 
	".ex":             map[int]string{2: "\ue62d"}, // 
	".exe":            map[int]string{2: "\uf17a"}, // 
	".exs":            map[int]string{2: "\ue62d"}, // 
	".fish":           map[int]string{2: "\uf489"}, // 
	".flac":           map[int]string{2: "\uf001"}, // 
	".flv":            map[int]string{2: "\uf03d"}, // 
	".font":           map[int]string{2: "\uf031"}, // 
	".fs":             map[int]string{2: "\ue7a7"}, // 
	".fsi":            map[int]string{2: "\ue7a7"}, // 
	".fsx":            map[int]string{2: "\ue7a7"}, // 
	".gdoc":           map[int]string{2: "\uf1c2"}, // 
	".gem":            map[int]string{2: "\ue21e"}, // 
	".gemfile":        map[int]string{2: "\ue21e"}, // 
	".gemspec":        map[int]string{2: "\ue21e"}, // 
	".gform":          map[int]string{2: "\uf298"}, // 
	".gif":            map[int]string{2: "\uf1c5"}, // 
	".git":            map[int]string{2: "\uf1d3"}, // 
	".gitattributes":  map[int]string{2: "\uf1d3"}, // 
	".gitignore":      map[int]string{2: "\uf1d3"}, // 
	".gitmodules":     map[int]string{2: "\uf1d3"}, // 
	".go":             map[int]string{2: "\ue626"}, // 
	".gradle":         map[int]string{2: "\ue256"}, // 
	".groovy":         map[int]string{2: "\ue775"}, // 
	".gsheet":         map[int]string{2: "\uf1c3"}, // 
	".gslides":        map[int]string{2: "\uf1c4"}, // 
	".guardfile":      map[int]string{2: "\ue21e"}, // 
	".gz":             map[int]string{2: "\uf410"}, // 
	".h":              map[int]string{2: "\uf0fd"}, // 
	".hbs":            map[int]string{2: "\ue60f"}, // 
	".hpp":            map[int]string{2: "\uf0fd"}, // 
	".hs":             map[int]string{2: "\ue777"}, // 
	".htm":            map[int]string{2: "\uf13b"}, // 
	".html":           map[int]string{2: "\uf13b"}, // 
	".hxx":            map[int]string{2: "\uf0fd"}, // 
	".ico":            map[int]string{2: "\uf1c5"}, // 
	".image":          map[int]string{2: "\uf1c5"}, // 
	".iml":            map[int]string{2: "\ue7b5"}, // 
	".ini":            map[int]string{2: "\uf17a"}, // 
	".ipynb":          map[int]string{2: "\ue606"}, // 
	".iso":            map[int]string{2: "\ue271"}, // 
	".j2c":            map[int]string{2: "\uf1c5"}, // 
	".j2k":            map[int]string{2: "\uf1c5"}, // 
	".jad":            map[int]string{2: "\ue256"}, // 
	".jar":            map[int]string{2: "\ue256"}, // 
	".java":           map[int]string{2: "\ue256"}, // 
	".jfi":            map[int]string{2: "\uf1c5"}, // 
	".jfif":           map[int]string{2: "\uf1c5"}, // 
	".jif":            map[int]string{2: "\uf1c5"}, // 
	".jl":             map[int]string{2: "\ue624"}, // 
	".jmd":            map[int]string{2: "\uf48a"}, // 
	".jp2":            map[int]string{2: "\uf1c5"}, // 
	".jpe":            map[int]string{2: "\uf1c5"}, // 
	".jpeg":           map[int]string{2: "\uf1c5"}, // 
	".jpg":            map[int]string{2: "\uf1c5"}, // 
	".jpx":            map[int]string{2: "\uf1c5"}, // 
	".js":             map[int]string{2: "\ue74e"}, // 
	".json":           map[int]string{2: "\ue60b"}, // 
	".jsx":            map[int]string{2: "\ue7ba"}, // 
	".jxl":            map[int]string{2: "\uf1c5"}, // 
	".ksh":            map[int]string{2: "\uf489"}, // 
	".kt":             map[int]string{2: "\ue634"}, // 
	".kts":            map[int]string{2: "\ue634"}, // 
	".latex":          map[int]string{2: "\uf034"}, // 
	".less":           map[int]string{2: "\ue758"}, // 
	".lhs":            map[int]string{2: "\ue777"}, // 
    ".license":        map[int]string{2: "\uf718", 3: "\U000f0219"}, //   󰈙
	".localized":      map[int]string{2: "\uf179"}, // 
	".lock":           map[int]string{2: "\uf023"}, // 
	".log":            map[int]string{2: "\uf18d"}, // 
	".lua":            map[int]string{2: "\ue620"}, // 
	".lz":             map[int]string{2: "\uf410"}, // 
	".lz4":            map[int]string{2: "\uf410"}, // 
	".lzh":            map[int]string{2: "\uf410"}, // 
	".lzma":           map[int]string{2: "\uf410"}, // 
	".lzo":            map[int]string{2: "\uf410"}, // 
	".m":              map[int]string{2: "\ue61e"}, // 
	".mm":             map[int]string{2: "\ue61d"}, // 
	".m4a":            map[int]string{2: "\uf001"}, // 
	".markdown":       map[int]string{2: "\uf48a"}, // 
	".md":             map[int]string{2: "\uf48a"}, // 
	".mjs":            map[int]string{2: "\ue74e"}, // 
	".mk":             map[int]string{2: "\uf489"}, // 
	".mkd":            map[int]string{2: "\uf48a"}, // 
	".mkv":            map[int]string{2: "\uf03d"}, // 
	".mobi":           map[int]string{2: "\ue28b"}, // 
	".mov":            map[int]string{2: "\uf03d"}, // 
	".mp3":            map[int]string{2: "\uf001"}, // 
	".mp4":            map[int]string{2: "\uf03d"}, // 
	".msi":            map[int]string{2: "\ue70f"}, // 
	".mustache":       map[int]string{2: "\ue60f"}, // 
	".nix":            map[int]string{2: "\uf313"}, // 
	".node":           map[int]string{2: "\uf898"}, // 
	".npmignore":      map[int]string{2: "\ue71e"}, // 
	".odp":            map[int]string{2: "\uf1c4"}, // 
	".ods":            map[int]string{2: "\uf1c3"}, // 
	".odt":            map[int]string{2: "\uf1c2"}, // 
	".ogg":            map[int]string{2: "\uf001"}, // 
	".ogv":            map[int]string{2: "\uf03d"}, // 
	".otf":            map[int]string{2: "\uf031"}, // 
	".part":           map[int]string{2: "\uf43a"}, // 
	".patch":          map[int]string{2: "\uf440"}, // 
	".pdf":            map[int]string{2: "\uf1c1"}, // 
	".php":            map[int]string{2: "\ue73d"}, // 
	".pl":             map[int]string{2: "\ue769"}, // 
	".png":            map[int]string{2: "\uf1c5"}, // 
	".ppt":            map[int]string{2: "\uf1c4"}, // 
	".pptx":           map[int]string{2: "\uf1c4"}, // 
	".procfile":       map[int]string{2: "\ue21e"}, // 
	".properties":     map[int]string{2: "\ue60b"}, // 
	".ps1":            map[int]string{2: "\uf489"}, // 
	".psd":            map[int]string{2: "\ue7b8"}, // 
	".pxm":            map[int]string{2: "\uf1c5"}, // 
	".py":             map[int]string{2: "\ue606"}, // 
	".pyc":            map[int]string{2: "\ue606"}, // 
	".r":              map[int]string{2: "\uf25d"}, // 
	".rakefile":       map[int]string{2: "\ue21e"}, // 
	".rar":            map[int]string{2: "\uf410"}, // 
	".razor":          map[int]string{2: "\uf1fa"}, // 
	".rb":             map[int]string{2: "\ue21e"}, // 
	".rdata":          map[int]string{2: "\uf25d"}, // 
	".rdb":            map[int]string{2: "\ue76d"}, // 
	".rdoc":           map[int]string{2: "\uf48a"}, // 
	".rds":            map[int]string{2: "\uf25d"}, // 
	".readme":         map[int]string{2: "\uf48a"}, // 
	".rlib":           map[int]string{2: "\ue7a8"}, // 
	".rmd":            map[int]string{2: "\uf48a"}, // 
	".rpm":            map[int]string{2: "\ue7bb"}, // 
	".rs":             map[int]string{2: "\ue7a8"}, // 
	".rspec":          map[int]string{2: "\ue21e"}, // 
	".rspec_parallel": map[int]string{2: "\ue21e"}, // 
	".rspec_status":   map[int]string{2: "\ue21e"}, // 
	".rss":            map[int]string{2: "\uf09e"}, // 
    ".rtf":            map[int]string{2: "\uf718", 3: "\U000f0219"}, //    󰈙
	".ru":             map[int]string{2: "\ue21e"}, // 
	".rubydoc":        map[int]string{2: "\ue73b"}, // 
	".sass":           map[int]string{2: "\ue603"}, // 
	".scala":          map[int]string{2: "\ue737"}, // 
	".scss":           map[int]string{2: "\ue749"}, // 
	".sh":             map[int]string{2: "\uf489"}, // 
	".shell":          map[int]string{2: "\uf489"}, // 
	".slim":           map[int]string{2: "\ue73b"}, // 
	".sln":            map[int]string{2: "\ue70c"}, // 
	".so":             map[int]string{2: "\uf17c"}, // 
	".sql":            map[int]string{2: "\uf1c0"}, // 
	".sqlite3":        map[int]string{2: "\ue7c4"}, // 
	".sty":            map[int]string{2: "\uf034"}, // 
	".styl":           map[int]string{2: "\ue600"}, // 
	".stylus":         map[int]string{2: "\ue600"}, // 
	".svg":            map[int]string{2: "\uf1c5"}, // 
	".swift":          map[int]string{2: "\ue755"}, // 
	".tar":            map[int]string{2: "\uf410"}, // 
	".taz":            map[int]string{2: "\uf410"}, // 
	".tbz":            map[int]string{2: "\uf410"}, // 
	".tbz2":           map[int]string{2: "\uf410"}, // 
	".tex":            map[int]string{2: "\uf034"}, // 
	".tgz":            map[int]string{2: "\uf410"}, // 
	".tiff":           map[int]string{2: "\uf1c5"}, // 
	".tlz":            map[int]string{2: "\uf410"}, // 
	".toml":           map[int]string{2: "\ue615"}, // 
	".torrent":        map[int]string{2: "\ue275"}, // 
	".ts":             map[int]string{2: "\ue628"}, // 
	".tsv":            map[int]string{2: "\uf1c3"}, // 
	".tsx":            map[int]string{2: "\ue7ba"}, // 
	".ttf":            map[int]string{2: "\uf031"}, // 
	".twig":           map[int]string{2: "\ue61c"}, // 
	".txt":            map[int]string{2: "\uf15c"}, // 
	".txz":            map[int]string{2: "\uf410"}, // 
	".tz":             map[int]string{2: "\uf410"}, // 
	".tzo":            map[int]string{2: "\uf410"}, // 
	".video":          map[int]string{2: "\uf03d"}, // 
	".vim":            map[int]string{2: "\ue62b"}, // 
	".vue":            map[int]string{2: "\ufd42"}, // ﵂
	".war":            map[int]string{2: "\ue256"}, // 
	".wav":            map[int]string{2: "\uf001"}, // 
	".webm":           map[int]string{2: "\uf03d"}, // 
	".webp":           map[int]string{2: "\uf1c5"}, // 
	".windows":        map[int]string{2: "\uf17a"}, // 
	".woff":           map[int]string{2: "\uf031"}, // 
	".woff2":          map[int]string{2: "\uf031"}, // 
	".xhtml":          map[int]string{2: "\uf13b"}, // 
	".xls":            map[int]string{2: "\uf1c3"}, // 
	".xlsx":           map[int]string{2: "\uf1c3"}, // 
	".xml":            map[int]string{2: "\uf121"}, // 
	".xul":            map[int]string{2: "\uf121"}, // 
	".xz":             map[int]string{2: "\uf410"}, // 
	".yaml":           map[int]string{2: "\uf481"}, // 
	".yml":            map[int]string{2: "\uf481"}, // 
	".zip":            map[int]string{2: "\uf410"}, // 
	".zsh":            map[int]string{2: "\uf489"}, // 
	".zsh-theme":      map[int]string{2: "\uf489"}, // 
	".zshrc":          map[int]string{2: "\uf489"}, // 
	".zst":            map[int]string{2: "\uf410"}, // 
}

func IconForFile(name string, isSubmodule bool, isDirectory bool) string {
	base := filepath.Base(name)
	if elm, ok := nameIconMap[base]; ok {
        if icon, ok := elm[GetNerdFontsVersion()]; ok {
		    return icon
        }
        return elm[2]
	}

	ext := filepath.Ext(name)
	if elm, ok := extIconMap[ext]; ok {
        if icon, ok := elm[GetNerdFontsVersion()]; ok {
		    return icon
        }
        return elm[2]
	}

	if isSubmodule {
		return DEFAULT_SUBMODULE_ICON
	} else if isDirectory {
		return DEFAULT_DIRECTORY_ICON
	}
	return DEFAULT_FILE_ICON
}
