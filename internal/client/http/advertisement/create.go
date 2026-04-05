package advertisement

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/shottaneflow/avito/internal/constants/path"
	"github.com/shottaneflow/avito/internal/managers/advertisement/models"
	"github.com/shottaneflow/avito/internal/runner"
)

func HttpPostCreateAdvertisement(t testing.TB, request models.CreateAdvertisementRequest) *http.Response {
	body, err := json.Marshal(request)
	require.NoError(t, err)
	url := runner.GetRunner().BaseURL + path.CreateAdvertisement
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	return resp
}

func HttpPostCreateAdvertisementRaw(t testing.TB, rawBody []byte) *http.Response {
	url := runner.GetRunner().BaseURL + path.CreateAdvertisement
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(rawBody))
	require.NoError(t, err)
	return resp
}
