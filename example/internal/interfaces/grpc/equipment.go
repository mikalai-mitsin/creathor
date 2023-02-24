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
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type EquipmentServiceServer struct {
	examplepb.UnimplementedEquipmentServiceServer
	equipmentInterceptor interceptors.EquipmentInterceptor
	logger               log.Logger
}

func NewEquipmentServiceServer(
	equipmentInterceptor interceptors.EquipmentInterceptor,
	logger log.Logger,
) examplepb.EquipmentServiceServer {
	return &EquipmentServiceServer{
		equipmentInterceptor: equipmentInterceptor,
		logger:               logger,
	}
}

func (s *EquipmentServiceServer) Create(
	ctx context.Context,
	input *examplepb.EquipmentCreate,
) (*examplepb.Equipment, error) {
	equipment, err := s.equipmentInterceptor.Create(
		ctx,
		encodeEquipmentCreate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeEquipment(equipment), nil
}

func (s *EquipmentServiceServer) Get(
	ctx context.Context,
	input *examplepb.EquipmentGet,
) (*examplepb.Equipment, error) {
	equipment, err := s.equipmentInterceptor.Get(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeEquipment(equipment), nil
}

func (s *EquipmentServiceServer) List(
	ctx context.Context,
	filter *examplepb.EquipmentFilter,
) (*examplepb.ListEquipment, error) {
	listEquipment, count, err := s.equipmentInterceptor.List(
		ctx,
		encodeEquipmentFilter(filter),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	header := metadata.Pairs("count", fmt.Sprint(count))
	_ = grpc.SendHeader(ctx, header)
	return decodeListEquipment(listEquipment, count), nil
}

func (s *EquipmentServiceServer) Update(
	ctx context.Context,
	input *examplepb.EquipmentUpdate,
) (*examplepb.Equipment, error) {
	equipment, err := s.equipmentInterceptor.Update(
		ctx,
		encodeEquipmentUpdate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeEquipment(equipment), nil
}

func (s *EquipmentServiceServer) Delete(
	ctx context.Context,
	input *examplepb.EquipmentDelete,
) (*emptypb.Empty, error) {
	if err := s.equipmentInterceptor.Delete(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}

func encodeEquipmentUpdate(input *examplepb.EquipmentUpdate) *models.EquipmentUpdate {
	update := &models.EquipmentUpdate{
		ID:          models.UUID(input.GetId()),
		Title:       nil,
		Description: nil,
		Weight:      nil,
		Versions:    nil,
		Release:     nil,
		Tested:      nil,
	}
	if input.GetTitle() != nil {
		update.Title = utils.Pointer(string(input.GetTitle().GetValue()))
	}
	if input.GetDescription() != nil {
		update.Description = utils.Pointer(string(input.GetDescription().GetValue()))
	}
	if input.GetWeight() != nil {
		update.Weight = utils.Pointer(uint64(input.GetWeight().GetValue()))
	}
	if input.GetVersions() != nil {
		var params []uint64
		for _, item := range input.GetVersions().GetValues() {
			params = append(params, uint64(item.GetNumberValue()))
		}
		update.Versions = &params
	}
	if input.GetRelease() != nil {
		update.Release = utils.Pointer(input.GetRelease().AsTime())
	}
	if input.GetTested() != nil {
		update.Tested = utils.Pointer(input.GetTested().AsTime())
	}
	return update
}

func decodeListEquipment(listEquipment []*models.Equipment, count uint64) *examplepb.ListEquipment {
	response := &examplepb.ListEquipment{
		Items: make([]*examplepb.Equipment, 0, len(listEquipment)),
		Count: count,
	}
	for _, equipment := range listEquipment {
		response.Items = append(response.Items, decodeEquipment(equipment))
	}
	return response
}

func encodeEquipmentFilter(input *examplepb.EquipmentFilter) *models.EquipmentFilter {
	filter := &models.EquipmentFilter{
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

func encodeEquipmentCreate(input *examplepb.EquipmentCreate) *models.EquipmentCreate {
	create := &models.EquipmentCreate{
		Title:       string(input.GetTitle()),
		Description: string(input.GetDescription()),
		Weight:      uint64(input.GetWeight()),
		Versions:    nil,
		Release:     input.GetRelease().AsTime(),
		Tested:      input.GetTested().AsTime(),
	}
	for _, param := range input.GetVersions() {
		create.Versions = append(create.Versions, uint64(param))
	}
	return create
}

func decodeEquipment(equipment *models.Equipment) *examplepb.Equipment {
	response := &examplepb.Equipment{
		Id:          string(equipment.ID),
		UpdatedAt:   timestamppb.New(equipment.UpdatedAt),
		CreatedAt:   timestamppb.New(equipment.CreatedAt),
		Title:       string(equipment.Title),
		Description: string(equipment.Description),
		Weight:      uint64(equipment.Weight),
		Versions:    nil,
		Release:     timestamppb.New(equipment.Release),
		Tested:      timestamppb.New(equipment.Tested),
	}
	for _, param := range equipment.Versions {
		response.Versions = append(response.Versions, uint64(param))
	}
	return response
}

func decodeEquipmentUpdate(update *models.EquipmentUpdate) *examplepb.EquipmentUpdate {
	result := &examplepb.EquipmentUpdate{
		Id:          string(update.ID),
		Title:       wrapperspb.String(string(*update.Title)),
		Description: wrapperspb.String(string(*update.Description)),
		Weight:      wrapperspb.UInt64(uint64(*update.Weight)),
		Versions:    nil,
		Release:     timestamppb.New(*update.Release),
		Tested:      timestamppb.New(*update.Tested),
	}
	if update.Versions != nil {
		params, err := structpb.NewList(utils.ToAnySlice(*update.Versions))
		if err != nil {
			return nil
		}
		result.Versions = params
	}
	return result
}
