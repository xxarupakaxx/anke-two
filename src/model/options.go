package model

// Options 選択肢のテーブル
type Options struct {
	ID         int    `gorm:"type:int(11);AUTO_INCREMENT;not null;primaryKey"`
	QuestionID int    `gorm:"type:int(11);not null"`
	OptionNum  int    `gorm:"type:int(11);not null"`
	Body       string `gorm:"type:text;default:NULL;"`
}