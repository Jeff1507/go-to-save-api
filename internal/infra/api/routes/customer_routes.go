package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/customer/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type CustomerRoutes struct {
	auth              *middleware.Authentication
	get               *handler.GetProfileHandler
	update            *handler.UpdateProfileHandler
	signup            *handler.SignupHandler
	login             *handler.LoginHandler
	emailVerification *handler.RequestEmailVerificationHandler
	verifyEmail       *handler.VerifyEmailHandler
	passwordReset     *handler.RequestPasswordResetHandler
	resetPassword     *handler.ResetPasswordHandler
	logout            *handler.LogoutHandler
	logoutAll         *handler.LogoutAllHandler
	updateEmail       *handler.UpdateEmailHandler
	updatePassword    *handler.UpdatePasswordHandler
}

func NewCustomerRoutes(
	auth *middleware.Authentication,
	get *handler.GetProfileHandler,
	update *handler.UpdateProfileHandler,
	signup *handler.SignupHandler,
	login *handler.LoginHandler,
	emailVerification *handler.RequestEmailVerificationHandler,
	verifyEmail *handler.VerifyEmailHandler,
	passwordReset *handler.RequestPasswordResetHandler,
	resetPassword *handler.ResetPasswordHandler,
	logout *handler.LogoutHandler,
	logoutAll *handler.LogoutAllHandler,
	updateEmail *handler.UpdateEmailHandler,
	updatePassword *handler.UpdatePasswordHandler,
) *CustomerRoutes {
	return &CustomerRoutes{
		auth:              auth,
		get:               get,
		update:            update,
		signup:            signup,
		login:             login,
		emailVerification: emailVerification,
		verifyEmail:       verifyEmail,
		passwordReset:     passwordReset,
		resetPassword:     resetPassword,
		logout:            logout,
		logoutAll:         logoutAll,
		updateEmail:       updateEmail,
		updatePassword:    updatePassword,
	}
}

func (r *CustomerRoutes) Register(router *gin.Engine) {
	customers := router.Group("/customers")
	customers.Use(handler.ErrorHandler)

	{
		// Rotas públicas
		customers.POST("/signup", r.signup.Handle)
		customers.POST("/login", r.login.Handle)
		customers.POST("/email-verification", r.emailVerification.Handle)
		customers.PATCH("/verify-email", r.verifyEmail.Handle)
		customers.POST("/password-reset", r.passwordReset.Handle)
		customers.PATCH("/reset-password", r.resetPassword.Handle)
	}

	{
		// Rotas privadas
		customers.Use(r.auth.VerifyCustomer)

		customers.GET("/me", r.get.Handle)
		customers.PUT("/me", r.update.Handle)
		customers.PATCH("/me/email", r.updateEmail.Handle)
		customers.PATCH("/me/password", r.updatePassword.Handle)
		customers.POST("/logout", r.logout.Handle)
		customers.POST("/logout-all", r.logoutAll.Handle)
	}
}
