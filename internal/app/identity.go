package app

const (
	name    = "mc-pacman"
	version = "0.1"
)

// UserAgent returns the user agent string for http requests
func UserAgent() string {
	return name + "-v" + version
}
