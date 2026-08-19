package grpc

import (
	"context"
	"fmt"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/gen/metadatapb"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetadataService interface {
	GetMetadataByNoradID(ctx context.Context, noradID int) (*models.SatelliteMetadata, error)
	GetMetadataByName(ctx context.Context, name string) (*models.SatelliteMetadata, error)
	ListSatellites(ctx context.Context, filter *models.ListFilter) ([]*models.SatelliteMetadata, string, uint32, error)
}

type MetadataServer struct {
	metadatapb.UnimplementedSatelliteMetadataServiceServer

	service         MetadataService
	logger          applog.Logger
	maxPageSize     uint32
	defaultPageSize uint32
}

func NewMetadataServer(s MetadataService, logger applog.Logger, maxPageSize, defaultPageSize uint32) *MetadataServer {
	return &MetadataServer{
		service:         s,
		logger:          logger,
		maxPageSize:     maxPageSize,
		defaultPageSize: defaultPageSize,
	}
}

func (s *MetadataServer) GetSatelliteMetadata(
	ctx context.Context,
	req *metadatapb.GetMetadataRequest,
) (*metadatapb.GetMetadataResponse, error) {

	if req == nil || req.Identifier == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"either norad_id or name must be set",
		)
	}

	var (
		metadata   *models.SatelliteMetadata
		err        error
		identifier string
	)

	switch id := req.Identifier.Kind.(type) {

	case *metadatapb.SatelliteIdentifier_NoradId:
		identifier = fmt.Sprintf("norad_id=%d", id.NoradId)
		metadata, err = s.service.GetMetadataByNoradID(ctx, int(id.NoradId))

	case *metadatapb.SatelliteIdentifier_SatelliteName:
		identifier = fmt.Sprintf("satellite_name=%s", id.SatelliteName)
		metadata, err = s.service.GetMetadataByName(ctx, id.SatelliteName)

	default:
		return nil, status.Error(codes.InvalidArgument, "either norad_id or name must be set")
	}

	if err != nil {
		s.logger.Error("failed to get metadata by identifier",
			applog.NewField("identifier", identifier),
			applog.NewErrorField(err),
		)
		return nil, status.Error(codes.Internal, "failed to get metadata")
	}
	if metadata == nil {
		return nil, status.Error(codes.NotFound, "metadata not found")
	}

	return &metadatapb.GetMetadataResponse{
		Metadata: toProtoSatelliteMetadata(metadata),
	}, nil
}

func (s *MetadataServer) ListSatelliteMetadata(
	ctx context.Context,
	req *metadatapb.ListSatelliteMetadataRequest,
) (*metadatapb.ListSatelliteMetadataResponse, error) {

	if req.PageSize == 0 {
		req.PageSize = s.defaultPageSize
	}

	if req.PageSize > s.maxPageSize {
		return nil, status.Error(codes.InvalidArgument, "page_size too large")
	}

	filter := &models.ListFilter{
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	}

	if req.ObjectType != nil {
		v := models.ObjectType(*req.ObjectType)
		filter.ObjectType = &v
	}

	if req.MissionType != nil {
		v := models.MissionType(*req.MissionType)
		filter.MissionType = &v
	}

	if req.OperationalStatus != nil {
		v := models.OperationalStatus(*req.OperationalStatus)
		filter.OperationalStatus = &v
	}

	if req.OrbitRegime != nil {
		v := models.OrbitRegime(*req.OrbitRegime)
		filter.OrbitRegime = &v
	}

	if req.Constellation != nil && *req.Constellation != "" {
		filter.Constellation = req.Constellation
	}

	items, nextToken, total, err := s.service.ListSatellites(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list satellite metadata", applog.NewErrorField(err))
		return nil, status.Error(codes.Internal, "failed to list satellite metadata")
	}

	respItems := make([]*metadatapb.SatelliteMetadata, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, toProtoSatelliteMetadata(item))
	}

	return &metadatapb.ListSatelliteMetadataResponse{
		Items:         respItems,
		NextPageToken: nextToken,
		Total:         total,
	}, nil
}
