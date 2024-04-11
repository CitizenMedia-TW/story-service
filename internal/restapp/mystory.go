package restapp

import (
	"encoding/json"
	"net/http"
	"story-service/internal/restapp/contextkeys"
)

func (s RestApp) MyStoryRoute(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		s.GetMyStory(writer, request)
		return
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s RestApp) GetMyStory(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(contextkeys.UserIdContextKey{})
	if userId == nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	stories, err := s.helper.GetMyStoryIds(request.Context(), userId.(string))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(stories)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
