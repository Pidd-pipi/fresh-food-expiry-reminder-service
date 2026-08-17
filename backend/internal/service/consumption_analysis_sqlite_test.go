package service

import (
    "context"
    "testing"

    "github.com/blueship581/cyfreshfood/internal/repository"
)

func TestConsumptionRecordService_Analysis_OnSQLite(t *testing.T) {
    db := newTestDB(t)
    groupRepo := repository.NewFamilyGroupRepository(db)
    memberRepo := repository.NewFamilyMemberRepository(db)
    familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
    consumeRepo := repository.NewConsumptionRecordRepository(db)
    foodRepo := repository.NewFoodItemRepository(db)
    svc := NewConsumptionRecordService(consumeRepo, foodRepo, familySvc, testLogger())

    group, err := familySvc.Create(context.Background(), 1, "测试家庭")
    if err != nil { t.Fatalf("Create() error = %v", err) }
    _, err = svc.Analysis(context.Background(), 1, group.ID, "2026-08")
    if err != nil { t.Fatalf("Analysis() error = %v", err) }
}
