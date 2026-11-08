package adapter

import (
	"fmt"

	infra "github.com/SteeperMold/Orbitalik/api-gateway-service/internal/infrastructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Connections struct {
	Metadata   *grpc.ClientConn
	TLE        *grpc.ClientConn
	Trajectory *grpc.ClientConn
	Auth       *grpc.ClientConn
}

func NewConnections(cfg *infra.Config) (*Connections, error) {
	metadataConn, err := dial(cfg.MetadataAddress, cfg.UseTLS)
	if err != nil {
		return nil, fmt.Errorf("connect to metadata service: %w", err)
	}

	tleConn, err := dial(cfg.TLEAddress, cfg.UseTLS)
	if err != nil {
		_ = metadataConn.Close()
		return nil, fmt.Errorf("connect to TLE service: %w", err)
	}

	trajectoryConn, err := dial(cfg.TrajectoryAddress, cfg.UseTLS)
	if err != nil {
		_ = metadataConn.Close()
		_ = tleConn.Close()
		return nil, fmt.Errorf("connect to trajectory service: %w", err)
	}
	authConn, err := dial(cfg.AuthAddress, cfg.UseTLS)
	if err != nil {
		_ = metadataConn.Close()
		_ = tleConn.Close()
		_ = trajectoryConn.Close()
		return nil, fmt.Errorf("connect to auth service: %w", err)
	}

	return &Connections{
		Metadata:   metadataConn,
		TLE:        tleConn,
		Trajectory: trajectoryConn,
		Auth:       authConn,
	}, nil
}

func dial(address string, useTLS bool) (*grpc.ClientConn, error) {
	var opts grpc.DialOption

	if useTLS {
		opts = grpc.WithTransportCredentials(
			credentials.NewTLS(nil),
		)
	} else {
		opts = grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		)
	}

	return grpc.NewClient(address, opts)
}

func (c *Connections) Close() error {
	var firstErr error

	for _, conn := range []*grpc.ClientConn{
		c.Metadata,
		c.TLE,
		c.Trajectory,
		c.Auth,
	} {
		if conn == nil {
			continue
		}

		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
