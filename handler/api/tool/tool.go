package tool

type CreateRequest struct {
}

type CreateResponse struct {
	List  []map[string]interface{}
	File  string
	Files []string
}

type ListRequest struct {
}

type ListResponse struct {
	Data []map[string]interface{}
}
