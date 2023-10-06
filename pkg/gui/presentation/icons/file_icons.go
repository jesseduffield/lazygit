package icons

import (
	"path/filepath"
)

// https://github.com/ogham/exa/blob/master/src/output/icons.rs
var (
	DEFAULT_FILE_ICON      = IconProperties{"\uf15b", 239} // 
	DEFAULT_SUBMODULE_ICON = IconProperties{"\uf1d3", 239} // 
	DEFAULT_DIRECTORY_ICON = IconProperties{"\uf114", 239} // 
)

var nameIconMap = map[string]IconProperties{
	".Trash":             {"\uf1f8", 239}, // 
	".atom":              {"\ue764", 239}, // 
	".bashprofile":       {"\ue615", 239}, // 
	".bashrc":            {"\uf489", 239}, // 
	".idea":              {"\ue7b5", 239}, // 
	".git":               {"\uf1d3", 239}, // 
	".gitattributes":     {"\uf1d3", 239}, // 
	".gitconfig":         {"\uf1d3", 239}, // 
	".github":            {"\uf408", 239}, // 
	".gitignore":         {"\uf1d3", 239}, // 
	".gitmodules":        {"\uf1d3", 239}, // 
	".rvm":               {"\ue21e", 239}, // 
	".vimrc":             {"\ue62b", 239}, // 
	".vscode":            {"\ue70c", 239}, // 
	".zshrc":             {"\uf489", 239}, // 
	"Cargo.lock":         {"\ue7a8", 239}, // 
	"Cargo.toml":         {"\ue7a8", 239}, // 
	"bin":                {"\ue5fc", 239}, // 
	"config":             {"\ue5fc", 239}, // 
	"docker-compose.yml": {"\uf308", 239}, // 
	"Dockerfile":         {"\uf308", 239}, // 
	"ds_store":           {"\uf179", 239}, // 
	"gitignore_global":   {"\uf1d3", 239}, // 
	"go.mod":             {"\ue626", 239}, // 
	"go.sum":             {"\ue626", 239}, // 
	"gradle":             {"\ue256", 239}, // 
	"gruntfile.coffee":   {"\ue611", 239}, // 
	"gruntfile.js":       {"\ue611", 239}, // 
	"gruntfile.ls":       {"\ue611", 239}, // 
	"gulpfile.coffee":    {"\ue610", 239}, // 
	"gulpfile.js":        {"\ue610", 239}, // 
	"gulpfile.ls":        {"\ue610", 239}, // 
	"hidden":             {"\uf023", 239}, // 
	"include":            {"\ue5fc", 239}, // 
	"lib":                {"\uf121", 239}, // 
	"localized":          {"\uf179", 239}, // 
	"Makefile":           {"\uf489", 239}, // 
	"node_modules":       {"\ue718", 239}, // 
	"npmignore":          {"\ue71e", 239}, // 
	"PKGBUILD":           {"\uf303", 239}, // 
	"rubydoc":            {"\ue73b", 239}, // 
	"yarn.lock":          {"\ue718", 239}, // 
}

