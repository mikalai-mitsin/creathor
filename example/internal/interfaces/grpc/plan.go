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

type PlanServiceServer struct {
	examplepb.UnimplementedPlanServiceServer
	planInterceptor interceptors.PlanInterceptor
	logger          log.Logger
}

func NewPlanServiceServer(
	planInterceptor interceptors.PlanInterceptor,
	logger log.Logger,
) examplepb.PlanServiceServer {
	return &PlanServiceServer{planInterceptor: planInterceptor, logger: logger}
}

func (s *PlanServiceServer) Create(
	ctx context.Context,
	input *examplepb.PlanCreate,
) (*examplepb.Plan, error) {
	plan, err := s.planInterceptor.Create(
		ctx,
		encodePlanCreate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodePlan(plan), nil
}

func (s *PlanServiceServer) Get(
	ctx context.Context,
	input *examplepb.PlanGet,
) (*examplepb.Plan, error) {
	plan, err := s.planInterceptor.Get(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodePlan(plan), nil
}

func (s *PlanServiceServer) List(
	ctx context.Context,
	filter *examplepb.PlanFilter,
) (*examplepb.ListPlan, error) {
	listPlans, count, err := s.planInterceptor.List(
		ctx,
		encodePlanFilter(filter),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	header := metadata.Pairs("count", fmt.Sprint(count))
	_ = grpc.SendHeader(ctx, header)
	return decodeListPlan(listPlans, count), nil
}

func (s *PlanServiceServer) Update(
	ctx context.Context,
	input *examplepb.PlanUpdate,
) (*examplepb.Plan, error) {
	plan, err := s.planInterceptor.Update(
		ctx,
		encodePlanUpdate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodePlan(plan), nil
}

func (s *PlanServiceServer) Delete(
	ctx context.Context,
	input *examplepb.PlanDelete,
) (*emptypb.Empty, error) {
	if err := s.planInterceptor.Delete(ctx, models.UUID(input.GetId()), ctx.Value(UserKey).(*models.User)); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodePlanCreate(input *examplepb.PlanCreate) *models.PlanCreate {
	create := &models.PlanCreate{
		Name:        input.GetName(),
		Repeat:      input.GetRepeat(),
		EquipmentID: input.GetEquipmentId(),
	}
	return create
}
func encodePlanFilter(input *examplepb.PlanFilter) *models.PlanFilter {
	filter := &models.PlanFilter{
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
func encodePlanUpdate(input *examplepb.PlanUpdate) *models.PlanUpdate {
	update := &models.PlanUpdate{ID: models.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = utils.Pointer(input.GetName().GetValue())
	}
	if input.GetRepeat() != nil {
		update.Repeat = utils.Pointer(input.GetRepeat().GetValue())
	}
	if input.GetEquipmentId() != nil {
		update.EquipmentID = utils.Pointer(input.GetEquipmentId().GetValue())
	}
	return update
}
func decodePlan(plan *models.Plan) *examplepb.Plan {
	response := &examplepb.Plan{
		Id:          string(plan.ID),
		UpdatedAt:   timestamppb.New(plan.UpdatedAt),
		CreatedAt:   timestamppb.New(plan.CreatedAt),
		Name:        plan.Name,
		Repeat:      plan.Repeat,
		EquipmentId: plan.EquipmentID,
	}
	return response
}
func decodeListPlan(listPlans []*models.Plan, count uint64) *examplepb.ListPlan {
	response := &examplepb.ListPlan{Items: make([]*examplepb.Plan, 0, len(listPlans)), Count: count}
	for _, plan := range listPlans {
		response.Items = append(response.Items, decodePlan(plan))
	}
	return response
}
func decodePlanUpdate(update *models.PlanUpdate) *examplepb.PlanUpdate {
	result := &examplepb.PlanUpdate{
		Id:          string(update.ID),
		Name:        wrapperspb.String(*update.Name),
		Repeat:      wrapperspb.UInt64(*update.Repeat),
		EquipmentId: wrapperspb.String(*update.EquipmentID),
	}
	return result
}
