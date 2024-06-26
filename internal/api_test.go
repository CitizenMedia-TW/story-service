package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"story-service/protobuffs/auth-service"
	"story-service/protobuffs/story-service"
	"testing"
)

var testUserId = "user1@example.com"

func GetAuthToken(t *testing.T) string {
	grpcClient, err := grpc.Dial("157.230.46.45:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	token, err := auth.NewAuthServiceClient(grpcClient).GenerateToken(context.TODO(), &auth.GenerateTokenRequest{
		Mail: "user1@example.com",
		Name: "Irrelevant",
	})
	assert.NoError(t, err)
	return token.Token
}

func createStory(t *testing.T, token string) string {
	body := &story.CreateStoryRequest{
		Tags:    []string{"test1", "test2"},
		Content: "test content",
		Title:   "test title",
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	assert.NoError(t, err)
	request, err := http.NewRequest("POST", "http://localhost:50051/story", b)
	request.Header.Set("Authorization", token)
	assert.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	assert.NoError(t, err)

	resBody := &story.CreateStoryResponse{}
	err = json.NewDecoder(response.Body).Decode(resBody)
	assert.NoError(t, err)
	return resBody.StoryId
}

func getStory(t *testing.T, token string, storyId string) *story.GetOneStoryResponse {
	req2Body := &story.GetOneStoryRequest{
		StoryId: storyId,
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(req2Body)
	assert.NoError(t, err)
	request, err := http.NewRequest("GET", "http://localhost:50051/story", b)
	assert.NoError(t, err)
	request.Header.Set("Authorization", token)
	response, err := http.DefaultClient.Do(request)
	assert.NoError(t, err)
	if response.StatusCode == http.StatusNotFound {
		return nil
	}
	res2Body := &story.GetOneStoryResponse{}
	err = json.NewDecoder(response.Body).Decode(res2Body)
	return res2Body
}

func TestCreateAndGetStory(t *testing.T) {
	token := GetAuthToken(t)
	println(token)
	storyId := createStory(t, token)

	assert.NotEmpty(t, storyId)

	storyData := getStory(t, token, storyId)

	assert.Equal(t, "Success", storyData.Message)
	assert.Equal(t, "www", storyData.Story.Author)
	assert.Equal(t, testUserId, storyData.Story.AuthorId)
	assert.Equal(t, "test content", storyData.Story.Content)
	assert.Equal(t, "test title", storyData.Story.Title)

	deleteReq := &story.DeleteStoryRequest{DeleterId: testUserId, StoryId: storyId}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(deleteReq)
	assert.NoError(t, err)
	request, err := http.NewRequest("DELETE", "http://localhost:50051/story", b)
	assert.NoError(t, err)
	request.Header.Set("Authorization", token)
	response, err := http.DefaultClient.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	storyData = getStory(t, token, storyId)
	assert.Nil(t, storyData)
}

func createComment(t *testing.T, token string, storyId string) string {
	request := &story.CreateCommentRequest{
		CommenterId:      testUserId,
		Comment:          "test comment",
		CommentedStoryId: storyId,
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(request)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "http://localhost:50051/comment", b)
	req.Header.Set("Authorization", token)
	assert.NoError(t, err)
	response, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	resBody := &story.CreateCommentResponse{}
	err = json.NewDecoder(response.Body).Decode(resBody)
	assert.NoError(t, err)
	return resBody.CommentId
}

func TestCreateComment(t *testing.T) {
	token := GetAuthToken(t)
	storyId := createStory(t, token)
	assert.NotEmpty(t, storyId)
	commentId := createComment(t, token, storyId)

	storyData := getStory(t, token, storyId)
	assert.Equal(t, storyData.Story.Comments[0].CommenterId, testUserId)
	assert.Equal(t, storyData.Story.Comments[0].Content, "test comment")
	assert.Equal(t, storyData.Story.Comments[0].Id, commentId)

	deleteReq := &story.DeleteCommentRequest{DeleterId: testUserId, CommentId: commentId}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(deleteReq)
	assert.NoError(t, err)
	request, err := http.NewRequest("DELETE", "http://localhost:50051/comment", b)
	assert.NoError(t, err)
	request.Header.Set("Authorization", token)
	response, err := http.DefaultClient.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode, response.Status)

	storyData = getStory(t, token, storyId)
	assert.Empty(t, storyData.Story.Comments)
}
