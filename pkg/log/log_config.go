package log

type LogConfig struct {
	AccessFile string `koanf:"access_file"`
	ErrorFile  string `koanf:"error_file"`
}
