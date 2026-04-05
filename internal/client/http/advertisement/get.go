package advertisement

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/shottaneflow/avito/internal/constants/path"
	"github.com/shottaneflow/avito/internal/runner"
)

func HttpGetAdvertisement(t testing.TB, id string) *http.Response {
	url := runner.GetRunner().BaseURL + fmt.Sprintf(path.GetAdvertisement, id)
	resp, err := http.Get(url)
	require.NoError(t, err)
	return resp
}

func HttpGetStatistic(t testing.TB, id string) *http.Response {
	url := runner.GetRunner().BaseURL + fmt.Sprintf(path.GetStatistic, id)
	resp, err := http.Get(url)
	require.NoError(t, err)
	return resp
}

func HttpGetSellerAdvertisements(t testing.TB, sellerID int64) *http.Response {
	url := runner.GetRunner().BaseURL + fmt.Sprintf(path.GetSellerAdvertisements, sellerID)
	resp, err := http.Get(url)
	require.NoError(t, err)
	return resp
}

func HttpGetSellerAdvertisementsRaw(t testing.TB, sellerID string) *http.Response {
	url := runner.GetRunner().BaseURL + fmt.Sprintf("/api/1/%s/item", sellerID)
	resp, err := http.Get(url)
	require.NoError(t, err)
	return resp
}
