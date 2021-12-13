package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/xxarupkaxx/anke-two/config"
	"github.com/xxarupkaxx/anke-two/src/model"
	"gorm.io/gorm"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
)

var allTables = []interface{}{
	model.Question{},
	model.QuestionType{},
	model.Questionnaire{},
	model.ResSharedTo{},
	model.Response{},
	model.Respondent{},
	model.Option{},
	model.Validation{},
	model.Respondent{},
	model.Target{},
	model.Administrator{},
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(c *config.Config) *GormRepository {
	db, err := connectDB(c)
	if err != nil {
		log.Panicf("failed to connect db:%w", err)
	}

	return &GormRepository{db: db}
}

func connectDB(c *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", c.MariaDBUsername, c.MariaDBPassword, c.MariaDBHostname, c.MariaDBDatabase) + "?parseTime=true&loc=Local&charset=utf8mb4"
	logLevel := logger.Info

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logLevel)})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB : %w", err)
	}

	db = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci")

	err = db.AutoMigrate(allTables...)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}

//getDB DBをコンテキストから取得
func (repo *GormRepository) getDB(ctx context.Context) (db *gorm.DB, err error) {
	iDB := ctx.Value(txKey)
	if iDB == nil {
		return repo.db.WithContext(ctx), nil
	}

	gormDB, ok := iDB.(*gorm.DB)
	if !ok {
		return nil, errors.New("failed to get gorm.DB")
	}

	return gormDB.WithContext(ctx), nil
}
