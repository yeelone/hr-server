package permission

type CreateRequest struct {
	Role     uint64                         `json:"role"`
	RoleName string                         `json:"role_name"`
	Fields   map[string]map[string]Resource `json:"fields"`
}

type CreateResponse struct {
	Fields map[string]map[string]Resource `json:"fields"`
}

type Resource struct {
	ID      string `json:"id"`
	Checked bool   `json:"checked"`
}
