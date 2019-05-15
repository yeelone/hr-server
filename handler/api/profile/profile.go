package profile

import (
	"hrgdrc/model"
)

type CreateRequest struct {
	Profile model.Profile `json:"profile"`
	Remark string `json:"remark"`
}

type RelateTagsRequest struct {
	Profile uint64 `json:"profile"`
	Tags    []uint64 `json:"tags"`
}

type DeleteRequest struct{
	Profiles []uint64 `json:"profiles"`
	Remark string `json:"remark"`
}

type FreezeRequest struct {
	Profiles []uint64 `json:"profiles"`
	Remark string `json:"remark"`
}

type CreateResponse struct {
	Name string `json:"name"`
	Profile model.Profile `json:"profile"`
	File string `json:"file"`
	Error string `json:"error"`
}

type ListTagsResponse struct {
	Profile model.Profile `json:"profile"`
	Tags    []TagResponse `json:"tags"`
}

type TagResponse struct {
	Tag  model.Tag `json:"tag"`
	Children []model.Tag `json:"children"`
}

type ListRequest struct {
	Name   string `form:"name"`
	IDCard string `form:"id_card"`
	Key    string `form:"key"`
	Value  string `form:"value"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
	Freezed bool `form:"freezed"`
}

type ListResponse struct {
	TotalCount  uint64           `json:"totalCount"`
	ProfileList []*model.Profile `json:"profileList"`
}

type TransferResponse struct {
	Transfer  []model.GroupTransfer `json:"transfer"`
}