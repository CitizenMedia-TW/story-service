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

func (db *SQLDatabase) InsertStory(ctx context.Context, story NewStory) (string, error) {
	storyId := uuid.New()
	_, err := db.database.ExecContext(ctx, `
		INSERT INTO story_t (id, user_mail, content, title, subtitle, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		storyId, story.UserEmail, story.Content, story.Title, story.SubTitle, time.Now(),
	)

	if err != nil {
		log.Println("Error in push")
		return "", err
	}

	//todo: use a better logging/tracing system
	return storyId.String(), nil
}

func (db *SQLDatabase) NewComment(ctx context.Context, commentedStoryId string, commenterMail string, content string) (string, error) {
	//probably should check if story and commenter exist, but since it's nosql database, and it does not affect the query outcome, we'll skip it for now.
	commentId := uuid.New()
	_, err := db.database.ExecContext(ctx, `
	INSERT INTO comment_t (id, story_id, content, created_at, user_mail) VALUES ($1, $2, $3, $4, $5)`,
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
		INSERT INTO subcomment_t (id, comment_id, content, created_at, user_mail) VALUES ($1, $2, $3, $4, $5)`,
		subCommentId, repliedCommentId, content, time.Now(), replierId,
	)

	if err != nil {
		return "", err
	}
	return subCommentId.String(), nil
}
