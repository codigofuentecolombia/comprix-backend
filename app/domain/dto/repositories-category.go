package dto

type CategoryRepositoryFindParams struct {
	ID          *string
	Name        *string
	OnlyParents *bool
	ParentID    *string

	Preloads *[]RepositoryGormParams
}
