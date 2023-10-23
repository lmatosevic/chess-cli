package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"github.com/r3labs/sse/v2"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpResponse[T any] struct {
	Data       T
	Error      model.ErrorResponse
	StatusCode int
}

type HttpClient struct {
	BaseUrl     string
	AccessToken string
	Client      *http.Client
}

var httpClient *HttpClient

func Init(baseUrl string) {
	if httpClient != nil {
		return
	}

	if !strings.HasPrefix(baseUrl, "http") {
		baseUrl = fmt.Sprintf("http://%s", baseUrl)
	}

	httpClient = &HttpClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		BaseUrl:     baseUrl,
		AccessToken: "",
	}
}

func SetAccessToken(token string) {
	httpClient.AccessToken = token
}

func GetAccessToken() string {
	return httpClient.AccessToken
}

func GetBaseUrl() string {
	return httpClient.BaseUrl
}

func SendRequest[T any](method string, path string, params *map[string]string, body any) (*HttpResponse[T], error) {
	if httpClient == nil {
		return nil, errors.New("HTTP client is not initialized")
	}

	var bodyReader io.Reader
	if body != nil {
		json, err := utils.ConvertJson(body)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(json)
	}

	sep := "/"
	if strings.HasPrefix(path, "/") || strings.HasSuffix(httpClient.BaseUrl, "/") {
		sep = ""
	}

	query := ""
	if params != nil {
		var entries []string
		for key, value := range *params {
			entries = append(entries, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
		}
		query = fmt.Sprintf("?%s", strings.Join(entries, "&"))
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s%s%s", httpClient.BaseUrl, sep, path, query), bodyReader)
	if err != nil {
		return nil, err
	}

	if httpClient.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", httpClient.AccessToken))
	}

	resp, err := httpClient.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var respModel T
	var errorModel model.ErrorResponse
	if resp.StatusCode >= 400 {
		errorModel, err = utils.ParseJson[model.ErrorResponse](resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		respModel, err = utils.ParseJson[T](resp.Body)
	}

	return &HttpResponse[T]{StatusCode: resp.StatusCode, Data: respModel, Error: errorModel}, err
}

func SubscribeOnEvent(eventType string, gameId int64, ctx context.Context, end func(),
	onEvent func(event *model.Event, end func())) error {
	client := sse.NewClient(fmt.Sprintf("%s/v1/events/subscribe?token=%s&event=%s&gameId=%d",
		httpClient.BaseUrl, httpClient.AccessToken, eventType, gameId))

	err := client.SubscribeWithContext(ctx, "message", func(msg *sse.Event) {
		event, err := utils.ParseJson[model.Event](bytes.NewReader(msg.Data))
		if err != nil {
			return
		}
		onEvent(&event, end)
	})

	return err
}
