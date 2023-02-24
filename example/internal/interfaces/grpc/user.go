package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/018bf/example/pkg/log"
	"github.com/018bf/example/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServiceServer struct {
	examplepb.UnimplementedUserServiceServer
	userInterceptor interceptors.UserInterceptor
	logger          log.Logger
}

func NewUserServiceServer(
	userInterceptor interceptors.UserInterceptor,
	logger log.Logger,
) examplepb.UserServiceServer {
	return &UserServiceServer{userInterceptor: userInterceptor, logger: logger}
}

func (u UserServiceServer) Signup(
	ctx context.Context,
	input *examplepb.Signup,
) (*examplepb.User, error) {
	signup := &models.UserCreate{
		Email:    strings.ToLower(input.GetEmail()),
		Password: input.GetPassword(),
	}
	user, err := u.userInterceptor.Create(ctx, signup, ctx.Value(UserKey).(*models.User))
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeUser(user), nil
}

func (u UserServiceServer) Update(
	ctx context.Context,
	input *examplepb.UserUpdate,
) (*examplepb.User, error) {
	update := &models.UserUpdate{
		ID: models.UUID(input.GetId()),
	}
	if input.GetFirstName() != nil {
		update.FirstName = utils.Pointer(input.GetFirstName().GetValue())
	}
	if input.GetLastName() != nil {
		update.LastName = utils.Pointer(input.GetLastName().GetValue())
	}
	if input.GetEmail() != nil {
		update.Email = utils.Pointer(strings.ToLower(input.GetEmail().GetValue()))
	}
	if input.GetPassword() != nil {
		update.Password = utils.Pointer(input.GetPassword().GetValue())
	}
	user, err := u.userInterceptor.Update(ctx, update, ctx.Value(UserKey).(*models.User))
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeUser(user), nil
}

func (u UserServiceServer) Get(
	ctx context.Context,
	input *examplepb.UserGet,
) (*examplepb.User, error) {
	user, err := u.userInterceptor.Get(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeUser(user), nil
}

func (u UserServiceServer) List(
	ctx context.Context,
	input *examplepb.UserFilter,
) (*examplepb.Users, error) {
	filter := encodeUserFilter(input)
	users, count, err := u.userInterceptor.List(ctx, filter, ctx.Value(UserKey).(*models.User))
	if err != nil {
		return nil, decodeError(err)
	}
	response := &examplepb.Users{
		Items: make([]*examplepb.User, 0, len(users)),
		Count: count,
	}
	for _, user := range users {
		response.Items = append(response.Items, decodeUser(user))
	}
	header := metadata.Pairs("count", fmt.Sprint(count))
	_ = grpc.SendHeader(ctx, header)
	return response, nil
}

func encodeUserFilter(input *examplepb.UserFilter) *models.UserFilter {
	filter := &models.UserFilter{}
	if input.PageNumber != nil {
		filter.PageNumber = utils.Pointer(input.GetPageNumber().GetValue())
	}
	if input.PageSize != nil {
		filter.PageSize = utils.Pointer(input.GetPageSize().GetValue())
	}
	if input.Search != nil {
		filter.Search = utils.Pointer(input.GetSearch().GetValue())
	}
	if len(input.OrderBy) > 0 {
		filter.OrderBy = input.OrderBy
	}
	return filter
}

func (u UserServiceServer) Delete(
	ctx context.Context,
	input *examplepb.UserDelete,
) (*emptypb.Empty, error) {
	if err := u.userInterceptor.Delete(ctx, models.UUID(input.GetId()), ctx.Value(UserKey).(*models.User)); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}

func decodeUser(user *models.User) *examplepb.User {
	u := &examplepb.User{
		Id:        string(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		RegionId:  nil,
	}
	return u
}
