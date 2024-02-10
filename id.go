package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

var _ Identifier = (*identifier)(nil)

type Identifier interface {
	driver.Valuer
	schema.GormDataTypeInterface
	migrator.GormDataTypeInterface
	Object() string
	Uid() int64
	Timestamp() time.Time
	String() string
}

type identifier struct {
	object    string
	uid       int64
	timestamp time.Time
}

func NewIdentifier(object string, uid uint64, timestamp time.Time) Identifier {
	return &identifier{
		object:    object,
		uid:       int64(uid),
		timestamp: timestamp,
	}
}

func (id identifier) Object() string {
	return id.object
}

func (id identifier) Uid() int64 {
	return id.uid
}

func (id identifier) Timestamp() time.Time {
	return id.timestamp
}

func (id identifier) String() string {
	return fmt.Sprintf("%s_%d", id.object, id.uid)
}

func (id identifier) Value() (driver.Value, error) {
	return id.uid, nil
}

func (id identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *identifier) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	_id, err := FromObjectId(string(data))
	if err != nil {
		return err
	}
	*id = *_id

	return nil
}

func (identifier) GormDataType() string {
	return "bigint"
}

func (identifier) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "bigint"
}

func FromObjectId(str string) (*identifier, error) {
	idParts := strings.Split(str, "_")
	if len(idParts) != 2 {
		return nil, errors.New("invalid str")
	}

	object := idParts[0]
	uid, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return FromUid(object, uid), nil
}

func FromUid(object string, uid int64) *identifier {
	timestamp := Timestamp(uid)
	return &identifier{
		uid:       uid,
		object:    object,
		timestamp: timestamp,
	}
}
