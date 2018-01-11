package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"os"
)

type App struct {
	Lout io.WriteCloser
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

func Initiate(cfile string, lenabled, clear bool) *App {
	app := App{}

	app.SetConf(cfile)
	app.SetDB(clear)
	app.SetRouter()

	app.Lout = os.Stdout

	if lenabled {
		lfile, err := os.Create(app.Conf.Ptolog)
		if err != nil {
			log.Fatal("Cannot create the logfile:", err)
		}
		app.Lout = lfile
	}
	return &app
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

func (app *App) SetDB(clear bool) {
	cs := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True&loc=Local", app.Conf.DBuser, app.Conf.DBname)
	db, err := gorm.Open(app.Conf.Driver, cs)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	app.DB = db.Debug()

	if clear {
		app.DB.DropTableIfExists(&Post{}, &Tag{}, &PostTag{})
	}

	app.DB.AutoMigrate(&Post{}, &Tag{}, &PostTag{})
}

func (app *App) SetRouter() {
	app.R = mux.NewRouter()
	// app.R.StrictSlash(true)

	app.R.HandleFunc("/health", app.CheckHealth).Methods("GET")

	app.R.HandleFunc("/posts", app.GetPostList).Methods("GET")
	app.R.HandleFunc("/posts", app.CreatePost).Methods("POST")
	app.R.HandleFunc("/posts/{id}", app.DeletePost).Methods("DELETE")

	app.R.HandleFunc("/tags", app.GetTagList).Methods("GET")
	app.R.HandleFunc("/tags", app.CreateTag).Methods("POST")
	app.R.HandleFunc("/tags/{id}", app.UpdateTag).Methods("PATCH")
	app.R.HandleFunc("/tags/{id}", app.DeleteTag).Methods("DELETE")

}

func (app *App) clearDB() {
	app.DB.DropTableIfExists(&Post{}, &Tag{}, &PostTag{})
}

func (app *App) Run() {
	var h http.Handler = app.R
	h = handlers.LoggingHandler(app.Lout, h)

	http.ListenAndServe(app.Conf.Addr, h)
}

func (app *App) Exit() {
	app.DB.Close()
	app.Lout.Close()
}
