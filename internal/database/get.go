package database

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"story-service/internal/model"
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
func (db *SQLDatabase) GetStoryById(ctx context.Context, storyId uuid.UUID) (model.Story, error) {
	storyQuery, err := getStoryQuery(ctx, db.database, storyId)
	if err != nil {
		return model.Story{}, err
	}
	comments, err := getStoryComments(ctx, db.database, storyId)
	if err != nil {
		return model.Story{}, err
	}
	storyQuery.Comments = comments
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

func getStories(ctx context.Context, database *sql.DB, skip int32, count int32) ([]StoryQuery, error) {
	rows, err := database.QueryContext(ctx, `
		SELECT s.id, s.title, s.subtitle, s.content, s.created_at, ut.name FROM story_t s
		LEFT JOIN user_t ut on ut.mail = s.user_mail
		ORDER BY s.created_at DESC
		LIMIT $1 OFFSET $2
	`)
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
		err := rows.Scan(&entity.Id, &entity.Title, &entity.SubTitle, &entity.Content, &entity.CreatedAt, &entity.AuthorName)
		if err != nil {
			return results, err
		}
		results = append(results, entity)
	}
	return results, nil
}

func getStoryQuery(ctx context.Context, database *sql.DB, storyId uuid.UUID) (StoryQuery, error) {
	logger := ctx.Value("logger").(*zap.Logger)
	statement, err := database.PrepareContext(ctx, `
		SELECT s.title, s.subtitle, s.content, s.created_at, u.name  FROM story_t s 
		JOIN user_t u ON s.user_mail = u.mail WHERE s.id = $1`,
	)

	if err != nil {
		return StoryQuery{}, err
	}
	defer statement.Close()
	rows, err := statement.QueryContext(ctx, storyId)
	if err != nil {
		return StoryQuery{}, err
	}
	rows.Close()
	entity := StoryQuery{}
	if !rows.Next() {
		return entity, ErrNotFound
	}

	err = rows.Scan(&entity.Title, &entity.SubTitle, &entity.Content, &entity.CreatedAt, &entity.AuthorName)
	if err != nil {
		logger.Log(zap.ErrorLevel, "Error scanning row", zap.Error(err))
		return StoryQuery{}, err
	}
	return entity, err
}

func getStoryComments(ctx context.Context, database *sql.DB, storyId uuid.UUID) ([]CommentQuery, error) {
	logger := ctx.Value("logger").(*zap.Logger)
	statement, err := database.PrepareContext(ctx, `
		SELECT c.id,c.content,cu.name,sc.id,sc.content,scu.name
FROM
    story_t s
    LEFT JOIN comment_t c ON s.id = c.story_id
    LEFT JOIN user_t cu ON c.user_mail = cu.mail
    LEFT JOIN subcomment_t sc ON c.id = sc.comment_id
    LEFT JOIN user_t scu ON sc.user_mail = scu.mail
WHERE
    s.id = $1
`)
	if err != nil {
		logger.Log(zap.ErrorLevel, "Error preparing statement", zap.Error(err))
		return []CommentQuery{}, nil
	}
	defer statement.Close()
	rows, err := statement.QueryContext(ctx, storyId)
	if err != nil {
		logger.Log(zap.ErrorLevel, "Error querying database", zap.Error(err))
		return []CommentQuery{}, nil
	}
	defer rows.Close()

	var comments map[uuid.UUID]CommentQuery
	for rows.Next() {
		var commentId uuid.UUID
		var commentContent string
		var commenterName string
		var subCommentId *uuid.UUID
		var subCommentContent *string
		var subCommenterName *string
		err = rows.Scan(&commentId, &commentContent, &commenterName, &subCommentId, &subCommentContent, &subCommenterName)
		if err != nil {
			logger.Log(zap.ErrorLevel, "Error scanning row", zap.Error(err))
			return []CommentQuery{}, nil
		}
		c, ok := comments[commentId]
		if !ok {
			c = CommentQuery{
				CommentEntity: CommentEntity{
					Id:      commentId,
					Content: commentContent,
				},
				CommenterName: commenterName,
			}
		}
		if subCommentId != nil {
			c.SubComments = append(c.SubComments, SubCommentQuery{
				SubCommentEntity: SubCommentEntity{
					Id:      *subCommentId,
					Content: *subCommentContent,
				},
				ReplierName: *subCommenterName,
			})
		}
		comments[commentId] = c
	}
	var result []CommentQuery = make([]CommentQuery, len(comments))
	for _, comment := range comments {
		result = append(result, comment)
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
		Id:         q.Id.String(),
		AuthorId:   q.AuthorId.String(),
		AuthorName: q.AuthorName,
		Content:    q.Content,
		Title:      q.Title,
		SubTitle:   q.SubTitle,
		CreatedAt:  q.CreatedAt,
		Tags:       q.Tags,
		Comments:   comments,
	}
}
