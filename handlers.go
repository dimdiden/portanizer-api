package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *App) CheckHealth(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{"Message": "OK"}
	ResponseWithJSON(w, &body, http.StatusOK)
}

func (app *App) GetPostList(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	app.DB.Preload("Tags").Find(&posts)

	ResponseWithJSON(w, &posts, http.StatusOK)
}

func (app *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	var rsvpost Post
	if err := json.NewDecoder(r.Body).Decode(&rsvpost); err != nil {
		ErrorWithJSON(w, "Can not decode json", http.StatusBadRequest)
		return
	}

	if rsvpost.Name == "" {
		ErrorWithJSON(w, "Failed. Name field is empty", http.StatusBadRequest)
		return
	}

	var post Post

	app.DB.First(&post, "name = ?", rsvpost.Name)
	if !app.DB.NewRecord(post) {
		ErrorWithJSON(w, "Failed. This post name already exists", http.StatusBadRequest)
		return
	}

	post = Post{Name: rsvpost.Name, Body: rsvpost.Body, Tags: []Tag{}}
	app.DB.Create(&post)

	for _, t := range rsvpost.Tags {
		app.DB.First(&t, "name = ?", t.Name)
		if app.DB.NewRecord(t) {
			app.DB.Create(&t)
		}
		app.DB.Model(&post).Association("Tags").Append(t)
	}

	ResponseWithJSON(w, &post, http.StatusOK)
}

func (app *App) UpdatePost(w http.ResponseWriter, r *http.Request) {
	pid := mux.Vars(r)["id"]

	var post Post
	app.DB.First(&post, "id = ?", pid)

	if app.DB.NewRecord(post) {
		ErrorWithJSON(w, "Failed. This post doesn't exist", http.StatusBadRequest)
		return
	}

	var rsvpost Post
	if err := json.NewDecoder(r.Body).Decode(&rsvpost); err != nil {
		ErrorWithJSON(w, "Can not decode json", http.StatusBadRequest)
		return
	}

	if rsvpost.Name == "" {
		ErrorWithJSON(w, "Failed. Name field is empty", http.StatusBadRequest)
		return
	}

	var chk_post Post
	app.DB.First(&chk_post, "name = ?", rsvpost.Name)

	if !app.DB.NewRecord(chk_post) && chk_post.ID != post.ID {
		ErrorWithJSON(w, "Failed. This post name already exists", http.StatusBadRequest)
		return
	}

	post.Name = rsvpost.Name
	post.Body = rsvpost.Body
	post.Tags = []Tag{}

	for _, t := range rsvpost.Tags {
		app.DB.First(&t, "name = ?", t.Name)
		if app.DB.NewRecord(t) {
			app.DB.Create(&t)
		}
		app.DB.Model(&post).Association("Tags").Append(t)
	}

	app.DB.Save(&post)

	ResponseWithJSON(w, &post, http.StatusOK)

}

func (app *App) DeletePost(w http.ResponseWriter, r *http.Request) {
	var post Post

	vars := mux.Vars(r)
	post_id := vars["id"]

	app.DB.Where("id = ?", post_id).Delete(&Post{})
	app.DB.Unscoped().Where("id = ?", post_id).First(&post)

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

func (app *App) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var tag Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		ErrorWithJSON(w, "Can not decode json", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	tag_id := vars["id"]

	app.DB.Table("tags").Where("id = ?", tag_id).Update(&tag)

	ResponseWithJSON(w, &tag, http.StatusOK)
}

func (app *App) DeleteTag(w http.ResponseWriter, r *http.Request) {
	var tag Tag

	vars := mux.Vars(r)
	tag_id := vars["id"]

	app.DB.Where("id = ?", tag_id).Delete(&Tag{})
	app.DB.Unscoped().Where("id = ?", tag_id).First(&tag)

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
