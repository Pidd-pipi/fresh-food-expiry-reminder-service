package handler

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blueship581/cyfreshfood/internal/middleware"
	"github.com/blueship581/cyfreshfood/internal/model"
	"github.com/blueship581/cyfreshfood/internal/repository"
	"github.com/blueship581/cyfreshfood/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.FamilyGroup{}, &model.FamilyMember{},
		&model.FoodItem{}, &model.ConsumptionRecord{}, &model.Notification{}, &model.Recipe{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func handlerLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestFamilyGroupDetail404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newHandlerTestDB(t)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	familySvc := service.NewFamilyGroupService(groupRepo, memberRepo, handlerLogger())
	memberSvc := service.NewFamilyMemberService(memberRepo, handlerLogger())
	h := NewFamilyGroupHandler(familySvc, memberSvc, handlerLogger())

	r := gin.New()
	r.Use(middleware.ErrorHandler(handlerLogger()))
	r.GET("/family-groups/:id", h.Detail)

	req := httptest.NewRequest(http.MethodGet, "/family-groups/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", w.Code, w.Body.String())
	}
}
