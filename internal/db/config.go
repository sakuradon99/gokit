package db

type Config struct {
	host     string `value:"db.host"`
	port     int    `value:"db.port"`
	user     string `value:"db.user"`
	password string `value:"db.password"`
	database string `value:"db.database"`
}
