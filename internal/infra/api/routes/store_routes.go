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
	updateAddress     *handler.UpdateAddressHandler
	listNearby        *handler.ListNearbyStoresHandler
	getWithActiveBags *handler.GetStoreWithActiveBagsHandler
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
	updateAddress *handler.UpdateAddressHandler,
	listNearby *handler.ListNearbyStoresHandler,
	getWithActiveBags *handler.GetStoreWithActiveBagsHandler,
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
		updateAddress:     updateAddress,
		listNearby:        listNearby,
		getWithActiveBags: getWithActiveBags,
	}
}

func (r *StoreRoutes) Register(router *gin.Engine) {
	stores := router.Group("/stores")
	stores.Use(handler.ErrorHandler)
	stores.Use(middleware.MaxRequestBodySize(5 << 20))

	// Rotas públicas
	public := stores.Group("")
	public.POST("/signup", r.signup.Handle)
	public.POST("/login", r.login.Handle)
	public.POST("/email-verification", r.emailVerification.Handle)
	public.PATCH("/verify-email", r.verifyEmail.Handle)
	public.POST("/password-reset", r.passwordReset.Handle)
	public.PATCH("/reset-password", r.resetPassword.Handle)

	// Rotas privadas
	private := stores.Group("")
	private.Use(r.auth.AuthenticateStore)
	private.GET("/me", r.getMe.Handle)
	private.PATCH("/me", r.updateMe.Handle)
	private.PUT("/me/avatar", r.updateAvatar.Handle)
	private.DELETE("/me/avatar", r.removeAvatar.Handle)
	private.PATCH("/me/email", r.updateEmail.Handle)
	private.PATCH("/me/password", r.updatePassword.Handle)
	private.PATCH("/me/address", r.updateAddress.Handle)
	private.POST("/logout", r.logout.Handle)
	private.POST("/logout-all", r.logoutAll.Handle)

	// Rotas do app mobile
	reservation := stores.Group("")
	reservation.Use(r.auth.AuthenticateCustomer)
	reservation.GET("", r.listNearby.Handle)
	reservation.GET("/:id", r.getWithActiveBags.Handle)
}
