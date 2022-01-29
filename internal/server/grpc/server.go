//go:generate -command PROTOC protoc -I../../../api/proto ../../../api/proto/metric.proto ../../../api/proto/service.proto
//go:generate PROTOC --go_out=../../../api/ --go-grpc_out=../../../api/

package internalgrpc

import (
	"context"
	"github.com/lomoval/otus-golang-project-sysmon/api"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics/calculator"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"strconv"
	"time"
)

const (
	minInterval = 1      // seconds
	maxInterval = 60 * 5 // seconds
)

type Config struct {
	Host string
	Port int
}

type Server struct {
	api.UnimplementedMetricsServer
	grpcServer *grpc.Server
	addr       string
	collectors []metric.Collector
}

func NewServer(config Config, collectors []metric.Collector) *Server {
	return &Server{addr: net.JoinHostPort(config.Host, strconv.Itoa(config.Port)), collectors: collectors}
}

func (s *Server) Start() error {
	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(loggingHandler))
	api.RegisterMetricsServer(s.grpcServer, s)

	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	log.Printf("starting grpc server on %s", s.addr)
	err = s.grpcServer.Serve(lsn)
	return err
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}

func (s *Server) GetMetrics(request *api.GetMetricsRequest, stream api.Metrics_GetMetricsServer) error {
	if request.GetNotifyInterval() < minInterval || request.GetNotifyInterval() > maxInterval {
		return status.Errorf(codes.InvalidArgument, "incorrect calc interval")
	}
	if request.GetAverageCalcInterval() < minInterval || request.GetAverageCalcInterval() > maxInterval {
		return status.Errorf(codes.InvalidArgument, "incorrect calc interval")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := metriccalc.Start(
		ctx,
		s.collectors,
		time.Duration(request.GetNotifyInterval())*time.Second,
		time.Duration(request.GetAverageCalcInterval())*time.Second,
	)

	for {
		select {
		case <-stream.Context().Done():
			err := stream.Context().Err()
			if err != nil {
				log.Debugf("stream ends: %v", err)
				return err
			}
		case metrics := <-ch:
			err := stream.Send(&api.GetMetricsResponse{Groups: toApiMetrics(metrics)})
			if err != nil {
				log.Errorf("stream send error: %v", err)
				return err
			}
		}
	}
}

func toApiMetrics(groups []metric.Group) []*api.MetricGroup {
	apiGroups := make([]*api.MetricGroup, 0, len(groups))
	for _, group := range groups {
		apiGroup := api.MetricGroup{
			Name:    group.Name,
			Metrics: make([]*api.Metric, len(group.Metrics)),
		}
		apiGroups = append(apiGroups, &apiGroup)
		for i, m := range group.Metrics {
			apiGroup.Metrics[i] = &api.Metric{
				Name:      m.Name,
				Timestamp: timestamppb.New(m.Time),
				Value:     m.Value,
			}
		}
	}
	return apiGroups
}
