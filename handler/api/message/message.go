package message

type CreateRequest struct {
	SendId uint64 `json:"sendId"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	MType  string `json:"messageType"`
	Group  uint64 `json:"groupId"`
	Role   uint64 `json:"roleId"`
	RecId  uint64 `json:"recId"`
}

type InboxRequest struct {
	RecId uint64 `json:"recId"`
}

type UpdateStatusRequest struct {
	MsgIds []uint64 `json:"msgIds"`
	Status int      `json:"status"`
}

type InboxCountResponse struct {
	Private int `json:"private"`
	Public  int `json:"public"`
	Global  int `json:"global"`
}

type CreateResponse struct {
	Id     uint64 `json:"id"`
	SenderId uint64  `json:"senderId"`
	SenderName string  `json:"senderName"`
	Title  string      `json:"title"`
	Text  string      `json:"text"`
	MType  string      `json:"messageType"`
	GroupName string `json:"groupName"`
	RoleName  string  `json:"roleName"`
	Date   string      `json:"date"`
}

type ListResponse struct {
	List []CreateResponse `json:"list"`
	Total int `json:"total"`
}

type ListRequest struct {
	Offset     int    `json:"offset" form:"offset"`
	Limit      int    `json:"limit" form:"limit"`
}