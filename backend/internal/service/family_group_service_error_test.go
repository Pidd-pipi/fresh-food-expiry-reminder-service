package service

import (
	"context"
	"errors"
	"testing"

	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
)

func TestFamilyGroupGetMissing(t *testing.T) {
	db := newTestDB(t)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	svc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 999)
	if err == nil {
		t.Fatal("expected error for missing family group")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.HTTPStatus != 404 {
		t.Fatalf("HTTPStatus = %d, want 404", appErr.HTTPStatus)
	}
}

func TestFamilyGroupJoinBadCode(t *testing.T) {
	db := newTestDB(t)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	svc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	ctx := context.Background()

	_, err := svc.JoinByCode(ctx, 1, "NO_SUCH_CODE")
	if err == nil {
		t.Fatal("expected error for invalid invite code")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.HTTPStatus != 404 {
		t.Fatalf("HTTPStatus = %d, want 404", appErr.HTTPStatus)
	}
}
