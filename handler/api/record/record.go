package record

import (
	"hr-server/model"
)

type CreateResponse struct {
}

type ListRequest struct {
	QueryDate string `form:"date"`
	Profile   uint64 `form:"profile"`
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit"`
}

type ListResponse struct {
	TotalCount    uint64                `json:"totalCount"`
	List          []model.Record        `json:"recordList"`
	OperationList []model.OperateRecord `json:"operateRecordList"`
}
