package model

// Validation Numberの値制限、Textの正規表現によるパターンマッチングのテーブル
type Validation struct {
	QuestionID   int    `gorm:"type:int(11);not null;primaryKey"`
	RegexPattern string `gorm:"type:text;default:NULL"`
	MinBound     string `gorm:"type:text;default:NULL"`
	MaxBound     string `gorm:"type:text;default:NULL"`
}
