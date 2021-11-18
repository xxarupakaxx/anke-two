package infrastructure

import (
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
	"os"
)

var (
	db        *gorm.DB
	allTables = []interface{}{
		model.Questionnaires{},
		model.Questions{},
		model.Respondents{},
		model.Responses{},
		model.Administrators{},
		model.Options{},
		model.ScaleLabels{},
		model.ResShareTypes{},
		model.QuestionType{},
		model.Targets{},
		model.Validations{},
	}
)

func EstablishConnection(isProduction bool) error {
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
	db.Use(prometheus.New(prometheus.Config{
		DBName:          "anke-two",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&usecase.MetricsCollector{},
		},
	}))
	db = _db
	return err
}

func Migrate() error {
	err := db.AutoMigrate(allTables...)
	if err != nil {
		return fmt.Errorf("failed in table's migration: %w", err)
	}

	return nil
}
