package wrk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/davidkarolyi/server-comparison/benchmark/wrk/types"
)

// Client is an HTTP cient for the wrk server.
type Client struct {
	wrkHostURL      string
	benchmarkParams *types.BenchmarkParams
}

// NewClient  creates a new wrk server client.
func NewClient(wrkHostURL string, benchmarkParams *types.BenchmarkParams) *Client {
	return &Client{
		wrkHostURL:      formatURL(wrkHostURL),
		benchmarkParams: benchmarkParams,
	}
}

// CheckConnection checks if the wrk server is available.
func (client *Client) CheckConnection() error {
	resp, err := http.Get(client.wrkHostURL + "/check")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"Wrk server responded with %s, instead of %d %s",
			resp.Status,
			http.StatusOK,
			http.StatusText(http.StatusOK),
		)
	}
	return nil
}

// RunBenchmark requests a benchmark,
// and returns with the result once it's done.
func (client *Client) RunBenchmark() (*types.WRKResult, error) {
	paramsJSON, err := json.Marshal(client.benchmarkParams)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		client.wrkHostURL+"/benchmark",
		"application/json",
		bytes.NewBuffer(paramsJSON),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return extractResultFromResponse(resp)
}

func extractResultFromResponse(resp *http.Response) (*types.WRKResult, error) {
	resultBuffer := new(bytes.Buffer)
	io.Copy(resultBuffer, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wrk server responded with error: '%s'", resultBuffer.String())
	}

	result := &types.WRKResult{}
	err := json.Unmarshal(resultBuffer.Bytes(), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func formatURL(url string) string {
	if strings.HasSuffix(url, "/") {
		return strings.TrimSuffix(url, "/")
	}
	return url
}
