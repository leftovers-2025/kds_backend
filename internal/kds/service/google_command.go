package service

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type GoogleCommandService struct {
	userCmdService   *UserCommandService
	authCmdService   *AuthCommandService
	userRepository   port.UserRepository
	googleRepository port.GoogleRepository
}

func NewGoogleCommandService(
	userCmdService *UserCommandService,
	authCmdService *AuthCommandService,
	userRepository port.UserRepository,
	googleRepository port.GoogleRepository,
) *GoogleCommandService {
	if userCmdService == nil {
		panic("nil UserCommandService")
	}
	if authCmdService == nil {
		panic("nil AuthCommandService")
	}
	if userRepository == nil {
		panic("nil UserRepository")
	}
	if googleRepository == nil {
		panic("nil GoogleRepository")
	}
	return &GoogleCommandService{
		userCmdService:   userCmdService,
		authCmdService:   authCmdService,
		userRepository:   userRepository,
		googleRepository: googleRepository,
	}
}

type GoogleOauthLoginCommandInput struct {
	Code string
}

type GoogleOauthLoginCommandOutput struct {
	Id           uuid.UUID
	Name         string
	Email        string
	AccessToken  string
	RefreshToken string
}

// Google OAuthを利用してログインする
func (s *GoogleCommandService) OauthLogin(input GoogleOauthLoginCommandInput) (*GoogleOauthLoginCommandOutput, error) {
	// グーグル認証
	googleUser, err := s.googleRepository.CodeAuthorization(input.Code)
	if err != nil {
		return nil, err
	}
	// GoogleのIdからユーザー取得
	user, err := s.userRepository.FindByGoogleId(googleUser.Id())
	// 既存ユーザー確認
	if err != nil {
		// 存在しない場合ユーザー作成
		userOutput, err := s.userCmdService.CreateUser(UserCreateCommandInput{
			Name:     googleUser.Name(),
			Email:    googleUser.Email().String(),
			GoogleId: googleUser.Id(),
		})
		if err != nil {
			return nil, err
		}
		// トークン生成
		tokenOutput, err := s.authCmdService.GenerateToken(AuthTokenCommandInput{
			UserId: userOutput.Id,
		})
		if err != nil {
			return nil, err
		}
		return &GoogleOauthLoginCommandOutput{
			Id:           userOutput.Id,
			Name:         userOutput.Name,
			Email:        userOutput.Email,
			AccessToken:  tokenOutput.AccessToken,
			RefreshToken: tokenOutput.RefreshToken,
		}, nil
	} else {
		// トークン生成
		tokenOutput, err := s.authCmdService.GenerateToken(AuthTokenCommandInput{
			UserId: user.Id(),
		})
		if err != nil {
			return nil, err
		}
		return &GoogleOauthLoginCommandOutput{
			Id:           user.Id(),
			Name:         user.Name(),
			Email:        user.Email().String(),
			AccessToken:  tokenOutput.AccessToken,
			RefreshToken: tokenOutput.RefreshToken,
		}, nil
	}
}
