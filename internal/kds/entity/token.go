package entity

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenRequire     = errors.New("access token is required")
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenIdRequired  = errors.New("token id is required")
	ErrTokenSubRequired = errors.New("token sub is required")
	ErrTokenExpRequired = errors.New("token exp is required")
)

type AccessToken struct {
	sub      uuid.UUID
	token    string
	expireIn time.Time
}

func NewAccessToken(sub uuid.UUID, token string, expireIn time.Time) (*AccessToken, error) {
	if token == "" {
		return nil, ErrTokenRequire
	}
	if err := isExpired(expireIn); err != nil {
		return nil, err
	}
	return &AccessToken{
		sub:      sub,
		token:    token,
		expireIn: expireIn,
	}, nil
}

// アクセストークンを生成
func GenerateAccessToken(sub uuid.UUID, expireIn time.Time) (*AccessToken, error) {
	if err := isExpired(expireIn); err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{
		"sub": sub.String(),
		"exp": expireIn.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(getJwtSecretKey())
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		sub:      sub,
		token:    signedStr,
		expireIn: expireIn,
	}, nil
}

// トークンをデコードしアクセストークンを取得
func AccessTokenFromToken(accessToken string) (*AccessToken, error) {
	claims := jwt.MapClaims{}
	// トークンデコード
	_, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (any, error) {
		return getJwtSecretKey(), nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	// sub取得
	sub, err := claims.GetSubject()
	if err != nil {
		return nil, ErrTokenSubRequired
	}
	// 期限取得
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, ErrTokenExpRequired
	}
	// subをUUIDに変換
	userId, err := uuid.Parse(sub)
	if err != nil {
		return nil, err
	}
	return &AccessToken{
		sub:      userId,
		token:    accessToken,
		expireIn: exp.Time,
	}, nil
}

func (t *AccessToken) IsExpired() error {
	err := isExpired(t.expireIn)
	if err != nil {
		return err
	}
	return nil
}

func (t *AccessToken) Sub() uuid.UUID {
	return t.sub
}

func (t *AccessToken) Token() string {
	return t.token
}

type RefreshToken struct {
	id       uuid.UUID
	userId   uuid.UUID
	token    string
	expireIn time.Time
}

func NewRefreshToken(id, userId uuid.UUID, token string, expireIn time.Time) (*RefreshToken, error) {
	if token == "" {
		return nil, ErrTokenRequire
	}
	if err := isExpired(expireIn); err != nil {
		return nil, err
	}
	return &RefreshToken{
		id:       id,
		userId:   userId,
		token:    token,
		expireIn: expireIn,
	}, nil
}

// リフレッシュトークンを生成
func GenerateRefreshToken(sub uuid.UUID, expireIn time.Time) (*RefreshToken, error) {
	tokenId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{
		"sub": sub.String(),
		"tid": tokenId.String(),
		"exp": expireIn.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(getJwtSecretKey())
	if err != nil {
		return nil, err
	}
	return &RefreshToken{
		id:       tokenId,
		userId:   sub,
		token:    signedStr,
		expireIn: expireIn,
	}, nil
}

// トークンをデコードしリフレッシュトークンを取得
func RefreshTokenFromToken(refreshToken string) (*RefreshToken, error) {
	claims := jwt.MapClaims{}
	// トークンデコード
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		return getJwtSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}
	// tid取得
	tid, ok := claims["tid"].(string)
	if !ok {
		return nil, ErrTokenIdRequired
	}
	// sub取得
	sub, err := claims.GetSubject()
	if err != nil {
		return nil, ErrTokenSubRequired
	}
	// 期限取得
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, ErrTokenExpRequired
	}
	// tidをUUIDに変換
	tokenId, err := uuid.Parse(tid)
	if err != nil {
		return nil, err
	}
	// subをUUIDに変換
	userId, err := uuid.Parse(sub)
	if err != nil {
		return nil, err
	}
	return &RefreshToken{
		id:       tokenId,
		userId:   userId,
		token:    refreshToken,
		expireIn: exp.Time,
	}, nil
}

func (t *RefreshToken) IsExpired() error {
	err := isExpired(t.expireIn)
	if err != nil {
		return err
	}
	return nil
}

func (t *RefreshToken) Id() uuid.UUID {
	return t.id
}

func (t *RefreshToken) UserId() uuid.UUID {
	return t.userId
}

func (t *RefreshToken) Token() string {
	return t.token
}

func (t *RefreshToken) ExpireIn() time.Time {
	return t.expireIn
}

// 期限が切れているかを確認
func isExpired(expireIn time.Time) error {
	if expireIn.After(time.Now()) {
		return nil
	}
	return ErrTokenExpired
}

// JWTの暗号化鍵
func getJwtSecretKey() []byte {
	secret, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		panic("\"JWT_SECRET_KEY\" is not set")
	}
	return []byte(secret)
}
