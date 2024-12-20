package models

import (
	"time"
)

type Post struct {
	ID          uint      `gorm:"primary_key"`
	AuthorLogin string    `gorm:"type:varchar(128);not null"`
	Title       string    `gorm:"type:varchar(256);not null"`
	Text        string    `gorm:"type:text;not null"`
	Timestamp   time.Time `gorm:"type:time;not null"`
	Comments    []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	ID          uint      `gorm:"primary_key"`
	AuthorLogin string    `gorm:"type:varchar(128);not null"`
	Text        string    `gorm:"type:text;not null"`
	Timestamp   time.Time `gorm:"type:time;not null"`
	PostID      uint      `gorm:"not null"`
}
