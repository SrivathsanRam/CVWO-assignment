package routes

import (
	"encoding/json"
	"net/http"

	"github.com/SrivathsanRam/CVWO_project/backend/internal/api"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/handlers/comments"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/handlers/posts"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/handlers/topics"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/handlers/users"
	"github.com/go-chi/chi/v5"
)

func GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		writeResponse := func(w http.ResponseWriter, response *api.Response) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		// get users
		r.Get("/users", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleList(w, req)
			writeResponse(w, response)
		})
		// get topics
		r.Get("/topics", func(w http.ResponseWriter, req *http.Request) {
			response, _ := topics.HandleList(w, req)
			writeResponse(w, response)
		})

		// create topics
		r.Post("/topics", func(w http.ResponseWriter, req *http.Request) {
			response, _ := topics.HandleCreate(w, req)
			writeResponse(w, response)
		})

		//get topic by id
		r.Get("/topics/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := topics.HandleGet(w, req)
			writeResponse(w, response)
		})
		//update topic by id
		r.Put("/topics/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := topics.HandleUpdate(w, req)
			writeResponse(w, response)
		})
		//delete topic by id
		r.Delete("/topics/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := topics.HandleDelete(w, req)
			writeResponse(w, response)
		})
		//get all posts
		r.Get("/posts", func(w http.ResponseWriter, req *http.Request) {
			response, _ := posts.HandleList(w, req)
			writeResponse(w, response)
		})
		//create post
		r.Post("/posts", func(w http.ResponseWriter, req *http.Request) {
			response, _ := posts.HandleCreate(w, req)
			writeResponse(w, response)
		})
		//get post by id
		r.Get("/posts/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := posts.HandleGet(w, req)
			writeResponse(w, response)
		})
		//update post by id
		r.Put("/posts/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := posts.HandleUpdate(w, req)
			writeResponse(w, response)
		})
		//delete post by id
		r.Delete("/posts/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := posts.HandleDelete(w, req)
			writeResponse(w, response)
		})
		// comments routes
		r.Get("/comments", func(w http.ResponseWriter, req *http.Request) {
			response, _ := comments.HandleList(w, req)
			writeResponse(w, response)
		})
		// create comment
		r.Post("/comments", func(w http.ResponseWriter, req *http.Request) {
			response, _ := comments.HandleCreate(w, req)
			writeResponse(w, response)
		})
		// get comment by id
		r.Get("/comments/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := comments.HandleGet(w, req)
			writeResponse(w, response)
		})

		// update comment by id
		r.Put("/comments/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := comments.HandleUpdate(w, req)
			writeResponse(w, response)
		})
		// delete comment by id
		r.Delete("/comments/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := comments.HandleDelete(w, req)
			writeResponse(w, response)
		})
	}
}
