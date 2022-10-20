package conf

type (
	DB struct {
		Driver   string // sqlite, mysql, postgres
		MySQL    mysql
		Postgres postgres
		Sqlite   sqlite
	}

	mysql struct {
		Host       string
		Port       int
		User       string
		Password   string
		DBName     string
		Parameters string
	}
	postgres struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	sqlite struct {
		Path string
	}
)

func init() {
	DefaultConf = append(DefaultConf, DB{
		Driver: "postgres",
		Sqlite: sqlite{
			Path: "db.db",
		},
		MySQL: mysql{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "root",
			Password: "666666",
			DBName:   "zls",
		},
		Postgres: postgres{
			Host:     "127.0.0.1",
			Port:     5432,
			User:     "root",
			Password: "666666",
			DBName:   "zls",
		},
	})
}
