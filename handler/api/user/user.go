package user

import (
	"hr-server/model"
)

type CreateRequest struct {
	ID           uint64   `json:"id"`
	Email        string   `json:"email"`
	Username     string   `json:"username"`
	IDCard       string   `json:"id_card"`
	Nickname     string   `json:"nickname"`
	Password     string   `json:"password"`
	OldPassword  string   `json:"oldPassword"`
	IsSuper      bool     `json:"is_super"`
	Picture      string   `json:"picture"`
	State        int      `json:"state"`
	Group        uint64   `json:"group"`
	Users        []uint64 `json:"users"`
	CaptchaId    string   `json:"captchaId"`
	CaptchaValue string   `json:"captchaValue"`
	// ProfileID uint64 `json:"profile_id"`
}

type CreateResponse struct {
	Username string `json:"username"`
}

type ListRequest struct {
	Username string `form:"username"`
	Offset   int    `form:"offset"`
	Limit    int    `form:"limit"`
}

type ListResponse struct {
	TotalCount uint64        `json:"totalCount"`
	UserList   []*model.User `json:"userList"`
}

type SwaggerListResponse struct {
	TotalCount uint64           `json:"totalCount"`
	UserList   []model.UserInfo `json:"userList"`
}
