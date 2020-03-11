package sd

type Health struct {
	UsedMB int `json:"usedMB"`
	UsedGB int `json:"usedGB"`
	TotalMB int `json:"totalMB"`
	TotalGB int `json:"totalGB"`
	UsedPercent int `json:"usedPercent"`
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15  float64 `json:"load15"`
	CpuUsage    float64 `json:"cpuUsage"`
	CpuTotal    float64 `json:"cpuTotal"`
}

type CreateResponse struct {
	Status  int `json:"status"`
	Disk Health `json:"disk"`
	RAM Health `json:"ram"`
	CPU Health `json:"cpu"`
}


