package restapp

import (
	"encoding/json"
	"net/http"
	"story-service/protobuffs/story-service"
	"strconv"
)

func (s RestApp) SearchRoute(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		s.Search(writer, request)
		return

	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s RestApp) Search(writer http.ResponseWriter, request *http.Request) {
	tag := request.URL.Query().Get("tag")
	count, _ := strconv.Atoi(request.URL.Query().Get("count"))
	if count == 0 {
		count = 10
	}
	skip, _ := strconv.Atoi(request.URL.Query().Get("skip"))
	storyIds, err := s.helper.GetStoryByTag(request.Context(), tag, int64(count), int64(skip))
	if err != nil {
		return
	}
	res := story.GetStoriesByTagResponse{
		StoryIdList: storyIds,
	}
	err = json.NewEncoder(writer).Encode(&res)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
