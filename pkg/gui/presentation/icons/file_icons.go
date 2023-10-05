package icons

import (
	"path/filepath"
)

type iconProperties struct {
    icon string
    color uint8 
}

// https://github.com/ogham/exa/blob/master/src/output/icons.rs
var (
	DEFAULT_FILE_ICON      = iconProperties {"\uf15b", 239} // 
	DEFAULT_SUBMODULE_ICON = iconProperties {"\uf1d3", 239} // 
	DEFAULT_DIRECTORY_ICON = iconProperties {"\uf114", 239} // 
)

var nameIconMap = map[string]iconProperties{
	".Trash":            iconProperties {"\uf1f8", 239},// 
	".atom":             iconProperties { "\ue764",239}, // 
	".bashprofile":      iconProperties { "\ue615",239}, // 
	".bashrc":           iconProperties { "\uf489",239}, // 
	".idea":             iconProperties { "\ue7b5",239}, // 
	".git":              iconProperties { "\uf1d3",239}, // 
	".gitattributes":    iconProperties { "\uf1d3",239}, // 
	".gitconfig":        iconProperties { "\uf1d3",239}, // 
	".github":           iconProperties { "\uf408",239}, // 
	".gitignore":        iconProperties { "\uf1d3",239}, // 
	".gitmodules":       iconProperties { "\uf1d3",239}, // 
	".rvm":              iconProperties { "\ue21e",239}, // 
	".vimrc":            iconProperties { "\ue62b",239}, // 
	".vscode":           iconProperties { "\ue70c",239}, // 
	".zshrc":            iconProperties { "\uf489",239}, // 
	"Cargo.lock":        iconProperties { "\ue7a8",239}, // 
	"Cargo.toml":        iconProperties {"\ue7a8", 239},// 
	"bin":               iconProperties { "\ue5fc",239}, // 
	"config":            iconProperties { "\ue5fc",239}, // 
    "docker-compose.yml":iconProperties { "\uf308",239}, // 
	"Dockerfile":        iconProperties { "\uf308",239}, // 
	"ds_store":          iconProperties { "\uf179",239}, // 
	"gitignore_global":  iconProperties { "\uf1d3",239}, // 
	"go.mod":            iconProperties { "\ue626",239}, // 
	"go.sum":            iconProperties { "\ue626",239}, // 
	"gradle":            iconProperties { "\ue256",239}, // 
	"gruntfile.coffee":  iconProperties {"\ue611", 239},// 
	"gruntfile.js":      iconProperties { "\ue611",239}, // 
	"gruntfile.ls":      iconProperties { "\ue611",239}, // 
	"gulpfile.coffee":   iconProperties { "\ue610",239}, // 
	"gulpfile.js":       iconProperties { "\ue610",239}, // 
	"gulpfile.ls":       iconProperties { "\ue610",239}, // 
	"hidden":            iconProperties { "\uf023",239}, // 
	"include":           iconProperties { "\ue5fc",239}, // 
	"lib":               iconProperties { "\uf121",239}, // 
	"localized":         iconProperties { "\uf179",239}, // 
	"Makefile":          iconProperties { "\uf489",239}, // 
	"node_modules":      iconProperties { "\ue718",239}, // 
	"npmignore":         iconProperties { "\ue71e",239}, // 
	"PKGBUILD":          iconProperties { "\uf303",239}, // 
	"rubydoc":           iconProperties { "\ue73b",239}, // 
	"yarn.lock":         iconProperties { "\ue718",239}, // 
}

