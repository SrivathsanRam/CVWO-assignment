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

type apiHandler func(http.ResponseWriter, *http.Request) (*api.Response, error)

func handle(fn apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		response, err := fn(w, req)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.Response{
				Messages:  []string{err.Error()},
				ErrorCode: http.StatusInternalServerError,
			})
			return
		}
		if response == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.Response{
				Messages:  []string{"Empty response"},
				ErrorCode: http.StatusInternalServerError,
			})
			return
		}
		json.NewEncoder(w).Encode(response)
	}
}

func GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		// get users
		r.Get("/users", handle(users.HandleList))
		// auth
		r.Post("/auth/login", handle(users.HandleLogin))
		r.Post("/auth/logout", handle(users.HandleLogout))
		r.Get("/users/me", handle(users.HandleGetCurrentUser))
		// get topics
		r.Get("/topics", handle(topics.HandleList))

		// create topics
		r.Post("/topics", handle(topics.HandleCreate))

		//get topic by id
		r.Get("/topics/{id}", handle(topics.HandleGet))
		//update topic by id
		r.Put("/topics/{id}", handle(topics.HandleUpdate))
		//delete topic by id
		r.Delete("/topics/{id}", handle(topics.HandleDelete))
		//get all posts
		r.Get("/posts", handle(posts.HandleList))
		//create post
		r.Post("/posts", handle(posts.HandleCreate))
		//get post by id
		r.Get("/posts/{id}", handle(posts.HandleGet))
		//update post by id
		r.Put("/posts/{id}", handle(posts.HandleUpdate))
		//delete post by id
		r.Delete("/posts/{id}", handle(posts.HandleDelete))
		// get comments
		r.Get("/comments", handle(comments.HandleList))
		// create comment
		r.Post("/comments", handle(comments.HandleCreate))
		// get comment by id
		r.Get("/comments/{id}", handle(comments.HandleGet))

		// update comment by id
		r.Put("/comments/{id}", handle(comments.HandleUpdate))
		// delete comment by id
		r.Delete("/comments/{id}", handle(comments.HandleDelete))
	}
}
