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
	"github.com/blueship581/cyfreshfood/internal/util"
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

func TestRecipeDetail404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newHandlerTestDB(t)
	recipeRepo := repository.NewRecipeRepository(db)
	foodRepo := repository.NewFoodItemRepository(db)
	groupRepo := repository.NewFamilyGroupRepository(db)
	memberRepo := repository.NewFamilyMemberRepository(db)
	familySvc := service.NewFamilyGroupService(groupRepo, memberRepo, handlerLogger())
	recipeSvc := service.NewRecipeService(recipeRepo, foodRepo, familySvc, util.NewFoodCalculator(), handlerLogger())
	h := NewRecipeHandler(recipeSvc, handlerLogger())

	r := gin.New()
	r.Use(middleware.ErrorHandler(handlerLogger()))
	r.GET("/recipes/:id", h.Detail)

	req := httptest.NewRequest(http.MethodGet, "/recipes/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", w.Code, w.Body.String())
	}
}
