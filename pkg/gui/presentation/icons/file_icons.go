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
	".Trash":             {2: "\uf1f8"}, // 
	".atom":              {2: "\ue764"}, // 
	".bashprofile":       {2: "\ue615"}, // 
	".bashrc":            {2: "\uf489"}, // 
	".idea":              {2: "\ue7b5"}, // 
	".git":               {2: "\uf1d3"}, // 
	".gitattributes":     {2: "\uf1d3"}, // 
	".gitconfig":         {2: "\uf1d3"}, // 
	".github":            {2: "\uf408"}, // 
	".gitignore":         {2: "\uf1d3"}, // 
	".gitmodules":        {2: "\uf1d3"}, // 
	".rvm":               {2: "\ue21e"}, // 
	".vimrc":             {2: "\ue62b"}, // 
	".vscode":            {2: "\ue70c"}, // 
	".zshrc":             {2: "\uf489"}, // 
	"Cargo.lock":         {2: "\ue7a8"}, // 
	"Cargo.toml":         {2: "\ue7a8"}, // 
	"bin":                {2: "\ue5fc"}, // 
	"config":             {2: "\ue5fc"}, // 
	"docker-compose.yml": {2: "\uf308"}, // 
	"Dockerfile":         {2: "\uf308"}, // 
	"ds_store":           {2: "\uf179"}, // 
	"gitignore_global":   {2: "\uf1d3"}, // 
	"go.mod":             {2: "\ue626"}, // 
	"go.sum":             {2: "\ue626"}, // 
	"gradle":             {2: "\ue256"}, // 
	"gruntfile.coffee":   {2: "\ue611"}, // 
	"gruntfile.js":       {2: "\ue611"}, // 
	"gruntfile.ls":       {2: "\ue611"}, // 
	"gulpfile.coffee":    {2: "\ue610"}, // 
	"gulpfile.js":        {2: "\ue610"}, // 
	"gulpfile.ls":        {2: "\ue610"}, // 
	"hidden":             {2: "\uf023"}, // 
	"include":            {2: "\ue5fc"}, // 
	"lib":                {2: "\uf121"}, // 
	"localized":          {2: "\uf179"}, // 
	"Makefile":           {2: "\uf489"}, // 
	"node_modules":       {2: "\ue718"}, // 
	"npmignore":          {2: "\ue71e"}, // 
	"PKGBUILD":           {2: "\uf303"}, // 
	"rubydoc":            {2: "\ue73b"}, // 
	"yarn.lock":          {2: "\ue718"}, // 
}

