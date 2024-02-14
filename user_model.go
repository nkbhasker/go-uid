package main

type User struct {
	ID    Identifier `gorm:"primaryKey;type:bigint;serializer:id;" kind:"user" json:"id"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
}
