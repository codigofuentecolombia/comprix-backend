package dto

import "comprix/app/domain/dao"

type AuthSocialiteResponse struct {
	Url string `json:"url"`
}

type AuthRequestCode struct {
	Code string `json:"code"`
}

type AuthLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type AuthRegister struct {
	Username    string `form:"username"     json:"username"     binding:"required,email"`
	Password    string `form:"password"     json:"password"     binding:"required,min=5"`
	LastName    string `form:"last_name"    json:"last_name"    binding:"required"`
	FirstName   string `form:"first_name"   json:"first_name"   binding:"required"`
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required,min=8"`
}

type AuthSendVerificationCode struct {
	Email string `form:"email" json:"email" binding:"required"`
}

type AuthVerificationCode struct {
	Email string `form:"email" json:"email" binding:"required"`
	Code  string `uri:"code"   json:"code"  binding:"required"`
}

type AuthResponse struct {
	User  dao.User `json:"user"`
	Token string   `json:"token"`
}

type SocialiteType string

const SocialiteGoogleType SocialiteType = "google"
const SocialiteFacebookType SocialiteType = "facebook"

type AuthFacebookUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthGoogleUser struct {
	ID    string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
