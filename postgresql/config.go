package postgresql

type Config struct {
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	SSLMode    string
	AutoCreate bool
}
