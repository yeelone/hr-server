package templateaccount

import (
	"hr-server/model"
	"hr-server/pkg/template"
)

type CreateRequest struct {
	ID        uint64   `json:"id"`
	Name      string   `json:"name"`
	Groups    []uint64 `json:"groups" ` //指定群组的人进行计算
	Order     []int64  `json:"order"`
	Templates []uint64 `json:"templates"`
	Remark    string   `json:"remark"`
}

type CreateResponse struct {
	TemplateAccount *model.TemplateAccount `json:"template_account"`
}

type ListRequest struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
}

type ListResponse struct {
	TotalCount int                      `json:"totalCount"`
	List       []*model.TemplateAccount `json:"List"`
	Templates  []TemplateResponse       `json:"templates"`
}

type TemplateResponse struct {
	model.Template
	Fields    []template.Field `json:"fields"`
	Templates []model.Template `json:"templates"`
}

type ListTemplateResponse struct {
	List []model.TemplateAccount
}

type TemplateFieldsResponse struct {
	Fields map[string]map[string][]string `json:"fields"`
}
