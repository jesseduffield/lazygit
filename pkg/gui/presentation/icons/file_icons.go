package icons

import (
	"path/filepath"
)

// https://github.com/ogham/exa/blob/master/src/output/icons.rs
var (
	DEFAULT_FILE_ICON      = IconProperties{"\uf15b", 241} // 
	DEFAULT_SUBMODULE_ICON = IconProperties{"\uf1d3", 202} // 
	DEFAULT_DIRECTORY_ICON = IconProperties{"\uf07b", 241} // 
)

var nameIconMap = map[string]IconProperties{
	".Trash":             {"\uf1f8", 241}, // 
	".atom":              {"\ue764", 241}, // 
	".bashprofile":       {"\ue615", 113}, // 
	".bashrc":            {"\ue795", 113}, // 
	".idea":              {"\ue7b5", 241}, // 
	".git":               {"\uf1d3", 202}, // 
	".gitattributes":     {"\uf1d3", 202}, // 
	".gitconfig":         {"\uf1d3", 202}, // 
	".github":            {"\uf408", 241}, // 
	".gitignore":         {"\uf1d3", 202}, // 
	".gitmodules":        {"\uf1d3", 202}, // 
	".rvm":               {"\ue21e", 160}, // 
	".vimrc":             {"\ue62b", 28},  // 
	".vscode":            {"\ue70c", 39},  // 
	".zshrc":             {"\ue795", 113}, // 
	"Cargo.lock":         {"\ue7a8", 216}, // 
	"Cargo.toml":         {"\ue7a8", 216}, // 
	"bin":                {"\ue5fc", 241}, // 
	"config":             {"\ue5fc", 241}, // 
	"docker-compose.yml": {"\uf308", 68},  // 
	"Dockerfile":         {"\uf308", 68},  // 
	"ds_store":           {"\uf179", 15},  // 
	"gitignore_global":   {"\uf1d3", 202}, // 
	"go.mod":             {"\ue627", 74},  // 
	"go.sum":             {"\ue627", 74},  // 
	"gradle":             {"\ue256", 168}, // 
	"gruntfile.coffee":   {"\ue611", 166}, // 
	"gruntfile.js":       {"\ue611", 166}, // 
	"gruntfile.ls":       {"\ue611", 166}, // 
	"gulpfile.coffee":    {"\ue610", 167}, // 
	"gulpfile.js":        {"\ue610", 167}, // 
	"gulpfile.ls":        {"\ue610", 168}, // 
	"hidden":             {"\uf023", 241}, // 
	"include":            {"\ue5fc", 241}, // 
	"lib":                {"\uf121", 241}, // 
	"localized":          {"\uf179", 15},  // 
	"Makefile":           {"\ue975", 241}, // 
	"node_modules":       {"\ue718", 197}, // 
	"npmignore":          {"\ue71e", 197}, // 
	"PKGBUILD":           {"\uf303", 38},  // 
	"rubydoc":            {"\ue73b", 160}, // 
	"yarn.lock":          {"\ue6a7", 74},  // 
}

