package config

type Config struct {
	ServerAddress  string
	AllowTypes     []string
	DBConnect      string
	FilmDecWrCount int
}

func NewConfig() *Config {
	return &Config{
		ServerAddress:  ":8080",
		DBConnect:      "usermysql:fmPWlkBjjLgSeapA@tcp(192.168.0.118:3306)/media_stream?parseTime=true",
		AllowTypes:     []string{"video/mp4", "audio/mpeg", "image/png"},
		FilmDecWrCount: 2,
	}
}
