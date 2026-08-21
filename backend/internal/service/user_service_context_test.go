package service

import (
	"context"
	"testing"

	"github.com/blueship581/cyfreshfood/internal/repository"
)

func TestUserGetCanceled(t *testing.T) {
	db := newTestDB(t)
	svc := NewUserService(repository.NewUserRepository(db), "test-secret", 24, testLogger())
	ctx := context.Background()

	user, _, err := svc.Register(ctx, "13900000001", "pass123", "测试用户")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := svc.GetByID(cancelled, user.ID); err == nil {
		t.Fatal("expected context canceled error for GetByID")
	}
}

func TestUserLoginCanceled(t *testing.T) {
	db := newTestDB(t)
	svc := NewUserService(repository.NewUserRepository(db), "test-secret", 24, testLogger())
	ctx := context.Background()

	if _, _, err := svc.Register(ctx, "13900000002", "pass123", "测试用户2"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := svc.Login(cancelled, "13900000002", "pass123"); err == nil {
		t.Fatal("expected context canceled error for Login")
	}
}
