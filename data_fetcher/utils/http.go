package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type IHttpClient interface {
	Get(url string, responseBody any, opts ...*HttpOptions) (*HttpResponse, error)
	Post(url string, body any, responseBody any, opts ...*HttpOptions) (*HttpResponse, error)
}

type HttpClient struct{}

type HttpResponse struct {
	StatusCode         int
	Body               *[]byte
	ContentType        *string
	ContentDisposition *string
	Length             *int64
}

type HttpOptions struct {
	Headers *map[string]string
}

func (h *HttpClient) Get(url string, responseBody any, opts ...*HttpOptions) (*HttpResponse, error) {
	options := HttpOptions{}
	if len(opts) > 0 {
		options = *opts[0]
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if options.Headers != nil {
		for key, value := range *options.Headers {
			req.Header.Set(key, value)
		}
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	r := HttpResponse{
		StatusCode: response.StatusCode,
	}

	if (response.ContentLength == -1 || response.ContentLength > 0) && responseBody != nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		if response.StatusCode > 299 {
			zap.L().Warn("Got error code with response from http request", zap.String("response", string(bodyBytes)))
			return &HttpResponse{
				StatusCode: response.StatusCode,
			}, nil
		}

		r.Body = &bodyBytes

		contentType := response.Header.Get("Content-Type")
		if len(contentType) > 0 {
			r.ContentType = &contentType
		}

		contentDisposition := response.Header.Get("Content-Disposition")
		if len(contentDisposition) > 0 {
			r.ContentDisposition = &contentDisposition
		}

		r.Length = &response.ContentLength

		switch {
		case strings.Contains(contentType, "application/json"):
			if err := json.Unmarshal(bodyBytes, responseBody); err != nil {
				return nil, err
			}
			break
		case strings.Contains(contentType, "plain/text"):
			s := string(bodyBytes)
			responseBody = &s
			break
		default:
			zap.L().Warn("Unsupported response type, not decoding", zap.String("type", contentType))
			break
		}
	}

	return &r, nil
}

func (h *HttpClient) Post(url string, body any, responseBody any, opts ...*HttpOptions) (*HttpResponse, error) {
	return h.performRequestWithBody("POST", url, body, responseBody, opts...)
}

// private

func (h *HttpClient) performRequestWithBody(verb string, url string, body any, responseBody any, opts ...*HttpOptions) (*HttpResponse, error) {
	options := HttpOptions{}
	if len(opts) > 0 {
		options = *opts[0]
	}

	var reader io.Reader
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reader = bytes.NewBuffer(jsonBytes)

	req, err := http.NewRequest(verb, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if options.Headers != nil {
		for key, value := range *options.Headers {
			req.Header.Set(key, value)
		}
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	r := HttpResponse{
		StatusCode: response.StatusCode,
	}

	if (response.ContentLength == -1 || response.ContentLength > 0) && responseBody != nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		if response.StatusCode > 299 {
			zap.L().Warn("Got error code with response from http request", zap.String("response", string(bodyBytes)))
			return &HttpResponse{
				StatusCode: response.StatusCode,
			}, nil
		}

		r.Body = &bodyBytes

		contentType := response.Header.Get("Content-Type")
		if len(contentType) > 0 {
			r.ContentType = &contentType
		}

		contentDisposition := response.Header.Get("Content-Disposition")
		if len(contentDisposition) > 0 {
			r.ContentDisposition = &contentDisposition
		}

		r.Length = &response.ContentLength

		switch {
		case strings.Contains(contentType, "application/json"):
			if err := json.Unmarshal(bodyBytes, responseBody); err != nil {
				return nil, err
			}
			break
		case strings.Contains(contentType, "plain/text"):
			s := string(bodyBytes)
			responseBody = &s
			break
		default:
			zap.L().Warn("Unsupported response type, not decoding", zap.String("type", contentType))
			break
		}
	}

	return &r, nil
}
