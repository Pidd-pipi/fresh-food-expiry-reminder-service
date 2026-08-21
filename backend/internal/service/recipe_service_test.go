package service

import (
	"context"
	"errors"
	"testing"

	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
)

func TestRecipeGetMissing(t *testing.T) {
	db := newTestDB(t)
	recipeRepo := repository.NewRecipeRepository(db)
	foodRepo := repository.NewFoodItemRepository(db)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	svc := NewRecipeService(recipeRepo, foodRepo, familySvc, util.NewFoodCalculator(), testLogger())
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 999)
	if err == nil {
		t.Fatal("expected error for missing recipe")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.HTTPStatus != 404 {
		t.Fatalf("HTTPStatus = %d, want 404", appErr.HTTPStatus)
	}
}
