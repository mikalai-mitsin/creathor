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

type ArchServiceServer struct {
	examplepb.UnimplementedArchServiceServer
	archInterceptor interceptors.ArchInterceptor
	logger          log.Logger
}

func NewArchServiceServer(
	archInterceptor interceptors.ArchInterceptor,
	logger log.Logger,
) examplepb.ArchServiceServer {
	return &ArchServiceServer{
		archInterceptor: archInterceptor,
		logger:          logger,
	}
}

func (s *ArchServiceServer) Create(
	ctx context.Context,
	input *examplepb.ArchCreate,
) (*examplepb.Arch, error) {
	arch, err := s.archInterceptor.Create(
		ctx,
		encodeArchCreate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeArch(arch), nil
}

func (s *ArchServiceServer) Get(
	ctx context.Context,
	input *examplepb.ArchGet,
) (*examplepb.Arch, error) {
	arch, err := s.archInterceptor.Get(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeArch(arch), nil
}

func (s *ArchServiceServer) List(
	ctx context.Context,
	filter *examplepb.ArchFilter,
) (*examplepb.ListArch, error) {
	listArches, count, err := s.archInterceptor.List(
		ctx,
		encodeArchFilter(filter),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeListArch(listArches, count), nil
}

func (s *ArchServiceServer) Update(
	ctx context.Context,
	input *examplepb.ArchUpdate,
) (*examplepb.Arch, error) {
	arch, err := s.archInterceptor.Update(
		ctx,
		encodeArchUpdate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeArch(arch), nil
}

func (s *ArchServiceServer) Delete(
	ctx context.Context,
	input *examplepb.ArchDelete,
) (*emptypb.Empty, error) {
	if err := s.archInterceptor.Delete(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	); err != nil {
		return nil, decodeError(err)
	}
	return &emptypb.Empty{}, nil
}

func encodeArchUpdate(input *examplepb.ArchUpdate) *models.ArchUpdate {
	update := &models.ArchUpdate{
		ID:      models.UUID(input.GetId()),
		Name:    nil,
		Release: nil,
		Tested:  nil,
	}
	if input.GetName() != nil {
		update.Name = utils.Pointer(string(input.GetName().GetValue()))
	}
	if input.GetRelease() != nil {
		update.Release = utils.Pointer(input.GetRelease().AsTime())
	}
	if input.GetTested() != nil {
		update.Tested = utils.Pointer(input.GetTested().AsTime())
	}
	return update
}

func decodeListArch(listArches []*models.Arch, count uint64) *examplepb.ListArch {
	response := &examplepb.ListArch{
		Items: make([]*examplepb.Arch, 0, len(listArches)),
		Count: count,
	}
	for _, arch := range listArches {
		response.Items = append(response.Items, decodeArch(arch))
	}
	return response
}

func encodeArchFilter(input *examplepb.ArchFilter) *models.ArchFilter {
	filter := &models.ArchFilter{
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

func encodeArchCreate(input *examplepb.ArchCreate) *models.ArchCreate {
	create := &models.ArchCreate{
		Name:    string(input.GetName()),
		Release: input.GetRelease().AsTime(),
		Tested:  input.GetTested().AsTime(),
	}
	return create
}

func decodeArch(arch *models.Arch) *examplepb.Arch {
	response := &examplepb.Arch{
		Id:        string(arch.ID),
		UpdatedAt: timestamppb.New(arch.UpdatedAt),
		CreatedAt: timestamppb.New(arch.CreatedAt),
		Name:      string(arch.Name),
		Release:   timestamppb.New(arch.Release),
		Tested:    timestamppb.New(arch.Tested),
	}
	return response
}
