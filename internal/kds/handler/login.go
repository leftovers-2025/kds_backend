package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

type GoogleHandler struct {
	googleCmdService *service.GoogleCommandService
}

func NewGoogleHandler(googleCmdService *service.GoogleCommandService) *GoogleHandler {
	if googleCmdService == nil {
		panic("nil GoogleCommandService")
	}
	return &GoogleHandler{
		googleCmdService,
	}
}

type GoogleLoginResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Redirect godoc
//
//	@Summary		Google OAuth redirect
//	@Description	Handle Google OAuth redirect
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"Authorization code"
//	@Success		200		{object}	GoogleLoginResponse
//	@Router			/oauth/google/redirect [get]
//
// Google OAuth認証時リダイレクト先
func (h *GoogleHandler) Redirect(ctx echo.Context) error {
	if !ctx.QueryParams().Has("code") {
		return common.NewValidationError(errors.New("code is not set"))
	}
	output, err := h.googleCmdService.OauthLogin(service.GoogleOauthLoginCommandInput{
		Code: ctx.QueryParam("code"),
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, GoogleLoginResponse{
		Id:           output.Id.String(),
		Name:         output.Name,
		Email:        output.Email,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}
