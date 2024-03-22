package database

import (
	"context"
	"github.com/google/uuid"
	"log"
	"time"
)

type NewStory struct {
	userEmail string
	Content   string
	Title     string
	SubTitle  string
	Tags      []string
}

func (db *SQLDatabase) NewStory(ctx context.Context, story NewStory) (uuid.UUID, error) {
	//probably should check if author exist, but since it's nosql database, and it does not affect the query outcome, we'll skip it for now.
	storyEntity := StoryEntity{
		Id:          uuid.New(),
		AuthorEmail: story.userEmail,
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
		return uuid.Nil, err
	}

	//todo: use a better logging/tracing system
	return storyEntity.Id, nil
}

func (db *SQLDatabase) NewComment(ctx context.Context, commentedStoryId uuid.UUID, commenterMail string, content string) (uuid.UUID, error) {
	//probably should check if story and commenter exist, but since it's nosql database, and it does not affect the query outcome, we'll skip it for now.
	commentId := uuid.New()
	_, err := db.database.ExecContext(ctx, `
	INSERT INTO comment_t (id, story_id, content, time, user_mail) VALUES ($1, $2, $3, $4, $5)`,
		commentId, commentedStoryId, content, time.Now(), commenterMail,
	)

	if err != nil {
		return uuid.Nil, err
	}

	return commentId, nil
}

func (db *SQLDatabase) NewSubComment(ctx context.Context, repliedCommentId uuid.UUID, replierId string, content string) (string, error) {
	subCommentId := uuid.New()

	db.database.ExecContext(ctx, `
		INSERT INTO subcomment_t (id, comment_id, content, time, user_mail) VALUES ($1, $2, $3, $4, $5)`,
		subCommentId, repliedCommentId, content, time.Now(), replierId,
	)

	_, err = db.database.Collection(SubCommentCollection).InsertOne(ctx, reply)

	if err != nil {
		return "", err
	}

	return reply.Id.Hex(), nil
}