var extIconMap = map[string]iconProperties{
	".ai":            iconProperties {"\ue7b4", 239},     // 
	".android":       iconProperties {"\ue70e", 239},     // 
	".apk":           iconProperties {"\ue70e", 239},     // 
	".apple":         iconProperties {"\uf179", 239},     // 
	".avi":           iconProperties {"\uf03d", 239},     // 
	".avif":          iconProperties {"\uf1c5", 239},     // 
	".avro":          iconProperties {"\ue60b", 239},     // 
	".awk":           iconProperties {"\uf489", 239},     // 
	".bash":          iconProperties {"\uf489", 239},     // 
	".bash_history":  iconProperties {"\uf489", 239},     // 
	".bash_profile":  iconProperties {"\uf489", 239},     // 
	".bashrc":        iconProperties {"\uf489", 239},     // 
	".bat":           iconProperties {"\uf17a", 239},     // 
	".bats":          iconProperties {"\uf489", 239},     // 
	".bmp":           iconProperties {"\uf1c5", 239},     // 
	".bz":            iconProperties {"\uf410", 239},     // 
	".bz2":           iconProperties {"\uf410", 239},     // 
	".c":             iconProperties {"\ue61e", 239},     // 
	".c++":           iconProperties {"\ue61d", 239},     // 
	".cab":           iconProperties {"\ue70f", 239},     // 
	".cc":            iconProperties {"\ue61d", 239},     // 
	".cfg":           iconProperties {"\ue615", 239},     // 
	".class":         iconProperties {"\ue256", 239},     // 
	".clj":           iconProperties {"\ue768", 239},     // 
	".cljs":          iconProperties {"\ue76a", 239},     // 
	".cls":           iconProperties {"\uf034", 239},     // 
	".cmd":           iconProperties {"\ue70f", 239},     // 
	".coffee":        iconProperties {"\uf0f4", 239},     // 
	".conf":          iconProperties {"\ue615", 239},     // 
	".cp":            iconProperties {"\ue61d", 239},     // 
	".cpio":          iconProperties {"\uf410", 239},     // 
	".cpp":           iconProperties {"\ue61d", 239},     // 
    ".cs":            iconProperties {"\U000f031b", 239}, // 󰌛
	".csh":           iconProperties {"\uf489", 239},     // 
	".cshtml":        iconProperties {"\uf1fa", 239},     // 
    ".csproj":        iconProperties {"\U000f031b", 239}, // 󰌛
	".css":           iconProperties {"\ue749", 239},     // 
	".csv":           iconProperties {"\uf1c3", 239},     // 
    ".csx":           iconProperties {"\U000f031b", 239}, // 󰌛
	".cxx":           iconProperties {"\ue61d", 239},     // 
	".d":             iconProperties {"\ue7af", 239},     // 
	".dart":          iconProperties {"\ue798", 239},     // 
	".db":            iconProperties {"\uf1c0", 239},     // 
	".deb":           iconProperties {"\ue77d", 239},     // 
	".diff":          iconProperties {"\uf440", 239},     // 
	".djvu":          iconProperties {"\uf02d", 239},     // 
	".dll":           iconProperties {"\ue70f", 239},     // 
	".doc":           iconProperties {"\uf1c2", 239},     // 
	".docx":          iconProperties {"\uf1c2", 239},     // 
	".ds_store":      iconProperties {"\uf179", 239},     // 
	".DS_store":      iconProperties {"\uf179", 239},     // 
	".dump":          iconProperties {"\uf1c0", 239},     // 
	".ebook":         iconProperties {"\ue28b", 239},     // 
	".ebuild":        iconProperties {"\uf30d", 239},     // 
	".editorconfig":  iconProperties {"\ue615", 239},     // 
	".ejs":           iconProperties {"\ue618", 239},     // 
	".elm":           iconProperties {"\ue62c", 239},     // 
	".env":           iconProperties {"\uf462", 239},     // 
	".eot":           iconProperties {"\uf031", 239},     // 
	".epub":          iconProperties {"\ue28a", 239},     // 
	".erb":           iconProperties {"\ue73b", 239},     // 
	".erl":           iconProperties {"\ue7b1", 239},     // 
	".ex":            iconProperties {"\ue62d", 239},     // 
	".exe":           iconProperties {"\uf17a", 239},     // 
	".exs":           iconProperties {"\ue62d", 239},     // 
	".fish":          iconProperties {"\uf489", 239},     // 
	".flac":          iconProperties {"\uf001", 239},     // 
	".flv":           iconProperties {"\uf03d", 239},     // 
	".font":          iconProperties {"\uf031", 239},     // 
	".fs":            iconProperties {"\ue7a7", 239},     // 
	".fsi":           iconProperties {"\ue7a7", 239},     // 
	".fsx":           iconProperties {"\ue7a7", 239},     // 
	".gdoc":          iconProperties {"\uf1c2", 239},     // 
	".gem":           iconProperties {"\ue21e", 239},     // 
	".gemfile":       iconProperties {"\ue21e", 239},     // 
	".gemspec":       iconProperties {"\ue21e", 239},     // 
	".gform":         iconProperties {"\uf298", 239},     // 
	".gif":           iconProperties {"\uf1c5", 239},     // 
	".git":           iconProperties {"\uf1d3", 239},     // 
	".gitattributes": iconProperties {"\uf1d3", 239},     // 
	".gitignore":     iconProperties {"\uf1d3", 239},     // 
	".gitmodules":    iconProperties {"\uf1d3", 239},     // 
	".go":            iconProperties {"\ue626", 239},     // 
	".gradle":        iconProperties {"\ue256", 239},     // 
	".groovy":        iconProperties {"\ue775", 239},     // 
	".gsheet":        iconProperties {"\uf1c3", 239},     // 
	".gslides":       iconProperties {"\uf1c4", 239},     // 
	".guardfile":     iconProperties {"\ue21e", 239},     // 
	".gz":            iconProperties {"\uf410", 239},     // 
	".h":             iconProperties {"\uf0fd", 239},     // 
	".hbs":           iconProperties {"\ue60f", 239},     // 
	".hpp":           iconProperties {"\uf0fd", 239},     // 
	".hs":            iconProperties {"\ue777", 239},     // 
	".htm":           iconProperties {"\uf13b", 239},     // 
	".html":          iconProperties {"\uf13b", 239},     // 
	".hxx":           iconProperties {"\uf0fd", 239},     // 
	".ico":           iconProperties {"\uf1c5", 239},     // 
	".image":         iconProperties {"\uf1c5", 239},     // 
	".iml":           iconProperties {"\ue7b5", 239},     // 
	".ini":           iconProperties {"\uf17a", 239},     // 
	".ipynb":         iconProperties {"\ue606", 239},     // 
	".iso":           iconProperties {"\ue271", 239},     // 
	".j2c":           iconProperties {"\uf1c5", 239},     // 
	".j2k":           iconProperties {"\uf1c5", 239},     // 
	".jad":           iconProperties {"\ue256", 239},     // 
	".jar":           iconProperties {"\ue256", 239},     // 
	".java":          iconProperties {"\ue256", 239},     // 
	".jfi":           iconProperties {"\uf1c5", 239},     // 
	".jfif":          iconProperties {"\uf1c5", 239},     // 
	".jif":           iconProperties {"\uf1c5", 239},     // 
	".jl":            iconProperties {"\ue624", 239},     // 
	".jmd":           iconProperties {"\uf48a", 239},     // 
	".jp2":           iconProperties {"\uf1c5", 239},     // 
	".jpe":           iconProperties {"\uf1c5", 239},     // 
	".jpeg":          iconProperties {"\uf1c5", 239},     // 
	".jpg":           iconProperties {"\uf1c5", 239},     // 
	".jpx":           iconProperties {"\uf1c5", 239},     // 
	".js":            iconProperties {"\ue74e", 239},     // 
	".json":          iconProperties {"\ue60b", 239},     // 
	".jsx":           iconProperties {"\ue7ba", 239},     // 
	".jxl":           iconProperties {"\uf1c5", 239},     // 
	".ksh":           iconProperties {"\uf489", 239},     // 
	".kt":            iconProperties {"\ue634", 239},     // 
	".kts":           iconProperties {"\ue634", 239},     // 
	".latex":         iconProperties {"\uf034", 239},     // 
	".less":          iconProperties {"\ue758", 239},     // 
	".lhs":           iconProperties {"\ue777", 239},     // 
	".license":       iconProperties {"\U000f0219", 239}, // 󰈙
	".localized":     iconProperties {"\uf179", 239},     // 
	".lock":          iconProperties {"\uf023", 239},     // 
	".log":           iconProperties {"\uf18d", 239},     // 
	".lua":           iconProperties {"\ue620", 239},     // 
	".lz":            iconProperties {"\uf410", 239},     // 
	".lz4":           iconProperties {"\uf410", 239},     // 
	".lzh":           iconProperties {"\uf410", 239},     // 
	".lzma":          iconProperties {"\uf410", 239},     // 
	".lzo":           iconProperties {"\uf410", 239},     // 
	".m":             iconProperties {"\ue61e", 239},     // 
	".mm":            iconProperties {"\ue61d", 239},     // 
	".m4a":           iconProperties {"\uf001", 239},     // 
	".markdown":      iconProperties {"\uf48a", 239},     // 
	".md":            iconProperties {"\uf48a", 239},     // 
	".mdx":           iconProperties {"\uf48a", 239},     // 
	".mjs":           iconProperties {"\ue74e", 239},     // 
	".mk":            iconProperties {"\uf489", 239},     // 
	".mkd":           iconProperties {"\uf48a", 239},     // 
	".mkv":           iconProperties {"\uf03d", 239},     // 
	".mobi":          iconProperties {"\ue28b", 239},     // 
	".mov":           iconProperties {"\uf03d", 239},     // 
	".mp3":           iconProperties {"\uf001", 239},     // 
	".mp4":           iconProperties {"\uf03d", 239},     // 
	".msi":           iconProperties {"\ue70f", 239},     // 
	".mustache":      iconProperties {"\ue60f", 239},     // 
	".nix":           iconProperties {"\uf313", 239},     // 
	".node":          iconProperties {"\U000f0399", 239}, // 󰎙
	".npmignore":     iconProperties {"\ue71e", 239},     // 
	".odp":           iconProperties {"\uf1c4", 239},     // 
	".ods":           iconProperties {"\uf1c3", 239},     // 
	".odt":           iconProperties {"\uf1c2", 239},     // 
	".ogg":           iconProperties {"\uf001", 239},     // 
	".ogv":           iconProperties {"\uf03d", 239},     // 
	".otf":           iconProperties {"\uf031", 239},     // 
	".part":          iconProperties {"\uf43a", 239},     // 
	".patch":         iconProperties {"\uf440", 239},     // 
	".pdf":           iconProperties {"\uf1c1", 239},     // 
	".php":           iconProperties {"\ue73d", 239},     // 
	".pl":            iconProperties {"\ue769", 239},     // 
	".png":           iconProperties {"\uf1c5", 239},     // 
	".ppt":           iconProperties {"\uf1c4", 239},     // 
	".pptx":          iconProperties {"\uf1c4", 239},     // 
	".procfile":      iconProperties {"\ue21e", 239},     // 
	".properties":    iconProperties {"\ue60b", 239},     // 
	".ps1":           iconProperties {"\uf489", 239},     // 
	".psd":           iconProperties {"\ue7b8", 239},     // 
	".pxm":           iconProperties {"\uf1c5", 239},     // 
	".py":            iconProperties {"\ue606", 239},     // 
	".pyc":           iconProperties {"\ue606", 239},     // 
	".r":             iconProperties {"\uf25d", 239},     // 
	".rakefile":      iconProperties {"\ue21e", 239},     // 
	".rar":           iconProperties {"\uf410", 239},     // 
	".razor":         iconProperties {"\uf1fa", 239},     // 
	".rb":            iconProperties {"\ue21e", 239},     // 
	".rdata":         iconProperties {"\uf25d", 239},     // 
	".rdb":           iconProperties {"\ue76d", 239},     // 
	".rdoc":          iconProperties {"\uf48a", 239},     // 
	".rds":           iconProperties {"\uf25d", 239},     // 
	".readme":        iconProperties {"\uf48a", 239},     // 
	".rlib":          iconProperties {"\ue7a8", 239},     // 
	".rmd":           iconProperties {"\uf48a", 239},     // 
	".rpm":           iconProperties {"\ue7bb", 239},     // 
	".rs":            iconProperties {"\ue7a8", 239},     // 
	".rspec":         iconProperties {"\ue21e", 239},     // 
	".rspec_parallel":iconProperties {"\ue21e", 239},     // 
	".rspec_status":  iconProperties {"\ue21e", 239},     // 
	".rss":           iconProperties {"\uf09e", 239},     // 
    ".rtf":           iconProperties {"\U000f0219", 239}, // 󰈙
	".ru":            iconProperties {"\ue21e", 239},     // 
	".rubydoc":       iconProperties {"\ue73b", 239},     // 
	".sass":          iconProperties {"\ue603", 239},     // 
	".scala":         iconProperties {"\ue737", 239},     // 
	".scss":          iconProperties {"\ue749", 239},     // 
	".sh":            iconProperties {"\uf489", 239},     // 
	".shell":         iconProperties {"\uf489", 239},     // 
	".slim":          iconProperties {"\ue73b", 239},     // 
	".sln":           iconProperties {"\ue70c", 239},     // 
	".so":            iconProperties {"\uf17c", 239},     // 
	".sql":           iconProperties {"\uf1c0", 239},     // 
	".sqlite3":       iconProperties {"\ue7c4", 239},     // 
	".sty":           iconProperties {"\uf034", 239},     // 
	".styl":          iconProperties {"\ue600", 239},     // 
	".stylus":        iconProperties {"\ue600", 239},     // 
	".svelte":        iconProperties {"\ue697", 239},     // 
	".svg":           iconProperties {"\uf1c5", 239},     // 
	".swift":         iconProperties {"\ue755", 239},     // 
	".tar":           iconProperties {"\uf410", 239},     // 
	".taz":           iconProperties {"\uf410", 239},     // 
	".tbz":           iconProperties {"\uf410", 239},     // 
	".tbz2":          iconProperties {"\uf410", 239},     // 
	".tex":           iconProperties {"\uf034", 239},     // 
	".tgz":           iconProperties {"\uf410", 239},     // 
	".tiff":          iconProperties {"\uf1c5", 239},     // 
	".tlz":           iconProperties {"\uf410", 239},     // 
	".toml":          iconProperties {"\ue615", 239},     // 
	".torrent":       iconProperties {"\ue275", 239},     // 
	".ts":            iconProperties {"\ue628", 239},     // 
	".tsv":           iconProperties {"\uf1c3", 239},     // 
	".tsx":           iconProperties {"\ue7ba", 239},     // 
	".ttf":           iconProperties {"\uf031", 239},     // 
	".twig":          iconProperties {"\ue61c", 239},     // 
	".txt":           iconProperties {"\uf15c", 239},     // 
	".txz":           iconProperties {"\uf410", 239},     // 
	".tz":            iconProperties {"\uf410", 239},     // 
	".tzo":           iconProperties {"\uf410", 239},     // 
	".video":         iconProperties {"\uf03d", 239},     // 
	".vim":           iconProperties {"\ue62b", 239},     // 
	".vue":           iconProperties {"\U000f0844", 239}, // 󰡄
	".war":           iconProperties {"\ue256", 239},     // 
	".wav":           iconProperties {"\uf001", 239},     // 
	".webm":          iconProperties {"\uf03d", 239},     // 
	".webp":          iconProperties {"\uf1c5", 239},     // 
	".windows":       iconProperties {"\uf17a", 239},     // 
	".woff":          iconProperties {"\uf031", 239},     // 
	".woff2":         iconProperties {"\uf031", 239},     // 
	".xhtml":         iconProperties {"\uf13b", 239},     // 
	".xls":           iconProperties {"\uf1c3", 239},     // 
	".xlsx":          iconProperties {"\uf1c3", 239},     // 
	".xml":           iconProperties {"\uf121", 239},     // 
	".xul":           iconProperties {"\uf121", 239},     // 
	".xz":            iconProperties {"\uf410", 239},     // 
	".yaml":          iconProperties {"\uf481", 239},     // 
	".yml":           iconProperties {"\uf481", 239},     // 
	".zip":           iconProperties {"\uf410", 239},     // 
	".zsh":           iconProperties {"\uf489", 239},     // 
	".zsh-theme":     iconProperties {"\uf489", 239},     // 
	".zshrc":         iconProperties {"\uf489", 239},     // 
	".zst":           iconProperties {"\uf410", 239},     // 
}

func patchFileIconsForNerdFontsV2() {
	extIconMap[".cs"] = iconProperties {"\uf81a", 239}      // 
	extIconMap[".csproj"] = iconProperties {"\uf81a", 239}  // 
	extIconMap[".csx"] = iconProperties {"\uf81a", 239}     // 
	extIconMap[".license"] = iconProperties {"\uf718", 239} // 
	extIconMap[".node"] = iconProperties {"\uf898", 239}    // 
	extIconMap[".rtf"] = iconProperties {"\uf718", 239}     // 
	extIconMap[".vue"] = iconProperties {"\ufd42", 239}     // ﵂
}

func IconForFile(name string, isSubmodule bool, isLinkedWorktree bool, isDirectory bool) iconProperties {
	base := filepath.Base(name)
	if icon, ok := nameIconMap[base]; ok {
		return icon
	}

	ext := filepath.Ext(name)
	if icon, ok := extIconMap[ext]; ok {
		return icon
	}

	if isSubmodule {
		return DEFAULT_SUBMODULE_ICON
	} else if isLinkedWorktree {
		return LINKED_WORKTREE_ICON
	} else if isDirectory {
		return DEFAULT_DIRECTORY_ICON
	}
	return DEFAULT_FILE_ICON
}
