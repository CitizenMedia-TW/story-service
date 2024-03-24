package database

import (
	"context"
	"github.com/google/uuid"
	"log"
	"time"
)

type NewStory struct {
	UserEmail string
	Content   string
	Title     string
	SubTitle  string
	Tags      []string
}

func (db *SQLDatabase) NewStory(ctx context.Context, story NewStory) (string, error) {
	//probably should check if author exist, but since it's nosql database, and it does not affect the query outcome, we'll skip it for now.
	storyEntity := StoryEntity{
		Id:          uuid.New(),
		AuthorEmail: story.UserEmail,
		Content:     story.Content,
		Title:       story.Title,
		SubTitle:    story.SubTitle,
		CreatedAt:   time.Now(),
		Tags:        story.Tags,
	}

	_, err := db.database.ExecContext(ctx, `
		INSERT INTO story_t (id, user_mail, content, title, subtitle, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		storyEntity.Id, storyEntity.AuthorEmail, storyEntity.Content, storyEntity.Title, storyEntity.SubTitle, storyEntity.CreatedAt,
	)

	if err != nil {
		log.Println("Error in push")
		return "", err
	}

	//todo: use a better logging/tracing system
	return storyEntity.Id.String(), nil
}

func (db *SQLDatabase) NewComment(ctx context.Context, commentedStoryId string, commenterMail string, content string) (string, error) {
	//probably should check if story and commenter exist, but since it's nosql database, and it does not affect the query outcome, we'll skip it for now.
	commentId := uuid.New()
	_, err := db.database.ExecContext(ctx, `
	INSERT INTO comment_t (id, story_id, content, time, user_mail) VALUES ($1, $2, $3, $4, $5)`,
		commentId, commentedStoryId, content, time.Now(), commenterMail,
	)

	if err != nil {
		return "", err
	}

	return commentId.String(), nil
}

func (db *SQLDatabase) NewSubComment(ctx context.Context, repliedCommentId string, replierId string, content string) (string, error) {
	subCommentId := uuid.New()

	_, err := db.database.ExecContext(ctx, `
		INSERT INTO subcomment_t (id, comment_id, content, time, user_mail) VALUES ($1, $2, $3, $4, $5)`,
		subCommentId, repliedCommentId, content, time.Now(), replierId,
	)

	if err != nil {
		return "", err
	}
	return subCommentId.String(), nil
}
