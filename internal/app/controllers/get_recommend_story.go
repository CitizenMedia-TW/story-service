package controllers

import (
	"encoding/json"
	"grpc-story-service/protobuffs/story-service"
	"net/http"
)

func (routes HttpRoutes) GetRecommendStory(writer http.ResponseWriter, request *http.Request) {
	in := &story.GetRecommendedRequest{}
	err := json.NewDecoder(request.Body).Decode(in)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := routes.app.GetRecommended(request.Context(), in)
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
