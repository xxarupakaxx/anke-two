package model

type ScaleLabels struct {
	QuestionID      int    `json:"questionID"`
	ScaleLabelRight string `json:"scale_label_right"`
	ScaleLabelLeft  string `json:"scale_label_left"`
	ScaleMin        int    `json:"scale_min"`
	ScaleMax        int    `json:"scale_max"`
}
