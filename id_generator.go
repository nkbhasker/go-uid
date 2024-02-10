package main

import (
	"time"

	"github.com/sony/sonyflake"
)

var startTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

type IdGenerator interface {
	Next(object string) (Identifier, error)
}

type idGenerator struct {
	flake *sonyflake.Sonyflake
}

func NewIdGenerator() IdGenerator {
	return &idGenerator{
		flake: sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: startTime,
		}),
	}
}

func (i *idGenerator) Next(object string) (Identifier, error) {
	uid, err := i.flake.NextID()
	if err != nil {
		return nil, err
	}

	return NewIdentifier(
		object,
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
