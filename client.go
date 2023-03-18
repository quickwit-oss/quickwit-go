package quickwit

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

// Error implements go error interface.
func (msg *ErrorMessage) Error() string {
	return fmt.Sprintf("API Error: %s", msg.Message)
}

type QuickwitClient struct {
	*req.Client
}

// endpoint is the root url of the quickwit server.
func NewQuickwitClient(endpoint string) *QuickwitClient {
	return &QuickwitClient{
		Client: req.C().
			SetBaseURL(endpoint).
			SetCommonErrorResult(&ErrorMessage{}).
			EnableDumpEachRequest().
			OnAfterResponse(func(client *req.Client, resp *req.Response) error {
				if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
					return nil
				}
				if errMsg, ok := resp.ErrorResult().(*ErrorMessage); ok {
					resp.Err = errMsg // Convert api error into go error
					return nil
				}
				if !resp.IsSuccessState() {
					// Neither a success response nor a error response, record details to help troubleshooting
					resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
				}
				return nil
			}),
	}
}

type SearchRequest struct {
	Query          string   `json:"query"`
	SearchFields   []string `json:"search_field,omitempty"`
	StartTimestamp int64    `json:"start_timestamp,omitempty"`
	EndTimestamp   int64    `json:"end_timestamp,omitempty"`
	MaxHits        uint64   `json:"max_hits"`
	StartOffset    uint64   `json:"start_offset,omitempty"`
	SortByField    string   `json:"sort_by_field,omitempty"`
}

type SearchResponse struct {
	NumHits           uint64        `json:"num_hits"`
	Hits              []interface{} `json:"hits"`
	ElapsedTimeMicros uint64        `json:"elapsed_time_micros"`
	Aggregations      interface{}   `json:"aggregations,omitempty"`
}

func (c *QuickwitClient) Search(indexId string, searchRequest SearchRequest) (searchResponse *SearchResponse, err error) {
	_, err = c.R().
		SetPathParam("indexId", indexId).
		SetBody(searchRequest).
		SetSuccessResult(&searchResponse).
		Post("/api/v1/{indexId}/search")
	return
}
