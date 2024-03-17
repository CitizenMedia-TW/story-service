package restapp

import (
	"encoding/json"
	"net/http"
)

func (s RestApp) LatestRoute(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		s.GetLatestStories(writer, request)
		return
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s RestApp) GetLatestStories(writer http.ResponseWriter, request *http.Request) {
	stories, err := s.helper.GetLatestStories(request.Context())
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
