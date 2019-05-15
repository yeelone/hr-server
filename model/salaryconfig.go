package model

type SalaryConfig struct {
	BaseModel
	Base         float64 `json:"base"`          //基本工资设定
	TaxThreshold float64 `json:"tax_threshold"` //个税起征点
	Reference    string  //根据发文
}

var SALARYCONFIGTABLENAME = "tb_salary_config"

// TableName :
func (s *SalaryConfig) TableName() string {
	return SALARYCONFIGTABLENAME
}

// Create creates a new salary config .
func (s *SalaryConfig) Create() error {
	return DB.Self.Create(&s).Error
}

func GetSalaryConfig() (s SalaryConfig, err error) {
	err = DB.Self.Model(SalaryConfig{}).Last(&s).Error
	return s, err
}
