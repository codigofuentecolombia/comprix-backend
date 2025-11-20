package routes

import (
	auth_controller "comprix/app/controller/auth"
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/auth")
	controller := auth_controller.InitAuthController(config)
	//
	// Login
	router.POST("login", controller.Login)
	router.POST("register", controller.Register)
	// Facebook
	router.GET("facebook", controller.Socialite(dto.SocialiteFacebookType))
	router.POST("facebook", controller.SocialiteCallback(dto.SocialiteFacebookType))
	// Google
	router.GET("google", controller.Socialite(dto.SocialiteGoogleType))
	router.POST("google", controller.SocialiteCallback(dto.SocialiteGoogleType))
	//
	router.GET("", server.ValidateJwt(config), controller.GetDetail)
	router.POST("update", server.ValidateJwt(config), controller.Update)
	//
	router.POST("verification-code", controller.SendVerificationCode)
	router.POST("verification-code/:code", controller.VerificationCode)

}
