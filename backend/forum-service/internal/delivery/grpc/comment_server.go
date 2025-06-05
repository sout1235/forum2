package grpc

import (
	"context"
	"strconv"

	"github.com/sout1235/forum2/backend/forum-service/api/proto"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommentServer struct {
	proto.UnimplementedCommentServiceServer
	commentUseCase usecase.CommentUseCase
}

func NewCommentServer(commentUseCase usecase.CommentUseCase) *CommentServer {
	return &CommentServer{
		commentUseCase: commentUseCase,
	}
}

func (s *CommentServer) GetCommentsByTopic(ctx context.Context, req *proto.GetCommentsByTopicRequest) (*proto.GetCommentsByTopicResponse, error) {
	topicID, err := strconv.ParseInt(req.TopicId, 10, 64)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentUseCase.GetCommentsByTopicID(ctx, topicID)
	if err != nil {
		return nil, err
	}

	var protoComments []*proto.Comment
	for _, comment := range comments {
		protoComments = append(protoComments, &proto.Comment{
			Id:        comment.ID,
			Content:   comment.Content,
			AuthorId:  comment.AuthorID,
			TopicId:   comment.TopicID,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			UpdatedAt: timestamppb.New(comment.UpdatedAt),
		})
	}

	return &proto.GetCommentsByTopicResponse{
		Comments: protoComments,
	}, nil
}

func (s *CommentServer) GetCommentByID(ctx context.Context, req *proto.GetCommentByIDRequest) (*proto.GetCommentByIDResponse, error) {
	commentID, err := strconv.ParseInt(req.CommentId, 10, 64)
	if err != nil {
		return nil, err
	}

	comment, err := s.commentUseCase.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return &proto.GetCommentByIDResponse{
		Comment: &proto.Comment{
			Id:        comment.ID,
			Content:   comment.Content,
			AuthorId:  comment.AuthorID,
			TopicId:   comment.TopicID,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			UpdatedAt: timestamppb.New(comment.UpdatedAt),
		},
	}, nil
}

func (s *CommentServer) CreateComment(ctx context.Context, req *proto.CreateCommentRequest) (*proto.CreateCommentResponse, error) {
	topicID, err := strconv.ParseInt(req.TopicId, 10, 64)
	if err != nil {
		return nil, err
	}

	comment := &entity.Comment{
		Content:  req.Content,
		AuthorID: req.AuthorId,
		TopicID:  topicID,
	}

	err = s.commentUseCase.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &proto.CreateCommentResponse{
		Comment: &proto.Comment{
			Id:        comment.ID,
			Content:   comment.Content,
			AuthorId:  comment.AuthorID,
			TopicId:   comment.TopicID,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			UpdatedAt: timestamppb.New(comment.UpdatedAt),
		},
	}, nil
}

func (s *CommentServer) DeleteComment(ctx context.Context, req *proto.DeleteCommentRequest) (*proto.DeleteCommentResponse, error) {
	commentID, err := strconv.ParseInt(req.CommentId, 10, 64)
	if err != nil {
		return nil, err
	}

	err = s.commentUseCase.DeleteComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteCommentResponse{}, nil
}

func (s *CommentServer) LikeComment(ctx context.Context, req *proto.LikeCommentRequest) (*proto.LikeCommentResponse, error) {
	commentID, err := strconv.ParseInt(req.CommentId, 10, 64)
	if err != nil {
		return nil, err
	}

	err = s.commentUseCase.LikeComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return &proto.LikeCommentResponse{}, nil
}
