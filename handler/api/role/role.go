package role

import (
	"hrgdrc/model"
)

type CreateRequest struct {
	ID    uint64   `json:"id"`
	Name  string   `json:"name"`
	Users []uint64 `json:"user_id_list"`
	Remark string `json:"remark"`
}

type CreateResponse struct {
	Role *model.Role `json:"role"`
}

type ListResponse struct {
	TotalCount int           `json:"totalCount"`
	RoleList   []*model.Role `json:"roleList"`
	UserList   []*model.User `json:"userList"`
}

type ListRequest struct {
	Offset     int    `json:"offset" form:"offset"`
	Limit      int    `json:"limit" form:"limit"`
	WhereKey   string `json:"where_key" form:"where_key"`
	WhereValue string `json:"where_value" form:"where_value"`
}

type SwaggerListResponse struct {
	TotalCount uint64       `json:"totalCount"`
	RoleList   []model.Role `json:"roleList"`
}
