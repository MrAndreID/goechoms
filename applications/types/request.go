package types

type PaginatorRequest struct {
	Page                  int    `query:"page" json:"page"`
	Limit                 int    `query:"limit" json:"limit"`
	OrderBy               string `query:"orderBy" json:"orderBy"`
	SortBy                string `query:"sortBy" json:"sortBy"`
	Search                string `query:"search" json:"search"`
	DisableCalculateTotal string `query:"disableCalculateTotal" json:"disableCalculateTotal"`
}

type GetUserRequest struct {
	PaginatorRequest
	ID string `query:"id" json:"id"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type EditUserRequest struct {
	ID   string `param:"id" json:"id"`
	Name string `json:"name"`
}

type DeleteUserRequest struct {
	ID string `param:"id" json:"id"`
}
