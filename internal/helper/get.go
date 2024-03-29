package helper

import (
	"context"
	"log"
	"story-service/protobuffs/story-service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Helper) GetOneStory(ctx context.Context, in *story.GetOneStoryRequest) (*story.GetOneStoryResponse, error) {
	result, err := h.database.GetStoryById(ctx, in.StoryId)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	var comments = make([]*story.Comment, len(result.Comments))

	for i, c := range result.Comments {
		var subComments = make([]*story.SubComment, len(c.SubComments))
		for j, sc := range c.SubComments {
			subComments[j] = &story.SubComment{
				Id:               sc.Id,
				Content:          sc.Content,
				Commenter:        sc.CommenterName,
				CommenterId:      sc.CommenterId,
				Time:             timestamppb.New(sc.CreatedAt),
				RepliedCommentId: c.Id,
			}
		}
		comments[i] = &story.Comment{
			Id:          c.Id,
			Content:     c.Content,
			Commenter:   c.CommenterName,
			CommenterId: c.CommenterId,
			Time:        timestamppb.New(c.CreatedAt),
			SubComments: subComments,
		}
	}

	res := &story.GetOneStoryResponse{
		Message: "Success",
		Story: &story.StoryContent{
			Author:    result.AuthorName,
			AuthorId:  result.AuthorId,
			Content:   result.Content,
			Comments:  comments,
			Title:     result.Title,
			SubTitle:  result.SubTitle,
			CreatedAt: timestamppb.New(result.CreatedAt),
			Tags:      result.Tags,
		},
	}

	return res, nil
}

func (h *Helper) GetMyStoryIds(ctx context.Context, user string) (*story.GetMyStoriesResponse, error) {
	result, err := h.database.GetUserStoryId(ctx, user)
	if err != nil {
		return nil, err
	}
	return &story.GetMyStoriesResponse{Message: "Success", StoryIdList: result}, nil
}

func (h *Helper) GetRecommended(ctx context.Context, in *story.GetRecommendedRequest) (*story.GetRecommendedResponse, error) {
	result, err := h.database.GetStories(ctx, in.Skip, in.Count)
	if err != nil {
		return nil, err
	}
	var ids []string = make([]string, len(result))

	for i, r := range result {
		ids[i] = r.Id
	}

	return &story.GetRecommendedResponse{Message: "Success", StoryIdList: ids}, nil
}

func (h *Helper) GetLatestStories(ctx context.Context) (*story.GetLatestStoriesResponse, error) {
	result, err := h.database.GetStories(ctx, 0, 10)
	if err != nil {
		return nil, err
	}
	var ids []string = make([]string, len(result))

	for i, r := range result {
		ids[i] = r.Id
	}

	return &story.GetLatestStoriesResponse{Message: "Success", StoryIdList: ids}, nil
}
