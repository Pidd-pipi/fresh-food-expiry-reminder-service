package service

import (
    "context"
    "testing"

    "github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
)

func TestFoodItemService_ImportCSV_AllOrNothingOnMalformedRow(t *testing.T) {
    db := newTestDB(t)
    groupRepo := repository.NewFamilyGroupRepository(db)
    memberRepo := repository.NewFamilyMemberRepository(db)
    foodRepo := repository.NewFoodItemRepository(db)
    consumeRepo := repository.NewConsumptionRecordRepository(db)
    familySvc := NewFamilyGroupService(groupRepo, memberRepo, testLogger())
    svc := NewFoodItemService(foodRepo, consumeRepo, familySvc, util.NewFoodCalculator(), testLogger())

    group, err := familySvc.Create(context.Background(), 1, "测试家庭")
    if err != nil { t.Fatalf("Create() error = %v", err) }

    csvText := "name,category,quantity,unit,shelf_life_days,storage_location\n牛奶,dairy,2,盒,7,fridge\n\"未闭合,other,1,份,7,fridge\n"
    if _, _, err := svc.ImportCSV(context.Background(), 1, group.ID, csvText); err == nil {
        t.Fatal("expected error for malformed CSV row")
    }
    count, err := foodRepo.CountByFamily(group.ID)
    if err != nil { t.Fatalf("CountByFamily() error = %v", err) }
    if count != 0 {
        t.Fatalf("food count = %d, want 0 after failed import", count)
    }
}
