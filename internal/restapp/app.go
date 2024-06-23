package restapp

import (
	"context"
	"net/http"
	"story-service/internal/helper"
	"story-service/internal/restapp/contextkeys"
	"story-service/protobuffs/jwt-service"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RestApp struct {
	helper helper.Helper
	logger *zap.Logger
}

func New(authClient jwt.JWTServiceClient) RestApp {
	h := helper.New(authClient)
	logger, _ := zap.NewDevelopment()

	return RestApp{
		helper: h,
		logger: logger,
	}
}

func (s RestApp) jwtProtect(next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(contextkeys.LoggerContextKey{}).(*zap.Logger)

	res, err := s.helper.JWTClient.VerifyToken(context.Background(), &jwt.VerifyTokenRequest{Token: r.Header.Get("Authorization")})
	if err != nil || res.Message == "Failed" {
		logger.Log(zap.DebugLevel, "verify failed", zap.String("err", err.Error()))
		http.Error(w, "UnAuthorized", http.StatusUnauthorized)
		return
	}

	logger.Log(zap.DebugLevel, "verify success", zap.String("userId", res.JwtContent.Mail))
	r = r.WithContext(context.WithValue(r.Context(), contextkeys.UserIdContextKey{}, res.JwtContent.Mail))
	next(w, r)
	return
}

func (s RestApp) middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.With(zap.String("requestId", uuid.New().String()))
		logger.Log(zap.DebugLevel, "Executing middleware")
		r = r.WithContext(context.WithValue(r.Context(), contextkeys.LoggerContextKey{}, logger))

		// Allow CORS here By * or specific origin
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
	mux.HandleFunc("/story", s.middleware(s.StoryRoute))
	mux.HandleFunc("/search", s.middleware(s.SearchRoute))
	mux.HandleFunc("/recommend", s.middleware(s.RecommendStoryRoute))
	mux.HandleFunc("/comment", s.middleware(s.CommentRoute))
	mux.HandleFunc("/subComment", s.middleware(s.SubCommentRoute))
	mux.HandleFunc("/mystory", s.middleware(s.MyStoryRoute))
	mux.HandleFunc("/latest", s.middleware(s.LatestRoute))

	return mux
}
