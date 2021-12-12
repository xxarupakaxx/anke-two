package model

// Administrators アンケートの管理者 (編集等ができる人)のテーブル
type Administrators struct {
	QuestionnaireID int    `gorm:"type:int(11);not null;primaryKey"`
	TraqID          string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}
