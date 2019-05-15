package statistics

type CreateRequest struct {
	Year string `json:"year" form:"year"`
	Month string `json:"month" form:"month"`
	ProfileID []uint64 `json:"profiles" form:"profiles"`
	DepartmentID []uint64 `json:"departments" form:"departments"`
}

type CreateResponse struct {
	File string `json:"file"`
}

type DetailRequest struct {
	Year   string `json:"year" form:"year"`
	Account uint64 `json:"account" form:"account"`
	Templates []DetailMap `json:"templates" form:"templates"`
}

type DetailMap struct {
	Template string  `json:"template" form:"template"`
	Fields   []string `json:"fields" form:"fields"`
}