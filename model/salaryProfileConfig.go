package model

import "hr-server/pkg/template"

// 计算工资之前 对员工的一些特殊配置
type SalaryProfileConfig struct {
	BaseModel
	TemplateFieldId string         `json:"template_field_id"`
	TemplateField   template.Field `json:"template_field" gorm:"-"`
	ProfileId       uint64
	Profile         Profile `json:"profile"`
	Operate         string  `json:"operate"` // +  - * /
	Value           float64 `json:"value"`
	Description     string  `json:"description"`
}

var SALARYProfileCONFIGTABLENAME = "tb_salary_profile_config"

// TableName :
func (s *SalaryProfileConfig) TableName() string {
	return SALARYProfileCONFIGTABLENAME
}

// Create creates a new salary config .
func (s *SalaryProfileConfig) Create() error {
	return DB.Self.Create(&s).Error
}

// Create creates a new salary config .
func (s *SalaryProfileConfig) Delete() error {
	return DB.Self.Delete(&s).Error
}

func GetSalaryProfileConfig() (list []SalaryProfileConfig, err error) {
	err = DB.Self.Model(SalaryProfileConfig{}).Preload("Profile").Find(&list).Error
	return list, err
}
