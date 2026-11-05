package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/auth"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/middleware"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type AuthHandler struct {
	service      *service.AuthService
	cookieSecure bool
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(
	service *service.AuthService,
	cookieSecure bool,
) *AuthHandler {
	return &AuthHandler{
		service:      service,
		cookieSecure: cookieSecure,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request loginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, token, err := h.service.Login(
		c.Request.Context(),
		request.Email,
		request.Password,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to login",
		})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     auth.AccessTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(auth.AccessTokenDuration.Seconds()),
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     auth.AccessTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDContextKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	user, err := h.service.GetCurrentUser(
		c.Request.Context(),
		userID.(string),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
