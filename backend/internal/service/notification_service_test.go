package service

import (
	"context"
	"testing"
	"time"

	"github.com/blueship581/cyfreshfood/internal/constants"
	"github.com/blueship581/cyfreshfood/internal/model"
	"github.com/blueship581/cyfreshfood/internal/repository"
	"gorm.io/gorm"
)

type notifFixture struct {
	familyID  uint
	n1, n2    uint
	repo      *repository.NotificationRepository
	memberRepo *repository.FamilyMemberRepository
	svc       *NotificationService
}

func seedNotifs(t *testing.T, db *gorm.DB, adminID uint) *notifFixture {
	t.Helper()
	group := &model.FamilyGroup{Name: "测试家庭", OwnerID: adminID, InviteCode: "INV008"}
	if err := db.Create(group).Error; err != nil {
		t.Fatalf("create group: %v", err)
	}
	admin := &model.FamilyMember{FamilyID: group.ID, UserID: adminID, Role: constants.FamilyRoleAdmin}
	member := &model.FamilyMember{FamilyID: group.ID, UserID: adminID + 1, Role: constants.FamilyRoleMember}
	if err := db.Create(&[]model.FamilyMember{*admin, *member}).Error; err != nil {
		t.Fatalf("create members: %v", err)
	}
	food := &model.FoodItem{FamilyID: group.ID, Name: "测试食品", Category: constants.FoodCategoryDairy, Quantity: 1, Unit: "盒", ShelfLifeDays: 5, StorageLocation: constants.StorageFridge, Status: constants.FreshnessFresh, CreatorID: adminID}
	if err := db.Create(food).Error; err != nil {
		t.Fatalf("create food: %v", err)
	}
	n1 := &model.Notification{FamilyID: group.ID, FoodItemID: food.ID, Type: constants.NotificationExpiring, Title: "临期", Content: "c1", SendAt: time.Now()}
	n2 := &model.Notification{FamilyID: group.ID, FoodItemID: food.ID, Type: constants.NotificationExpired, Title: "过期", Content: "c2", SendAt: time.Now()}
	created := []model.Notification{*n1, *n2}
	if err := db.Create(&created).Error; err != nil {
		t.Fatalf("create notifications: %v", err)
	}
	n1.ID = created[0].ID
	n2.ID = created[1].ID
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	notifyRepo := repository.NewNotificationRepository(db)
	familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
	svc := NewNotificationService(notifyRepo, familySvc, testLogger())
	return &notifFixture{familyID: group.ID, n1: n1.ID, n2: n2.ID, repo: notifyRepo, memberRepo: memberRepo, svc: svc}
}

func TestNotifyMarkReadTarget(t *testing.T) {
	db := newTestDB(t)
	fx := seedNotifs(t, db, 1)
	ctx := context.Background()

	if err := fx.svc.MarkRead(ctx, 1, fx.n1); err != nil {
		t.Fatalf("MarkRead() error = %v", err)
	}

	n1, err := fx.repo.FindByID(fx.n1)
	if err != nil {
		t.Fatalf("find n1: %v", err)
	}
	n2, err := fx.repo.FindByID(fx.n2)
	if err != nil {
		t.Fatalf("find n2: %v", err)
	}
	if !n1.IsRead {
		t.Fatalf("n1 should be read")
	}
	if n2.IsRead {
		t.Fatalf("n2 should remain unread")
	}
}

func TestNotifyMarkAllReadState(t *testing.T) {
	db := newTestDB(t)
	fx := seedNotifs(t, db, 1)
	ctx := context.Background()

	if err := fx.svc.MarkAllRead(ctx, 1, fx.familyID); err != nil {
		t.Fatalf("MarkAllRead() error = %v", err)
	}

	items, _, err := fx.repo.ListByFamily(fx.familyID, false, 1, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, n := range items {
		if !n.IsRead {
			t.Fatalf("notification %d should be read", n.ID)
		}
		if n.ReadAt == nil {
			t.Fatalf("notification %d read_at should be set", n.ID)
		}
	}
}

func TestNotifyUnreadCount(t *testing.T) {
	db := newTestDB(t)
	fx := seedNotifs(t, db, 1)
	ctx := context.Background()

	count, err := fx.svc.UnreadCount(ctx, 2, fx.familyID)
	if err != nil {
		t.Fatalf("UnreadCount() error = %v", err)
	}
	if count != 2 {
		t.Fatalf("unread count = %d, want 2", count)
	}
}
