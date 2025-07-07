package handler

import (
	"errors"
	"net/http"
	"strings"

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
	authQueryService *service.AuthQueryService
}

func NewAuthHandler(authQueryService *service.AuthQueryService) *AuthHandler {
	if authQueryService == nil {
		panic("nil AuthQueryService")
	}
	return &AuthHandler{
		authQueryService: authQueryService,
	}
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
			return err
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
