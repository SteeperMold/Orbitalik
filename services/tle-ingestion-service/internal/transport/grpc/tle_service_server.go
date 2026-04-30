package grpc

import (
	"context"
	"fmt"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/gen/tlepb"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TLEServiceServer struct {
	tlepb.UnimplementedTleServiceServer

	service domain.TLEService
	logger  applog.Logger
}

func NewTLEServiceServer(s domain.TLEService, logger applog.Logger) *TLEServiceServer {
	return &TLEServiceServer{
		service: s,
		logger:  logger,
	}
}

func (s *TLEServiceServer) GetTle(ctx context.Context, req *tlepb.GetTleRequest) (*tlepb.GetTleResponse, error) {
	var tle *models.TLE
	var err error
	var identifier string

	switch id := req.Identifier.Kind.(type) {

	case *tlepb.SatelliteIdentifier_NoradId:
		identifier = fmt.Sprintf("norad_id=%d", id.NoradId)
		tle, err = s.service.GetTLEByNoradID(ctx, int(id.NoradId))

	case *tlepb.SatelliteIdentifier_SatelliteName:
		identifier = fmt.Sprintf("satellite_name=%s", id.SatelliteName)
		tle, err = s.service.GetTLEBySatelliteName(ctx, id.SatelliteName)

	default:
		return nil, status.Error(codes.InvalidArgument, "either norad_id or name must be set")
	}

	if err != nil {
		s.logger.Error("failed to get tle by identifier",
			applog.NewField("identifier", identifier),
			applog.NewErrorField(err),
		)
		return nil, status.Error(codes.Internal, "failed to get TLE")
	}
	if tle == nil {
		return nil, status.Error(codes.NotFound, "TLE not found")
	}

	return &tlepb.GetTleResponse{
		Tle: &tlepb.Tle{
			NoradId:       uint32(tle.NoradID),
			SatelliteName: tle.SatelliteName,
			Line1:         tle.Line1,
			Line2:         tle.Line2,
			Epoch:         timestamppb.New(tle.Epoch),
		},
	}, nil
}

func (s *TLEServiceServer) ListTles(ctx context.Context, _ *tlepb.ListTlesRequest) (*tlepb.ListTlesResponse, error) {
	tles, err := s.service.GetAllTLEs(ctx)
	if err != nil {
		s.logger.Error("failed to get all tles", applog.NewErrorField(err))
		return nil, status.Error(codes.Internal, "failed to get TLEs")
	}

	pbTLEs := make([]*tlepb.Tle, 0, len(tles))
	for _, tle := range tles {
		pbTLEs = append(pbTLEs, &tlepb.Tle{
			NoradId:       uint32(tle.NoradID),
			SatelliteName: tle.SatelliteName,
			Line1:         tle.Line1,
			Line2:         tle.Line2,
			Epoch:         timestamppb.New(tle.Epoch),
		})
	}

	return &tlepb.ListTlesResponse{Tles: pbTLEs}, nil
}
