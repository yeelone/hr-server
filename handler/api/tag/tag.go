package tag

import (
	"hr-server/model"
)

type CreateRequest struct {
	IDS         []uint64 `json:"tag_id_list"`
	ID          uint64   `json:"id"`
	Parent      uint64   `json:"parent"`
	Name        string   `json:"name"`
	Coefficient float64  `json:"coefficient"`
	Users       []uint64 `json:"user_id_list"`
	Profiles    []uint64 `json:"profile_id_list"`
	Remark      string   `json:"remark"`
	GroupIds    []uint64 `json:"commensalism_group_ids"`
}

type CreateResponse struct {
	Tag   *model.Tag `json:"tag"`
	File  string     `json:"file"`
	Error string     `json:"error"`
}

type ListResponse struct {
	TotalCount  uint64          `json:"totalCount"`
	TagList     []*model.Tag    `json:"tagList"`
	UserList    []*model.User   `json:"userList"`
	ProfileList []model.Profile `json:"profileList"`
}

type ListRequest struct {
	Offset     int    `json:"offset" form:"offset"`
	Limit      int    `json:"limit" form:"limit"`
	WhereKey   string `json:"where_key" form:"where_key"`
	WhereValue string `json:"where_value" form:"where_value"`
}

type SwaggerListResponse struct {
}
