package router

import (
	"Test_App/internal/transport/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.UserAuthHandler) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		auth.POST("/resend-code", authHandler.ResendCode)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	return r
}
