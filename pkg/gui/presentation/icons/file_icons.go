package icons

import (
	"path/filepath"
)

// NOTE: Visit next links for inspiration:
// https://github.com/eza-community/eza/blob/main/src/output/icons.rs
// https://github.com/nvim-tree/nvim-web-devicons/blob/master/lua/nvim-web-devicons/icons-default.lua

var (
	DEFAULT_FILE_ICON      = IconProperties{Icon: "\uf15b", Color: 241} // 
	DEFAULT_SUBMODULE_ICON = IconProperties{Icon: "\uf1d3", Color: 202} // 
	DEFAULT_DIRECTORY_ICON = IconProperties{Icon: "\uf07b", Color: 241} // 
)

var nameIconMap = map[string]IconProperties{
	".Trash":             {Icon: "\uf1f8", Color: 241}, // 
	".atom":              {Icon: "\ue764", Color: 241}, // 
	".bashprofile":       {Icon: "\ue615", Color: 113}, // 
	".bashrc":            {Icon: "\ue795", Color: 113}, // 
	".idea":              {Icon: "\ue7b5", Color: 241}, // 
	".git":               {Icon: "\uf1d3", Color: 202}, // 
	".gitattributes":     {Icon: "\uf1d3", Color: 202}, // 
	".gitconfig":         {Icon: "\uf1d3", Color: 202}, // 
	".github":            {Icon: "\uf408", Color: 241}, // 
	".gitignore":         {Icon: "\uf1d3", Color: 202}, // 
	".gitmodules":        {Icon: "\uf1d3", Color: 202}, // 
	".mailmap":           {Icon: "\uf1d3", Color: 202}, // 
	".npmrc":             {Icon: "\ue71e", Color: 197}, // 
	".prettierrc":        {Icon: "\ue6b4", Color: 33},  // 
	".rvm":               {Icon: "\ue21e", Color: 160}, // 
	".SRCINFO":           {Icon: "\uf129", Color: 230}, // 
	".vimrc":             {Icon: "\ue62b", Color: 28},  // 
	".vscode":            {Icon: "\ue70c", Color: 39},  // 
	".Xauthority":        {Icon: "\uf369", Color: 196}, // 
	".xinitrc":           {Icon: "\uf369", Color: 196}, // 
	".Xresources":        {Icon: "\uf369", Color: 196}, // 
	".xsession":          {Icon: "\uf369", Color: 196}, // 
	".zshrc":             {Icon: "\ue795", Color: 113}, // 
	"bin":                {Icon: "\ue5fc", Color: 241}, // 
	"Cargo.lock":         {Icon: "\ue7a8", Color: 216}, // 
	"Cargo.toml":         {Icon: "\ue7a8", Color: 216}, // 
	"COMMIT_EDITMSG":     {Icon: "\ue702", Color: 239}, // 
	"compose.yaml":       {Icon: "\uf308", Color: 68},  // 
	"compose.yml":        {Icon: "\uf308", Color: 68},  // 
	"config":             {Icon: "\ue5fc", Color: 241}, // 
	"docker-compose.yml": {Icon: "\uf308", Color: 68},  // 
	"docker-compose.yaml":{Icon: "\uf308", Color: 68},  // 
	"Dockerfile":         {Icon: "\uf308", Color: 68},  // 
	"ds_store":           {Icon: "\uf179", Color: 15},  // 
	"favicon.ico":        {Icon: "\ue623", Color: 185}, // 
	"fp-info-cache":      {Icon: "\uf49b", Color: 231}, // 
	"fp-lib-table":       {Icon: "\uf34c", Color: 231}, // 
	"gitignore_global":   {Icon: "\uf1d3", Color: 202}, // 
	"GNUmakefile":        {Icon: "\ue779", Color: 66},  // 
	"go.mod":             {Icon: "\ue627", Color: 74},  // 
	"go.sum":             {Icon: "\ue627", Color: 74},  // 
	"gradle":             {Icon: "\ue256", Color: 168}, // 
	"gruntfile.coffee":   {Icon: "\ue611", Color: 166}, // 
	"gruntfile.js":       {Icon: "\ue611", Color: 166}, // 
	"gruntfile.ls":       {Icon: "\ue611", Color: 166}, // 
	"gtkrc":              {Icon: "\uf362", Color: 231}, // 
	"gulpfile.coffee":    {Icon: "\ue610", Color: 167}, // 
	"gulpfile.js":        {Icon: "\ue610", Color: 167}, // 
	"gulpfile.ls":        {Icon: "\ue610", Color: 168}, // 
	"hidden":             {Icon: "\uf023", Color: 241}, // 
	"include":            {Icon: "\ue5fc", Color: 241}, // 
	"lib":                {Icon: "\uf121", Color: 241}, // 
	"localized":          {Icon: "\uf179", Color: 15},  // 
	"Makefile":           {Icon: "\ue975", Color: 241}, // 
	"node_modules":       {Icon: "\ue718", Color: 197}, // 
	"npmignore":          {Icon: "\ue71e", Color: 197}, // 
	"PKGBUILD":           {Icon: "\uf303", Color: 38},  // 
	"rubydoc":            {Icon: "\ue73b", Color: 160}, // 
	"sym-lib-table":      {Icon: "\uf34c", Color: 231}, // 
	"xorg.conf":          {Icon: "\uf369", Color: 196}, // 
	"xsettingsd.conf":    {Icon: "\uf369", Color: 196}, // 
	"yarn.lock":          {Icon: "\ue6a7", Color: 74},  // 
}

