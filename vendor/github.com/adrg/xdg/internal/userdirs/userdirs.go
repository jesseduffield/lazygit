package userdirs

// XDG user directories environment variables.
const (
	EnvDesktopDir     = "XDG_DESKTOP_DIR"
	EnvDownloadDir    = "XDG_DOWNLOAD_DIR"
	EnvDocumentsDir   = "XDG_DOCUMENTS_DIR"
	EnvMusicDir       = "XDG_MUSIC_DIR"
	EnvPicturesDir    = "XDG_PICTURES_DIR"
	EnvVideosDir      = "XDG_VIDEOS_DIR"
	EnvTemplatesDir   = "XDG_TEMPLATES_DIR"
	EnvPublicShareDir = "XDG_PUBLICSHARE_DIR"
)

// Directories defines the locations of well known user directories.
type Directories struct {
	// Desktop defines the location of the user's desktop directory.
	Desktop string

	// Download defines a suitable location for user downloaded files.
	Download string

	// Documents defines a suitable location for user document files.
	Documents string

	// Music defines a suitable location for user audio files.
	Music string

	// Pictures defines a suitable location for user image files.
	Pictures string

	// VideosDir defines a suitable location for user video files.
	Videos string

	// Templates defines a suitable location for user template files.
	Templates string

	// PublicShare defines a suitable location for user shared files.
	PublicShare string
}
