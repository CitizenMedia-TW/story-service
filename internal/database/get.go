package database

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"story-service/internal/model"
	"story-service/internal/restapp/contextkeys"
	"time"
)

type StoryQuery struct {
	StoryEntity
	AuthorName string
	Comments   []CommentQuery
}

type CommentQuery struct {
	CommentEntity
	CommenterName string
	SubComments   []SubCommentQuery
}

type SubCommentQuery struct {
	SubCommentEntity
	ReplierName string
}

// GetStoryById todo: aggregate author name , commenter name
func (db *SQLDatabase) GetStoryById(ctx context.Context, storyId string) (model.Story, error) {
	storyIdUUID, err := uuid.Parse(storyId)
	if err != nil {
		return model.Story{}, err
	}
	storyQuery, err := getStoryQuery(ctx, db.database, storyIdUUID)
	if err != nil {
		return model.Story{}, err
	}
	comments, err := getStoryComments(ctx, db.database, storyIdUUID)
	if err != nil {
		return model.Story{}, err
	}
	tags, err := getStoryTags(ctx, db.database, storyIdUUID)
	if err != nil {
		return model.Story{}, err
	}
	storyQuery.Comments = comments
	storyQuery.Tags = tags
	return storyQuery.ToDomain(), nil
}

func (db *SQLDatabase) GetStories(ctx context.Context, skip int32, count int32) ([]model.Story, error) {
	stories, err := getStories(ctx, db.database, skip, count)
	if err != nil {
		return nil, err
	}
	var results []model.Story
	for _, story := range stories {
		comments, err := getStoryComments(ctx, db.database, story.Id)
		if err != nil {
			return nil, err
		}
		story.Comments = comments
		results = append(results, story.ToDomain())
	}
	return results, nil
}

func (db *SQLDatabase) GetUserStoryId(ctx context.Context, userEmail string) ([]string, error) {
	var storyIds []string
	stmt, err := db.database.PrepareContext(ctx, `SELECT id FROM story_t WHERE user_mail = $1`)
	if err != nil {
		return storyIds, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, userEmail)
	if err != nil {
		return storyIds, err
	}

	for rows.Next() {
		var storyId string
		err = rows.Scan(&storyId)
		if err != nil {
			return storyIds, err
		}
		storyIds = append(storyIds, storyId)
	}
	return storyIds, err
}

