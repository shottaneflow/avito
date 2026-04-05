package advertisement

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	client "github.com/shottaneflow/avito/internal/client/http/advertisement"
	"github.com/shottaneflow/avito/internal/managers/advertisement/models"
)

func CreateAdvertisement(t testing.TB, req models.CreateAdvertisementRequest) (models.Advertisement, int) {
	resp := client.HttpPostCreateAdvertisement(t, req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	code := resp.StatusCode
	if code != 200 {
		return models.Advertisement{}, code
	}

	// BUG-01: сервер возвращает {"status": "Сохранили объявление - <uuid>"}
	// вместо полного объекта. Парсим id из строки статуса.
	var statusResp struct {
		Status string `json:"status"`
	}
	err = json.Unmarshal(body, &statusResp)
	require.NoError(t, err)

	parts := strings.Split(statusResp.Status, " - ")
	if len(parts) < 2 {
		t.Logf("Неожиданный формат ответа: %s", statusResp.Status)
		return models.Advertisement{}, code
	}
	id := strings.TrimSpace(parts[len(parts)-1])

	ad := models.Advertisement{ID: id}
	return ad, code
}

func CreateAdvertisementRaw(t testing.TB, rawBody []byte) int {
	resp := client.HttpPostCreateAdvertisementRaw(t, rawBody)
	defer resp.Body.Close()
	return resp.StatusCode
}

func GetAdvertisement(t testing.TB, id string) (models.Advertisement, int) {
	resp := client.HttpGetAdvertisement(t, id)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	code := resp.StatusCode
	if code != 200 {
		return models.Advertisement{}, code
	}

	var result []models.Advertisement
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	return result[0], code
}

func GetStatistic(t testing.TB, id string) (models.StatisticResponse, int) {
	resp := client.HttpGetStatistic(t, id)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	code := resp.StatusCode
	if code != 200 {
		return models.StatisticResponse{}, code
	}

	var result []models.StatisticResponse
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	return result[0], code
}

func GetSellerAdvertisements(t testing.TB, sellerID int64) ([]models.Advertisement, int) {
	resp := client.HttpGetSellerAdvertisements(t, sellerID)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	code := resp.StatusCode
	if code != 200 {
		return nil, code
	}

	var result []models.Advertisement
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	return result, code
}

func DeleteAdvertisement(t testing.TB, id string) int {
	resp := client.HttpDeleteAdvertisement(t, id)
	defer resp.Body.Close()
	return resp.StatusCode
}
