package e2e

import (
	"encoding/json"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	manager "github.com/shottaneflow/avito/internal/managers/advertisement"
	"github.com/shottaneflow/avito/internal/managers/advertisement/models"
)

type E2ESuite struct {
	suite.Suite
}

func TestE2ESuite(t *testing.T) {
	suite.RunSuite(t, new(E2ESuite))
}

func (s *E2ESuite) TestTC_E2E_01_FullLifecycle(t provider.T) {
	t.Title("TC-E2E-01: Полный жизненный цикл объявления")
	t.Description("1. Создание  2. GET по ID 3. проверка в списке продавца 4. проверка статистики 5. удаление")
	t.Tags("e2e")
	t.Severity(allure.CRITICAL)

	var createdID string
	const sellerID = int64(111289)

	req := models.CreateAdvertisementRequest{
		SellerID:  sellerID,
		Name:      "e2e-test",
		Price:     15000,
		Statistic: models.Statistic{Likes: 5, ViewCount: 12, Contacts: 10},
	}

	t.WithNewStep("Шаг 1: POST /api/1/item ", func(sCtx provider.StepCtx) {
		b, _ := json.MarshalIndent(req, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, b))

		ad, code := manager.CreateAdvertisement(t, req)
		sCtx.Require().Equal(200, code, "Ожидался 200 при создании")
		sCtx.Require().NotEmpty(ad.ID, "id должен присутствовать")
		createdID = ad.ID

		b, _ = json.MarshalIndent(ad, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Response body", allure.JSON, b))
		sCtx.Logf("Создано объявление id=%s", createdID)
	})

	t.WithNewStep("Шаг 2: GET /api/1/item/{id} ", func(sCtx provider.StepCtx) {
		ad, code := manager.GetAdvertisement(t, createdID)
		sCtx.Require().Equal(200, code)

		sCtx.Assert().Equal(req.SellerID, ad.SellerID)
		sCtx.Assert().Equal(req.Name, ad.Name)
		sCtx.Assert().Equal(req.Price, ad.Price)
		sCtx.Assert().Equal(req.Statistic.Likes, ad.Statistic.Likes)
		sCtx.Assert().Equal(req.Statistic.ViewCount, ad.Statistic.ViewCount)
		sCtx.Assert().Equal(req.Statistic.Contacts, ad.Statistic.Contacts)
	})

	t.WithNewStep("Шаг 3: GET /api/1/{sellerID}/item ", func(sCtx provider.StepCtx) {
		items, code := manager.GetSellerAdvertisements(t, sellerID)
		sCtx.Require().Equal(200, code)

		found := false
		for _, item := range items {
			if item.ID == createdID {
				found = true
				break
			}
		}
		sCtx.Assert().True(found, "Объявление %s должно быть в списке продавца", createdID)
	})

	t.WithNewStep("Шаг 4: GET /api/1/statistic/{id} ", func(sCtx provider.StepCtx) {
		stat, code := manager.GetStatistic(t, createdID)
		b, _ := json.MarshalIndent(stat, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Statistic response", allure.JSON, b))
		sCtx.Require().Equal(200, code)

		sCtx.Assert().Equal(req.Statistic.Likes, stat.Likes)
		sCtx.Assert().Equal(req.Statistic.ViewCount, stat.ViewCount)
		sCtx.Assert().Equal(req.Statistic.Contacts, stat.Contacts)
	})

	t.WithNewStep("Шаг 5: DELETE /api/2/item/{id} ", func(sCtx provider.StepCtx) {
		code := manager.DeleteAdvertisement(t, createdID)
		sCtx.Assert().Equal(200, code, "Ожидался 200 при удалении")
		sCtx.Logf("Удалено объявление %s", createdID)
	})
}
