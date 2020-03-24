package tool

type CreateRequest struct {
}

type CustomFile struct {
	Size int64 `json:"size"`
	Path string `json:"path"`
	Date string `json:"date"`
}

type CreateResponse struct {
	List  []map[string]interface{}
	File  string
	Files []CustomFile
}

type ListRequest struct {
}

type ListResponse struct {
	Data []map[string]interface{}
}
