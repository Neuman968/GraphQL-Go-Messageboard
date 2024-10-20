package dataloaders

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/vikstrous/dataloadgen"
	dm "messageboard.example.graphql/.gen/messageboardDB/public/model"
	. "messageboard.example.graphql/.gen/messageboardDB/public/table"
	gqlmodel "messageboard.example.graphql/graph/model"
)

type ctxKey string

const (
	LoadersKey    = ctxKey("dataloaders")
	CommentsLimit = ctxKey("commentsLimit")
)

type dataLoadersService struct {
	db     *sql.DB
	logger *log.Logger
}

func NewLoaderService(logger *log.Logger, db *sql.DB) *Loaders {
	service := &dataLoadersService{
		db:     db,
		logger: logger,
	}

	return &Loaders{
		CommentsLoader: dataloadgen.NewLoader(service.PostComments, dataloadgen.WithWait(time.Millisecond)),
		UserLoader:     dataloadgen.NewLoader(service.Users, dataloadgen.WithWait(time.Millisecond)),
	}
}

type Loaders struct {
	CommentsLoader *dataloadgen.Loader[string, []*gqlmodel.Comment]
	UserLoader     *dataloadgen.Loader[string, *gqlmodel.User]
}

// Middleware for injecting data loaders into graphql http handler.
func DataLoaderMiddleware(logger *log.Logger, db *sql.DB, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaderService(logger, db)
		r = r.WithContext(context.WithValue(r.Context(), LoadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// Users dataloader implementation.
func (s *dataLoadersService) Users(ctx context.Context, userIds []string) ([]*gqlmodel.User, []error) {

	s.logger.Println("$$ Get Users Data Loader $$")

	var (
		jetIn       []Expression
		errors      []error
		queryResult []*dm.Users
	)

	// Dataloaders require the order of result and ids returned match, This lookup tracks
	// The index of the
	idOrderLookup := make(map[string]int)

	for idx, userId := range userIds {

		userIntId, err := strconv.Atoi(userId)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		// Update index lookup.
		idOrderLookup[userId] = idx
		// Add to Jet Predicates (Where clause).
		jetIn = append(jetIn, Int(int64(userIntId)))
	}

	if len(errors) > 0 {
		return nil, errors
	}

	// Database query.
	stmt := SELECT(Users.AllColumns).FROM(Users).
		WHERE(Users.ID.IN(jetIn...))

	if err := stmt.Query(s.db, &queryResult); err != nil {
		return nil, []error{err}
	}

	// Result is an array matching the size of userIds. Any Elements that do not exist in DB will be marked as "null"
	result := make([]*gqlmodel.User, len(userIds))

	for _, dmUser := range queryResult {

		modelUser := gqlmodel.User{
			ID:   fmt.Sprint(dmUser.ID),
			Name: *dmUser.Name,
		}

		idx, present := idOrderLookup[modelUser.ID]
		if present {
			result[idx] = &modelUser
		}
	}

	return result, nil
}

// PostComments Data Loader.
func (s *dataLoadersService) PostComments(ctx context.Context, postIds []string) ([][]*gqlmodel.Comment, []error) {

	s.logger.Println("$$ Post Comments Data Loader $$")

	var (
		jetIn       []Expression
		errors      []error
		queryResult []*dm.Comment
	)

	postIdIdx := make(map[string]int)

	// Gather postIds for IN(..) Query.
	// Constructs postIdIdx map for tracking order of postIds for result.
	for idx, postId := range postIds {

		postIdInt, err := strconv.Atoi(postId)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		postIdIdx[postId] = idx
		jetIn = append(jetIn, Int(int64(postIdInt)))
	}

	if len(errors) > 0 {
		return nil, errors
	}

	postCommentsWith, rowNumberAlias := CTE("PostCommentsWith"), "rowNumber"

	// Get Limit param from context.
	limit := ctx.Value(CommentsLimit).(int)

	// Query statement
	stmt := WITH(postCommentsWith.AS(
		SELECT(Comment.AllColumns, ROW_NUMBER().OVER(PARTITION_BY(Comment.PostID).ORDER_BY(Comment.CreatedTime)).AS(rowNumberAlias)).
			FROM(Comment).
			WHERE(Comment.PostID.IN(jetIn...))))(
		SELECT(postCommentsWith.AllColumns()).
			FROM(Post.LEFT_JOIN(postCommentsWith, Post.ID.EQ(Comment.PostID.From(postCommentsWith)))).
			WHERE(Post.ID.IN(jetIn...).AND(IntegerColumn(rowNumberAlias).LT_EQ(Int(int64(limit))))))

	// Execute Query
	if err := stmt.Query(s.db, &queryResult); err != nil {
		return nil, []error{err}
	}

	// size of result must match post ids. If there are no comments, the post still must be
	// represented as a null value in the array.
	result := make([][]*gqlmodel.Comment, len(postIds))

	for _, comment := range queryResult {
		modelComment := gqlmodel.Comment{
			ID:           fmt.Sprint(comment.ID),
			PostID:       fmt.Sprint(comment.PostID),
			AuthorUserID: int(comment.AuthorUsersID),
			Text:         *comment.Text,
		}

		if idx, present := postIdIdx[fmt.Sprint(comment.PostID)]; present {
			result[idx] = append(result[idx], &modelComment)
		}
	}

	return result, nil
}

// Returns the dataloader from the context
func For(ctx context.Context) *Loaders {
	return ctx.Value(LoadersKey).(*Loaders)
}

func LoadUser(ctx context.Context, userId string) (*gqlmodel.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userId)
}

func LoadPostComment(ctx context.Context, postId string, limit int) ([]*gqlmodel.Comment, error) {
	loaders := For(ctx)
	return loaders.CommentsLoader.Load(context.WithValue(ctx, CommentsLimit, limit), postId)
}
