package admin

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rostmebel/backend/internal/domain/admin"
	"github.com/rostmebel/backend/internal/domain/apperror"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo       admin.Repository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewUseCase(repo admin.Repository, jwtSecret string, accessTTL, refreshTTL time.Duration) *UseCase {
	return &UseCase{
		repo:       repo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *UseCase) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	a, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, apperror.New(apperror.CodeAuthInvalidCredentials, "Invalid credentials", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)); err != nil {
		return nil, apperror.New(apperror.CodeAuthInvalidCredentials, "Invalid credentials", nil)
	}

	now := time.Now()
	a.LastLoginAt = &now
	if err := u.repo.Update(ctx, a); err != nil {
		return nil, err
	}

	tokens, err := u.GenerateTokens(a.ID)
	if err != nil {
		return nil, err
	}

	refreshTokenHash, err := hashRefreshToken(tokens.RefreshToken)
	if err != nil {
		return nil, err
	}
	if err := u.repo.UpdateRefreshToken(ctx, a.ID, &refreshTokenHash); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *UseCase) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return nil, apperror.New(apperror.CodeAuthInvalidRefreshToken, "Invalid refresh token", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperror.New(apperror.CodeAuthInvalidRefreshToken, "Invalid refresh token claims", nil)
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return nil, apperror.New(apperror.CodeAuthInvalidRefreshToken, "Invalid refresh token subject", nil)
	}
	tokenType, ok := claims["typ"].(string)
	if !ok || tokenType != "refresh" {
		return nil, apperror.New(apperror.CodeAuthInvalidRefreshToken, "Invalid refresh token type", nil)
	}
	adminID := int64(sub)
	a, err := u.repo.GetByID(ctx, adminID)
	if err != nil || a == nil || a.RefreshToken == nil || !compareRefreshToken(*a.RefreshToken, refreshToken) {
		return nil, apperror.New(apperror.CodeAuthInvalidRefreshToken, "Invalid refresh token", nil)
	}

	tokens, err := u.GenerateTokens(adminID)
	if err != nil {
		return nil, err
	}

	refreshTokenHash, err := hashRefreshToken(tokens.RefreshToken)
	if err != nil {
		return nil, err
	}
	if err := u.repo.UpdateRefreshToken(ctx, adminID, &refreshTokenHash); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *UseCase) GenerateTokens(adminID int64) (*TokenPair, error) {
	now := time.Now()
	accessJTI, err := randomTokenID()
	if err != nil {
		return nil, err
	}
	refreshJTI, err := randomTokenID()
	if err != nil {
		return nil, err
	}

	accessClaims := jwt.MapClaims{
		"sub": adminID,
		"exp": now.Add(u.accessTTL).Unix(),
		"iat": now.Unix(),
		"jti": accessJTI,
		"typ": "access",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"sub": adminID,
		"exp": now.Add(u.refreshTTL).Unix(),
		"iat": now.Unix(),
		"jti": refreshJTI,
		"typ": "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessString,
		RefreshToken: refreshString,
	}, nil
}

func (u *UseCase) GetStats(ctx context.Context) (*admin.Stats, error) {
	return u.repo.GetStats(ctx)
}

func (u *UseCase) Logout(ctx context.Context, adminID int64) error {
	return u.repo.UpdateRefreshToken(ctx, adminID, nil)
}

func hashRefreshToken(token string) (string, error) {
	digest := refreshTokenDigest(token)
	hash, err := bcrypt.GenerateFromPassword([]byte(digest), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func compareRefreshToken(hash, token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(refreshTokenDigest(token))) == nil
}

func refreshTokenDigest(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func randomTokenID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
