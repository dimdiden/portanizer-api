package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler for checking the api status
func (app *App) CheckHealth(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{"Message": "OK"}
	ResponseWithJSON(w, &body, http.StatusOK)
}

// Handler for getting all posts
func (app *App) GetPostList(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	// Find all posts and populate their tags
	app.DB.Preload("Tags").Order("ID ASC").Find(&posts)

	ResponseWithJSON(w, &posts, http.StatusOK)
}

func (app *App) GetPost(w http.ResponseWriter, r *http.Request) {
	var post Post

	pid := mux.Vars(r)["id"]
	if app.DB.First(&post, "id = ?", pid).RecordNotFound() {
		// (w, "Failed. This post doesn't exist", http.StatusNotFound)
		ResponseWithJSON(w, "Failed. This post doesn't exist", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, &post, http.StatusOK)
}

// Handler for creating post
func (app *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	var rsvpost Post
	if err := json.NewDecoder(r.Body).Decode(&rsvpost); err != nil {
		ResponseWithJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	// Empty post name validation
	if rsvpost.Name == "" {
		ResponseWithJSON(w, "Failed. Name field is empty", http.StatusBadRequest)
		return
	}

	var post Post
	// Unique post name validation
	app.DB.First(&post, "name = ?", rsvpost.Name)
	if !app.DB.NewRecord(post) {
		ResponseWithJSON(w, "Failed. This post name already exists", http.StatusBadRequest)
		return
	}
	// Creating post
	post = Post{Name: rsvpost.Name, Body: rsvpost.Body, Tags: []Tag{}}
	app.DB.Create(&post)
	// Create tag if doesn't exist and assign tags to post
	for _, t := range rsvpost.Tags {
		app.DB.FirstOrCreate(&t, Tag{Name: t.Name})
		app.DB.Model(&post).Association("Tags").Append(t)
	}

	ResponseWithJSON(w, &post, http.StatusOK)
}

// Handler for updating post
func (app *App) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	// Query the post by id and validate existance
	pid := mux.Vars(r)["id"]
	if app.DB.First(&post, "id = ?", pid).RecordNotFound() {
		ResponseWithJSON(w, "Failed. This post doesn't exist", http.StatusBadRequest)
		return
	}
	// Read the request body
	var rsvpost Post
	if err := json.NewDecoder(r.Body).Decode(&rsvpost); err != nil {
		ResponseWithJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	// Empty post name validation
	if rsvpost.Name == "" {
		ResponseWithJSON(w, "Failed. Name field is empty", http.StatusBadRequest)
		return
	}
	// Unique post name validation
	var chk_post Post
	if !app.DB.First(&chk_post, "name = ?", rsvpost.Name).RecordNotFound() && chk_post.ID != post.ID {
		ResponseWithJSON(w, "Failed. This post name already exists", http.StatusBadRequest)
		return
	}
	// Update post with new parameters
	post.Name, post.Body, post.Tags = rsvpost.Name, rsvpost.Body, []Tag{}
	app.DB.Save(&post)
	// Create tag if doesn't exist and assign tags to post
	for _, t := range rsvpost.Tags {
		app.DB.FirstOrCreate(&t, Tag{Name: t.Name})
		app.DB.Model(&post).Association("Tags").Append(t)
	}

	ResponseWithJSON(w, &post, http.StatusOK)
}

func (app *App) DeletePost(w http.ResponseWriter, r *http.Request) {
	var post Post

	id := mux.Vars(r)["id"]

	// app.DB.Where("id = ?", id).Delete(&Tag{})
	// app.DB.Unscoped().Where("id = ?", id).First(&tag)

	if app.DB.First(post, id).RecordNotFound() {
		ResponseWithJSON(w, "Failed. This post doesn't exist", http.StatusNotFound)
		return
	}
	app.DB.Model(&post).Association("Tags").Clear()

	app.DB.Delete(&post)

	// app.DB.Where("id = ?", pid).Delete(&Post{})
	// app.DB.Unscoped().Where("id = ?", pid).First(&post)

	ResponseWithJSON(w, &post, http.StatusOK)
}

func (app *App) GetTagList(w http.ResponseWriter, r *http.Request) {
	var tags []Tag
	app.DB.Order("ID ASC").Find(&tags)

	ResponseWithJSON(w, &tags, http.StatusOK)
}

func (app *App) GetTag(w http.ResponseWriter, r *http.Request) {
	var tag Tag

	pid := mux.Vars(r)["id"]
	if app.DB.First(&tag, "id = ?", pid).RecordNotFound() {
		ResponseWithJSON(w, "Failed. This tag doesn't exist", http.StatusNotFound)
		return
	}

	ResponseWithJSON(w, &tag, http.StatusOK)
}

func (app *App) CreateTag(w http.ResponseWriter, r *http.Request) {
	var tag Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		ResponseWithJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
		return
	}
	app.DB.Save(&tag)
	ResponseWithJSON(w, &tag, http.StatusOK)
}

func (app *App) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var tag Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		ResponseWithJSON(w, "Failed. Please check json syntax", http.StatusBadRequest)
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
	id := vars["id"]

	// app.DB.Unscoped().Where("id = ?", id).First(&tag)

	if app.DB.First(tag, id).RecordNotFound() {
		ResponseWithJSON(w, "Failed. This tag doesn't exist", http.StatusNotFound)
		return
	}

	// app.DB.Where("id = ?", id).Delete(&Tag{})
	app.DB.Delete(tag)
	ResponseWithJSON(w, "Tag "+id+" has been deleted successfully", http.StatusOK)
	// // app.DB.Model(&tag).Clear()
	//
	// app.DB.Delete(&tag)
	// ResponseWithJSON(w, &tag, http.StatusOK)
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
