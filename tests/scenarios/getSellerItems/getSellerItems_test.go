package getselleritems

import (
	"encoding/json"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	advertisementClient "github.com/shottaneflow/avito/internal/client/http/advertisement"
	manager "github.com/shottaneflow/avito/internal/managers/advertisement"
	"github.com/shottaneflow/avito/internal/managers/advertisement/models"
)

const testSellerID = int64(111299)

type GetSellerItemsSuite struct {
	suite.Suite
	advertisementID string
}

func TestGetSellerItemsSuite(t *testing.T) {
	suite.RunSuite(t, new(GetSellerItemsSuite))
}

func (s *GetSellerItemsSuite) BeforeAll(t provider.T) {
	t.WithNewStep("создание объявления для продавца", func(sCtx provider.StepCtx) {
		req := models.CreateAdvertisementRequest{
			SellerID:  testSellerID,
			Name:      "seller-test",
			Price:     750,
			Statistic: models.Statistic{Likes: 1, ViewCount: 1, Contacts: 1},
		}
		ad, code := manager.CreateAdvertisement(t, req)
		sCtx.Require().Equal(200, code, "Не удалось создать объявление")
		s.advertisementID = ad.ID
		sCtx.Logf("Создано объявление с id=%s", s.advertisementID)
	})
}

func (s *GetSellerItemsSuite) AfterAll(t provider.T) {
	t.WithNewStep("Teardown: удаление тестового объявления", func(sCtx provider.StepCtx) {
		if s.advertisementID != "" {
			code := manager.DeleteAdvertisement(t, s.advertisementID)
			sCtx.Logf("Удалено %s, статус %d", s.advertisementID, code)
		}
	})
}

func (s *GetSellerItemsSuite) TestTC_SELLER_01_ExistingSeller(t provider.T) {
	t.Title("TC-SELLER-01: Получение объявлений продавца")
	t.Description("GET /api/1/{sellerID}/item")
	t.Tags("positive", "smoke")
	t.Severity(allure.CRITICAL)

	var items []models.Advertisement
	var code int

	t.WithNewStep("Отправка GET /api/1/{sellerID}/item", func(sCtx provider.StepCtx) {
		items, code = manager.GetSellerAdvertisements(t, testSellerID)
		b, _ := json.MarshalIndent(items, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Response body", allure.JSON, b))
		sCtx.Require().Equal(200, code)
		sCtx.Require().NotEmpty(items)
	})

	t.WithNewStep("Проверка структуры элементов массива", func(sCtx provider.StepCtx) {
		first := items[0]
		sCtx.Assert().NotEmpty(first.ID)
		sCtx.Assert().NotZero(first.SellerID)
		sCtx.Assert().NotEmpty(first.Name)
		sCtx.Assert().NotZero(first.Price)
		sCtx.Assert().NotEmpty(first.CreatedAt)
	})

	t.WithNewStep("Проверка наличия созданного объявления в списке", func(sCtx provider.StepCtx) {
		found := false
		for _, item := range items {
			if item.ID == s.advertisementID {
				found = true
				break
			}
		}
		sCtx.Assert().True(found, "Созданное объявление должно быть в списке")
	})
}

func (s *GetSellerItemsSuite) TestTC_SELLER_02_SellerNoItems(t provider.T) {
	t.Title("TC-SELLER-02: Продавец без объявлений")
	t.Description("GET /api/1/{sellerID}/item для продавца без объявлений")
	t.Tags("positive")
	t.Severity(allure.NORMAL)

	t.WithNewStep("Отправка запроса для продавца без объявлений", func(sCtx provider.StepCtx) {
		emptySellerID := int64(999999999)
		items, code := manager.GetSellerAdvertisements(t, emptySellerID)
		if len(items) != 0 {
			for _, item := range items {
				_ = manager.DeleteAdvertisement(t, item.ID)
			}
			items, code = manager.GetSellerAdvertisements(t, emptySellerID)
		}
		sCtx.Assert().Equal(200, code)
		sCtx.Assert().Empty(items)
	})
}

func (s *GetSellerItemsSuite) TestTC_SELLER_03_InvalidSellerID(t provider.T) {
	t.Title("TC-SELLER-03: Некорректный sellerID")
	t.Description("GET /api/1/abc/item ")
	t.Tags("negative", "seller")
	t.Severity(allure.NORMAL)

	t.WithNewStep("Отправка запроса с sellerID=abc", func(sCtx provider.StepCtx) {
		resp := advertisementClient.HttpGetSellerAdvertisementsRaw(t, "abc")
		defer resp.Body.Close()
		sCtx.Assert().Equal(400, resp.StatusCode)
	})
}
