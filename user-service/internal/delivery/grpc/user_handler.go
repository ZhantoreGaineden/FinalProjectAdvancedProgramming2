package grpcdelivery

import (
	"context"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/userpb"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/entity"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/usecase"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	usecase *usecase.UserUsecase
}

func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.UserResponse, error) {
	user, err := h.usecase.RegisterUser(ctx, req.GetFullName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{User: toProtoUser(user)}, nil
}

func (h *UserHandler) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.LoginResponse, error) {
	token, user, err := h.usecase.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{
		Token: token,
		User:  toProtoUser(user),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	user, err := h.usecase.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{User: toProtoUser(user)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	user := entity.User{
		ID:       req.GetId(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	}

	updated, err := h.usecase.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{User: toProtoUser(updated)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	err := h.usecase.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userpb.DeleteUserResponse{
		Success: true,
		Message: "user deleted successfully",
	}, nil
}

func toProtoUser(user entity.User) *userpb.User {
	createdAt := ""
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt.Format(time.RFC3339)
	}

	return &userpb.User{
		Id:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: createdAt,
	}
}
