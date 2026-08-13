package integrationtest

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sdinsure/agent/pkg/logger"
	storagetestutils "github.com/sdinsure/agent/pkg/storage/testutils"
	"github.com/stretchr/testify/assert"

	restcolgohttpclient "github.com/footprintai/restcol/api/go-http-client"
	restcolopenapi "github.com/footprintai/restcol/api/go-openapiv2/client"
	bootstrap "github.com/footprintai/restcol/pkg/bootstrap"
	serverapp "github.com/footprintai/restcol/pkg/server/app"
)

type SuiteCloser interface {
	Close() error
}

type suite struct {
	svr *serverapp.Server
}

func (s *suite) Close() error {
	return s.svr.Stop()
}

func (s *suite) NewClient() *restcolopenapi.RestColAPIDocumentations {
	return restcolgohttpclient.MustNewClient("localhost:50051", nil)
}

func SetupTest(t *testing.T) *suite {
	log := logger.NewLogger(false)
	postgresDb, err := storagetestutils.NewTestPostgresCli(log)
	if err != nil {
		assert.NoError(t, err)
	}
	svr, err := serverapp.NewServer(50050, 50051, postgresDb, log)
	if err != nil {
		log.Fatal("%+v", err)
	}

	fmt.Print("integrationtest about to start\n")
	go svr.Start()

	waitGatewayReady(t)

	return &suite{
		svr: svr,
	}
}

// waitGatewayReady blocks until the HTTP gateway can actually reach the gRPC
// server behind it.
//
// A fixed sleep is not enough, and not because the server is slow: the gateway
// registers its client connection before the gRPC listener binds, so that
// connection starts in failure and only recovers on gRPC's own reconnect
// backoff. The listener is up in ~100ms but the gateway keeps answering 503
// (code 14, "connection refused") for a second or more afterwards, which is
// exactly where the old one-second sleep landed. Poll the real thing instead.
func waitGatewayReady(t *testing.T) {
	t.Helper()

	const timeout = 30 * time.Second
	probe := "http://localhost:50051/v1/projects/" + bootstrap.DefaultModelProject.ID.String() + "/collections"

	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		resp, err := http.Get(probe)
		if err != nil {
			last = err.Error()
			time.Sleep(50 * time.Millisecond)
			continue
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		resp.Body.Close()

		// Any answer that came from the application means the gateway reached
		// gRPC. Only an unavailable-transport reply means it did not.
		if !(resp.StatusCode == http.StatusServiceUnavailable &&
			(strings.Contains(string(body), `"code":14`) || strings.Contains(string(body), "transport:"))) {
			return
		}
		last = fmt.Sprintf("%d %s", resp.StatusCode, body)
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("gateway never reached the grpc server within %s; last response: %s", timeout, last)
}
