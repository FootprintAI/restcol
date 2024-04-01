package storage

import (
	sderrors "github.com/sdinsure/agent/pkg/errors"
	storageerrors "github.com/sdinsure/agent/pkg/storage/errors"
)

func wrap(e *sderrors.Error) error {
	if e == (*sderrors.Error)(nil) {
		return nil
	}
	return e
}

func WrapStorageError(e error) error {
	return wrap(
		storageerrors.WrapStorageError(e),
	)
}
