package service

import (
	"context"
	"testing"
	"time"

	"github.com/blueship581/cyfreshfood/internal/constants"
	"github.com/blueship581/cyfreshfood/internal/model"
	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
)

func TestFoodItemService_ListRefreshesStatusBeforeFiltering(t *testing.T) {
	db := newTestDB(t)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	foodRepo := repository.NewFoodItemRepository(db)
	consumeRepo := repository.NewConsumptionRecordRepository(db)
	familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	svc := NewFoodItemService(foodRepo, consumeRepo, familySvc, util.NewFoodCalculator(), testLogger())
	ctx := context.Background()

	group, err := familySvc.Create(ctx, 1, "测试家庭")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	expiry := time.Now().AddDate(0, 0, -2)
	item := &model.FoodItem{
		FamilyID: group.ID, Name: "过期牛奶", Category: constants.FoodCategoryDairy,
		Quantity: 1, Unit: "盒", ShelfLifeDays: 7, StorageLocation: constants.StorageFridge,
		Status: constants.FreshnessFresh, CreatorID: 1, ExpiryDate: &expiry,
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("seed food: %v", err)
	}

	items, _, err := svc.List(ctx, 1, group.ID, "", constants.FreshnessFresh, "", "", 1, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("fresh items = %d, want 0 after refreshing status", len(items))
	}
}
