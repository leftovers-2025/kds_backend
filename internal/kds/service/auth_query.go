package service

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

type AuthQueryService struct {
}

func NewAuthQueryService() *AuthQueryService {
	return &AuthQueryService{}
}

type AuthAccessTokenQueryInput struct {
	AccessToken string
}

type AuthIdQueryOutput struct {
	UserId uuid.UUID
}

// アクセストークンからIDを取得する
func (s *AuthQueryService) IdFromAccessToken(input AuthAccessTokenQueryInput) (*AuthIdQueryOutput, error) {
	accessToken, err := entity.AccessTokenFromToken(input.AccessToken)
	if err != nil {
		return nil, err
	}
	err = accessToken.IsExpired()
	if err != nil {
		return nil, err
	}
	return &AuthIdQueryOutput{
		UserId: accessToken.Sub(),
	}, nil
}
