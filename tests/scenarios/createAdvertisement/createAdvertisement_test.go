package createadvertisement

import (
	"encoding/json"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	manager "github.com/shottaneflow/avito/internal/managers/advertisement"
	"github.com/shottaneflow/avito/internal/managers/advertisement/models"
)

type CreateAdvertisementSuite struct {
	suite.Suite
	createdIDs []string
}

func TestCreateAdvertisementSuite(t *testing.T) {
	suite.RunSuite(t, new(CreateAdvertisementSuite))
}

func (s *CreateAdvertisementSuite) AfterEach(t provider.T) {
	t.WithNewStep("Teardown: удаление созданных объявлений", func(sCtx provider.StepCtx) {
		for _, id := range s.createdIDs {
			code := manager.DeleteAdvertisement(t, id)
			sCtx.Logf("Удалено объявление %s, статус %d", id, code)
		}
		s.createdIDs = nil
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_01_SuccessAllFields(t provider.T) {
	t.Title("TC-POST-01: Успешное создание объявления со всеми полями")
	t.Description("POST /api/1/item с валидными данными  , в ответе id, данные совпадают")
	t.Tags("positive", "post")
	t.Severity(allure.CRITICAL)

	var req models.CreateAdvertisementRequest
	t.WithNewStep("Подготовка тела запроса", func(sCtx provider.StepCtx) {
		req = models.CreateAdvertisementRequest{
			SellerID:  111299,
			Name:      "test",
			Price:     1000,
			Statistic: models.Statistic{Likes: 10, ViewCount: 10, Contacts: 10},
		}
		b, _ := json.MarshalIndent(req, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, b))
	})

	var ad models.Advertisement
	var code int
	t.WithNewStep("Отправка POST /api/1/item", func(sCtx provider.StepCtx) {
		ad, code = manager.CreateAdvertisement(t, req)
		b, _ := json.MarshalIndent(ad, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Response body", allure.JSON, b))
		sCtx.Require().Equal(200, code)
		sCtx.Require().NotEmpty(ad.ID)
	})
	t.WithNewStep("Проверка данных ответа", func(sCtx provider.StepCtx) {
		sCtx.Assert().Equal(int64(111289), ad.SellerID)
		sCtx.Assert().Equal("test", ad.Name)
		sCtx.Assert().Equal(int64(1000), ad.Price)
		sCtx.Assert().Equal(int64(10), ad.Statistic.Likes)
		sCtx.Assert().Equal(int64(10), ad.Statistic.ViewCount)
		sCtx.Assert().Equal(int64(10), ad.Statistic.Contacts)
	})

	s.createdIDs = append(s.createdIDs, ad.ID)
}

func (s *CreateAdvertisementSuite) TestTC_POST_02_MissingSellerID(t provider.T) {
	t.Title("TC-POST-02: Отсутствие sellerID")
	t.Description("Запрос без sellerID ")
	t.Tags("negative", "post")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"name":"test","price":1000,"statistics":{"likes":10,"viewCount":10,"contacts":10}}`)
	t.WithNewStep("Отправка запроса без sellerID", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_03_InvalidPriceType(t provider.T) {
	t.Title("TC-POST-03: Неверный тип данных price")
	t.Description("price как строка")
	t.Tags("negative", "post")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"sellerID":111289,"name":"test","price":"1sb","statistics":{"likes":10,"viewCount":10,"contacts":10}}`)
	t.WithNewStep("Отправка запроса с price=string", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_04_ZeroPrice(t provider.T) {
	t.Title("TC-POST-04: Нулевая цена")
	t.Description("price=0 ")
	t.Tags("negative", "post")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"sellerID":111289,"name":"test","price":0,"statistics":{"likes":10,"viewCount":10,"contacts":10}}`)
	t.WithNewStep("Отправка запроса с price=0", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_05_PriceOverflow(t provider.T) {
	t.Title("TC-POST-05: Переполнение price")
	t.Description("price больше int64 max ")
	t.Tags("negative", "post")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"sellerID":111289,"name":"test","price":9223372036854775808,"statistics":{"likes":10,"viewCount":10,"contacts":10}}`)
	t.WithNewStep("Отправка запроса с переполненным price", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_06_EmptyName(t provider.T) {
	t.Title("TC-POST-06: Пустое имя")
	t.Description("name=\"\" ")
	t.Tags("negative", "post")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"sellerID":111289,"name":"","price":100,"statistics":{"likes":10,"viewCount":10,"contacts":10}}`)
	t.WithNewStep("Отправка запроса с пустым name", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_07_ZeroStatistics(t provider.T) {
	t.Title("TC-POST-07: Нулевая статистика")
	t.Description("statistics с нулями")
	t.Tags("positive", "post")
	t.Severity(allure.NORMAL)

	var req models.CreateAdvertisementRequest
	t.WithNewStep("Подготовка запроса с нулевой статистикой", func(sCtx provider.StepCtx) {
		req = models.CreateAdvertisementRequest{
			SellerID:  111299,
			Name:      "test",
			Price:     100,
			Statistic: models.Statistic{Likes: 0, ViewCount: 0, Contacts: 0},
		}
		b, _ := json.MarshalIndent(req, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, b))
	})

	var ad models.Advertisement
	var code int
	t.WithNewStep("Отправка post", func(sCtx provider.StepCtx) {
		ad, code = manager.CreateAdvertisement(t, req)
		b, _ := json.MarshalIndent(ad, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Response body", allure.JSON, b))
		sCtx.Require().Equal(200, code)
		sCtx.Require().NotEmpty(ad.ID)
	})

	t.WithNewStep("Проверка статистики в ответе", func(sCtx provider.StepCtx) {
		sCtx.Assert().Equal(int64(0), ad.Statistic.Likes)
		sCtx.Assert().Equal(int64(0), ad.Statistic.ViewCount)
		sCtx.Assert().Equal(int64(0), ad.Statistic.Contacts)
	})

	s.createdIDs = append(s.createdIDs, ad.ID)
}

func (s *CreateAdvertisementSuite) TestTC_POST_08_NegativePrice(t provider.T) {
	t.Title("TC-POST-08: Отрицательный price")
	t.Description("price=-1")
	t.Tags("negative", "post")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"sellerID":111289,"name":"test","price":-1,"statistics":{"likes":10,"viewCount":10,"contacts":10}}`)
	t.WithNewStep("Отправка запроса с price=-1", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_09_NegativeStatistics(t provider.T) {
	t.Title("TC-POST-09: Отрицательные значения статистики")
	t.Description("likes/viewCount/contacts < 0 ")
	t.Tags("negative", "post")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"sellerID":111289,"name":"test","price":100,"statistics":{"likes":-1,"viewCount":-1,"contacts":-1}}`)
	t.WithNewStep("Отправка запроса с отрицательной статистикой", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_10_StatisticsOverflow(t provider.T) {
	t.Title("TC-POST-10: Переполнение значений статистики")
	t.Description("likes/viewCount/contacts > int64 max")
	t.Tags("negative", "post", "validation", "boundary")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"sellerID":111289,"name":"test","price":100,"statistics":{"likes":9223372036854775808,"viewCount":9223372036854775808,"contacts":9223372036854775808}}`)
	t.WithNewStep("Отправка запроса с переполненной статистикой", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_11_MissingName(t provider.T) {
	t.Title("TC-POST-11: Отсутствие обязательного поля name")
	t.Description("Запрос без name ")
	t.Tags("negative", "post", "validation")
	t.Severity(allure.NORMAL)

	rawBody := []byte(`{"sellerID":111289,"price":100,"statistics":{"likes":1,"viewCount":1,"contacts":1}}`)
	t.WithNewStep("Отправка запроса без name", func(sCtx provider.StepCtx) {
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, rawBody))
		code := manager.CreateAdvertisementRaw(t, rawBody)
		sCtx.Assert().Equal(400, code)
	})
}

func (s *CreateAdvertisementSuite) TestTC_POST_12_MinPrice(t provider.T) {
	t.Title("TC-POST-12: Минимальная цена price=1")
	t.Description("price=1 ")
	t.Tags("positive", "post", "boundary")
	t.Severity(allure.NORMAL)

	var req models.CreateAdvertisementRequest
	t.WithNewStep("Подготовка запроса с price=1", func(sCtx provider.StepCtx) {
		req = models.CreateAdvertisementRequest{
			SellerID:  111299,
			Name:      "test",
			Price:     1,
			Statistic: models.Statistic{Likes: 1, ViewCount: 1, Contacts: 1},
		}
		b, _ := json.MarshalIndent(req, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, b))
	})

	var ad models.Advertisement
	var code int
	t.WithNewStep("Отправка post", func(sCtx provider.StepCtx) {
		ad, code = manager.CreateAdvertisement(t, req)
		sCtx.Require().Equal(200, code)
		sCtx.Require().NotEmpty(ad.ID)
	})

	t.WithNewStep("Проверка price в ответе", func(sCtx provider.StepCtx) {
		sCtx.Assert().Equal(int64(1), ad.Price)
	})

	s.createdIDs = append(s.createdIDs, ad.ID)
}

