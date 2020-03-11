package statistics

import "hr-server/model"

type CreateRequest struct {
	Year         string   `json:"year" form:"year"`
	Month        string   `json:"month" form:"month"`
	ProfileID    []uint64 `json:"profiles" form:"profiles"`
	DepartmentID []uint64 `json:"departments" form:"departments"`
}

type CreateResponse struct {
	File string `json:"file"`
}

type DetailRequest struct {
	Year      string      `json:"year" form:"year"`
	Account   uint64      `json:"account" form:"account"`
	Templates []DetailMap `json:"templates" form:"templates"`
}

type DetailMap struct {
	Template string   `json:"template" form:"template"`
	Fields   []string `json:"fields" form:"fields"`
}

type ProfileRequest struct {
	GetYear    bool   `json:"getYear"`
	GetMonth   bool   `json:"getMonth"`
	GetDay     bool  `json:"getDay"`
	Amount      int   `json:"amount"`
}

type ProfileResponse struct {
	Data []model.ProfileIncrease `json:"data"`
}

