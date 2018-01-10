package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"io"
)

type App struct {
	Lfile io.WriteCloser
	Conf *Conf
	R    *mux.Router
	DB   *gorm.DB
}

type Conf struct {
	Addr   string
	Driver string
	DBuser string
	DBname string
	Ptolog string `json:"logfile"`
}

func (app *App) Initiate(cfile string) {
	app.SetConf(cfile)
	app.SetDB()
	app.SetRouter()

	lfile, err := os.Create(app.Conf.Ptolog)
	if err != nil {
		log.Fatal("Cannot create the logfile:", err)
	}
	app.Lfile = lfile
}

func (app *App) SetConf(cfile string) {
	file, err := os.Open(cfile)
	if err != nil {
		log.Fatal("Error opening conf file:", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&app.Conf)
	if err != nil {
		log.Fatal("Error decoding conf file:", err)
	}
	return
}

func (app *App) SetDB() {
	cs := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True&loc=Local", app.Conf.DBuser, app.Conf.DBname)
	db, err := gorm.Open(app.Conf.Driver, cs)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	app.DB = db.Debug()
}

func (app *App) SetRouter() {
	app.R = mux.NewRouter()

	app.R.HandleFunc("/health", Health).Methods("GET")
	app.R.HandleFunc("/posts", app.GetPostList).Methods("GET")
	app.R.HandleFunc("/posts", app.CreatePost).Methods("POST")

	app.R.HandleFunc("/tags", app.GetTagList).Methods("GET")
	app.R.HandleFunc("/tags", app.CreateTag).Methods("POST")

}

func (app *App) CleanDB() {
	app.DB.DropTableIfExists(&Post{}, &Tag{}, &PostTag{})
	app.DB.AutoMigrate(&Post{}, &Tag{}, &PostTag{})
}

func (app *App) Run(lenabled bool) {
	var h http.Handler = app.R

	if lenabled {
		h = handlers.LoggingHandler(app.Lfile, h)
	}

	http.ListenAndServe(app.Conf.Addr, h)
}

func (app *App) Exit() {
	app.DB.Close()
	app.Lfile.Close()
}
