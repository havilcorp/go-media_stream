package config

type Config struct {
	AllowTypes []string
	DBConnect  string
}

func NewConfig() *Config {
	return &Config{
		DBConnect:  "usermysql:fmPWlkBjjLgSeapA@tcp(192.168.0.118:3306)/media_stream?parseTime=true",
		AllowTypes: []string{"video/mp4", "audio/mpeg", "image/png"},
	}
}
