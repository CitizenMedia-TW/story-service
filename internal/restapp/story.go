package restapp

import (
	"encoding/json"
	"net/http"
	"story-service/internal/database"
	"story-service/internal/restapp/contextkeys"
	"story-service/protobuffs/story-service"
)

func (s RestApp) StoryRoute(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		s.jwtProtect(s.CreateStory, writer, request)
		return
	case "GET":
		s.GetOneStory(writer, request)
		return
	case "DELETE":
		s.jwtProtect(s.DeleteStory, writer, request)
		return
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s RestApp) CreateStory(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(contextkeys.UserIdContextKey{})

	if userId == nil {
		http.Error(writer, "UnAuthorized", http.StatusUnauthorized)
		return
	}

	in := &story.CreateStoryRequest{}
	err := json.NewDecoder(request.Body).Decode(in)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	in.AuthorId = userId.(string)

	res, err := s.helper.CreateStory(request.Context(), in)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(res)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (s RestApp) GetOneStory(writer http.ResponseWriter, request *http.Request) {
	storyId := request.URL.Query().Get("storyId")
	in := &story.GetOneStoryRequest{
		StoryId: storyId,
	}
	// err := json.NewDecoder(request.Body).Decode(in)
	// if err != nil {
	// 	http.Error(writer, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	res, err := s.helper.GetOneStory(request.Context(), in)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(writer, "Story not found", http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(res)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (s RestApp) DeleteStory(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(contextkeys.UserIdContextKey{})
	if userId == nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	in := &story.DeleteStoryRequest{}
	err := json.NewDecoder(request.Body).Decode(in)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	in.DeleterId = userId.(string)

	res, err := s.helper.DeleteStory(request.Context(), in)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(res)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