func getStoryTags(ctx context.Context, database *sql.DB, storyId uuid.UUID) ([]string, error) {
	rows, err := database.QueryContext(ctx, `
		SELECT tag FROM story_tag_t WHERE story_id = $1
	`, storyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
func getStories(ctx context.Context, database *sql.DB, skip int32, count int32) ([]StoryQuery, error) {
	rows, err := database.QueryContext(ctx, `
		SELECT s.id, s.title, s.subtitle, s.content, s.created_at, ut.name, s.user_mail FROM story_t s
		LEFT JOIN user_t ut on ut.mail = s.user_mail
		ORDER BY s.created_at DESC
		LIMIT $1 OFFSET $2
	`, count, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []StoryQuery
	if err != nil {
		return results, err
	}
	for rows.Next() {
		var entity StoryQuery
		err := rows.Scan(&entity.Id, &entity.Title, &entity.SubTitle, &entity.Content, &entity.CreatedAt, &entity.AuthorName, &entity.AuthorEmail)
		if err != nil {
			return results, err
		}
		results = append(results, entity)
	}
	return results, nil
}

func getStoryQuery(ctx context.Context, database *sql.DB, storyId uuid.UUID) (StoryQuery, error) {
	logger := ctx.Value(contextkeys.LoggerContextKey{}).(*zap.Logger)
	row := database.QueryRowContext(ctx, `
		SELECT s.title, s.subtitle, s.content, s.created_at, u.name, u.mail  FROM story_t s 
		LEFT JOIN user_t u ON s.user_mail = u.mail WHERE s.id = $1`,
		storyId,
	)
	if row.Err() != nil {
		return StoryQuery{}, row.Err()
	}
	entity := StoryQuery{}
	err := row.Scan(&entity.Title, &entity.SubTitle, &entity.Content, &entity.CreatedAt, &entity.AuthorName, &entity.AuthorEmail)
	if err != nil {
		logger.Log(zap.ErrorLevel, "Error scanning row", zap.Error(err))
		return StoryQuery{}, err
	}
	return entity, err
}

func getStoryComments(ctx context.Context, database *sql.DB, storyId uuid.UUID) ([]CommentQuery, error) {
	logger := ctx.Value(contextkeys.LoggerContextKey{}).(*zap.Logger)
	rows, err := database.QueryContext(ctx, `
		SELECT c.id, c.content, c.created_at, cu.name, cu.mail, sc.id, sc.content, sc.created_at, scu.name, scu.mail
        FROM
            story_t as s
            LEFT JOIN comment_t as c ON s.id = c.story_id
            LEFT JOIN user_t as cu ON c.user_mail = cu.mail
            LEFT JOIN subcomment_t as sc ON c.id = sc.comment_id
            LEFT JOIN user_t as scu ON sc.user_mail = scu.mail
        WHERE s.id = $1`,
		storyId,
	)
	if err != nil {
		logger.Log(zap.ErrorLevel, "Error preparing statement", zap.Error(err))
		return []CommentQuery{}, nil
	}
	defer rows.Close()

	comments := make(map[uuid.UUID]CommentQuery)
	for rows.Next() {
		var commentId uuid.UUID
		var commentContent sql.NullString
		var commentTime sql.NullTime
		var commenterName sql.NullString
		var commenterId sql.NullString
		var subCommentId *uuid.UUID
		var subCommentContent *string
		var subCommentTime *time.Time
		var subCommenterId *string
		var subCommenterName *string
		err = rows.Scan(&commentId, &commentContent, &commentTime, &commenterName, &commenterId, &subCommentId, &subCommentContent, &subCommentTime, &subCommenterName, &subCommenterId)
		if err != nil {
			logger.Log(zap.ErrorLevel, "Error scanning row", zap.Error(err))
			return []CommentQuery{}, nil
		}
		c, ok := comments[commentId]
		if !ok {
			c = CommentQuery{
				CommentEntity: CommentEntity{
					Id:          commentId,
					Content:     commentContent.String,
					CreatedAt:   commentTime.Time,
					CommenterId: commenterId.String,
				},
				CommenterName: commenterName.String,
			}
		}
		if subCommentId != nil {
			c.SubComments = append(c.SubComments, SubCommentQuery{
				SubCommentEntity: SubCommentEntity{
					Id:        *subCommentId,
					Content:   *subCommentContent,
					CreatedAt: *subCommentTime,
					ReplierId: *subCommenterId,
				},
				ReplierName: *subCommenterName,
			})
		}
		comments[commentId] = c
	}
	var result []CommentQuery = make([]CommentQuery, len(comments))
	i := 0
	for _, comment := range comments {
		result[i] = comment
		i++
	}
	return result, nil
}

func (q StoryQuery) ToDomain() model.Story {
	comments := make([]model.Comment, len(q.Comments))
	for i, comment := range q.Comments {
		comments[i] = comment.ToDomain()
		comments[i].SubComments = make([]model.SubComment, len(comment.SubComments))
		for j, subComment := range comment.SubComments {
			comments[i].SubComments[j] = subComment.ToDomain()
		}
	}
	return model.Story{
		Id:          q.Id.String(),
		AuthorEmail: q.AuthorEmail,
		AuthorName:  q.AuthorName,
		Content:     q.Content,
		Title:       q.Title,
		SubTitle:    q.SubTitle,
		CreatedAt:   q.CreatedAt,
		Tags:        q.Tags,
		Comments:    comments,
	}
}

func (db *SQLDatabase) GetStoryIdsByTag(ctx context.Context, tag string, limit int64, offset int64) ([]string, error) {
	var storyIds []string
	rows, err := db.database.QueryContext(ctx, `SELECT story_id FROM story_tag_t WHERE tag = $1 LIMIT $2 OFFSET $3`, tag, limit, offset)
	if err != nil {
		return storyIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var storyId string
		err = rows.Scan(&storyId)
		if err != nil {
			return storyIds, err
		}
		storyIds = append(storyIds, storyId)
	}
	return storyIds, nil

}