func (s *CreateAdvertisementSuite) TestTC_POST_13_MaxPrice(t provider.T) {
	t.Title("TC-POST-13: Максимальное значение price ")
	t.Description("price=9223372036854775807")
	t.Tags("positive", "post")
	t.Severity(allure.NORMAL)

	var req models.CreateAdvertisementRequest
	t.WithNewStep("Подготовка запроса с максимальным price", func(sCtx provider.StepCtx) {
		req = models.CreateAdvertisementRequest{
			SellerID:  111299,
			Name:      "test",
			Price:     9223372036854775807,
			Statistic: models.Statistic{Likes: 1, ViewCount: 1, Contacts: 1},
		}
		b, _ := json.MarshalIndent(req, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Request body", allure.JSON, b))
	})

	var ad models.Advertisement
	var code int
	t.WithNewStep("Отправка post", func(sCtx provider.StepCtx) {
		ad, code = manager.CreateAdvertisement(t, req)
		sCtx.Require().Equal(200, code)
		sCtx.Require().NotEmpty(ad.ID)
	})

	t.WithNewStep("Проверка price в ответе", func(sCtx provider.StepCtx) {
		sCtx.Assert().Equal(int64(9223372036854775807), ad.Price)
	})

	s.createdIDs = append(s.createdIDs, ad.ID)
}

func (s *CreateAdvertisementSuite) TestTC_NF_01_Stability(t provider.T) {
	t.Title("TC-NF-01: Стабильность POST /item (10 запросов)")
	t.Description("10 последовательных POST запросов  ни один не должен вернуть 5xx")
	t.Tags("nonfunctional", "post")
	t.Severity(allure.NORMAL)

	req := models.CreateAdvertisementRequest{
		SellerID:  111299,
		Name:      "stability-test",
		Price:     100,
		Statistic: models.Statistic{Likes: 1, ViewCount: 1, Contacts: 1},
	}

	t.WithNewStep("Выполнение 10 последовательных POST запросов", func(sCtx provider.StepCtx) {
		for i := 0; i < 10; i++ {
			ad, code := manager.CreateAdvertisement(t, req)
			sCtx.Assert().True(code == 200,
				"Итерация %d: ожидался 200 , получен %d", i+1, code)
			if code == 200 {
				s.createdIDs = append(s.createdIDs, ad.ID)
			}
		}
	})
}
