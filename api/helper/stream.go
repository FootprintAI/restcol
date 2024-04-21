package apihelper

import (
	"errors"
	"fmt"
	"time"

	openapiclientdocument "github.com/footprintai/restcol/api/go-openapiv2/client/document"
	sdinsureerrors "github.com/sdinsure/agent/pkg/errors"
	sdinsureopenapistream "github.com/sdinsure/agent/pkg/http/openapi/stream"
)

var (
	_ sdinsureopenapistream.ClientResponseSinkCloser = &DocumentStreamSinkCloser{}
)

type DocumentStreamSinkCloser struct {
	recvChan                 chan *openapiclientdocument.RestColServiceQueryDocumentsStreamOK
	readWriteTimeoutDuration time.Duration
}

func NewDocumentStreamSinkCloser() *DocumentStreamSinkCloser {
	return &DocumentStreamSinkCloser{
		recvChan:                 make(chan *openapiclientdocument.RestColServiceQueryDocumentsStreamOK, 10),
		readWriteTimeoutDuration: 10 * time.Second,
	}
}

var (
	errTimeout = sdinsureerrors.NewTimeoutError(errors.New("timeout when reading/writing to streamsinker"))
)

func (d *DocumentStreamSinkCloser) DefaultResponse() interface{} {
	return &openapiclientdocument.RestColServiceQueryDocumentsStreamOK{}
}

func (d *DocumentStreamSinkCloser) SinkResponse(v interface{}) error {
	fmt.Printf("sinkresponse is called:+%v\n", v)
	_, isTypeMatched := v.(*openapiclientdocument.RestColServiceQueryDocumentsStreamOK)
	if !isTypeMatched {
		return errors.New("type are not matched")
	}
	select {
	case d.recvChan <- v.(*openapiclientdocument.RestColServiceQueryDocumentsStreamOK):
		return nil
	case <-time.After(d.readWriteTimeoutDuration):
		return errTimeout
	}
}

// return nil, nil when no more data
func (d *DocumentStreamSinkCloser) Recv() (*openapiclientdocument.RestColServiceQueryDocumentsStreamOK, error) {
	fmt.Printf("recv is called\n")
	select {
	case ret, more := <-d.recvChan:
		if !more {
			return nil, nil
		}
		return ret, nil
	case <-time.After(d.readWriteTimeoutDuration):
		return nil, errTimeout
	}
}

var (
	ErrEOF = errors.New("stream: eof")
)

func WithError(in *openapiclientdocument.RestColServiceQueryDocumentsStreamOK, e error) (*openapiclientdocument.RestColServiceQueryDocumentsStreamOK, error) {
	if in == nil && e == nil {
		return nil, ErrEOF
	}
	return in, e

}

func (d *DocumentStreamSinkCloser) Close() error {
	close(d.recvChan)
	return nil

}
