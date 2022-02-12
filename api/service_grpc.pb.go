// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: service.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	GetMetrics(ctx context.Context, in *GetMetricsRequest, opts ...grpc.CallOption) (Metrics_GetMetricsClient, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) GetMetrics(ctx context.Context, in *GetMetricsRequest, opts ...grpc.CallOption) (Metrics_GetMetricsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Metrics_ServiceDesc.Streams[0], "/Metrics/GetMetrics", opts...)
	if err != nil {
		return nil, err
	}
	x := &metricsGetMetricsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Metrics_GetMetricsClient interface {
	Recv() (*GetMetricsResponse, error)
	grpc.ClientStream
}

type metricsGetMetricsClient struct {
	grpc.ClientStream
}

func (x *metricsGetMetricsClient) Recv() (*GetMetricsResponse, error) {
	m := new(GetMetricsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility
type MetricsServer interface {
	GetMetrics(*GetMetricsRequest, Metrics_GetMetricsServer) error
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsServer struct {
}

func (UnimplementedMetricsServer) GetMetrics(*GetMetricsRequest, Metrics_GetMetricsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetMetrics not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_GetMetrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetMetricsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetricsServer).GetMetrics(m, &metricsGetMetricsServer{stream})
}

type Metrics_GetMetricsServer interface {
	Send(*GetMetricsResponse) error
	grpc.ServerStream
}

type metricsGetMetricsServer struct {
	grpc.ServerStream
}

func (x *metricsGetMetricsServer) Send(m *GetMetricsResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetMetrics",
			Handler:       _Metrics_GetMetrics_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}
