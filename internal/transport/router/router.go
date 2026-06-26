package router

import (
	"Test_App/internal/transport/handler"
	"Test_App/internal/transport/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handler.UserAuthHandler,
	subjectHandler *handler.SubjectHandler,
	sectionHandler *handler.SectionHandler,
	testHandler *handler.TestHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	// ===================== публичные маршруты =====================
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		auth.POST("/resend-code", authHandler.ResendCode)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	// ===================== защищённые маршруты =====================
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// предметы
		protected.GET("/subjects", subjectHandler.GetAllSubjects)
		protected.GET("/subjects/:id", subjectHandler.GetSubjectByID)

		// разделы
		protected.GET("/subjects/:id/sections", sectionHandler.GetSectionsBySubjectID)

		// тесты
		protected.GET("/sections/:id/tests", testHandler.GetTestsBySection)
		protected.GET("/tests/:id", testHandler.GetTest)
		protected.POST("/tests/:id/answer", testHandler.SubmitAnswer)
		protected.POST("/tests/:id/restart", testHandler.Restart)

		// попытки
		protected.POST("/attempts/:id/finish", testHandler.Finish)
		protected.GET("/attempts/:id/result", testHandler.GetResult)
	}

	return r
}

//package router
//
//import (
//	"Test_App/internal/transport/handler"
//
//	"github.com/gin-gonic/gin"
//)
//
//func SetupRouter(authHandler *handler.UserAuthHandler) *gin.Engine {
//	r := gin.Default()
//
//	auth := r.Group("/auth")
//	{
//		auth.POST("/register", authHandler.Register)
//		auth.POST("/verify-email", authHandler.VerifyEmail)
//		auth.POST("/resend-code", authHandler.ResendCode)
//		auth.POST("/login", authHandler.Login)
//		auth.POST("/forgot-password", authHandler.ForgotPassword)
//		auth.POST("/reset-password", authHandler.ResetPassword)
//	}
//
//	return r
//}
