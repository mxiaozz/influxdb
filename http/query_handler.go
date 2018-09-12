package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/influxdata/platform"
	pcontext "github.com/influxdata/platform/context"
	"github.com/influxdata/platform/kit/errors"
	"github.com/influxdata/platform/query"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	fluxPath = "/v2/query"
)

// FluxHandler implements handling flux queries.
type FluxHandler struct {
	*httprouter.Router

	Logger *zap.Logger

	AuthorizationService platform.AuthorizationService
	OrganizationService  platform.OrganizationService
	ProxyQueryService    query.ProxyQueryService
}

// NewFluxHandler returns a new handler at /v2/query for flux queries.
func NewFluxHandler() *FluxHandler {
	h := &FluxHandler{
		Router: httprouter.New(),
		Logger: zap.NewNop(),
	}

	h.HandlerFunc("POST", "/v2/query", h.handlePostQuery)
	return h
}

func (h *FluxHandler) handlePostQuery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tok, err := pcontext.GetToken(ctx)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	auth, err := h.AuthorizationService.FindAuthorizationByToken(ctx, tok)
	if err != nil {
		EncodeError(ctx, errors.Wrap(err, "invalid token", errors.InvalidData), w)
		return
	}

	if !platform.IsActive(auth) {
		EncodeError(ctx, errors.Forbiddenf("insufficient permissions for write"), w)
		return
	}

	req, err := decodeProxyQueryRequest(ctx, r, h.OrganizationService)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	hd, ok := req.Dialect.(HTTPDialect)
	if !ok {
		EncodeError(ctx, fmt.Errorf("unsupported dialect over HTTP %T", req.Dialect), w)
		return
	}
	hd.SetHeaders(w)

	n, err := h.ProxyQueryService.Query(ctx, w, req)
	if err != nil {
		if n == 0 {
			// Only record the error headers IFF nothing has been written to w.
			EncodeError(ctx, err, w)
			return
		}
		h.Logger.Info("Error writing response to client",
			zap.String("handler", "flux"),
			zap.Error(err),
		)
	}
}

// PrometheusCollectors satisifies the prom.PrometheusCollector interface.
func (h *FluxHandler) PrometheusCollectors() []prometheus.Collector {
	// TODO: gather and return relevant metrics.
	return nil
}

// FluxService connects to Influx via HTTP using tokens to run queries.
type FluxService struct {
	URL                string
	Token              string
	InsecureSkipVerify bool
}

// Query runs a flux query against a influx server and sends the results to the io.Writer.
func (s *FluxService) Query(ctx context.Context, w io.Writer, req *query.ProxyRequest) (int64, error) {
	u, err := newURL(s.URL, fluxPath)
	if err != nil {
		return 0, err
	}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return 0, err
	}

	hreq, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return 0, err
	}
	SetToken(s.Token, hreq)
	hreq.Header.Set("Content-Type", "application/json")
	hreq = hreq.WithContext(ctx)

	hc := newClient(u.Scheme, s.InsecureSkipVerify)
	resp, err := hc.Do(hreq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if err := CheckError(resp); err != nil {
		return 0, err
	}
	return io.Copy(w, resp.Body)
}
