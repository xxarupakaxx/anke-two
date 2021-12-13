package model

// Target アンケートの対象者のテーブル
type Target struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

func (t *Target) TableName() string {
	return "target"
}