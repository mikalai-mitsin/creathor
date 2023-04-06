package grpc

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/018bf/example/pkg/log"
	"github.com/018bf/example/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
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
	return &EquipmentServiceServer{equipmentInterceptor: equipmentInterceptor, logger: logger}
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
	if err := s.equipmentInterceptor.Delete(ctx, models.UUID(input.GetId()), ctx.Value(UserKey).(*models.User)); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeEquipmentCreate(input *examplepb.EquipmentCreate) *models.EquipmentCreate {
	create := &models.EquipmentCreate{
		Name:   input.GetName(),
		Repeat: int(input.GetRepeat()),
		Weight: int(input.GetWeight()),
	}
	return create
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
	for _, id := range input.GetIds() {
		filter.IDs = append(filter.IDs, models.UUID(id))
	}
	if input.GetSearch() != nil {
		filter.Search = utils.Pointer(input.GetSearch().GetValue())
	}
	return filter
}
func encodeEquipmentUpdate(input *examplepb.EquipmentUpdate) *models.EquipmentUpdate {
	update := &models.EquipmentUpdate{ID: models.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = utils.Pointer(input.GetName().GetValue())
	}
	if input.GetRepeat() != nil {
		update.Repeat = utils.Pointer(int(input.GetRepeat().GetValue()))
	}
	if input.GetWeight() != nil {
		update.Weight = utils.Pointer(int(input.GetWeight().GetValue()))
	}
	return update
}
func decodeEquipment(equipment *models.Equipment) *examplepb.Equipment {
	response := &examplepb.Equipment{
		Id:        string(equipment.ID),
		UpdatedAt: timestamppb.New(equipment.UpdatedAt),
		CreatedAt: timestamppb.New(equipment.CreatedAt),
		Name:      equipment.Name,
		Repeat:    int32(equipment.Repeat),
		Weight:    int32(equipment.Weight),
	}
	return response
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
func decodeEquipmentUpdate(update *models.EquipmentUpdate) *examplepb.EquipmentUpdate {
	result := &examplepb.EquipmentUpdate{
		Id:     string(update.ID),
		Name:   wrapperspb.String(*update.Name),
		Repeat: wrapperspb.Int32(int32(*update.Repeat)),
		Weight: wrapperspb.Int32(int32(*update.Weight)),
	}
	return result
}
