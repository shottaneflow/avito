package getadvertisement

import (
	"encoding/json"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	manager "github.com/shottaneflow/avito/internal/managers/advertisement"
	"github.com/shottaneflow/avito/internal/managers/advertisement/models"
)

type GetAdvertisementSuite struct {
	suite.Suite
	advertisementID string
}

func TestGetAdvertisementSuite(t *testing.T) {
	suite.RunSuite(t, new(GetAdvertisementSuite))
}

func (s *GetAdvertisementSuite) BeforeAll(t provider.T) {
	t.WithNewStep("создание тестового объявления", func(sCtx provider.StepCtx) {
		req := models.CreateAdvertisementRequest{
			SellerID:  111289,
			Name:      "get-test",
			Price:     500,
			Statistic: models.Statistic{Likes: 3, ViewCount: 7, Contacts: 2},
		}
		ad, code := manager.CreateAdvertisement(t, req)
		sCtx.Require().Equal(200, code, "Не удалось создать объявление для теста")
		sCtx.Require().NotEmpty(ad.ID)
		s.advertisementID = ad.ID
		sCtx.Logf("Создано объявление с id=%s", s.advertisementID)
	})
}

func (s *GetAdvertisementSuite) AfterAll(t provider.T) {
	t.WithNewStep("Teardown: удаление тестового объявления", func(sCtx provider.StepCtx) {
		if s.advertisementID != "" {
			code := manager.DeleteAdvertisement(t, s.advertisementID)
			sCtx.Logf("Удалено объявление %s, статус %d", s.advertisementID, code)
		}
	})
}

func (s *GetAdvertisementSuite) TestTC_GET_01_ExistingID(t provider.T) {
	t.Title("TC-GET-01: Получение существующего объявления")
	t.Description("GET /api/1/item/{id} с существующим id")
	t.Tags("positive", "get")
	t.Severity(allure.CRITICAL)

	var ad models.Advertisement
	var code int

	t.WithNewStep("Отправка GET /api/1/item/{id}", func(sCtx provider.StepCtx) {
		sCtx.Logf("Запрос объявления с id=%s", s.advertisementID)
		ad, code = manager.GetAdvertisement(t, s.advertisementID)
		b, _ := json.MarshalIndent(ad, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Response body", allure.JSON, b))
		sCtx.Require().Equal(200, code)
	})

	t.WithNewStep("Проверка структуры ответа", func(sCtx provider.StepCtx) {
		sCtx.Assert().NotEmpty(ad.ID)
		sCtx.Assert().NotZero(ad.SellerID)
		sCtx.Assert().NotEmpty(ad.Name)
		sCtx.Assert().NotZero(ad.Price)
		sCtx.Assert().NotEmpty(ad.CreatedAt)
		sCtx.Assert().GreaterOrEqual(ad.Statistic.Likes, int64(0))
		sCtx.Assert().GreaterOrEqual(ad.Statistic.ViewCount, int64(0))
		sCtx.Assert().GreaterOrEqual(ad.Statistic.Contacts, int64(0))
	})
}

func (s *GetAdvertisementSuite) TestTC_GET_02_NonExistingID(t provider.T) {
	t.Title("TC-GET-02: Несуществующий id")
	t.Description("GET /api/1/item/{id} с несуществующим id ")
	t.Tags("negative", "get")
	t.Severity(allure.NORMAL)

	t.WithNewStep("Отправка запроса с несуществующим id", func(sCtx provider.StepCtx) {
		nonExistingID := "6eb9b015-b2cb-4502-8b61-c58903106074"
		manager.DeleteAdvertisement(t, nonExistingID)
		sCtx.Logf("id=%s", nonExistingID)
		_, code := manager.GetAdvertisement(t, nonExistingID)
		sCtx.Assert().Equal(404, code)
	})
}

func (s *GetAdvertisementSuite) TestTC_GET_03_InvalidID(t provider.T) {
	t.Title("TC-GET-03: Некорректный id")
	t.Description("GET /api/1/item/abc ")
	t.Tags("negative", "get")
	t.Severity(allure.NORMAL)

	t.WithNewStep("Отправка запроса с id=abc", func(sCtx provider.StepCtx) {
		_, code := manager.GetAdvertisement(t, "abc")
		sCtx.Assert().Equal(400, code)
	})
}

func (s *GetAdvertisementSuite) TestTC_GET_04_EmptyID(t provider.T) {
	t.Title("TC-GET-04: Пустой id")
	t.Description("GET /api/1/item/ без id")
	t.Tags("negative", "get")
	t.Severity(allure.NORMAL)

	t.WithNewStep("Отправка запроса с пустым id", func(sCtx provider.StepCtx) {
		_, code := manager.GetAdvertisement(t, "")
		sCtx.Assert().Equal(404, code)
	})
}

func (s *GetAdvertisementSuite) TestTC_GET_05_Idempotency(t provider.T) {
	t.Title("TC-GET-05: Идемпотентность GET")
	t.Description("Два одинаковых GET запроса возвращают идентичные данные")
	t.Tags("positive", "get", "idempotency")
	t.Severity(allure.NORMAL)

	var ad1, ad2 models.Advertisement

	t.WithNewStep("Первый GET запрос", func(sCtx provider.StepCtx) {
		var code int
		ad1, code = manager.GetAdvertisement(t, s.advertisementID)
		sCtx.Require().Equal(200, code)
	})

	t.WithNewStep("Второй GET запрос", func(sCtx provider.StepCtx) {
		var code int
		ad2, code = manager.GetAdvertisement(t, s.advertisementID)
		sCtx.Require().Equal(200, code)
	})

	t.WithNewStep("Сравнение ответов", func(sCtx provider.StepCtx) {
		sCtx.Assert().Equal(ad1, ad2, "Ответы должны быть идентичны")
	})
}
