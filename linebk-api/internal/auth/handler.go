package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/validator"
)

// Handler defines HTTP request handler methods for user authentication and system health checks.
type Handler interface {
	SignInWithPasscode(c *gin.Context)
	SignInWithPassword(c *gin.Context)
	SignOut(c *gin.Context)

	HealthCheck(c *gin.Context)
}

type handler struct {
	vdt *validator.Instance
	cfg *config.Config
	svc Service
}

// NewHandler creates a new instance of AuthHandler with the provided dependencies.
func NewHandler(
	vdt *validator.Instance,
	cfg *config.Config,
	svc Service,
) Handler {
	return &handler{vdt: vdt, cfg: cfg, svc: svc}
}

// SignInWithPasscode godoc
//
//	@Summary		Signin with id and passcode
//	@Description	Signin using id and passcode
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			req	body		SignInWithPasscodeRequest	true	"Signin with passscode payload"
//	@Success		200	{object}	ApiResponseWithData{data=TokenResponse}
//	@Failure		400	{object}	ApiResponse{errors=ErrResponse}
//	@Failure		401	{object}	ApiResponse
//	@Failure		500	{object}	ApiResponse
//	@Router			/v1/auth/signin/passcode [post]
func (h *handler) SignInWithPasscode(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req SignInWithPasscodeRequest

	if err := errs.GinBindAndValidate(c, &req, h.vdt); err != nil {
		errs.GinHandleError(c, err)
		return
	}

	token, err := h.svc.SignInWithPasscode(
		ctx,
		req.ID,
		req.Passcode,
		c.ClientIP(),
		c.Request.UserAgent(),
	)
	if err != nil {
		errs.GinHandleError(c, err)
		return
	}

	res := ApiResponseWithData[*TokenResponse]{
		Status:     api.StatusSuccess,
		StatusCode: http.StatusOK,
		Data: &TokenResponse{
			AccessToken: token.AccessToken,
		},
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		UserRefreshTokenCookieKey,
		token.RefreshToken,
		int(h.cfg.Auth.RefreshTTL.Seconds()),
		"/",
		h.cfg.Cookie.Domain,
		true,
		true,
	)

	c.JSON(http.StatusOK, res)
}

// SignInWithPassword godoc
//
//	@Summary		Signin with id and password
//	@Description	Signin using id and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			req	body		SignInWithPasswordRequest	true	"Signin with password payload"
//	@Success		200	{object}	ApiResponseWithData{data=TokenResponse}
//	@Failure		400	{object}	ApiResponse{errors=ErrResponse}
//	@Failure		401	{object}	ApiResponse
//	@Failure		500	{object}	ApiResponse
//	@Router			/v1/auth/signin/password [post]
func (h *handler) SignInWithPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var req SignInWithPasswordRequest

	if err := errs.GinBindAndValidate(c, &req, h.vdt); err != nil {
		errs.GinHandleError(c, err)
		return
	}

	token, err := h.svc.SignInWithPassword(
		ctx,
		req.ID,
		req.Password,
		c.ClientIP(),
		c.Request.UserAgent(),
	)
	if err != nil {
		errs.GinHandleError(c, err)
		return
	}

	res := ApiResponseWithData[*TokenResponse]{
		Status:     api.StatusSuccess,
		StatusCode: http.StatusOK,
		Data: &TokenResponse{
			AccessToken: token.AccessToken,
		},
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		UserRefreshTokenCookieKey,
		token.RefreshToken,
		int(h.cfg.Auth.RefreshTTL.Seconds()),
		"/",
		h.cfg.Cookie.Domain,
		true,
		true,
	)

	c.JSON(http.StatusOK, res)
}

// SignOut godoc
//
//	@Summary		Sign out user
//	@Description	Invalidate user session and clear authentication cookies
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ApiResponse
//	@Failure		400	{object}	ApiResponse
//	@Failure		401	{object}	ApiResponse	"Unauthorized"
//	@Failure		500	{object}	ApiResponse
//	@Router			/v1/auth/signout [post]
func (h *handler) SignOut(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	refreshToken, err := c.Cookie(UserRefreshTokenCookieKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			errs.GinHandleError(
				c,
				errs.NewUnauthorizedError(
					"Authentication session expired or missing.",
				),
			)
			return
		}

		errs.GinHandleError(c, err)
		return
	}

	if err := h.svc.SignOut(
		ctx,
		refreshToken,
	); err != nil {
		errs.GinHandleError(c, err)
		return
	}

	res := ApiResponse{
		Status:     api.StatusSuccess,
		StatusCode: http.StatusOK,
		Message:    "Sign Out Successfully.",
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(UserRefreshTokenCookieKey, "", -1, "/", h.cfg.Cookie.Domain, true, true)

	c.JSON(http.StatusOK, res)
}

// HealthCheck is a simple health check endpoint to verify that the service is running.
func (h *handler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
