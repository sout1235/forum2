package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sout1235/forum2/backend/forum-service/api/proto"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommentUseCase struct {
	mock.Mock
}

func (m *MockCommentUseCase) GetCommentsByTopicID(ctx context.Context, topicID int64) ([]*entity.Comment, error) {
	args := m.Called(ctx, topicID)
	return args.Get(0).([]*entity.Comment), args.Error(1)
}

func (m *MockCommentUseCase) GetCommentByID(ctx context.Context, commentID int64) (*entity.Comment, error) {
	args := m.Called(ctx, commentID)
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *MockCommentUseCase) CreateComment(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentUseCase) DeleteComment(ctx context.Context, commentID int64) error {
	args := m.Called(ctx, commentID)
	return args.Error(0)
}

func (m *MockCommentUseCase) LikeComment(ctx context.Context, commentID int64) error {
	args := m.Called(ctx, commentID)
	return args.Error(0)
}

func TestCommentServer_GetCommentsByTopic(t *testing.T) {
	muc := new(MockCommentUseCase)
	server := NewCommentServer(muc)
	ctx := context.Background()
	timeNow := time.Now()
	comments := []*entity.Comment{{ID: 1, Content: "c1", AuthorID: 2, TopicID: 3, CreatedAt: timeNow, UpdatedAt: timeNow}}
	muc.On("GetCommentsByTopicID", ctx, int64(3)).Return(comments, nil)

	resp, err := server.GetCommentsByTopic(ctx, &proto.GetCommentsByTopicRequest{TopicId: "3"})
	assert.NoError(t, err)
	assert.Len(t, resp.Comments, 1)
	assert.Equal(t, int64(1), resp.Comments[0].Id)

	// error case
	muc.On("GetCommentsByTopicID", ctx, int64(99)).Return([]*entity.Comment{}, errors.New("fail"))
	_, err = server.GetCommentsByTopic(ctx, &proto.GetCommentsByTopicRequest{TopicId: "99"})
	assert.Error(t, err)

	// invalid id
	_, err = server.GetCommentsByTopic(ctx, &proto.GetCommentsByTopicRequest{TopicId: "bad"})
	assert.Error(t, err)
}

func TestCommentServer_GetCommentByID(t *testing.T) {
	muc := new(MockCommentUseCase)
	server := NewCommentServer(muc)
	ctx := context.Background()
	timeNow := time.Now()
	comment := &entity.Comment{ID: 1, Content: "c1", AuthorID: 2, TopicID: 3, CreatedAt: timeNow, UpdatedAt: timeNow}
	muc.On("GetCommentByID", ctx, int64(1)).Return(comment, nil)

	resp, err := server.GetCommentByID(ctx, &proto.GetCommentByIDRequest{CommentId: "1"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Comment.Id)

	// error case
	muc.On("GetCommentByID", ctx, int64(99)).Return(&entity.Comment{}, errors.New("fail"))
	_, err = server.GetCommentByID(ctx, &proto.GetCommentByIDRequest{CommentId: "99"})
	assert.Error(t, err)

	// invalid id
	_, err = server.GetCommentByID(ctx, &proto.GetCommentByIDRequest{CommentId: "bad"})
	assert.Error(t, err)
}

func TestCommentServer_CreateComment(t *testing.T) {
	muc := new(MockCommentUseCase)
	server := NewCommentServer(muc)
	ctx := context.Background()
	muc.On("CreateComment", ctx, mock.AnythingOfType("*entity.Comment")).Return(nil)

	req := &proto.CreateCommentRequest{Content: "c1", AuthorId: 2, TopicId: "3"}
	resp, err := server.CreateComment(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "c1", resp.Comment.Content)

	// error case
	muc.ExpectedCalls = nil
	muc.On("CreateComment", ctx, mock.AnythingOfType("*entity.Comment")).Return(errors.New("fail"))
	_, err = server.CreateComment(ctx, req)
	assert.Error(t, err)

	// invalid id
	badReq := &proto.CreateCommentRequest{Content: "c1", AuthorId: 2, TopicId: "bad"}
	_, err = server.CreateComment(ctx, badReq)
	assert.Error(t, err)
}

func TestCommentServer_DeleteComment(t *testing.T) {
	muc := new(MockCommentUseCase)
	server := NewCommentServer(muc)
	ctx := context.Background()
	muc.On("DeleteComment", ctx, int64(1)).Return(nil)

	resp, err := server.DeleteComment(ctx, &proto.DeleteCommentRequest{CommentId: "1"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// error case
	muc.On("DeleteComment", ctx, int64(99)).Return(errors.New("fail"))
	_, err = server.DeleteComment(ctx, &proto.DeleteCommentRequest{CommentId: "99"})
	assert.Error(t, err)

	// invalid id
	_, err = server.DeleteComment(ctx, &proto.DeleteCommentRequest{CommentId: "bad"})
	assert.Error(t, err)
}

func TestCommentServer_LikeComment(t *testing.T) {
	muc := new(MockCommentUseCase)
	server := NewCommentServer(muc)
	ctx := context.Background()
	muc.On("LikeComment", ctx, int64(1)).Return(nil)

	resp, err := server.LikeComment(ctx, &proto.LikeCommentRequest{CommentId: "1"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// error case
	muc.On("LikeComment", ctx, int64(99)).Return(errors.New("fail"))
	_, err = server.LikeComment(ctx, &proto.LikeCommentRequest{CommentId: "99"})
	assert.Error(t, err)

	// invalid id
	_, err = server.LikeComment(ctx, &proto.LikeCommentRequest{CommentId: "bad"})
	assert.Error(t, err)
}
