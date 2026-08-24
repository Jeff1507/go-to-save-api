package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/store/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type StoreRoutes struct {
	auth              *middleware.Authentication
	getMe             *handler.GetProfileHandler
	signup            *handler.SignupHandler
	login             *handler.LoginHandler
	emailVerification *handler.RequestEmailVerificationHandler
	verifyEmail       *handler.VerifyEmailHandler
	passwordReset     *handler.RequestPasswordResetHandler
	resetPassword     *handler.ResetPasswordHandler
	logout            *handler.LogoutHandler
	logoutAll         *handler.LogoutAllHandler
	updateMe          *handler.UpdateProfileHandler
	updateAvatar      *handler.UpdateAvatarHandler
	removeAvatar      *handler.RemoveAvatarHandler
	updateEmail       *handler.UpdateEmailHandler
	updatePassword    *handler.UpdatePasswordHandler
}

func NewStoreRoutes(
	auth *middleware.Authentication,
	getMe *handler.GetProfileHandler,
	signup *handler.SignupHandler,
	login *handler.LoginHandler,
	emailVerification *handler.RequestEmailVerificationHandler,
	verifyEmail *handler.VerifyEmailHandler,
	passwordReset *handler.RequestPasswordResetHandler,
	resetPassword *handler.ResetPasswordHandler,
	logout *handler.LogoutHandler,
	logoutAll *handler.LogoutAllHandler,
	updateMe *handler.UpdateProfileHandler,
	updateAvatar *handler.UpdateAvatarHandler,
	removeAvatar *handler.RemoveAvatarHandler,
	updateEmail *handler.UpdateEmailHandler,
	updatePassword *handler.UpdatePasswordHandler,
) *StoreRoutes {
	return &StoreRoutes{
		auth:              auth,
		getMe:             getMe,
		signup:            signup,
		login:             login,
		emailVerification: emailVerification,
		verifyEmail:       verifyEmail,
		passwordReset:     passwordReset,
		resetPassword:     resetPassword,
		logout:            logout,
		logoutAll:         logoutAll,
		updateMe:          updateMe,
		updateAvatar:      updateAvatar,
		removeAvatar:      removeAvatar,
		updateEmail:       updateEmail,
		updatePassword:    updatePassword,
	}
}

func (r *StoreRoutes) Register(router *gin.Engine) {
	stores := router.Group("/stores")
	stores.Use(handler.ErrorHandler)
	stores.Use(middleware.MaxRequestBodySize(5 << 20))

	{
		// Rotas públicas
		stores.POST("/signup", r.signup.Handle)
		stores.POST("/login", r.login.Handle)
		stores.POST("/email-verification", r.emailVerification.Handle)
		stores.PATCH("/verify-email", r.verifyEmail.Handle)
		stores.POST("/password-reset", r.passwordReset.Handle)
		stores.PATCH("/reset-password", r.resetPassword.Handle)
	}

	{
		// Rotas privadas
		stores.Use(r.auth.AuthenticateStore)

		stores.GET("/me", r.getMe.Handle)
		stores.PATCH("/me", r.updateMe.Handle)
		stores.PUT("/me/avatar", r.updateAvatar.Handle)
		stores.DELETE("/me/avatar", r.removeAvatar.Handle)
		stores.PATCH("/me/email", r.updateEmail.Handle)
		stores.PATCH("/me/password", r.updatePassword.Handle)
		stores.POST("/logout", r.logout.Handle)
		stores.POST("/logout-all", r.logoutAll.Handle)
	}
}
