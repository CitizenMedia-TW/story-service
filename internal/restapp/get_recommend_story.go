package restapp

import (
	"encoding/json"
	"grpc-story-service/protobuffs/story-service"
	"net/http"
	"strconv"
)

func (s RestApp) GetRecommendStory(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value("userId")
	if userId == nil {
		http.Error(writer, "UnAuthorized", http.StatusUnauthorized)
		return
	}

	// in := &story.GetRecommendedRequest{}
	// err := json.NewDecoder(request.Body).Decode(in)
	// if err != nil {
	// 	http.Error(writer, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// userId := request.URL.Query().Get("userId")
	strCount := request.URL.Query().Get("count")
	count, err := strconv.ParseInt(strCount, 10, 32)
	strSkip := request.URL.Query().Get("skip")
	skip, err := strconv.ParseInt(strSkip, 10, 32)
	in := &story.GetRecommendedRequest{
		UserId: userId.(string),
		Count:  int32(count),
		Skip:   int32(skip),
	}

	res, err := s.helper.GetRecommended(request.Context(), in)
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
