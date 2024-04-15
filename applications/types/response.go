package types

type MainResponse struct {
	Code        string      `json:"code"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

type PaginatorResponse struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	NextPage bool        `json:"nextPage"`
}

type HTTPResponse struct {
	Headers    map[string][]string
	Body       interface{}
	StatusCode int
	Error      error
}
