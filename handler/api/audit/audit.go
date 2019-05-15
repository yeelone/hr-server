package audit

import (
	"hrgdrc/model"
)

type CreateRequest struct {
	ID     uint64
	Object string `json:"object"` //审核对象
	State  int    `json:"state"`
	Reply  string `json:"reply"`
}

type CreateResponse struct {
}

type ListRequest struct {
	State  int `form:"state"`
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type ListResponse struct {
	TotalCount uint64        `json:"totalCount"`
	List       []model.Audit `json:"auditList"`
}
