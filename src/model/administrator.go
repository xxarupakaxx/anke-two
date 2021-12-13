package model

// Administrator アンケートの管理者 (編集等ができる人)のテーブル
type Administrator struct {
	QuestionnaireID int    `gorm:"type:int(11);not null;primaryKey"`
	TraqID          string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

func (a *Administrator) TableName() string {
	return "administrator"
}