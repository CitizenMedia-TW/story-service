package database

import (
	"context"
	"log"
)

func (db *SQLDatabase) DeleteComment(ctx context.Context, commentId string) error {
	_, err := db.database.ExecContext(ctx, `DELETE FROM comment_t WHERE id = $1`, commentId)
	return err
}

func (db *SQLDatabase) DeleteSubComment(ctx context.Context, subCommentId string) error {
	_, err := db.database.ExecContext(ctx, `DELETE FROM subcomment_t WHERE id = $1`, subCommentId)
	return err
}

func (db *SQLDatabase) DeleteStory(ctx context.Context, storyId string) error {
	result, err := db.database.ExecContext(ctx, `DELETE FROM story_t WHERE id = $1`, storyId)
	if err != nil {
		log.Println("Error in DeleteStory")
		return err
	}
	if a, err := result.RowsAffected(); err == nil && a == 0 {
		return ErrNotFound
	}
	return nil
}
