package misskey

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MisskeyHttpClient struct {
	baseUrl     string
	timeout     time.Duration
	accessToken string
	client      *http.Client
}

func NewHttpClient(baseUrl string, timeoutSeconds int, accessToken string) *MisskeyHttpClient {
	timeout := time.Duration(timeoutSeconds) * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	return &MisskeyHttpClient{
		baseUrl:     baseUrl,
		timeout:     timeout,
		accessToken: accessToken,
		client:      client,
	}
}

func (c *MisskeyHttpClient) Do(method string, endpoint string, reqBody any) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, 0, err
	}

	fullUrl := fmt.Sprintf("%s/%s", c.baseUrl, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, fullUrl, bytes.NewBuffer(reqBodyData))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}

func (c *MisskeyHttpClient) Post(endpoint string, reqBody any) ([]byte, int, error) {
	return c.Do("POST", endpoint, reqBody)
}
