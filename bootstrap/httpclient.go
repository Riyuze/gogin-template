package bootstrap

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type WithHttpContext func(r *http.Request)

const (
	HTTP_METHOD_GET    string = "GET"
	HTTP_METHOD_POST   string = "POST"
	HTTP_METHOD_PUT    string = "PUT"
	HTTP_METHOD_DELETE string = "DELETE"
)

type HttpClient struct {
	Client *http.Client
	trace  trace.Tracer
	cfg    *Container
}

func NewHttpClient(cfg *Container, trace trace.Tracer) *HttpClient {
	propagate := otelhttp.WithPropagators(b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader | b3.B3SingleHeader)))

	timeout := viper.GetInt("client_configuration.timeout")
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport, propagate),
	}

	return &HttpClient{cfg: cfg, Client: client, trace: trace}
}

func (client *HttpClient) RequestWithLog(ctx context.Context, method, url string, body string, opts ...WithHttpContext) (int, string, string, error) {
	client.cfg.Logger().Info("HTTPCLIENT REQUEST: ", method, " ", url, ": ", string(body))

	res, err := client.Request(ctx, method, url, bytes.NewReader([]byte(body)), opts...)
	if err != nil {
		return 0, "", "", err
	}

	resBody := ""
	resBuf := new(strings.Builder)
	if res != nil {
		_, err := io.Copy(resBuf, res.Body)
		if err != nil {
			return 0, "", "", err
		}
		resBody = resBuf.String()
	}
	client.cfg.Logger().Info("HTTPCLIENT RESPONSE: ", method, " ", url, ": ", resBody)

	return res.StatusCode, res.Status, resBody, err
}

func (client *HttpClient) Request(ctx context.Context, method, url string, body io.Reader, opts ...WithHttpContext) (*http.Response, error) {
	ctx, span := client.trace.Start(ctx, url)
	defer span.End()

	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithRedactedHeaders()))
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("http.method", req.Method),
		attribute.String("http.url", req.URL.String()),
	)

	for _, o := range opts {
		o(req)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Client.Transport = tr
	res, err := client.Client.Do(req.WithContext(ctx))
	if err != nil && res != nil {
		if res.StatusCode >= 400 {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}

	return res, err
}

func (client *HttpClient) Get(ctx context.Context, url string, body io.Reader, opts ...WithHttpContext) (*http.Response, error) {
	return client.Request(ctx, HTTP_METHOD_GET, url, body, opts...)
}

func (client *HttpClient) Post(ctx context.Context, url string, body io.Reader, opts ...WithHttpContext) (*http.Response, error) {
	return client.Request(ctx, HTTP_METHOD_POST, url, body, opts...)
}

func (client *HttpClient) Put(ctx context.Context, url string, body io.Reader, opts ...WithHttpContext) (*http.Response, error) {
	return client.Request(ctx, HTTP_METHOD_PUT, url, body, opts...)
}

func (client *HttpClient) Delete(ctx context.Context, url string, body io.Reader, opts ...WithHttpContext) (*http.Response, error) {
	return client.Request(ctx, HTTP_METHOD_DELETE, url, body, opts...)
}
