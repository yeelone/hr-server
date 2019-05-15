package grouptransfer

import "hrgdrc/model"

type CreateRequest struct {
	UserID      uint64 `json:"user_id"`
	ProfileID   uint64 `json:"profile_id"`
	GroupID     uint64 `json:"group_id"`
	OldGroupID  uint64 `json:"old_group_id"`
	NewGroupID  uint64 `json:"new_group_id"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
}

type CreateResponse struct {
}

type ListRequest struct {
}

type ListResponse struct {
	Record []model.GroupTransfer
	File   string `json:"file"`
}

type UploadResponse struct {
}
