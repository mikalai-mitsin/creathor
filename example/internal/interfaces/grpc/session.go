package grpc

import (
	"context"
	"fmt"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/018bf/example/pkg/log"
	"github.com/018bf/example/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type SessionServiceServer struct {
	examplepb.UnimplementedSessionServiceServer
	sessionInterceptor interceptors.SessionInterceptor
	logger             log.Logger
}

func NewSessionServiceServer(
	sessionInterceptor interceptors.SessionInterceptor,
	logger log.Logger,
) examplepb.SessionServiceServer {
	return &SessionServiceServer{sessionInterceptor: sessionInterceptor, logger: logger}
}

func (s *SessionServiceServer) Create(
	ctx context.Context,
	input *examplepb.SessionCreate,
) (*examplepb.Session, error) {
	session, err := s.sessionInterceptor.Create(
		ctx,
		encodeSessionCreate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeSession(session), nil
}

func (s *SessionServiceServer) Get(
	ctx context.Context,
	input *examplepb.SessionGet,
) (*examplepb.Session, error) {
	session, err := s.sessionInterceptor.Get(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeSession(session), nil
}

func (s *SessionServiceServer) List(
	ctx context.Context,
	filter *examplepb.SessionFilter,
) (*examplepb.ListSession, error) {
	listSessions, count, err := s.sessionInterceptor.List(
		ctx,
		encodeSessionFilter(filter),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	header := metadata.Pairs("count", fmt.Sprint(count))
	_ = grpc.SendHeader(ctx, header)
	return decodeListSession(listSessions, count), nil
}

func (s *SessionServiceServer) Update(
	ctx context.Context,
	input *examplepb.SessionUpdate,
) (*examplepb.Session, error) {
	session, err := s.sessionInterceptor.Update(
		ctx,
		encodeSessionUpdate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeSession(session), nil
}

func (s *SessionServiceServer) Delete(
	ctx context.Context,
	input *examplepb.SessionDelete,
) (*emptypb.Empty, error) {
	if err := s.sessionInterceptor.Delete(ctx, models.UUID(input.GetId()), ctx.Value(UserKey).(*models.User)); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeSessionCreate(input *examplepb.SessionCreate) *models.SessionCreate {
	create := &models.SessionCreate{Title: input.GetTitle(), Description: input.GetDescription()}
	return create
}
func encodeSessionFilter(input *examplepb.SessionFilter) *models.SessionFilter {
	filter := &models.SessionFilter{
		IDs:        nil,
		PageSize:   nil,
		PageNumber: nil,
		OrderBy:    input.GetOrderBy(),
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = utils.Pointer(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = utils.Pointer(input.GetPageNumber().GetValue())
	}
	for _, id := range input.GetIds() {
		filter.IDs = append(filter.IDs, models.UUID(id))
	}
	if input.GetSearch() != nil {
		filter.Search = utils.Pointer(input.GetSearch().GetValue())
	}
	return filter
}
func encodeSessionUpdate(input *examplepb.SessionUpdate) *models.SessionUpdate {
	update := &models.SessionUpdate{ID: models.UUID(input.GetId())}
	if input.GetTitle() != nil {
		update.Title = utils.Pointer(input.GetTitle().GetValue())
	}
	if input.GetDescription() != nil {
		update.Description = utils.Pointer(input.GetDescription().GetValue())
	}
	return update
}
func decodeSession(session *models.Session) *examplepb.Session {
	response := &examplepb.Session{
		Id:          string(session.ID),
		UpdatedAt:   timestamppb.New(session.UpdatedAt),
		CreatedAt:   timestamppb.New(session.CreatedAt),
		Title:       session.Title,
		Description: session.Description,
	}
	return response
}
func decodeListSession(listSessions []*models.Session, count uint64) *examplepb.ListSession {
	response := &examplepb.ListSession{
		Items: make([]*examplepb.Session, 0, len(listSessions)),
		Count: count,
	}
	for _, session := range listSessions {
		response.Items = append(response.Items, decodeSession(session))
	}
	return response
}
func decodeSessionUpdate(update *models.SessionUpdate) *examplepb.SessionUpdate {
	result := &examplepb.SessionUpdate{
		Id:          string(update.ID),
		Title:       wrapperspb.String(*update.Title),
		Description: wrapperspb.String(*update.Description),
	}
	return result
}
