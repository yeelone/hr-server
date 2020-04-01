package salary

import (
	"hr-server/model"
	"hr-server/pkg/template"
)

type CreateRequest struct {
	ID                  uint64                    `json:"id"`
	BaseSalary          float64                   `json:"base"`
	TaxThreshold        float64                   `json:"tax_threshold"`         //个税起征点
	TemplateAccountID   uint64                    `json:"template_account_id"`   //账套名字
	TemplateAccountName string                    `json:"template_account_name"` //账套名字
	Name                string                    `json:"name"`
	Type                string                    `json:"type" `  //可以是普通模板和分摊模板
	Months              int                       `json:"months"` //如果分摊模板，可以设置分几个月摊完
	InitData            string                    `json:"file"`
	Startup             bool                      `json:"startup"`
	Year                string                    `json:"year"`
	Month               string                    `json:"month"`
	Password            string                    `json:"password"`
	Group               uint64                    `json:"group" `
	Body                map[string]template.Field `json:"body"` //模板配置内容,从客户端发送的json格式，在服务器端要转化为yaml并保存到template文件夹中
	Remark              string                    `json:"remark"`
}
type CreateProfileConfigCreateRequest struct {
	ProfileID       uint64  `json:"profile_id"`
	TemplateFieldID string  `json:"template_field_id"`
	Operate         string  `json:"operate"`
	Value           float64 `json:"value"`
	Description     string  `json:"description"`
}

type SalaryProfileConfigResponse struct {
	ConfigList []model.SalaryProfileConfig `json:"config_list"`
}
type DeleteRequest struct {
	ID uint64 `json:"id"`
}
type ProfileSalaryRequest struct {
	ProfileID uint64 `json:"profile_id"`
	Year      string `json:"year"`
	Month     string `json:"month"`
}

type ProfileSalaryResponse struct {
	Profile      model.Profile                 `json:"profile"`
	Department   string                        `json:"department"`
	Post         string                        `json:"post"`
	TemplateList []map[string][]template.Field `json:"template_list"`
}

type TemplateOrderRequest struct {
	Orders map[uint64]int `form:"orders"`
}

type TaxRequest struct {
	model.TaxConf
}

type TaxResponse struct {
	Conf model.TaxConf `json:"conf"`
}

type CreateResponse struct {
	SalaryData       map[uint64]*Row
	ErrorMessage     []string
	ErrorMessageFile string
	Base             float64
	TaxThreshold     float64 //个税起征点
}

type TemplateResponse struct {
	Template *model.Template
	Fields   []template.Field
}

type ListRequest struct {
}

type ListResponse struct {
	Data []map[string]interface{}
}

type ListTemplateResponse struct {
	List []model.Template
}

type UploadResponse struct {
	Data        *SalaryData
	UploadFile  string                   `json:"file"`
	DataPreview []map[string]interface{} `json:"preview"`
}

type SwaggerListResponse struct {
}

type ListCoefficient struct {
	Coefficient map[string]*model.Group
}

type Row struct {
	Data map[string]interface{}
}

type SalaryData struct {
	Base   float64
	Upload map[string]map[string]float64
	Rows   map[uint64]*Row
}

type ExportRequest struct {
	AccountID uint64 `form:"accountid"`
	Year      string `form:"year"`
	Month     string `form:"month"`
	Template  string `form:"template"`
}

type ExportResponse struct {
	File string `json:"file"`
}

var salaryDataMap *SalaryData
