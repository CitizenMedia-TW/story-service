package helper

import (
	"context"
	"story-service/internal/database"
	"story-service/internal/utils"
	"story-service/protobuffs/story-service"
)

func (h *Helper) CreateComment(ctx context.Context, in *story.CreateCommentRequest) (*story.CreateCommentResponse, error) {
	id, err := h.database.NewComment(ctx, in.CommentedStoryId, in.CommenterId, in.Comment)
	if err != nil {
		return nil, err
	}
	return &story.CreateCommentResponse{Message: "Success", CommentId: id}, nil
}

func (h *Helper) CreateStory(ctx context.Context, in *story.CreateStoryRequest) (*story.CreateStoryResponse, error) {
	//remove duplicate tags
	in.Tags = utils.RemoveDuplicate(in.Tags)

	id, err := h.database.InsertStory(ctx, database.NewStory{
		UserEmail: in.AuthorId,
		Content:   in.Content,
		Title:     in.Title,
		SubTitle:  in.Subtitle,
		Tags:      in.Tags,
	})

	if err != nil {
		return nil, err
	}

	return &story.CreateStoryResponse{Message: "Success", StoryId: id}, nil
}

func (h *Helper) CreateSubComment(ctx context.Context, in *story.CreateSubCommentRequest) (*story.CreateSubCommentResponse, error) {
	id, err := h.database.NewSubComment(ctx, in.RepliedCommentId, in.StoryId, in.CommenterId, in.Content)
	if err != nil {
		return nil, err
	}
	return &story.CreateSubCommentResponse{Message: "Success", SubCommentId: id}, nil
}
