package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/gormigrate.v1"
)

type Post struct {
	ID   uint
	Name string `gorm:"unique"`
	Body string
	Tags []Tag `gorm:"many2many:post_tags;"`
}

type Tag struct {
	ID   uint
	Name string `gorm:"unique"`
}

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				return tx.CreateTable(&Post{}, &Tag{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable(&Post{}, &Tag{}).Error
			},
		},
	})
	return m.Migrate()
}
