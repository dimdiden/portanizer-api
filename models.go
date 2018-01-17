package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Post struct {
	gorm.Model
	Name string
	Body string
	Tags []Tag `gorm:"many2many:post_tags;"`
}

type Tag struct {
	gorm.Model
	Name string
}

type PostTag struct {
	PostId int
	TagId  int
}

func (p Post) String() string {
	return "posts"
}

func (p Tag) String() string {
	return "tags"
}
