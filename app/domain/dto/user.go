package dto

import "mime/multipart"

type UpdateProfile struct {
	File        *multipart.FileHeader `form:"file"`
	FirstName   string                `form:"first_name"     binding:"required,max=250"`
	LastName    string                `form:"last_name"      binding:"required,max=250"`
	PhoneNumber string                `form:"phone_number"   binding:"required,max=20"`
}
