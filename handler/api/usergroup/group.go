package usergroup

import (
	"hrgdrc/model"
)

type CreateRequest struct {
	ID     uint64   `json:"id"`
	Name   string   `json:"name"`
	Parent uint64   `json:"parent"`
	Users  []uint64 `json:"user_id_list"`
}

type CreateResponse struct {
	Group *model.UserGroup `json:"group"`
}

type ListResponse struct {
	TotalCount int                `json:"totalCount"`
	GroupList  []*model.UserGroup `json:"groupList"`
	UserList   []*model.User      `json:"userList"`
}

type ListRequest struct {
	Offset     int    `json:"offset" form:"offset"`
	Limit      int    `json:"limit" form:"limit"`
	WhereKey   string `json:"where_key" form:"where_key"`
	WhereValue string `json:"where_value" form:"where_value"`
}

type SwaggerListResponse struct {
	TotalCount uint64            `json:"totalCount"`
	GroupList  []model.UserGroup `json:"groupList"`
}
