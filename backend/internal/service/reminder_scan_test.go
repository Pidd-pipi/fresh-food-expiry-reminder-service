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
