package grpc

import (
	"context"
	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/examplepb"
	"github.com/018bf/example/pkg/log"
	"github.com/018bf/example/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	return &DayServiceServer{
		dayInterceptor: dayInterceptor,
		logger:         logger,
	}
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
	if err := s.dayInterceptor.Delete(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}

func encodeDayUpdate(input *examplepb.DayUpdate) *models.DayUpdate {
	update := &models.DayUpdate{
		ID:          models.UUID(input.GetId()),
		Name:        nil,
		Repeat:      nil,
		EquipmentID: nil,
	}
	if input.GetName() != nil {
		update.Name = utils.Pointer(string(input.GetName().GetValue()))
	}
	if input.GetRepeat() != nil {
		update.Repeat = utils.Pointer(int(input.GetRepeat().GetValue()))
	}
	if input.GetEquipmentId() != nil {
		update.EquipmentID = utils.Pointer(string(input.GetEquipmentId().GetValue()))
	}
	return update
}

func decodeListDay(listDays []*models.Day, count uint64) *examplepb.ListDay {
	response := &examplepb.ListDay{
		Items: make([]*examplepb.Day, 0, len(listDays)),
		Count: count,
	}
	for _, day := range listDays {
		response.Items = append(response.Items, decodeDay(day))
	}
	return response
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
	if input.GetSearch() != nil {
		filter.Search = utils.Pointer(input.GetSearch().GetValue())
	}
	for _, id := range input.GetIds() {
		filter.IDs = append(filter.IDs, models.UUID(id))
	}
	return filter
}

func encodeDayCreate(input *examplepb.DayCreate) *models.DayCreate {
	create := &models.DayCreate{
		Name:        string(input.GetName()),
		Repeat:      int(input.GetRepeat()),
		EquipmentID: string(input.GetEquipmentId()),
	}
	return create
}

func decodeDay(day *models.Day) *examplepb.Day {
	response := &examplepb.Day{
		Id:          string(day.ID),
		UpdatedAt:   timestamppb.New(day.UpdatedAt),
		CreatedAt:   timestamppb.New(day.CreatedAt),
		Name:        string(day.Name),
		Repeat:      int32(day.Repeat),
		EquipmentId: string(day.EquipmentID),
	}
	return response
}
