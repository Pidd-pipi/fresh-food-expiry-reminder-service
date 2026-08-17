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

func TestStatsService_Dashboard_CountsAllFoods(t *testing.T) {
    db := newTestDB(t)
    groupRepo := repository.NewFamilyGroupRepository(db)
    memberRepo := repository.NewFamilyMemberRepository(db)
    foodRepo := repository.NewFoodItemRepository(db)
    consumeRepo := repository.NewConsumptionRecordRepository(db)
    notifyRepo := repository.NewNotificationRepository(db)
    familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
    memberSvc := NewFamilyMemberService(memberRepo, testLogger())
    svc := NewStatsService(foodRepo, consumeRepo, notifyRepo, familySvc, memberSvc, util.NewFoodCalculator(), testLogger())

    group, err := familySvc.Create(context.Background(), 1, "测试家庭")
    if err != nil { t.Fatalf("Create() error = %v", err) }

    expiry := time.Now().AddDate(0, 0, 10)
    for i := 0; i < 1001; i++ {
        item := &model.FoodItem{
            FamilyID: group.ID, Name: "食品", Category: constants.FoodCategoryDairy,
            Quantity: 1, Unit: "盒", ShelfLifeDays: 10, StorageLocation: constants.StorageFridge,
            Status: constants.FreshnessFresh, CreatorID: 1, ExpiryDate: &expiry,
        }
        if err := db.Create(item).Error; err != nil { t.Fatalf("seed food %d: %v", i, err) }
    }

    data, err := svc.Dashboard(context.Background(), 1, group.ID)
    if err != nil { t.Fatalf("Dashboard() error = %v", err) }
    if data.TotalItems != 1001 {
        t.Fatalf("TotalItems = %d, want 1001", data.TotalItems)
    }
    var dairyCount int64
    for _, c := range data.ByCategory {
        if c.Category == constants.FoodCategoryDairy {
            dairyCount = c.Count
        }
    }
    if dairyCount != 1001 {
        t.Fatalf("dairy count = %d, want 1001", dairyCount)
    }
}
