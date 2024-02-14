package main

import (
	"time"

	"github.com/sony/sonyflake"
)

var startTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

type IdGenerator interface {
	Next(kind string) (Identifier, error)
}

type idGenerator struct {
	*sonyflake.Sonyflake
}

func NewIdGenerator() IdGenerator {
	return &idGenerator{
		sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: startTime,
		}),
	}
}

func (i *idGenerator) Next(kind string) (Identifier, error) {
	uid, err := i.NextID()
	if err != nil {
		return nil, err
	}

	return NewIdentifier(
		kind,
		uid,
		timestamp(uid),
	), nil
}

func Timestamp(uid int64) time.Time {
	return timestamp(uint64(uid))
}

func timestamp(uid uint64) time.Time {
	elapsedTime := sonyflake.ElapsedTime(uid)
	timestamp := startTime.Add(elapsedTime)

	return timestamp
}
