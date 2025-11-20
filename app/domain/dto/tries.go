package dto

type HandleScrapperTotalTriesParams[ResponseType any] struct {
	Err             error
	Url             string
	Msg             string
	Callback        func(ScrapperParams) ResponseType
	CallbackArgs    ScrapperParams
	DefaultResponse ResponseType
}
