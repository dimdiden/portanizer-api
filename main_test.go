package main

import (
	"os"
	"testing"
)

const TESTCONFFILE = "./conf_test.json"

var app *App

func TestMain(m *testing.M) {

	app = Initiate(TESTCONFFILE, true)

	os.Exit(m.Run())
}
