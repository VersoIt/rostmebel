package admin

import (
	"context"
	"testing"
	"time"

	domAdmin "github.com/rostmebel/backend/internal/domain/admin"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginStoresHashedRefreshTokenAndRefreshRotates(t *testing.T) {
	ctx := context.Background()
	repo := newMemoryAdminRepo(t, "admin", "correct-password")
	useCase := NewUseCase(repo, "0123456789abcdef0123456789abcdef", time.Minute, time.Hour)

	tokens, err := useCase.Login(ctx, "admin", "correct-password")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if repo.admin.RefreshToken == nil {
		t.Fatal("refresh token hash was not stored")
	}
	if *repo.admin.RefreshToken == tokens.RefreshToken {
		t.Fatal("refresh token must not be stored in plaintext")
	}
	if !compareRefreshToken(*repo.admin.RefreshToken, tokens.RefreshToken) {
		t.Fatal("stored refresh hash does not match issued token")
	}

	if _, err := useCase.Refresh(ctx, tokens.AccessToken); err == nil {
		t.Fatal("access token must not be accepted as refresh token")
	}

	refreshed, err := useCase.Refresh(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if refreshed.RefreshToken == tokens.RefreshToken {
		t.Fatal("refresh token must rotate on refresh")
	}
	if compareRefreshToken(*repo.admin.RefreshToken, tokens.RefreshToken) {
		t.Fatal("old refresh token must stop matching stored hash after rotation")
	}
	if !compareRefreshToken(*repo.admin.RefreshToken, refreshed.RefreshToken) {
		t.Fatal("new refresh token hash was not stored")
	}
}

type memoryAdminRepo struct {
	admin *domAdmin.Admin
}

func newMemoryAdminRepo(t *testing.T, username, password string) *memoryAdminRepo {
	t.Helper()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	return &memoryAdminRepo{
		admin: &domAdmin.Admin{
			ID:           42,
			Username:     username,
			PasswordHash: string(passwordHash),
			CreatedAt:    time.Now(),
		},
	}
}

func (r *memoryAdminRepo) GetByUsername(_ context.Context, username string) (*domAdmin.Admin, error) {
	if r.admin == nil || r.admin.Username != username {
		return nil, nil
	}
	return cloneAdmin(r.admin), nil
}

func (r *memoryAdminRepo) GetByID(_ context.Context, id int64) (*domAdmin.Admin, error) {
	if r.admin == nil || r.admin.ID != id {
		return nil, nil
	}
	return cloneAdmin(r.admin), nil
}

func (r *memoryAdminRepo) Create(_ context.Context, admin *domAdmin.Admin) error {
	r.admin = cloneAdmin(admin)
	return nil
}

func (r *memoryAdminRepo) Update(_ context.Context, admin *domAdmin.Admin) error {
	currentRefreshToken := r.admin.RefreshToken
	r.admin = cloneAdmin(admin)
	r.admin.RefreshToken = currentRefreshToken
	return nil
}

func (r *memoryAdminRepo) UpdateRefreshToken(_ context.Context, _ int64, token *string) error {
	if token == nil {
		r.admin.RefreshToken = nil
		return nil
	}
	tokenCopy := *token
	r.admin.RefreshToken = &tokenCopy
	return nil
}

func (r *memoryAdminRepo) GetStats(_ context.Context) (*domAdmin.Stats, error) {
	return &domAdmin.Stats{}, nil
}

func cloneAdmin(admin *domAdmin.Admin) *domAdmin.Admin {
	if admin == nil {
		return nil
	}
	cp := *admin
	if admin.RefreshToken != nil {
		token := *admin.RefreshToken
		cp.RefreshToken = &token
	}
	return &cp
}
