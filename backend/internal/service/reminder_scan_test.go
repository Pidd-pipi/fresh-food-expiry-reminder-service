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

func TestReminderService_Scan_CreatesExpiringNotification(t *testing.T) {
	db := newTestDB(t)
	foodRepo := repository.NewFoodItemRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	svc := NewReminderService(foodRepo, notifyRepo, util.NewFoodCalculator(), testLogger())

	expiry := time.Now().AddDate(0, 0, 2)
	item := &model.FoodItem{
		FamilyID: 1, Name: "牛奶", Category: constants.FoodCategoryDairy,
		Quantity: 1, Unit: "盒", ShelfLifeDays: 2, StorageLocation: constants.StorageFridge,
		Status: constants.FreshnessFresh, ExpiryDate: &expiry,
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("seed food: %v", err)
	}

	created, err := svc.Scan(context.Background())
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if created != 1 {
		t.Fatalf("created notifications = %d, want 1", created)
	}

	var got model.FoodItem
	if err := db.First(&got, item.ID).Error; err != nil {
		t.Fatalf("reload food: %v", err)
	}
	if got.Status != constants.FreshnessExpiring {
		t.Fatalf("food status = %s, want %s", got.Status, constants.FreshnessExpiring)
	}

	var n model.Notification
	if err := db.Where("food_item_id = ?", item.ID).First(&n).Error; err != nil {
		t.Fatalf("load notification: %v", err)
	}
	if n.Type != constants.NotificationExpiring {
		t.Fatalf("notification type = %s, want %s", n.Type, constants.NotificationExpiring)
	}
}

// TestReminderService_Scan_CoversFoodCreatedExpiring 回归测试：
// 食品在录入时即已临期（剩余 2 天），FoodItemService.Create 会用 ComputeFreshness
// 把 status 落库为 expiring。ListReminderCandidates 必须仍能扫到它并补发临期通知，
// 否则「录入即临期」的食品永远收不到提醒。
func TestReminderService_Scan_CoversFoodCreatedExpiring(t *testing.T) {
	db := newTestDB(t)
	foodRepo := repository.NewFoodItemRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	calc := util.NewFoodCalculator()
	svc := NewReminderService(foodRepo, notifyRepo, calc, testLogger())

	// 模拟真实 Create：落库时 status 由 ComputeFreshness 计算得到 expiring。
	expiry := time.Now().AddDate(0, 0, 2)
	item := &model.FoodItem{
		FamilyID: 1, Name: "临期牛奶", Category: constants.FoodCategoryDairy,
		Quantity: 1, Unit: "盒", ShelfLifeDays: 2, StorageLocation: constants.StorageFridge,
		ExpiryDate: &expiry,
	}
	item.Status = calc.ComputeFreshness("", item.ExpiryDate)
	if item.Status != constants.FreshnessExpiring {
		t.Fatalf("seed status = %s, want %s", item.Status, constants.FreshnessExpiring)
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("seed food: %v", err)
	}

	created, err := svc.Scan(context.Background())
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if created != 1 {
		t.Fatalf("created notifications = %d, want 1", created)
	}
	var n model.Notification
	if err := db.Where("food_item_id = ?", item.ID).First(&n).Error; err != nil {
		t.Fatalf("load notification: %v", err)
	}
	if n.Type != constants.NotificationExpiring {
		t.Fatalf("notification type = %s, want %s", n.Type, constants.NotificationExpiring)
	}
	// 再次扫描应去重，不重复创建通知。
	created2, err := svc.Scan(context.Background())
	if err != nil {
		t.Fatalf("second Scan() error = %v", err)
	}
	if created2 != 0 {
		t.Fatalf("second scan created = %d, want 0 (去重)", created2)
	}
}
