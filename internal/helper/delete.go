package helper

import (
	"context"
	"log"
	"story-service/protobuffs/story-service"
)

func (h *Helper) DeleteComment(ctx context.Context, in *story.DeleteCommentRequest) (*story.DeleteCommentResponse, error) {

	err := h.database.DeleteComment(ctx, in.CommentId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &story.DeleteCommentResponse{Message: "Success"}, nil
}

func (h *Helper) DeleteStory(ctx context.Context, in *story.DeleteStoryRequest) (*story.DeleteStoryResponse, error) {

	err := h.database.DeleteStory(ctx, in.StoryId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &story.DeleteStoryResponse{Message: "Success"}, nil
}

func (h *Helper) DeleteSubComment(ctx context.Context, in *story.DeleteSubCommentRequest) (*story.DeleteSubCommentResponse, error) {

	err := h.database.DeleteSubComment(ctx, in.SubCommentId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &story.DeleteSubCommentResponse{Message: "Success"}, nil
}
