package restapp

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"story-service/internal/helper"
	"story-service/internal/restapp/contextkeys"
	"story-service/protobuffs/auth-service"
)

type RestApp struct {
	helper helper.Helper
	logger *zap.Logger
}

func New(authClient auth.AuthServiceClient) RestApp {
	h := helper.New(authClient)
	logger, _ := zap.NewDevelopment()

	return RestApp{
		helper: h,
		logger: logger,
	}
}

func (s RestApp) middlewares(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.With(zap.String("requestId", uuid.New().String()))
		logger.Log(zap.DebugLevel, "Executing middleware")
		r = r.WithContext(context.WithValue(r.Context(), contextkeys.LoggerContextKey{}, logger))

		res, err := s.helper.AuthClient.VerifyToken(context.Background(), &auth.VerifyTokenRequest{Token: r.Header.Get("Authorization")})
		logger.Log(zap.DebugLevel, "verify token", zap.String("status", res.Message))
		if err == nil && res.Message != "Failed" {
			logger.Log(zap.DebugLevel, "verify Success", zap.String("userId", res.JwtContent.Mail))
			r = r.WithContext(context.WithValue(r.Context(), contextkeys.UserIdContextKey{}, res.JwtContent.Mail))
		}

		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s RestApp) Routes() http.Handler {
	// Declare a new router
	mux := http.NewServeMux()

	// Use middleware for all routes
	mux.Handle("/story", s.middlewares(http.HandlerFunc(s.StoryRoute)))
	mux.Handle("/search", s.middlewares(http.HandlerFunc(s.SearchRoute)))
	mux.Handle("/recommend", s.middlewares(http.HandlerFunc(s.GetRecommendStory)))
	mux.Handle("/comment", s.middlewares(http.HandlerFunc(s.CommentRoute)))
	mux.Handle("/subComment", s.middlewares(http.HandlerFunc(s.SubCommentRoute)))
	mux.Handle("/mystory", s.middlewares(http.HandlerFunc(s.MyStoryRoute)))
	mux.Handle("/latest", s.middlewares(http.HandlerFunc(s.LatestRoute)))

	return mux
}
