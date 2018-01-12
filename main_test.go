package main

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"testing"
)

const TESTCONFFILE = "./conf_test.json"

func TestMain(m *testing.M) {

	app := Initiate(TESTCONFFILE, true)

	ensureTablesExist(app.DB)

	os.Exit(m.Run())
}

func ensureTablesExist(db *gorm.DB) {
	var tabs []interface{}
	tabs = append(tabs, Post{}, Tag{})

	for _, t := range tabs {
		var ok bool = true
		if !ok {
			log.Fatal("Table %s doesn't exist in database")
		}
		ok = db.HasTable(t)
	}
}