var extIconMap = map[string]IconProperties{
	".ai":             {"\ue7b4", 239},     // 
	".android":        {"\ue70e", 239},     // 
	".apk":            {"\ue70e", 239},     // 
	".apple":          {"\uf179", 239},     // 
	".avi":            {"\uf03d", 239},     // 
	".avif":           {"\uf1c5", 239},     // 
	".avro":           {"\ue60b", 239},     // 
	".awk":            {"\uf489", 239},     // 
	".bash":           {"\uf489", 239},     // 
	".bash_history":   {"\uf489", 239},     // 
	".bash_profile":   {"\uf489", 239},     // 
	".bashrc":         {"\uf489", 239},     // 
	".bat":            {"\uf17a", 239},     // 
	".bats":           {"\uf489", 239},     // 
	".bmp":            {"\uf1c5", 239},     // 
	".bz":             {"\uf410", 239},     // 
	".bz2":            {"\uf410", 239},     // 
	".c":              {"\ue61e", 239},     // 
	".c++":            {"\ue61d", 239},     // 
	".cab":            {"\ue70f", 239},     // 
	".cc":             {"\ue61d", 239},     // 
	".cfg":            {"\ue615", 239},     // 
	".class":          {"\ue256", 239},     // 
	".clj":            {"\ue768", 239},     // 
	".cljs":           {"\ue76a", 239},     // 
	".cls":            {"\uf034", 239},     // 
	".cmd":            {"\ue70f", 239},     // 
	".coffee":         {"\uf0f4", 239},     // 
	".conf":           {"\ue615", 239},     // 
	".cp":             {"\ue61d", 239},     // 
	".cpio":           {"\uf410", 239},     // 
	".cpp":            {"\ue61d", 239},     // 
	".cs":             {"\U000f031b", 239}, // 󰌛
	".csh":            {"\uf489", 239},     // 
	".cshtml":         {"\uf1fa", 239},     // 
	".csproj":         {"\U000f031b", 239}, // 󰌛
	".css":            {"\ue749", 239},     // 
	".csv":            {"\uf1c3", 239},     // 
	".csx":            {"\U000f031b", 239}, // 󰌛
	".cxx":            {"\ue61d", 239},     // 
	".d":              {"\ue7af", 239},     // 
	".dart":           {"\ue798", 239},     // 
	".db":             {"\uf1c0", 239},     // 
	".deb":            {"\ue77d", 239},     // 
	".diff":           {"\uf440", 239},     // 
	".djvu":           {"\uf02d", 239},     // 
	".dll":            {"\ue70f", 239},     // 
	".doc":            {"\uf1c2", 239},     // 
	".docx":           {"\uf1c2", 239},     // 
	".ds_store":       {"\uf179", 239},     // 
	".DS_store":       {"\uf179", 239},     // 
	".dump":           {"\uf1c0", 239},     // 
	".ebook":          {"\ue28b", 239},     // 
	".ebuild":         {"\uf30d", 239},     // 
	".editorconfig":   {"\ue615", 239},     // 
	".ejs":            {"\ue618", 239},     // 
	".elm":            {"\ue62c", 239},     // 
	".env":            {"\uf462", 239},     // 
	".eot":            {"\uf031", 239},     // 
	".epub":           {"\ue28a", 239},     // 
	".erb":            {"\ue73b", 239},     // 
	".erl":            {"\ue7b1", 239},     // 
	".ex":             {"\ue62d", 239},     // 
	".exe":            {"\uf17a", 239},     // 
	".exs":            {"\ue62d", 239},     // 
	".fish":           {"\uf489", 239},     // 
	".flac":           {"\uf001", 239},     // 
	".flv":            {"\uf03d", 239},     // 
	".font":           {"\uf031", 239},     // 
	".fs":             {"\ue7a7", 239},     // 
	".fsi":            {"\ue7a7", 239},     // 
	".fsx":            {"\ue7a7", 239},     // 
	".gdoc":           {"\uf1c2", 239},     // 
	".gem":            {"\ue21e", 239},     // 
	".gemfile":        {"\ue21e", 239},     // 
	".gemspec":        {"\ue21e", 239},     // 
	".gform":          {"\uf298", 239},     // 
	".gif":            {"\uf1c5", 239},     // 
	".git":            {"\uf1d3", 239},     // 
	".gitattributes":  {"\uf1d3", 239},     // 
	".gitignore":      {"\uf1d3", 239},     // 
	".gitmodules":     {"\uf1d3", 239},     // 
	".go":             {"\ue626", 239},     // 
	".gradle":         {"\ue256", 239},     // 
	".groovy":         {"\ue775", 239},     // 
	".gsheet":         {"\uf1c3", 239},     // 
	".gslides":        {"\uf1c4", 239},     // 
	".guardfile":      {"\ue21e", 239},     // 
	".gz":             {"\uf410", 239},     // 
	".h":              {"\uf0fd", 239},     // 
	".hbs":            {"\ue60f", 239},     // 
	".hpp":            {"\uf0fd", 239},     // 
	".hs":             {"\ue777", 239},     // 
	".htm":            {"\uf13b", 239},     // 
	".html":           {"\uf13b", 239},     // 
	".hxx":            {"\uf0fd", 239},     // 
	".ico":            {"\uf1c5", 239},     // 
	".image":          {"\uf1c5", 239},     // 
	".iml":            {"\ue7b5", 239},     // 
	".ini":            {"\uf17a", 239},     // 
	".ipynb":          {"\ue606", 239},     // 
	".iso":            {"\ue271", 239},     // 
	".j2c":            {"\uf1c5", 239},     // 
	".j2k":            {"\uf1c5", 239},     // 
	".jad":            {"\ue256", 239},     // 
	".jar":            {"\ue256", 239},     // 
	".java":           {"\ue256", 239},     // 
	".jfi":            {"\uf1c5", 239},     // 
	".jfif":           {"\uf1c5", 239},     // 
	".jif":            {"\uf1c5", 239},     // 
	".jl":             {"\ue624", 239},     // 
	".jmd":            {"\uf48a", 239},     // 
	".jp2":            {"\uf1c5", 239},     // 
	".jpe":            {"\uf1c5", 239},     // 
	".jpeg":           {"\uf1c5", 239},     // 
	".jpg":            {"\uf1c5", 239},     // 
	".jpx":            {"\uf1c5", 239},     // 
	".js":             {"\ue74e", 239},     // 
	".json":           {"\ue60b", 239},     // 
	".jsx":            {"\ue7ba", 239},     // 
	".jxl":            {"\uf1c5", 239},     // 
	".ksh":            {"\uf489", 239},     // 
	".kt":             {"\ue634", 239},     // 
	".kts":            {"\ue634", 239},     // 
	".latex":          {"\uf034", 239},     // 
	".less":           {"\ue758", 239},     // 
	".lhs":            {"\ue777", 239},     // 
	".license":        {"\U000f0219", 239}, // 󰈙
	".localized":      {"\uf179", 239},     // 
	".lock":           {"\uf023", 239},     // 
	".log":            {"\uf18d", 239},     // 
	".lua":            {"\ue620", 239},     // 
	".lz":             {"\uf410", 239},     // 
	".lz4":            {"\uf410", 239},     // 
	".lzh":            {"\uf410", 239},     // 
	".lzma":           {"\uf410", 239},     // 
	".lzo":            {"\uf410", 239},     // 
	".m":              {"\ue61e", 239},     // 
	".mm":             {"\ue61d", 239},     // 
	".m4a":            {"\uf001", 239},     // 
	".markdown":       {"\uf48a", 239},     // 
	".md":             {"\uf48a", 239},     // 
	".mdx":            {"\uf48a", 239},     // 
	".mjs":            {"\ue74e", 239},     // 
	".mk":             {"\uf489", 239},     // 
	".mkd":            {"\uf48a", 239},     // 
	".mkv":            {"\uf03d", 239},     // 
	".mobi":           {"\ue28b", 239},     // 
	".mov":            {"\uf03d", 239},     // 
	".mp3":            {"\uf001", 239},     // 
	".mp4":            {"\uf03d", 239},     // 
	".msi":            {"\ue70f", 239},     // 
	".mustache":       {"\ue60f", 239},     // 
	".nix":            {"\uf313", 239},     // 
	".node":           {"\U000f0399", 239}, // 󰎙
	".npmignore":      {"\ue71e", 239},     // 
	".odp":            {"\uf1c4", 239},     // 
	".ods":            {"\uf1c3", 239},     // 
	".odt":            {"\uf1c2", 239},     // 
	".ogg":            {"\uf001", 239},     // 
	".ogv":            {"\uf03d", 239},     // 
	".otf":            {"\uf031", 239},     // 
	".part":           {"\uf43a", 239},     // 
	".patch":          {"\uf440", 239},     // 
	".pdf":            {"\uf1c1", 239},     // 
	".php":            {"\ue73d", 239},     // 
	".pl":             {"\ue769", 239},     // 
	".png":            {"\uf1c5", 239},     // 
	".ppt":            {"\uf1c4", 239},     // 
	".pptx":           {"\uf1c4", 239},     // 
	".procfile":       {"\ue21e", 239},     // 
	".properties":     {"\ue60b", 239},     // 
	".ps1":            {"\uf489", 239},     // 
	".psd":            {"\ue7b8", 239},     // 
	".pxm":            {"\uf1c5", 239},     // 
	".py":             {"\ue606", 239},     // 
	".pyc":            {"\ue606", 239},     // 
	".r":              {"\uf25d", 239},     // 
	".rakefile":       {"\ue21e", 239},     // 
	".rar":            {"\uf410", 239},     // 
	".razor":          {"\uf1fa", 239},     // 
	".rb":             {"\ue21e", 239},     // 
	".rdata":          {"\uf25d", 239},     // 
	".rdb":            {"\ue76d", 239},     // 
	".rdoc":           {"\uf48a", 239},     // 
	".rds":            {"\uf25d", 239},     // 
	".readme":         {"\uf48a", 239},     // 
	".rlib":           {"\ue7a8", 239},     // 
	".rmd":            {"\uf48a", 239},     // 
	".rpm":            {"\ue7bb", 239},     // 
	".rs":             {"\ue7a8", 239},     // 
	".rspec":          {"\ue21e", 239},     // 
	".rspec_parallel": {"\ue21e", 239},     // 
	".rspec_status":   {"\ue21e", 239},     // 
	".rss":            {"\uf09e", 239},     // 
	".rtf":            {"\U000f0219", 239}, // 󰈙
	".ru":             {"\ue21e", 239},     // 
	".rubydoc":        {"\ue73b", 239},     // 
	".sass":           {"\ue603", 239},     // 
	".scala":          {"\ue737", 239},     // 
	".scss":           {"\ue749", 239},     // 
	".sh":             {"\uf489", 239},     // 
	".shell":          {"\uf489", 239},     // 
	".slim":           {"\ue73b", 239},     // 
	".sln":            {"\ue70c", 239},     // 
	".so":             {"\uf17c", 239},     // 
	".sql":            {"\uf1c0", 239},     // 
	".sqlite3":        {"\ue7c4", 239},     // 
	".sty":            {"\uf034", 239},     // 
	".styl":           {"\ue600", 239},     // 
	".stylus":         {"\ue600", 239},     // 
	".svelte":         {"\ue697", 239},     // 
	".svg":            {"\uf1c5", 239},     // 
	".swift":          {"\ue755", 239},     // 
	".tar":            {"\uf410", 239},     // 
	".taz":            {"\uf410", 239},     // 
	".tbz":            {"\uf410", 239},     // 
	".tbz2":           {"\uf410", 239},     // 
	".tex":            {"\uf034", 239},     // 
	".tgz":            {"\uf410", 239},     // 
	".tiff":           {"\uf1c5", 239},     // 
	".tlz":            {"\uf410", 239},     // 
	".toml":           {"\ue615", 239},     // 
	".torrent":        {"\ue275", 239},     // 
	".ts":             {"\ue628", 239},     // 
	".tsv":            {"\uf1c3", 239},     // 
	".tsx":            {"\ue7ba", 239},     // 
	".ttf":            {"\uf031", 239},     // 
	".twig":           {"\ue61c", 239},     // 
	".txt":            {"\uf15c", 239},     // 
	".txz":            {"\uf410", 239},     // 
	".tz":             {"\uf410", 239},     // 
	".tzo":            {"\uf410", 239},     // 
	".video":          {"\uf03d", 239},     // 
	".vim":            {"\ue62b", 239},     // 
	".vue":            {"\U000f0844", 239}, // 󰡄
	".war":            {"\ue256", 239},     // 
	".wav":            {"\uf001", 239},     // 
	".webm":           {"\uf03d", 239},     // 
	".webp":           {"\uf1c5", 239},     // 
	".windows":        {"\uf17a", 239},     // 
	".woff":           {"\uf031", 239},     // 
	".woff2":          {"\uf031", 239},     // 
	".xhtml":          {"\uf13b", 239},     // 
	".xls":            {"\uf1c3", 239},     // 
	".xlsx":           {"\uf1c3", 239},     // 
	".xml":            {"\uf121", 239},     // 
	".xul":            {"\uf121", 239},     // 
	".xz":             {"\uf410", 239},     // 
	".yaml":           {"\uf481", 239},     // 
	".yml":            {"\uf481", 239},     // 
	".zip":            {"\uf410", 239},     // 
	".zsh":            {"\uf489", 239},     // 
	".zsh-theme":      {"\uf489", 239},     // 
	".zshrc":          {"\uf489", 239},     // 
	".zst":            {"\uf410", 239},     // 
}

func patchFileIconsForNerdFontsV2() {
	extIconMap[".cs"] = IconProperties{"\uf81a", 239}      // 
	extIconMap[".csproj"] = IconProperties{"\uf81a", 239}  // 
	extIconMap[".csx"] = IconProperties{"\uf81a", 239}     // 
	extIconMap[".license"] = IconProperties{"\uf718", 239} // 
	extIconMap[".node"] = IconProperties{"\uf898", 239}    // 
	extIconMap[".rtf"] = IconProperties{"\uf718", 239}     // 
	extIconMap[".vue"] = IconProperties{"\ufd42", 239}     // ﵂
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
