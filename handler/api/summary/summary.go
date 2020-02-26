package summary

type CreateResponse struct {
	ProfileCount int `json:"profile_count"`
	UserCount    int `json:"user_count"`
	GroupCount   int `json:"group_count"`
	AuditCount   int `json:"audit_count"`
}
