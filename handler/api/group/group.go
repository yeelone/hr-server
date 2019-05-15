package group

import (
	"hrgdrc/model"
)

type CreateRequest struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Parent      uint64   `json:"parent"`
	Code        int      `json:"code"`
	Coefficient float64  `json:"coefficient"`
	Profiles    []uint64 `json:"profile_id_list"`
	Remark string `json:"remark"`
}

type CreateResponse struct {
	Group *model.Group `json:"group"`
	Rules []string `json:"rules"`
	File string `json:"file"`
	Error string `json:"error"`
}

type ListResponse struct {
	TotalCount  int             `json:"totalCount"`
	GroupList   []*model.Group  `json:"groupList"`
	UserList    []*model.User   `json:"userList"`
	ProfileList []model.Profile `json:"profileList"`
}

type ListRequest struct {
	Offset     int    `json:"offset" form:"offset"`
	Limit      int    `json:"limit" form:"limit"`
	WhereKey   string `json:"where_key" form:"where_key"`
	WhereValue string `json:"where_value" form:"where_value"`
	Freezed    bool   `json:"freezed" form:"freezed"`
}

type SwaggerListResponse struct {
	TotalCount uint64        `json:"totalCount"`
	GroupList  []model.Group `json:"groupList"`
}

type RelateTagsRequest struct {
	Group uint64 `json:"group"`
	Tags    []uint64 `json:"tags"`
	Rules    []string `json:"rules"`
}
