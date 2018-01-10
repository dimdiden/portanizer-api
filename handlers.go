package main

import (
	"encoding/json"
	// "fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{"Status": "OK"}
	ResponseWithJSON(w, &body, http.StatusOK)
}

func (app *App) GetPostList(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	app.DB.Find(&posts)

	for i, _ := range posts {
		app.DB.Model(posts[i]).Related(&posts[i].Tags, "Tags")
	}
	ResponseWithJSON(w, &posts, http.StatusOK)
}

func (app *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		ErrorWithJSON(w, "Can not decode json", http.StatusBadRequest)
		return
	}
	app.DB.Save(&post)
	ResponseWithJSON(w, &post, http.StatusOK)
}

func (app *App) GetTagList(w http.ResponseWriter, r *http.Request) {
	var tags []Tag
	app.DB.Find(&tags)

	ResponseWithJSON(w, &tags, http.StatusOK)
}

func (app *App) CreateTag(w http.ResponseWriter, r *http.Request) {
	var tag Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		ErrorWithJSON(w, "Can not decode json", http.StatusBadRequest)
		return
	}
	app.DB.Save(&tag)
	ResponseWithJSON(w, &tag, http.StatusOK)
}

// {
//         "Name": "Hui",
//         "Body": "TestTest",
//         "Tags": [
//             {
//                 "Name": "mysql"
//             },
//             {
//                 "Name": "ntp"
//             }
//         ]
// }
