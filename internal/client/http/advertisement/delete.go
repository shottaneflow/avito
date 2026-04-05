package advertisement

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/shottaneflow/avito/internal/constants/path"
	"github.com/shottaneflow/avito/internal/runner"
)

func HttpDeleteAdvertisement(t testing.TB, id string) *http.Response {
	url := runner.GetRunner().BaseURL + fmt.Sprintf(path.DeleteAdvertisement, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	require.NoError(t, err)
	return resp
}