var extIconMap = map[string]IconProperties{
	".ai":             {"\ue7b4", 185},     // 
	".android":        {"\ue70e", 70},      // 
	".apk":            {"\ue70e", 70},      // 
	".apple":          {"\uf179", 15},      // 
	".avi":            {"\uf03d", 241},     // 
	".avif":           {"\uf1c5", 241},     // 
	".avro":           {"\ue60b", 241},     // 
	".awk":            {"\ue795", 241},     // 
	".bash":           {"\ue795", 113},     // 
	".bash_history":   {"\ue795", 113},     // 
	".bash_profile":   {"\ue795", 113},     // 
	".bashrc":         {"\ue795", 113},     // 
	".bat":            {"\uf17a", 81},      // 
	".bats":           {"\ue795", 241},     // 
	".bmp":            {"\uf1c5", 149},     // 
	".bz":             {"\uf410", 239},     // 
	".bz2":            {"\uf410", 239},     // 
	".c":              {"\ue61e", 111},     // 
	".c++":            {"\ue61d", 204},     // 
	".cab":            {"\ue70f", 241},     // 
	".cc":             {"\ue61d", 204},     // 
	".cfg":            {"\ue615", 255},     // 
	".class":          {"\ue256", 168},     // 
	".clj":            {"\ue768", 113},     // 
	".cljs":           {"\ue76a", 74},      // 
	".cls":            {"\uf034", 239},     // 
	".cmd":            {"\ue70f", 239},     // 
	".coffee":         {"\uf0f4", 185},     // 
	".conf":           {"\ue615", 66},      // 
	".cp":             {"\ue61d", 74},      // 
	".cpio":           {"\uf410", 239},     // 
	".cpp":            {"\ue61d", 74},      // 
	".cs":             {"\U000f031b", 58},  // 󰌛
	".csh":            {"\ue795", 240},     // 
	".cshtml":         {"\uf1fa", 239},     // 
	".csproj":         {"\U000f031b", 58},  // 󰌛
	".css":            {"\ue749", 75},      // 
	".csv":            {"\uf1c3", 113},     // 
	".csx":            {"\U000f031b", 58},  // 󰌛
	".cxx":            {"\ue61d", 74},      // 
	".d":              {"\ue7af", 28},      // 
	".dart":           {"\ue798", 25},      // 
	".db":             {"\uf1c0", 188},     // 
	".deb":            {"\ue77d", 88},      // 
	".diff":           {"\uf440", 241},     // 
	".djvu":           {"\uf02d", 241},     // 
	".dll":            {"\ue70f", 241},     // 
	".doc":            {"\uf0219", 26},     // 󰈙
	".docx":           {"\uf0219", 26},     // 󰈙
	".ds_store":       {"\uf179", 15},      // 
	".DS_store":       {"\uf179", 15},      // 
	".dump":           {"\uf1c0", 188},     // 
	".ebook":          {"\ue28b", 241},     // 
	".ebuild":         {"\uf30d", 241},     // 
	".editorconfig":   {"\ue615", 241},     // 
	".ejs":            {"\ue618", 185},     // 
	".elm":            {"\ue62c", 74},      // 
	".env":            {"\uf462", 227},     // 
	".eot":            {"\uf031", 241},     // 
	".epub":           {"\ue28a", 241},     // 
	".erb":            {"\ue73b", 160},     // 
	".erl":            {"\ue7b1", 163},     // 
	".ex":             {"\ue62d", 140},     // 
	".exe":            {"\uf17a", 81},      // 
	".exs":            {"\ue62d", 140},     // 
	".fish":           {"\ue795", 249},     // 
	".flac":           {"\uf001", 241},     // 
	".flv":            {"\uf03d", 241},     // 
	".font":           {"\uf031", 241},     // 
	".fs":             {"\ue7a7", 74},      // 
	".fsi":            {"\ue7a7", 74},      // 
	".fsx":            {"\ue7a7", 74},      // 
	".gdoc":           {"\uf1c2", 241},     // 
	".gem":            {"\ue21e", 241},     // 
	".gemfile":        {"\ue21e", 241},     // 
	".gemspec":        {"\ue21e", 241},     // 
	".gform":          {"\uf298", 241},     // 
	".gif":            {"\uf1c5", 140},     // 
	".git":            {"\uf1d3", 202},     // 
	".gitattributes":  {"\uf1d3", 202},     // 
	".gitignore":      {"\uf1d3", 202},     // 
	".gitmodules":     {"\uf1d3", 202},     // 
	".go":             {"\ue627", 74},      // 
	".gradle":         {"\ue256", 168},     // 
	".groovy":         {"\ue775", 241},     // 
	".gsheet":         {"\uf1c3", 241},     // 
	".gslides":        {"\uf1c4", 241},     // 
	".guardfile":      {"\ue21e", 241},     // 
	".gz":             {"\uf410", 241},     // 
	".h":              {"\uf0fd", 140},     // 
	".hbs":            {"\ue60f", 202},     // 
	".hpp":            {"\uf0fd", 140},     // 
	".hs":             {"\ue777", 140},     // 
	".htm":            {"\uf13b", 196},     // 
	".html":           {"\uf13b", 196},     // 
	".hxx":            {"\uf0fd", 140},     // 
	".ico":            {"\uf1c5", 185},     // 
	".image":          {"\uf1c5", 185},     // 
	".iml":            {"\ue7b5", 239},     // 
	".ini":            {"\uf17a", 81},      // 
	".ipynb":          {"\ue606", 214},     // 
	".iso":            {"\ue271", 239},     // 
	".j2c":            {"\uf1c5", 239},     // 
	".j2k":            {"\uf1c5", 239},     // 
	".jad":            {"\ue256", 168},     // 
	".jar":            {"\ue256", 168},     // 
	".java":           {"\ue256", 168},     // 
	".jfi":            {"\uf1c5", 241},     // 
	".jfif":           {"\uf1c5", 241},     // 
	".jif":            {"\uf1c5", 241},     // 
	".jl":             {"\ue624", 241},     // 
	".jmd":            {"\uf48a", 74},      // 
	".jp2":            {"\uf1c5", 241},     // 
	".jpe":            {"\uf1c5", 241},     // 
	".jpeg":           {"\uf1c5", 241},     // 
	".jpg":            {"\uf1c5", 241},     // 
	".jpx":            {"\uf1c5", 241},     // 
	".js":             {"\ue74e", 185},     // 
	".json":           {"\ue60b", 185},     // 
	".jsx":            {"\ue7ba", 45},      // 
	".jxl":            {"\uf1c5", 241},     // 
	".ksh":            {"\ue795", 241},     // 
	".kt":             {"\ue634", 99},      // 
	".kts":            {"\ue634", 99},      // 
	".latex":          {"\uf034", 241},     // 
	".less":           {"\ue758", 54},      // 
	".lhs":            {"\ue777", 140},     // 
	".license":        {"\U000f0219", 185}, // 󰈙
	".localized":      {"\uf179", 15},      // 
	".lock":           {"\uf023", 241},     // 
	".log":            {"\uf18d", 188},     // 
	".lua":            {"\ue620", 74},      // 
	".lz":             {"\uf410", 241},     // 
	".lz4":            {"\uf410", 241},     // 
	".lzh":            {"\uf410", 241},     // 
	".lzma":           {"\uf410", 241},     // 
	".lzo":            {"\uf410", 241},     // 
	".m":              {"\ue61e", 111},     // 
	".mm":             {"\ue61d", 111},     // 
	".m4a":            {"\uf001", 239},     // 
	".markdown":       {"\uf48a", 74},      // 
	".md":             {"\uf48a", 74},      // 
	".mdx":            {"\uf48a", 74},      // 
	".mjs":            {"\ue74e", 185},     // 
	".mk":             {"\ue795", 241},     // 
	".mkd":            {"\uf48a", 74},      // 
	".mkv":            {"\uf03d", 241},     // 
	".mobi":           {"\ue28b", 241},     // 
	".mov":            {"\uf03d", 241},     // 
	".mp3":            {"\uf001", 241},     // 
	".mp4":            {"\uf03d", 241},     // 
	".msi":            {"\ue70f", 241},     // 
	".mustache":       {"\ue60f", 241},     // 
	".nix":            {"\uf313", 241},     // 
	".node":           {"\U000f0399", 197}, // 󰎙
	".npmignore":      {"\ue71e", 197},     // 
	".odp":            {"\uf1c4", 241},     // 
	".ods":            {"\uf1c3", 241},     // 
	".odt":            {"\uf1c2", 241},     // 
	".ogg":            {"\uf001", 241},     // 
	".ogv":            {"\uf03d", 241},     // 
	".otf":            {"\uf031", 241},     // 
	".part":           {"\uf43a", 241},     // 
	".patch":          {"\uf440", 241},     // 
	".pdf":            {"\uf1c1", 241},     // 
	".php":            {"\ue73d", 61},      // 
	".pl":             {"\ue769", 241},     // 
	".png":            {"\uf1c5", 241},     // 
	".ppt":            {"\uf1c4", 241},     // 
	".pptx":           {"\uf1c4", 241},     // 
	".procfile":       {"\ue21e", 241},     // 
	".properties":     {"\ue60b", 185},     // 
	".ps1":            {"\ue795", 241},     // 
	".psd":            {"\ue7b8", 241},     // 
	".pxm":            {"\uf1c5", 241},     // 
	".py":             {"\ue606", 214},     // 
	".pyc":            {"\ue606", 214},     // 
	".r":              {"\uf25d", 241},     // 
	".rakefile":       {"\ue21e", 241},     // 
	".rar":            {"\uf410", 241},     // 
	".razor":          {"\uf1fa", 241},     // 
	".rb":             {"\ue21e", 241},     // 
	".rdata":          {"\uf25d", 241},     // 
	".rdb":            {"\ue76d", 241},     // 
	".rdoc":           {"\uf48a", 74},      // 
	".rds":            {"\uf25d", 241},     // 
	".readme":         {"\uf48a", 241},     // 
	".rlib":           {"\ue7a8", 241},     // 
	".rmd":            {"\uf48a", 74},      // 
	".rpm":            {"\ue7bb", 241},     // 
	".rs":             {"\ue7a8", 216},     // 
	".rspec":          {"\ue21e", 241},     // 
	".rspec_parallel": {"\ue21e", 241},     // 
	".rspec_status":   {"\ue21e", 241},     // 
	".rss":            {"\uf09e", 241},     // 
	".rtf":            {"\U000f0219", 241}, // 󰈙
	".ru":             {"\ue21e", 241},     // 
	".rubydoc":        {"\ue73b", 160},     // 
	".sass":           {"\ue603", 169},     // 
	".scala":          {"\ue737", 74},      // 
	".scss":           {"\ue749", 204},     // 
	".sh":             {"\ue795", 239},     // 
	".shell":          {"\ue795", 239},     // 
	".slim":           {"\ue73b", 160},     // 
	".sln":            {"\ue70c", 39},      // 
	".so":             {"\uf17c", 241},     // 
	".sql":            {"\uf1c0", 188},     // 
	".sqlite3":        {"\ue7c4", 25},      // 
	".sty":            {"\uf034", 239},     // 
	".styl":           {"\ue600", 241},     // 
	".stylus":         {"\ue600", 241},     // 
	".svelte":         {"\ue697", 208},     // 
	".svg":            {"\uf1c5", 241},     // 
	".swift":          {"\ue755", 208},     // 
	".tar":            {"\uf410", 241},     // 
	".taz":            {"\uf410", 241},     // 
	".tbz":            {"\uf410", 241},     // 
	".tbz2":           {"\uf410", 241},     // 
	".tex":            {"\uf034", 241},     // 
	".tgz":            {"\uf410", 241},     // 
	".tiff":           {"\uf1c5", 241},     // 
	".tlz":            {"\uf410", 241},     // 
	".toml":           {"\ue615", 241},     // 
	".torrent":        {"\ue275", 241},     // 
	".ts":             {"\ue628", 74},      // 
	".tsv":            {"\uf1c3", 241},     // 
	".tsx":            {"\ue7ba", 74},      // 
	".ttf":            {"\uf031", 241},     // 
	".twig":           {"\ue61c", 241},     // 
	".txt":            {"\uf15c", 241},     // 
	".txz":            {"\uf410", 241},     // 
	".tz":             {"\uf410", 241},     // 
	".tzo":            {"\uf410", 241},     // 
	".video":          {"\uf03d", 241},     // 
	".vim":            {"\ue62b", 28},      // 
	".vue":            {"\U000f0844", 113}, // 󰡄
	".war":            {"\ue256", 168},     // 
	".wav":            {"\uf001", 241},     // 
	".webm":           {"\uf03d", 241},     // 
	".webp":           {"\uf1c5", 241},     // 
	".windows":        {"\uf17a", 81},      // 
	".woff":           {"\uf031", 241},     // 
	".woff2":          {"\uf031", 241},     // 
	".xhtml":          {"\uf13b", 196},     // 
	".xls":            {"\uf1c3", 241},     // 
	".xlsx":           {"\uf1c3", 241},     // 
	".xml":            {"\uf121", 241},     // 
	".xul":            {"\uf121", 241},     // 
	".xz":             {"\uf410", 241},     // 
	".yaml":           {"\uf481", 241},     // 
	".yml":            {"\uf481", 241},     // 
	".zip":            {"\uf410", 241},     // 
	".zsh":            {"\ue795", 241},     // 
	".zsh-theme":      {"\ue795", 241},     // 
	".zshrc":          {"\ue795", 241},     // 
	".zst":            {"\uf410", 241},     // 
}

func patchFileIconsForNerdFontsV2() {
	extIconMap[".cs"] = IconProperties{"\uf81a", 58}       // 
	extIconMap[".csproj"] = IconProperties{"\uf81a", 58}   // 
	extIconMap[".csx"] = IconProperties{"\uf81a", 58}      // 
	extIconMap[".license"] = IconProperties{"\uf718", 241} // 
	extIconMap[".node"] = IconProperties{"\uf898", 197}    // 
	extIconMap[".rtf"] = IconProperties{"\uf718", 241}     // 
	extIconMap[".vue"] = IconProperties{"\ufd42", 113}     // ﵂
}

func IconForFile(name string, isSubmodule bool, isLinkedWorktree bool, isDirectory bool) IconProperties {
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
