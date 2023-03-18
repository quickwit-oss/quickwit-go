package quickwit

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestClient_Search(t *testing.T) {
	qClient := NewQuickwitClient("http://localhost:7280")
	httpmock.ActivateNonDefault(qClient.GetClient())
	httpmock.RegisterResponder("POST", "http://localhost:7280/api/v1/otel-logs-v0/search", func(request *http.Request) (*http.Response, error) {
		respBody := `{"hits": [], "num_hits": 0, "elapsed_time_micros": 100}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})
	searchResponse, err := qClient.Search("otel-logs-v0", SearchRequest{Query: "test"})
	if err != nil {
		t.Error(err)
	}
	require.NoError(t, err)
	require.Equal(t, uint64(0), searchResponse.NumHits)
	require.Equal(t, 0, len(searchResponse.Hits))
}
