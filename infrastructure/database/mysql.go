package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)



type SqlHandler struct {
	Db *gorm.DB
}

func EstablishConnection(isProduction bool) (*SqlHandler, error) {
	user, ok := os.LookupEnv("MARIADB_USERNAME")
	if !ok {
		user = "root"
	}

	pass, ok := os.LookupEnv("MARIADB_PASSWORD")
	if !ok {
		pass = "password"
	}

	host, ok := os.LookupEnv("MARIADB_HOSTNAME")
	if !ok {
		host = "localhost"
	}

	dbname, ok := os.LookupEnv("MARIADB_DATABASE")
	if !ok {
		dbname = "anke-two"
	}

	var loglevel logger.LogLevel
	if isProduction {
		loglevel = logger.Silent
	} else {
		loglevel = logger.Info
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname) + "?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4"
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(loglevel)})
	_db = _db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci")
	/*_db.Use(prometheus.New(prometheus.Config{
		DBName:          "anke-two",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&database.MetricsCollector{},
		},
	}))*/
	sqlHandler := new(SqlHandler)
	sqlHandler.Db = _db
	return sqlHandler, err
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(allTables...)
	if err != nil {
		return fmt.Errorf("failed in table's migration: %w", err)
	}

	return nil
}
