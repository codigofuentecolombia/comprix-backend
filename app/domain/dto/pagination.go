package dto

type Pagination[T any] struct {
	Items []T `json:"items"`

	Sort   string `json:"sort"  form:"sort"`
	Index  int    `json:"index" form:"page"`
	Limit  int    `json:"limit" form:"limit"`
	Search string `json:"-"     form:"search"`
	Offset int    `json:"-"`

	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func (p *Pagination[T]) Validate() {
	// Asegurarse de que el índice y el límite sean válidos
	if p.Index < 1 {
		p.Index = 1
	}
	// Asegurar que el limite sea valdio
	if p.Limit < 1 {
		p.Limit = 10
	}
	// Actualizar offset
	p.Offset = (p.Index - 1) * p.Limit
}
