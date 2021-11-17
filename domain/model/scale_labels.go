package model

type ScaleLabels struct {
	QuestionID      int    `json:"questionID" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	ScaleLabelRight string `json:"scale_label_right" gorm:"type:varchar(50);default:NULL;"`
	ScaleLabelLeft  string `json:"scale_label_left" gorm:"type:varchar(50);default:NULL;"`
	ScaleMin        int    `json:"scale_min" gorm:"type:int(11);default:NULL;"`
	ScaleMax        int    `json:"scale_max" gorm:"type:int(11);default:NULL;"`
}
