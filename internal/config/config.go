package config

type Config struct {
	HttpServer HttpServer `yaml:"HttpServer"`
	Database   Database   `yaml:"Database"`
	Auth       Auth       `yaml:"Auth"`
}

type Database struct {
	User       string `yaml:"User"`
	Password   string `yaml:"Password"`
	Port       string `yaml:"Port"`
	Host       string `yaml:"Host"`
	Database   string `yaml:"Database"`
	Collection string `yaml:"Collection"`
}

type HttpServer struct {
	Port int `yaml:"Port"`
}

type Auth struct {
	JwtSecretKey string `yaml:"JwtSecretKey"`
}
