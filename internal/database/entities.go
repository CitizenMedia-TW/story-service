package database

import (
	"github.com/google/uuid"
	"story-service/internal/model"
	"time"
)

var StoryTable = "story_t"
var CommentTable = "comment_t"
var SubCommentTable = "subcomment_t"

type StoryEntity struct {
	Id          uuid.UUID
	AuthorEmail string
	Content     string    `bson:"content"`
	Title       string    `bson:"title"`
	SubTitle    string    `bson:"subTitle"`
	CreatedAt   time.Time `bson:"createdAt"`
	Tags        []string  `bson:"tags"`
}

var CommentCollection = "StoryComments"

type CommentEntity struct {
	Id          uuid.UUID
	StoryId     uuid.UUID
	Content     string `bson:"content"`
	CreatedAt   time.Time
	CommenterId uuid.UUID
}

var SubCommentCollection = "StorySubComments"

type SubCommentEntity struct {
	Id        uuid.UUID
	ParentId  uuid.UUID
	Content   string `bson:"content"`
	CreatedAt time.Time
	ReplierId uuid.UUID
}

func (e CommentEntity) ToDomain() model.Comment {
	return model.Comment{
		Id:          e.Id.String(),
		Content:     e.Content,
		CreatedAt:   e.CreatedAt,
		CommenterId: e.CommenterId.String(),
	}
}

func (e SubCommentEntity) ToDomain() model.SubComment {
	return model.SubComment{
		Id:          e.Id.String(),
		Content:     e.Content,
		CreatedAt:   e.CreatedAt,
		CommenterId: e.ReplierId.String(),
	}
}

func (e StoryEntity) ToDomain() model.Story {
	return model.Story{
		Id:        e.Id.String(),
		AuthorId:  e.AuthorId.String(),
		Content:   e.Content,
		Title:     e.Title,
		SubTitle:  e.SubTitle,
		CreatedAt: e.CreatedAt,
		Tags:      e.Tags,
	}
}
