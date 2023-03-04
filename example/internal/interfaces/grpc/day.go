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

type DayServiceServer struct {
	examplepb.UnimplementedDayServiceServer
	dayInterceptor interceptors.DayInterceptor
	logger         log.Logger
}

func NewDayServiceServer(
	dayInterceptor interceptors.DayInterceptor,
	logger log.Logger,
) examplepb.DayServiceServer {
	return &DayServiceServer{dayInterceptor: dayInterceptor, logger: logger}
}

func (s *DayServiceServer) Create(
	ctx context.Context,
	input *examplepb.DayCreate,
) (*examplepb.Day, error) {
	day, err := s.dayInterceptor.Create(
		ctx,
		encodeDayCreate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeDay(day), nil
}

func (s *DayServiceServer) Get(
	ctx context.Context,
	input *examplepb.DayGet,
) (*examplepb.Day, error) {
	day, err := s.dayInterceptor.Get(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeDay(day), nil
}

func (s *DayServiceServer) List(
	ctx context.Context,
	filter *examplepb.DayFilter,
) (*examplepb.ListDay, error) {
	listDays, count, err := s.dayInterceptor.List(
		ctx,
		encodeDayFilter(filter),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	header := metadata.Pairs("count", fmt.Sprint(count))
	_ = grpc.SendHeader(ctx, header)
	return decodeListDay(listDays, count), nil
}

func (s *DayServiceServer) Update(
	ctx context.Context,
	input *examplepb.DayUpdate,
) (*examplepb.Day, error) {
	day, err := s.dayInterceptor.Update(
		ctx,
		encodeDayUpdate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeDay(day), nil
}

func (s *DayServiceServer) Delete(
	ctx context.Context,
	input *examplepb.DayDelete,
) (*emptypb.Empty, error) {
	if err := s.dayInterceptor.Delete(ctx, models.UUID(input.GetId()), ctx.Value(UserKey).(*models.User)); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeDayCreate(input *examplepb.DayCreate) *models.DayCreate {
	create := &models.DayCreate{
		Name:        input.GetName(),
		Repeat:      int(input.GetRepeat()),
		EquipmentID: input.GetEquipmentId(),
	}
	return create
}
func encodeDayFilter(input *examplepb.DayFilter) *models.DayFilter {
	filter := &models.DayFilter{
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
func encodeDayUpdate(input *examplepb.DayUpdate) *models.DayUpdate {
	update := &models.DayUpdate{ID: models.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = utils.Pointer(input.GetName().GetValue())
	}
	if input.GetRepeat() != nil {
		update.Repeat = utils.Pointer(int(input.GetRepeat().GetValue()))
	}
	if input.GetEquipmentId() != nil {
		update.EquipmentID = utils.Pointer(input.GetEquipmentId().GetValue())
	}
	return update
}
func decodeDay(day *models.Day) *examplepb.Day {
	response := &examplepb.Day{
		Id:          string(day.ID),
		UpdatedAt:   timestamppb.New(day.UpdatedAt),
		CreatedAt:   timestamppb.New(day.CreatedAt),
		Name:        day.Name,
		Repeat:      int32(day.Repeat),
		EquipmentId: day.EquipmentID,
	}
	return response
}
func decodeListDay(listDays []*models.Day, count uint64) *examplepb.ListDay {
	response := &examplepb.ListDay{Items: make([]*examplepb.Day, 0, len(listDays)), Count: count}
	for _, day := range listDays {
		response.Items = append(response.Items, decodeDay(day))
	}
	return response
}
func decodeDayUpdate(update *models.DayUpdate) *examplepb.DayUpdate {
	result := &examplepb.DayUpdate{
		Id:          string(update.ID),
		Name:        wrapperspb.String(*update.Name),
		Repeat:      wrapperspb.Int32(int32(*update.Repeat)),
		EquipmentId: wrapperspb.String(*update.EquipmentID),
	}
	return result
}
