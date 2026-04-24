package integrationtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/sdinsure/agent/pkg/logger"
	storagetestutils "github.com/sdinsure/agent/pkg/storage/testutils"
	"github.com/stretchr/testify/assert"

	restcolgohttpclient "github.com/footprintai/restcol/api/go-http-client"
	restcolopenapi "github.com/footprintai/restcol/api/go-openapiv2/client"
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
	log := logger.NewLogger()
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

	// wait for everything is ready
	time.Sleep(1 * time.Second)

	return &suite{
		svr: svr,
	}
}
