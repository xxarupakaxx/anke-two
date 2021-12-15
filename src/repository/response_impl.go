package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) CreateResponse(ctx context.Context, response *model.Response) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db :%w", err)
	}

	err = db.Create(&response).Error
	if err != nil {
		return fmt.Errorf("failed to create reponse:%w", err)
	}

	return nil
}

func (repo *GormRepository) DeleteResponse(ctx context.Context, responseID int) error {
	panic("implement me")
}
