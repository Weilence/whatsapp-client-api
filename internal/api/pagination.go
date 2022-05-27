package api

type Pagination struct {
	Page int `json:"page,omitempty" form:"page,omitempty"`
	Size int `json:"size,omitempty" form:"size,omitempty"`
}

func (p Pagination) Limit() int {
	return p.Size
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}