var extIconMap = map[string]map[int]string{
	".ai":             {2: "\ue7b4"},                  // 
	".android":        {2: "\ue70e"},                  // 
	".apk":            {2: "\ue70e"},                  // 
	".apple":          {2: "\uf179"},                  // 
	".avi":            {2: "\uf03d"},                  // 
	".avif":           {2: "\uf1c5"},                  // 
	".avro":           {2: "\ue60b"},                  // 
	".awk":            {2: "\uf489"},                  // 
	".bash":           {2: "\uf489"},                  // 
	".bash_history":   {2: "\uf489"},                  // 
	".bash_profile":   {2: "\uf489"},                  // 
	".bashrc":         {2: "\uf489"},                  // 
	".bat":            {2: "\uf17a"},                  // 
	".bats":           {2: "\uf489"},                  // 
	".bmp":            {2: "\uf1c5"},                  // 
	".bz":             {2: "\uf410"},                  // 
	".bz2":            {2: "\uf410"},                  // 
	".c":              {2: "\ue61e"},                  // 
	".c++":            {2: "\ue61d"},                  // 
	".cab":            {2: "\ue70f"},                  // 
	".cc":             {2: "\ue61d"},                  // 
	".cfg":            {2: "\ue615"},                  // 
	".class":          {2: "\ue256"},                  // 
	".clj":            {2: "\ue768"},                  // 
	".cljs":           {2: "\ue76a"},                  // 
	".cls":            {2: "\uf034"},                  // 
	".cmd":            {2: "\ue70f"},                  // 
	".coffee":         {2: "\uf0f4"},                  // 
	".conf":           {2: "\ue615"},                  // 
	".cp":             {2: "\ue61d"},                  // 
	".cpio":           {2: "\uf410"},                  // 
	".cpp":            {2: "\ue61d"},                  // 
	".cs":             {2: "\uf81a"},                  // 
	".csh":            {2: "\uf489"},                  // 
	".cshtml":         {2: "\uf1fa"},                  // 
	".csproj":         {2: "\uf81a"},                  // 
	".css":            {2: "\ue749"},                  // 
	".csv":            {2: "\uf1c3"},                  // 
	".csx":            {2: "\uf81a"},                  // 
	".cxx":            {2: "\ue61d"},                  // 
	".d":              {2: "\ue7af"},                  // 
	".dart":           {2: "\ue798"},                  // 
	".db":             {2: "\uf1c0"},                  // 
	".deb":            {2: "\ue77d"},                  // 
	".diff":           {2: "\uf440"},                  // 
	".djvu":           {2: "\uf02d"},                  // 
	".dll":            {2: "\ue70f"},                  // 
	".doc":            {2: "\uf1c2"},                  // 
	".docx":           {2: "\uf1c2"},                  // 
	".ds_store":       {2: "\uf179"},                  // 
	".DS_store":       {2: "\uf179"},                  // 
	".dump":           {2: "\uf1c0"},                  // 
	".ebook":          {2: "\ue28b"},                  // 
	".ebuild":         {2: "\uf30d"},                  // 
	".editorconfig":   {2: "\ue615"},                  // 
	".ejs":            {2: "\ue618"},                  // 
	".elm":            {2: "\ue62c"},                  // 
	".env":            {2: "\uf462"},                  // 
	".eot":            {2: "\uf031"},                  // 
	".epub":           {2: "\ue28a"},                  // 
	".erb":            {2: "\ue73b"},                  // 
	".erl":            {2: "\ue7b1"},                  // 
	".ex":             {2: "\ue62d"},                  // 
	".exe":            {2: "\uf17a"},                  // 
	".exs":            {2: "\ue62d"},                  // 
	".fish":           {2: "\uf489"},                  // 
	".flac":           {2: "\uf001"},                  // 
	".flv":            {2: "\uf03d"},                  // 
	".font":           {2: "\uf031"},                  // 
	".fs":             {2: "\ue7a7"},                  // 
	".fsi":            {2: "\ue7a7"},                  // 
	".fsx":            {2: "\ue7a7"},                  // 
	".gdoc":           {2: "\uf1c2"},                  // 
	".gem":            {2: "\ue21e"},                  // 
	".gemfile":        {2: "\ue21e"},                  // 
	".gemspec":        {2: "\ue21e"},                  // 
	".gform":          {2: "\uf298"},                  // 
	".gif":            {2: "\uf1c5"},                  // 
	".git":            {2: "\uf1d3"},                  // 
	".gitattributes":  {2: "\uf1d3"},                  // 
	".gitignore":      {2: "\uf1d3"},                  // 
	".gitmodules":     {2: "\uf1d3"},                  // 
	".go":             {2: "\ue626"},                  // 
	".gradle":         {2: "\ue256"},                  // 
	".groovy":         {2: "\ue775"},                  // 
	".gsheet":         {2: "\uf1c3"},                  // 
	".gslides":        {2: "\uf1c4"},                  // 
	".guardfile":      {2: "\ue21e"},                  // 
	".gz":             {2: "\uf410"},                  // 
	".h":              {2: "\uf0fd"},                  // 
	".hbs":            {2: "\ue60f"},                  // 
	".hpp":            {2: "\uf0fd"},                  // 
	".hs":             {2: "\ue777"},                  // 
	".htm":            {2: "\uf13b"},                  // 
	".html":           {2: "\uf13b"},                  // 
	".hxx":            {2: "\uf0fd"},                  // 
	".ico":            {2: "\uf1c5"},                  // 
	".image":          {2: "\uf1c5"},                  // 
	".iml":            {2: "\ue7b5"},                  // 
	".ini":            {2: "\uf17a"},                  // 
	".ipynb":          {2: "\ue606"},                  // 
	".iso":            {2: "\ue271"},                  // 
	".j2c":            {2: "\uf1c5"},                  // 
	".j2k":            {2: "\uf1c5"},                  // 
	".jad":            {2: "\ue256"},                  // 
	".jar":            {2: "\ue256"},                  // 
	".java":           {2: "\ue256"},                  // 
	".jfi":            {2: "\uf1c5"},                  // 
	".jfif":           {2: "\uf1c5"},                  // 
	".jif":            {2: "\uf1c5"},                  // 
	".jl":             {2: "\ue624"},                  // 
	".jmd":            {2: "\uf48a"},                  // 
	".jp2":            {2: "\uf1c5"},                  // 
	".jpe":            {2: "\uf1c5"},                  // 
	".jpeg":           {2: "\uf1c5"},                  // 
	".jpg":            {2: "\uf1c5"},                  // 
	".jpx":            {2: "\uf1c5"},                  // 
	".js":             {2: "\ue74e"},                  // 
	".json":           {2: "\ue60b"},                  // 
	".jsx":            {2: "\ue7ba"},                  // 
	".jxl":            {2: "\uf1c5"},                  // 
	".ksh":            {2: "\uf489"},                  // 
	".kt":             {2: "\ue634"},                  // 
	".kts":            {2: "\ue634"},                  // 
	".latex":          {2: "\uf034"},                  // 
	".less":           {2: "\ue758"},                  // 
	".lhs":            {2: "\ue777"},                  // 
	".license":        {2: "\uf718", 3: "\U000f0219"}, //   󰈙
	".localized":      {2: "\uf179"},                  // 
	".lock":           {2: "\uf023"},                  // 
	".log":            {2: "\uf18d"},                  // 
	".lua":            {2: "\ue620"},                  // 
	".lz":             {2: "\uf410"},                  // 
	".lz4":            {2: "\uf410"},                  // 
	".lzh":            {2: "\uf410"},                  // 
	".lzma":           {2: "\uf410"},                  // 
	".lzo":            {2: "\uf410"},                  // 
	".m":              {2: "\ue61e"},                  // 
	".mm":             {2: "\ue61d"},                  // 
	".m4a":            {2: "\uf001"},                  // 
	".markdown":       {2: "\uf48a"},                  // 
	".md":             {2: "\uf48a"},                  // 
	".mjs":            {2: "\ue74e"},                  // 
	".mk":             {2: "\uf489"},                  // 
	".mkd":            {2: "\uf48a"},                  // 
	".mkv":            {2: "\uf03d"},                  // 
	".mobi":           {2: "\ue28b"},                  // 
	".mov":            {2: "\uf03d"},                  // 
	".mp3":            {2: "\uf001"},                  // 
	".mp4":            {2: "\uf03d"},                  // 
	".msi":            {2: "\ue70f"},                  // 
	".mustache":       {2: "\ue60f"},                  // 
	".nix":            {2: "\uf313"},                  // 
	".node":           {2: "\uf898"},                  // 
	".npmignore":      {2: "\ue71e"},                  // 
	".odp":            {2: "\uf1c4"},                  // 
	".ods":            {2: "\uf1c3"},                  // 
	".odt":            {2: "\uf1c2"},                  // 
	".ogg":            {2: "\uf001"},                  // 
	".ogv":            {2: "\uf03d"},                  // 
	".otf":            {2: "\uf031"},                  // 
	".part":           {2: "\uf43a"},                  // 
	".patch":          {2: "\uf440"},                  // 
	".pdf":            {2: "\uf1c1"},                  // 
	".php":            {2: "\ue73d"},                  // 
	".pl":             {2: "\ue769"},                  // 
	".png":            {2: "\uf1c5"},                  // 
	".ppt":            {2: "\uf1c4"},                  // 
	".pptx":           {2: "\uf1c4"},                  // 
	".procfile":       {2: "\ue21e"},                  // 
	".properties":     {2: "\ue60b"},                  // 
	".ps1":            {2: "\uf489"},                  // 
	".psd":            {2: "\ue7b8"},                  // 
	".pxm":            {2: "\uf1c5"},                  // 
	".py":             {2: "\ue606"},                  // 
	".pyc":            {2: "\ue606"},                  // 
	".r":              {2: "\uf25d"},                  // 
	".rakefile":       {2: "\ue21e"},                  // 
	".rar":            {2: "\uf410"},                  // 
	".razor":          {2: "\uf1fa"},                  // 
	".rb":             {2: "\ue21e"},                  // 
	".rdata":          {2: "\uf25d"},                  // 
	".rdb":            {2: "\ue76d"},                  // 
	".rdoc":           {2: "\uf48a"},                  // 
	".rds":            {2: "\uf25d"},                  // 
	".readme":         {2: "\uf48a"},                  // 
	".rlib":           {2: "\ue7a8"},                  // 
	".rmd":            {2: "\uf48a"},                  // 
	".rpm":            {2: "\ue7bb"},                  // 
	".rs":             {2: "\ue7a8"},                  // 
	".rspec":          {2: "\ue21e"},                  // 
	".rspec_parallel": {2: "\ue21e"},                  // 
	".rspec_status":   {2: "\ue21e"},                  // 
	".rss":            {2: "\uf09e"},                  // 
	".rtf":            {2: "\uf718", 3: "\U000f0219"}, //    󰈙
	".ru":             {2: "\ue21e"},                  // 
	".rubydoc":        {2: "\ue73b"},                  // 
	".sass":           {2: "\ue603"},                  // 
	".scala":          {2: "\ue737"},                  // 
	".scss":           {2: "\ue749"},                  // 
	".sh":             {2: "\uf489"},                  // 
	".shell":          {2: "\uf489"},                  // 
	".slim":           {2: "\ue73b"},                  // 
	".sln":            {2: "\ue70c"},                  // 
	".so":             {2: "\uf17c"},                  // 
	".sql":            {2: "\uf1c0"},                  // 
	".sqlite3":        {2: "\ue7c4"},                  // 
	".sty":            {2: "\uf034"},                  // 
	".styl":           {2: "\ue600"},                  // 
	".stylus":         {2: "\ue600"},                  // 
	".svg":            {2: "\uf1c5"},                  // 
	".swift":          {2: "\ue755"},                  // 
	".tar":            {2: "\uf410"},                  // 
	".taz":            {2: "\uf410"},                  // 
	".tbz":            {2: "\uf410"},                  // 
	".tbz2":           {2: "\uf410"},                  // 
	".tex":            {2: "\uf034"},                  // 
	".tgz":            {2: "\uf410"},                  // 
	".tiff":           {2: "\uf1c5"},                  // 
	".tlz":            {2: "\uf410"},                  // 
	".toml":           {2: "\ue615"},                  // 
	".torrent":        {2: "\ue275"},                  // 
	".ts":             {2: "\ue628"},                  // 
	".tsv":            {2: "\uf1c3"},                  // 
	".tsx":            {2: "\ue7ba"},                  // 
	".ttf":            {2: "\uf031"},                  // 
	".twig":           {2: "\ue61c"},                  // 
	".txt":            {2: "\uf15c"},                  // 
	".txz":            {2: "\uf410"},                  // 
	".tz":             {2: "\uf410"},                  // 
	".tzo":            {2: "\uf410"},                  // 
	".video":          {2: "\uf03d"},                  // 
	".vim":            {2: "\ue62b"},                  // 
	".vue":            {2: "\ufd42"},                  // ﵂
	".war":            {2: "\ue256"},                  // 
	".wav":            {2: "\uf001"},                  // 
	".webm":           {2: "\uf03d"},                  // 
	".webp":           {2: "\uf1c5"},                  // 
	".windows":        {2: "\uf17a"},                  // 
	".woff":           {2: "\uf031"},                  // 
	".woff2":          {2: "\uf031"},                  // 
	".xhtml":          {2: "\uf13b"},                  // 
	".xls":            {2: "\uf1c3"},                  // 
	".xlsx":           {2: "\uf1c3"},                  // 
	".xml":            {2: "\uf121"},                  // 
	".xul":            {2: "\uf121"},                  // 
	".xz":             {2: "\uf410"},                  // 
	".yaml":           {2: "\uf481"},                  // 
	".yml":            {2: "\uf481"},                  // 
	".zip":            {2: "\uf410"},                  // 
	".zsh":            {2: "\uf489"},                  // 
	".zsh-theme":      {2: "\uf489"},                  // 
	".zshrc":          {2: "\uf489"},                  // 
	".zst":            {2: "\uf410"},                  // 
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
