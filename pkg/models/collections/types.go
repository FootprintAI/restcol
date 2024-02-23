package modelcollections

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type CollectionID uuid.UUID

var (
	_ fmt.Stringer = CollectionID{}
)

func (c CollectionID) String() string {
	// TODO maybe add some prefix like "c:"
	return uuid.UUID(c).String()
}

var (
	_ driver.Valuer = CollectionID{}
	_ sql.Scanner   = &CollectionID{}
)

func (c CollectionID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *CollectionID) Scan(value interface{}) error {
	_, isStringType := value.(string)
	if !isStringType {
		return fmt.Errorf("db.model: invalid type, expect string")
	}
	cid, err := NewCollectionIDFromStr(value.(string))
	if err != nil {
		return err
	}
	(*c) = cid
	return nil
}

func NewCollectionID() CollectionID {
	uid, _ := uuid.NewV7()
	return CollectionID(uid)
}

func NewCollectionIDFromStr(s string) (CollectionID, error) {
	innerUuid, err := uuid.Parse(s)
	if err != nil {
		return CollectionID{}, err
	}
	return CollectionID(innerUuid), nil
}