var extIconMap = map[string]IconProperties{
	".3gp":            {Icon: "\uf03d", Color: 208},     // 
	".3mf":            {Icon: "\U000f01a7", Color: 102}, // 󰆧
	".aac":            {Icon: "\uf001", Color: 45},      // 
	".ai":             {Icon: "\ue7b4", Color: 185},     // 
	".android":        {Icon: "\ue70e", Color: 70},      // 
	".apk":            {Icon: "\ue70e", Color: 70},      // 
	".app":            {Icon: "\ueae8", Color: 124},     // 
	".apple":          {Icon: "\uf179", Color: 15},      // 
	".applescript":    {Icon: "\uf179", Color: 66},      // 
	".ass":            {Icon: "\U000f0a16", Color: 214}, // 󰨖
	".avi":            {Icon: "\uf03d", Color: 140},     // 
	".avif":           {Icon: "\uf1c5", Color: 140},     // 
	".avro":           {Icon: "\ue60b", Color: 130},     // 
	".awk":            {Icon: "\ue795", Color: 140},     // 
	".azcli":          {Icon: "\uebe8", Color: 32},      // 
	".bak":            {Icon: "\U000f006f", Color: 66},  // 󰁯
	".bash":           {Icon: "\ue795", Color: 113},     // 
	".bash_history":   {Icon: "\ue795", Color: 113},     // 
	".bash_profile":   {Icon: "\ue795", Color: 113},     // 
	".bashrc":         {Icon: "\ue795", Color: 113},     // 
	".bat":            {Icon: "\uf17a", Color: 81},      // 
	".bats":           {Icon: "\ue795", Color: 241},     // 
	".bazel":          {Icon: "\ue63a", Color: 113},     // 
	".bib":            {Icon: "\U000f125f", Color: 185}, // 󱉟
	".blend":          {Icon: "\U000f00ab", Color: 208}, // 󰂫
	".blp":            {Icon: "\U000f0ebe", Color: 68},  // 󰺾
	".bmp":            {Icon: "\uf1c5", Color: 149},     // 
	".brep":           {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".bz":             {Icon: "\uf410", Color: 239},     // 
	".bz2":            {Icon: "\uf410", Color: 239},     // 
	".bzl":            {Icon: "\ue63a", Color: 113},     // 
	".c":              {Icon: "\ue61e", Color: 111},     // 
	".c++":            {Icon: "\ue61d", Color: 204},     // 
	".cab":            {Icon: "\ue70f", Color: 241},     // 
	".cache":          {Icon: "\uf49b", Color: 231},     // 
	".cast":           {Icon: "\uf03d", Color: 208},     // 
	".cc":             {Icon: "\ue61d", Color: 204},     // 
	".cfg":            {Icon: "\ue615", Color: 255},     // 
	".class":          {Icon: "\ue256", Color: 168},     // 
	".clj":            {Icon: "\ue768", Color: 113},     // 
	".cljs":           {Icon: "\ue76a", Color: 74},      // 
	".cls":            {Icon: "\ue69b", Color: 239},     // 
	".cmd":            {Icon: "\ue70f", Color: 239},     // 
	".coffee":         {Icon: "\uf0f4", Color: 185},     // 
	".conf":           {Icon: "\ue615", Color: 66},      // 
	".cp":             {Icon: "\ue61d", Color: 74},      // 
	".cpio":           {Icon: "\uf410", Color: 239},     // 
	".cpp":            {Icon: "\ue61d", Color: 74},      // 
	".cs":             {Icon: "\U000f031b", Color: 58},  // 󰌛
	".csh":            {Icon: "\ue795", Color: 240},     // 
	".cshtml":         {Icon: "\uf1fa", Color: 239},     // 
	".csproj":         {Icon: "\U000f031b", Color: 58},  // 󰌛
	".css":            {Icon: "\ue749", Color: 75},      // 
	".csv":            {Icon: "\uf1c3", Color: 113},     // 
	".csx":            {Icon: "\U000f031b", Color: 58},  // 󰌛
	".cue":            {Icon: "\U000f0cb9", Color: 211}, // 󰲹
	".cxx":            {Icon: "\ue61d", Color: 74},      // 
	".d":              {Icon: "\ue7af", Color: 28},      // 
	".dart":           {Icon: "\ue798", Color: 25},      // 
	".dconf":          {Icon: "\ue706", Color: 188},     // 
	".db":             {Icon: "\uf1c0", Color: 188},     // 
	".deb":            {Icon: "\ue77d", Color: 88},      // 
	".desktop":        {Icon: "\uf108", Color: 54},      // 
	".diff":           {Icon: "\uf440", Color: 241},     // 
	".djvu":           {Icon: "\uf02d", Color: 241},     // 
	".dll":            {Icon: "\ue70f", Color: 241},     // 
	".doc":            {Icon: "\U000f0219", Color: 26},  // 󰈙
	".docx":           {Icon: "\U000f0219", Color: 26},  // 󰈙
	".dot":            {Icon: "\U000f1049", Color: 24},  // 󱁉
	".ds_store":       {Icon: "\uf179", Color: 15},      // 
	".DS_store":       {Icon: "\uf179", Color: 15},      // 
	".dump":           {Icon: "\uf1c0", Color: 188},     // 
	".dwg":            {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".dxf":            {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".ebook":          {Icon: "\ue28b", Color: 241},     // 
	".ebuild":         {Icon: "\uf30d", Color: 56},      // 
	".editorconfig":   {Icon: "\ue615", Color: 241},     // 
	".ejs":            {Icon: "\ue618", Color: 185},     // 
	".elf":            {Icon: "\ueae8", Color: 124},     // 
	".elm":            {Icon: "\ue62c", Color: 74},      // 
	".env":            {Icon: "\uf462", Color: 227},     // 
	".eot":            {Icon: "\uf031", Color: 124},     // 
	".epub":           {Icon: "\ue28a", Color: 241},     // 
	".erb":            {Icon: "\ue73b", Color: 160},     // 
	".erl":            {Icon: "\ue7b1", Color: 163},     // 
	".ex":             {Icon: "\ue62d", Color: 140},     // 
	".exe":            {Icon: "\uf17a", Color: 81},      // 
	".exs":            {Icon: "\ue62d", Color: 140},     // 
	".f90":            {Icon: "\U000f121a", Color: 97},  // 󱈚
	".fbx":            {Icon: "\U000f01a7", Color: 102}, // 󰆧
	".fcbak":          {Icon: "\uf336", Color: 160},     // 
	".fcmacro":        {Icon: "\uf336", Color: 160},     // 
	".fcmat":          {Icon: "\uf336", Color: 160},     // 
	".fcparam":        {Icon: "\uf336", Color: 160},     // 
	".fcscript":       {Icon: "\uf336", Color: 160},     // 
	".fcstd":          {Icon: "\uf336", Color: 160},     // 
	".fcstd1":         {Icon: "\uf336", Color: 160},     // 
	".fctb":           {Icon: "\uf336", Color: 160},     // 
	".fctl":           {Icon: "\uf336", Color: 160},     // 
	".fish":           {Icon: "\ue795", Color: 249},     // 
	".flac":           {Icon: "\uf001", Color: 241},     // 
	".flc":            {Icon: "\uf031", Color: 255},     // 
	".flf":            {Icon: "\uf031", Color: 255},     // 
	".flv":            {Icon: "\uf03d", Color: 241},     // 
	".fnl":            {Icon: "\ue6af", Color: 230},     // 
	".font":           {Icon: "\uf031", Color: 241},     // 
	".fs":             {Icon: "\ue7a7", Color: 74},      // 
	".fsi":            {Icon: "\ue7a7", Color: 74},      // 
	".fsx":            {Icon: "\ue7a7", Color: 74},      // 
	".gcode":          {Icon: "\U000f0af4", Color: 234}, // 󰫴
	".gd":             {Icon: "\ue65f", Color: 66},      // 
	".gdoc":           {Icon: "\uf1c2", Color: 40},      // 
	".gem":            {Icon: "\ue21e", Color: 160},     // 
	".gemfile":        {Icon: "\ue21e", Color: 160},     // 
	".gemspec":        {Icon: "\ue21e", Color: 160},     // 
	".gform":          {Icon: "\uf298", Color: 40},      // 
	".gif":            {Icon: "\uf1c5", Color: 140},     // 
	".git":            {Icon: "\uf1d3", Color: 202},     // 
	".gitattributes":  {Icon: "\uf1d3", Color: 202},     // 
	".gitignore":      {Icon: "\uf1d3", Color: 202},     // 
	".gitmodules":     {Icon: "\uf1d3", Color: 202},     // 
	".go":             {Icon: "\ue627", Color: 74},      // 
	".godot":          {Icon: "\ue65f", Color: 66},      // 
	".gql":            {Icon: "\uf20e", Color: 199},     // 
	".graphql":        {Icon: "\uf20e", Color: 199},     // 
	".gradle":         {Icon: "\ue256", Color: 168},     // 
	".gresource":      {Icon: "\uf362", Color: 231},     // 
	".groovy":         {Icon: "\ue775", Color: 24},      // 
	".gsheet":         {Icon: "\uf1c3", Color: 10},      // 
	".gslides":        {Icon: "\uf1c4", Color: 226},     // 
	".guardfile":      {Icon: "\ue21e", Color: 241},     // 
	".gv":             {Icon: "\U000f1049", Color: 24},  // 󱁉
	".gz":             {Icon: "\uf410", Color: 241},     // 
	".h":              {Icon: "\uf0fd", Color: 140},     // 
	".hbs":            {Icon: "\ue60f", Color: 202},     // 
	".hc":             {Icon: "\U000f00a2", Color: 227}, // 󰂢
	".hex":            {Icon: "\U000f12a7", Color: 27},  // 󱊧
	".hh":             {Icon: "\uf0fd", Color: 140},     // 
	".hpp":            {Icon: "\uf0fd", Color: 140},     // 
	".hs":             {Icon: "\ue777", Color: 140},     // 
	".htm":            {Icon: "\uf13b", Color: 196},     // 
	".html":           {Icon: "\uf13b", Color: 196},     // 
	".hxx":            {Icon: "\uf0fd", Color: 140},     // 
	".ico":            {Icon: "\uf1c5", Color: 185},     // 
	".ifc":            {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".ige":            {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".iges":           {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".igs":            {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".image":          {Icon: "\uf1c5", Color: 185},     // 
	".iml":            {Icon: "\ue7b5", Color: 239},     // 
	".info":           {Icon: "\uf129", Color: 230},     // 
	".ini":            {Icon: "\uf17a", Color: 81},      // 
	".ino":            {Icon: "\uf34b", Color: 73},      // 
	".ipynb":          {Icon: "\ue606", Color: 214},     // 
	".iso":            {Icon: "\ue271", Color: 239},     // 
	".j2c":            {Icon: "\uf1c5", Color: 239},     // 
	".j2k":            {Icon: "\uf1c5", Color: 239},     // 
	".jad":            {Icon: "\ue256", Color: 168},     // 
	".jar":            {Icon: "\ue256", Color: 168},     // 
	".java":           {Icon: "\ue256", Color: 168},     // 
	".jfi":            {Icon: "\uf1c5", Color: 241},     // 
	".jfif":           {Icon: "\uf1c5", Color: 241},     // 
	".jif":            {Icon: "\uf1c5", Color: 241},     // 
	".jl":             {Icon: "\ue624", Color: 241},     // 
	".jmd":            {Icon: "\uf48a", Color: 74},      // 
	".jp2":            {Icon: "\uf1c5", Color: 241},     // 
	".jpe":            {Icon: "\uf1c5", Color: 241},     // 
	".jpeg":           {Icon: "\uf1c5", Color: 241},     // 
	".jpg":            {Icon: "\uf1c5", Color: 241},     // 
	".jpx":            {Icon: "\uf1c5", Color: 241},     // 
	".js":             {Icon: "\ue74e", Color: 185},     // 
	".json":           {Icon: "\ue60b", Color: 185},     // 
	".jsx":            {Icon: "\ue7ba", Color: 45},      // 
	".jxl":            {Icon: "\uf1c5", Color: 241},     // 
	".kdenlive":       {Icon: "\uf33c", Color: 110},     // 
	".kdenlivetitle":  {Icon: "\uf33c", Color: 110},     // 
	".kicad_dru":      {Icon: "\uf34c", Color: 231},     // 
	".kicad_mod":      {Icon: "\uf34c", Color: 231},     // 
	".kicad_pcb":      {Icon: "\uf34c", Color: 231},     // 
	".kicad_prl":      {Icon: "\uf34c", Color: 231},     // 
	".kicad_pro":      {Icon: "\uf34c", Color: 231},     // 
	".kicad_sch":      {Icon: "\uf34c", Color: 231},     // 
	".kicad_sym":      {Icon: "\uf34c", Color: 231},     // 
	".kicad_wks":      {Icon: "\uf34c", Color: 231},     // 
	".kpp":            {Icon: "\uf33d", Color: 201},     // 
	".kra":            {Icon: "\uf33d", Color: 201},     // 
	".krz":            {Icon: "\uf33d", Color: 201},     // 
	".ksh":            {Icon: "\ue795", Color: 241},     // 
	".kt":             {Icon: "\ue634", Color: 99},      // 
	".kts":            {Icon: "\ue634", Color: 99},      // 
	".latex":          {Icon: "\ue69b", Color: 241},     // 
	".lck":            {Icon: "\ue672", Color: 250},     // 
	".less":           {Icon: "\ue758", Color: 54},      // 
	".lhs":            {Icon: "\ue777", Color: 140},     // 
	".license":        {Icon: "\U000f0219", Color: 185}, // 󰈙
	".localized":      {Icon: "\uf179", Color: 15},      // 
	".lock":           {Icon: "\uf023", Color: 241},     // 
	".log":            {Icon: "\uf4ed", Color: 188},     // 
	".lrc":            {Icon: "\U000f0a16", Color: 214}, // 󰨖
	".lua":            {Icon: "\ue620", Color: 74},      // 
	".luac":           {Icon: "\ue620", Color: 74},      // 
	".luau":           {Icon: "\ue620", Color: 74},      // 
	".lz":             {Icon: "\uf410", Color: 241},     // 
	".lz4":            {Icon: "\uf410", Color: 241},     // 
	".lzh":            {Icon: "\uf410", Color: 241},     // 
	".lzma":           {Icon: "\uf410", Color: 241},     // 
	".lzo":            {Icon: "\uf410", Color: 241},     // 
	".m":              {Icon: "\ue61e", Color: 111},     // 
	".m3u":            {Icon: "\U000f0cb9", Color: 211}, // 󰲹
	".m3u8":           {Icon: "\U000f0cb9", Color: 211}, // 󰲹
	".m4a":            {Icon: "\uf001", Color: 239},     // 
	".m4v":            {Icon: "\uf03d", Color: 208},     // 
	".magnet":         {Icon: "\uf076", Color: 124},     // 
	".makefile":       {Icon: "\ue779", Color: 66},      // 
	".markdown":       {Icon: "\uf48a", Color: 74},      // 
	".md":             {Icon: "\uf48a", Color: 74},      // 
	".mdx":            {Icon: "\uf48a", Color: 74},      // 
	".mjs":            {Icon: "\ue74e", Color: 185},     // 
	".mk":             {Icon: "\ue795", Color: 241},     // 
	".mkd":            {Icon: "\uf48a", Color: 74},      // 
	".mkv":            {Icon: "\uf03d", Color: 241},     // 
	".mm":             {Icon: "\ue61d", Color: 111},     // 
	".mobi":           {Icon: "\ue28b", Color: 241},     // 
	".mov":            {Icon: "\uf03d", Color: 241},     // 
	".mp3":            {Icon: "\uf001", Color: 241},     // 
	".mp4":            {Icon: "\uf03d", Color: 241},     // 
	".msi":            {Icon: "\ue70f", Color: 241},     // 
	".mustache":       {Icon: "\ue60f", Color: 241},     // 
	".nfo":            {Icon: "\uf129", Color: 230},     // 
	".nix":            {Icon: "\uf313", Color: 111},     // 
	".node":           {Icon: "\U000f0399", Color: 197}, // 󰎙
	".npmignore":      {Icon: "\ue71e", Color: 197},     // 
	".obj":            {Icon: "\U000f01a7", Color: 102}, // 󰆧
	".odp":            {Icon: "\uf1c4", Color: 241},     // 
	".ods":            {Icon: "\uf1c3", Color: 241},     // 
	".odt":            {Icon: "\uf1c2", Color: 241},     // 
	".ogg":            {Icon: "\uf001", Color: 241},     // 
	".ogv":            {Icon: "\uf03d", Color: 241},     // 
	".opus":           {Icon: "\U000f0223", Color: 208}, // 󰈣
	".org":            {Icon: "\ue633", Color: 73},      // 
	".otf":            {Icon: "\uf031", Color: 241},     // 
	".out":            {Icon: "\ueae8", Color: 124},     // 
	".part":           {Icon: "\uf43a", Color: 241},     // 
	".patch":          {Icon: "\uf440", Color: 241},     // 
	".pdf":            {Icon: "\uf1c1", Color: 124},     // 
	".php":            {Icon: "\ue73d", Color: 61},      // 
	".pl":             {Icon: "\ue769", Color: 74},      // 
	".pls":            {Icon: "\U000f0cb9", Color: 211}, // 󰲹
	".ply":            {Icon: "\U000f01a7", Color: 102}, // 󰆧
	".png":            {Icon: "\uf1c5", Color: 241},     // 
	".po":             {Icon: "\U000f05ca", Color: 31},  // 󰗊
	".pot":            {Icon: "\U000f05ca", Color: 31},  // 󰗊
	".ppt":            {Icon: "\uf1c4", Color: 241},     // 
	".pptx":           {Icon: "\uf1c4", Color: 241},     // 
	".procfile":       {Icon: "\ue21e", Color: 241},     // 
	".properties":     {Icon: "\ue60b", Color: 185},     // 
	".ps1":            {Icon: "\ue795", Color: 241},     // 
	".psd":            {Icon: "\ue7b8", Color: 241},     // 
	".pxm":            {Icon: "\uf1c5", Color: 241},     // 
	".py":             {Icon: "\ue606", Color: 214},     // 
	".pyc":            {Icon: "\ue606", Color: 214},     // 
	".qm":             {Icon: "\U000f05ca", Color: 31},  // 󰗊
	".qml":            {Icon: "\uf375", Color: 77},      // 
	".qrc":            {Icon: "\uf375", Color: 77},      // 
	".qss":            {Icon: "\uf375", Color: 77},      // 
	".r":              {Icon: "\uf25d", Color: 68},      // 
	".rakefile":       {Icon: "\ue21e", Color: 160},     // 
	".rar":            {Icon: "\uf410", Color: 241},     // 
	".razor":          {Icon: "\uf1fa", Color: 81},      // 
	".rb":             {Icon: "\ue21e", Color: 160},     // 
	".rdata":          {Icon: "\uf25d", Color: 68},      // 
	".rdb":            {Icon: "\ue76d", Color: 160},     // 
	".rdoc":           {Icon: "\uf48a", Color: 74},      // 
	".rds":            {Icon: "\uf25d", Color: 68},      // 
	".readme":         {Icon: "\uf48a", Color: 74},      // 
	".rlib":           {Icon: "\ue7a8", Color: 216},     // 
	".rmd":            {Icon: "\uf48a", Color: 74},      // 
	".rpm":            {Icon: "\ue7bb", Color: 52},      // 
	".rs":             {Icon: "\ue7a8", Color: 216},     // 
	".rspec":          {Icon: "\ue21e", Color: 160},     // 
	".rspec_parallel": {Icon: "\ue21e", Color: 160},     // 
	".rspec_status":   {Icon: "\ue21e", Color: 160},     // 
	".rss":            {Icon: "\uf09e", Color: 130},     // 
	".rtf":            {Icon: "\U000f0219", Color: 241}, // 󰈙
	".ru":             {Icon: "\ue21e", Color: 160},     // 
	".rubydoc":        {Icon: "\ue73b", Color: 160},     // 
	".sass":           {Icon: "\ue603", Color: 169},     // 
	".scad":           {Icon: "\uf34e", Color: 220},     // 
	".scala":          {Icon: "\ue737", Color: 74},      // 
	".scss":           {Icon: "\ue749", Color: 204},     // 
	".sh":             {Icon: "\ue795", Color: 239},     // 
	".shell":          {Icon: "\ue795", Color: 239},     // 
	".skp":            {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".sldasm":         {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".sldprt":         {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".slim":           {Icon: "\ue73b", Color: 160},     // 
	".slvs":           {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".sln":            {Icon: "\ue70c", Color: 39},      // 
	".so":             {Icon: "\uf17c", Color: 241},     // 
	".sql":            {Icon: "\uf1c0", Color: 188},     // 
	".sqlite":         {Icon: "\ue7c4", Color: 25},      // 
	".sqlite3":        {Icon: "\ue7c4", Color: 25},      // 
	".srt":            {Icon: "\U000f0a16", Color: 214}, // 󰨖
	".ssa":            {Icon: "\U000f0a16", Color: 214}, // 󰨖
	".ste":            {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".step":           {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".stl":            {Icon: "\U000f01a7", Color: 102}, // 󰆧
	".sty":            {Icon: "\ue69b", Color: 239},     // 
	".styl":           {Icon: "\ue600", Color: 148},     // 
	".stylus":         {Icon: "\ue600", Color: 148},     // 
	".stp":            {Icon: "\U000f0eeb", Color: 101}, // 󰻫
	".sub":            {Icon: "\U000f0a16", Color: 214}, // 󰨖
	".sv":             {Icon: "\U000f035b", Color: 28},  // 󰍛
	".svh":            {Icon: "\U000f035b", Color: 28},  // 󰍛
	".svelte":         {Icon: "\ue697", Color: 208},     // 
	".svg":            {Icon: "\uf1c5", Color: 241},     // 
	".swift":          {Icon: "\ue755", Color: 208},     // 
	".tar":            {Icon: "\uf410", Color: 241},     // 
	".taz":            {Icon: "\uf410", Color: 241},     // 
	".tbz":            {Icon: "\uf410", Color: 241},     // 
	".tbz2":           {Icon: "\uf410", Color: 241},     // 
	".tex":            {Icon: "\ue69b", Color: 79},      // 
	".tgz":            {Icon: "\uf410", Color: 241},     // 
	".tiff":           {Icon: "\uf1c5", Color: 241},     // 
	".tlz":            {Icon: "\uf410", Color: 241},     // 
	".toml":           {Icon: "\ue615", Color: 241},     // 
	".torrent":        {Icon: "\ue275", Color: 76},      // 
	".tres":           {Icon: "\ue65f", Color: 66},      // 
	".tscn":           {Icon: "\ue65f", Color: 66},      // 
	".ts":             {Icon: "\ue628", Color: 74},      // 
	".tsv":            {Icon: "\uf1c3", Color: 241},     // 
	".tsx":            {Icon: "\ue7ba", Color: 74},      // 
	".ttf":            {Icon: "\uf031", Color: 241},     // 
	".twig":           {Icon: "\ue61c", Color: 241},     // 
	".txt":            {Icon: "\uf15c", Color: 241},     // 
	".txz":            {Icon: "\uf410", Color: 241},     // 
	".tz":             {Icon: "\uf410", Color: 241},     // 
	".tzo":            {Icon: "\uf410", Color: 241},     // 
	".ui":             {Icon: "\uf2d0", Color: 17},      // 
	".v":              {Icon: "\U000f035b", Color: 28},  // 󰍛
	".vh":             {Icon: "\U000f035b", Color: 28},  // 󰍛
	".vhd":            {Icon: "\U000f035b", Color: 28},  // 󰍛
	".vhdl":           {Icon: "\U000f035b", Color: 28},  // 󰍛
	".video":          {Icon: "\uf03d", Color: 241},     // 
	".vim":            {Icon: "\ue62b", Color: 28},      // 
	".vsix":           {Icon: "\ue70c", Color: 98},      // 
	".vue":            {Icon: "\U000f0844", Color: 113}, // 󰡄
	".war":            {Icon: "\ue256", Color: 168},     // 
	".wav":            {Icon: "\uf001", Color: 241},     // 
	".webm":           {Icon: "\uf03d", Color: 241},     // 
	".webp":           {Icon: "\uf1c5", Color: 241},     // 
	".windows":        {Icon: "\uf17a", Color: 81},      // 
	".woff":           {Icon: "\uf031", Color: 241},     // 
	".woff2":          {Icon: "\uf031", Color: 241},     // 
	".wrl":            {Icon: "\U000f01a7", Color: 102}, // 󰆧
	".wrz":            {Icon: "\U000f01a7", Color: 102}, // 󰆧
	".xcf":            {Icon: "\uf338", Color: 240},     // 
	".xhtml":          {Icon: "\uf13b", Color: 196},     // 
	".xls":            {Icon: "\uf1c3", Color: 34},      // 
	".xlsx":           {Icon: "\uf1c3", Color: 34},      // 
	".xml":            {Icon: "\uf121", Color: 160},     // 
	".xpi":            {Icon: "\ueae6", Color: 17},      // 
	".xul":            {Icon: "\uf121", Color: 166},     // 
	".xz":             {Icon: "\uf410", Color: 241},     // 
	".yaml":           {Icon: "\uf481", Color: 160},     // 
	".yml":            {Icon: "\uf481", Color: 160},     // 
	".zig":            {Icon: "\ue6a9", Color: 172},     // 
	".zip":            {Icon: "\uf410", Color: 241},     // 
	".zsh":            {Icon: "\ue795", Color: 241},     // 
	".zsh-theme":      {Icon: "\ue795", Color: 241},     // 
	".zshrc":          {Icon: "\ue795", Color: 241},     // 
	".zst":            {Icon: "\uf410", Color: 241},     // 
}

func patchFileIconsForNerdFontsV2() {
	extIconMap[".cs"] = IconProperties{Icon: "\uf81a", Color: 58}       // 
	extIconMap[".csproj"] = IconProperties{Icon: "\uf81a", Color: 58}   // 
	extIconMap[".csx"] = IconProperties{Icon: "\uf81a", Color: 58}      // 
	extIconMap[".license"] = IconProperties{Icon: "\uf718", Color: 241} // 
	extIconMap[".node"] = IconProperties{Icon: "\uf898", Color: 197}    // 
	extIconMap[".rtf"] = IconProperties{Icon: "\uf718", Color: 241}     // 
	extIconMap[".vue"] = IconProperties{Icon: "\ufd42", Color: 113}     // ﵂
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
		return IconProperties{LINKED_WORKTREE_ICON, 239}
	} else if isDirectory {
		return DEFAULT_DIRECTORY_ICON
	}
	return DEFAULT_FILE_ICON
}
