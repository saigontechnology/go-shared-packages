package newrelic

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
	nragent "github.com/newrelic/go-agent/v3/newrelic"
)

func HTTPClientWithNewRelic() *http.Client {
	client := &http.Client{}
	client.Transport = nragent.NewRoundTripper(client.Transport)

	return client
}

func RestyClientWithNewRelic() *resty.Client {
	client := &http.Client{}
	client.Transport = nragent.NewRoundTripper(client.Transport)

	return resty.NewWithClient(client)
}

func HTTPRequestWithNewRelic(ctx context.Context, request *http.Request) *http.Request {
	txn := nragent.FromContext(ctx)
	request = nragent.RequestWithTransactionContext(request, txn)

	return request
}

//nolint:contextcheck
func RestyRequestWithNewRelic(ctx context.Context, request *resty.Request) *resty.Request {
	txn := nragent.FromContext(ctx)
	request = RestyRequestWithTransactionContext(request, txn)

	return request
}

func RestyRequestWithTransactionContext(req *resty.Request, txn *nragent.Transaction) *resty.Request {
	ctx := req.Context()
	ctx = nragent.NewContext(ctx, txn)
	return req.SetContext(ctx)
}
