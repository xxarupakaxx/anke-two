package model

// ScaleLabel 目盛り (LinearScale) 形式の質問の左右のラベルのテーブル
type ScaleLabel struct {
	QuestionID      int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	ScaleLabelRight string `gorm:"type:varchar(50);default:NULL;"`
	ScaleLabelLeft  string `gorm:"type:varchar(50);default:NULL;"`
	ScaleMin        int    `gorm:"type:int(11);default:NULL;"`
	ScaleMax        int    `gorm:"type:int(11);default:NULL;"`
}
