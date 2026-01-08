package session

import (
	"context"
	"fmt"
	"os"
	"time"

	authcontext "github.com/AlphaCodinggroup/alpha_auth-api/pkg/context"
	"github.com/dgrijalva/jwt-go"
)

type AccessTokenClaim struct {
	AccessTokenInfo
	KeyType string
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	ID      uint
	Hash    string
	KeyType string
	jwt.StandardClaims
}

type jwtService struct {
	secretKey        string
	refreshSecretKey string
}

func NewService() *jwtService {
	return &jwtService{
		secretKey:        os.Getenv("SECRET"),
		refreshSecretKey: os.Getenv("REFRESH_SECRET"),
	}
}

func (service *jwtService) CreateAccessToken(tokenInfo AccessTokenInfo) (tokenString string, err error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &AccessTokenClaim{
		AccessTokenInfo: AccessTokenInfo{
			ID:       tokenInfo.ID,
			Rol:      tokenInfo.Rol,
			Username: tokenInfo.Username,
			Hash:     tokenInfo.Hash,
		},
		KeyType: "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "auth.service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(service.secretKey))
	if err != nil {
		return tokenString, fmt.Errorf("error signing token: %w", err)
	}

	return
}

func (service *jwtService) CreateRefreshToken(tokenInfo RefreshTokenInfo) (string, error) {
	expirationTime := time.Now().Add(6 * 30 * 24 * time.Hour)

	claims := RefreshTokenClaims{
		ID:      tokenInfo.ID,
		Hash:    tokenInfo.Hash,
		KeyType: "refresh",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "auth.service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(service.refreshSecretKey))
	if err != nil {
		return tokenString, fmt.Errorf("error signing refresh token: %w", err)
	}

	return tokenString, nil
}

func (service *jwtService) ValidateAccessToken(ctx context.Context, tokenString string) (AccessTokenInfo, error) {
	logger := authcontext.Logger(ctx)

	accesTokenInfo := AccessTokenInfo{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaim{},
		func(token *jwt.Token) (any, error) {
			return []byte(service.secretKey), nil
		})

	if err != nil {
		logger.Error("error in function ParseWithClaims: %w", err)
		return accesTokenInfo, fmt.Errorf("error parsing token")
	}

	if !token.Valid {
		return accesTokenInfo, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(*AccessTokenClaim)
	if !ok {
		return accesTokenInfo, fmt.Errorf("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return accesTokenInfo, fmt.Errorf("token expired")
	}

	accesTokenInfo.ID = claims.ID
	accesTokenInfo.Username = claims.Username
	accesTokenInfo.Rol = claims.Rol
	accesTokenInfo.Hash = claims.Hash
	accesTokenInfo.Exp = claims.ExpiresAt

	return accesTokenInfo, nil
}

func (service *jwtService) ValidateRefreshToken(ctx context.Context, tokenString string) (RefreshTokenInfo, error) {
	logger := authcontext.Logger(ctx)
	tokenInfo := RefreshTokenInfo{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&RefreshTokenClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(service.refreshSecretKey), nil
		})

	if err != nil {
		logger.Error("error in function ParseWithClaims: %w", err)
		return tokenInfo, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid || claims.ID == 0 || claims.KeyType != "refresh" {
		return tokenInfo, fmt.Errorf("invalid token: authentication failed")
	}

	tokenInfo.ID = claims.ID
	tokenInfo.Hash = claims.Hash

	return tokenInfo, nil
}
