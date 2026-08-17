package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/blueship581/cyfreshfood/internal/constants"
	"github.com/blueship581/cyfreshfood/internal/model"
	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/util"
	"gorm.io/gorm"
)

// ReminderService 临期/过期扫描：刷新食品状态并生成站内通知。
type ReminderService struct {
	foodRepo   *repository.FoodItemRepository
	notifyRepo *repository.NotificationRepository
	calculator *util.FoodCalculator
	log        *slog.Logger
}

// NewReminderService 构造临期扫描服务。
func NewReminderService(foodRepo *repository.FoodItemRepository, notifyRepo *repository.NotificationRepository, calculator *util.FoodCalculator, log *slog.Logger) *ReminderService {
	return &ReminderService{foodRepo: foodRepo, notifyRepo: notifyRepo, calculator: calculator, log: log}
}

// Scan 执行一轮临期/过期扫描（全家庭），返回新增通知数。
func (s *ReminderService) Scan(ctx context.Context) (int, error) {
	s.log.InfoContext(ctx, constants.LOG_EXPIRY_SCAN_STARTED)
	items, err := s.foodRepo.ListReminderCandidates(0)
	if err != nil {
		return 0, util.LogError(s.log, ctx, constants.LOG_EXPIRY_SCAN_FAILED, fmt.Errorf("scan foods: %w", err))
	}
	created := 0
	for _, item := range items {
		// 已消耗食品不再参与扫描。
		if item.Status == constants.FreshnessConsumed {
			continue
		}
		status := s.calculator.ComputeFreshness(item.Status, item.ExpiryDate)
		// 仅临期/过期需发通知；仍新鲜则跳过，避免无谓的状态更新与通知查重。
		var typ, title string
		switch status {
		case constants.FreshnessExpiring:
			typ, title = constants.NotificationExpiring, constants.MsgFoodExpiring
		case constants.FreshnessExpired:
			typ, title = constants.NotificationExpired, constants.MsgFoodExpired
		default:
			continue
		}
		// 通知去重：同一食品同一类型只发一次。注意这必须在状态判断之后、事务之前，
		// 否则「录入即临期」的食品（status 落库即 expiring、从未经历 fresh→expiring 扫描转换）
		// 会被 status==item.Status 早退逻辑跳过而永不发通知。
		exists, _ := s.notifyRepo.HasForFoodAndType(item.ID, typ)
		if exists {
			continue
		}
		content := fmt.Sprintf("%s 已%s，请及时处理。", item.Name, util.FreshnessStatusText(status))
		n := &model.Notification{
			FamilyID: item.FamilyID, FoodItemID: item.ID, Type: typ,
			Title: title, Content: content,
		}
		err = s.notifyRepo.Transaction(func(tx *gorm.DB) error {
			if err := s.foodRepo.WithTx(tx).UpdateStatus(item.ID, status); err != nil {
				return fmt.Errorf("update food status: %w", err)
			}
			if err := s.notifyRepo.WithTx(tx).Create(n); err != nil {
				return fmt.Errorf("create notification: %w", err)
			}
			return nil
		})
		if err != nil {
			return created, util.LogError(s.log, ctx, constants.LOG_EXPIRY_SCAN_FAILED, err)
		}
		s.log.InfoContext(ctx, constants.LOG_FOOD_STATUS_REFRESHED, "food_id", item.ID, "status", status)
		s.log.InfoContext(ctx, constants.LOG_NOTIFICATION_CREATED, "food_item_id", item.ID, "type", typ)
		created++
	}
	s.log.InfoContext(ctx, constants.LOG_EXPIRY_SCAN_FINISHED, "created", created)
	return created, nil
}
