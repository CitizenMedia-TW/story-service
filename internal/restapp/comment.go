package restapp

import (
	"encoding/json"
	"log"
	"net/http"
	"story-service/internal/restapp/contextkeys"
	"story-service/protobuffs/story-service"
)

func (s RestApp) CommentRoute(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		s.jwtProtect(s.CreateComment, writer, request)
		return
	case "DELETE":
		s.jwtProtect(s.DeleteComment, writer, request)
		return
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s RestApp) DeleteComment(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(contextkeys.UserIdContextKey{})
	if userId == nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	in := &story.DeleteCommentRequest{}
	err := json.NewDecoder(request.Body).Decode(in)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	in.DeleterId = userId.(string)

	res, err := s.helper.DeleteComment(request.Context(), in)

	if err != nil {
		log.Println("Error deleting comment", err.Error())
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(res)
	if err != nil {
		log.Println("Error encoding response")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (s RestApp) CreateComment(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(contextkeys.UserIdContextKey{})
	if userId == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	in := &story.CreateCommentRequest{}
	err := json.NewDecoder(r.Body).Decode(in)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}
	in.CommenterId = userId.(string)
	res, err := s.helper.CreateComment(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
