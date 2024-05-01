package modelcollections

import (
	"github.com/google/uuid"
)

type CollectionID string

func NewCollectionIDFromStr(s string) CollectionID {
	return CollectionID(s)
}

func (c CollectionID) String() string {
	// TODO maybe add some prefix like "c:"
	return string(c)
}

func NewCollectionID() CollectionID {
	uid, _ := uuid.NewV7()
	return CollectionID(uid.String())
}
