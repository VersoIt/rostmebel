package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rostmebel/backend/internal/domain/admin"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo      admin.Repository
	jwtSecret string
	accessTTL time.Duration
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
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
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

	if err := u.repo.UpdateRefreshToken(ctx, a.ID, &tokens.RefreshToken); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *UseCase) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	adminID := int64(claims["sub"].(float64))
	a, err := u.repo.GetByID(ctx, adminID)
	if err != nil || a == nil || a.RefreshToken == nil || *a.RefreshToken != refreshToken {
		return nil, fmt.Errorf("invalid refresh token or session expired")
	}

	tokens, err := u.GenerateTokens(adminID)
	if err != nil {
		return nil, err
	}

	if err := u.repo.UpdateRefreshToken(ctx, adminID, &tokens.RefreshToken); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *UseCase) GenerateTokens(adminID int64) (*TokenPair, error) {
	accessClaims := jwt.MapClaims{
		"sub": adminID,
		"exp": time.Now().Add(u.accessTTL).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"sub": adminID,
		"exp": time.Now().Add(u.refreshTTL).Unix(),
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
