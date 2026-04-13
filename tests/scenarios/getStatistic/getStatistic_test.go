package getstatistic

import (
	"encoding/json"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	manager "github.com/shottaneflow/avito/internal/managers/advertisement"
	"github.com/shottaneflow/avito/internal/managers/advertisement/models"
)

type GetStatisticSuite struct {
	suite.Suite
	advertisementID string
	expectedStat    models.Statistic
}

func TestGetStatisticSuite(t *testing.T) {
	suite.RunSuite(t, new(GetStatisticSuite))
}

func (s *GetStatisticSuite) BeforeAll(t provider.T) {
	t.WithNewStep("создание тестового объявления со статистикой", func(sCtx provider.StepCtx) {
		s.expectedStat = models.Statistic{Likes: 5, ViewCount: 12, Contacts: 8}
		req := models.CreateAdvertisementRequest{
			SellerID:  111299,
			Name:      "stat-test",
			Price:     300,
			Statistic: s.expectedStat,
		}
		ad, code := manager.CreateAdvertisement(t, req)
		sCtx.Require().Equal(200, code, "Не удалось создать объявление для теста")
		s.advertisementID = ad.ID
		sCtx.Logf("Создано объявление с id=%s", s.advertisementID)
	})
}

func (s *GetStatisticSuite) AfterAll(t provider.T) {
	t.WithNewStep("Teardown: удаление тестового объявления", func(sCtx provider.StepCtx) {
		if s.advertisementID != "" {
			code := manager.DeleteAdvertisement(t, s.advertisementID)
			sCtx.Logf("Удалено %s, статус %d", s.advertisementID, code)
		}
	})
}

func (s *GetStatisticSuite) TestTC_STAT_01_ExistingID(t provider.T) {
	t.Title("TC-STAT-01: Получение статистики по существующему id")
	t.Description("GET /api/1/statistic/{id}")
	t.Tags("positive", "smoke", "statistic")
	t.Severity(allure.CRITICAL)

	var stat models.StatisticResponse
	var code int

	t.WithNewStep("Отправка GET /api/1/statistic/{id}", func(sCtx provider.StepCtx) {
		stat, code = manager.GetStatistic(t, s.advertisementID)
		b, _ := json.MarshalIndent(stat, "", "  ")
		sCtx.WithAttachments(allure.NewAttachment("Response body", allure.JSON, b))
		sCtx.Require().Equal(200, code)
	})

	t.WithNewStep("Проверка соответствия данным при создании", func(sCtx provider.StepCtx) {
		sCtx.Assert().Equal(s.expectedStat.Likes, stat.Likes)
		sCtx.Assert().Equal(s.expectedStat.ViewCount, stat.ViewCount)
		sCtx.Assert().Equal(s.expectedStat.Contacts, stat.Contacts)
	})
}

func (s *GetStatisticSuite) TestTC_STAT_02_NonExistingID(t provider.T) {
	t.Title("TC-STAT-02: Несуществующий id")
	t.Description("GET /api/1/statistic/{id} с несуществующим id ")
	t.Tags("negative", "statistic")
	t.Severity(allure.NORMAL)

	t.WithNewStep("Отправка запроса с несуществующим id", func(sCtx provider.StepCtx) {
		nonExistingID := "6eb9b015-b2cb-4502-8b61-c58903106074"
		manager.DeleteAdvertisement(t, nonExistingID)
		_, code := manager.GetStatistic(t, nonExistingID)
		sCtx.Assert().Equal(404, code)
	})
}
