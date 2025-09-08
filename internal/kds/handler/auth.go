package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

const (
	AUTHORIZATION_HEADER_PREFIX        = "Bearer "
	AUTHORIZATION_USER_ID_CONTEXT_NAME = "userId"
)

var (
	ErrInvalidUserIdType = errors.New("invalid userId type")
)

type AuthHandler struct {
	authCmdService   *service.AuthCommandService
	authQueryService *service.AuthQueryService
}

func NewAuthHandler(authCmdService *service.AuthCommandService, authQueryService *service.AuthQueryService) *AuthHandler {
	if authCmdService == nil {
		panic("nil AuthCommandService")
	}
	if authQueryService == nil {
		panic("nil AuthQueryService")
	}
	return &AuthHandler{
		authCmdService:   authCmdService,
		authQueryService: authQueryService,
	}
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenRefreshResponse struct {
	AccessToken string    `json:"accessToken"`
	ExpiresIn   time.Time `json:"expiresIn"`
}

// @Summary		Refresh access token
// @Description	Refresh the access token using the provided refresh token
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			refreshToken	body		TokenRefreshRequest	true	"Refresh Token"
// @Success		200				{object}	TokenRefreshResponse
// @Failure		400				{object}	ErrorResponse
// @Failure		401				{object}	ErrorResponse
// @Failure		500				{object}	ErrorResponse
// @Router			/refreshToken [post]
func (h *AuthHandler) RefreshToken(ctx echo.Context) error {
	request := TokenRefreshRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	output, err := h.authCmdService.RefreshToken(service.TokenRefreshCommandInput{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, TokenRefreshResponse{
		AccessToken: output.AccessToken,
		ExpiresIn:   output.ExpiresIn,
	})
}

// JWT認証を行うミドルウェア
func (h *AuthHandler) JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// 認証ヘッダー取得
		authHeader := ctx.Request().Header.Get("Authorization")
		// prefix確認
		if !strings.HasPrefix(authHeader, AUTHORIZATION_HEADER_PREFIX) {
			return ctx.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization header is required",
			})
		}
		// トークンのみ取得
		token := strings.TrimPrefix(authHeader, AUTHORIZATION_HEADER_PREFIX)
		// ID取得
		idOutput, err := h.authQueryService.IdFromAccessToken(service.AuthAccessTokenQueryInput{
			AccessToken: token,
		})
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
		}
		// idを設定
		ctx.Set(AUTHORIZATION_USER_ID_CONTEXT_NAME, idOutput.UserId)

		// エンドポイント
		if err := next(ctx); err != nil {
			ctx.Error(err)
		}
		return nil
	}
}

// コンテキストからユーザーIDを取得
func getUserIdFromCtx(ctx echo.Context) (uuid.UUID, error) {
	userId := ctx.Get(AUTHORIZATION_USER_ID_CONTEXT_NAME)
	if sub, ok := userId.(uuid.UUID); ok {
		return sub, nil
	}
	return uuid.Nil, ErrInvalidUserIdType
}
