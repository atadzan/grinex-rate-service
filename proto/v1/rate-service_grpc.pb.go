// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: proto/v1/rate-service.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RateService_GetRates_FullMethodName    = "/rateservice.v1.RateService/GetRates"
	RateService_Healthcheck_FullMethodName = "/rateservice.v1.RateService/Healthcheck"
)

// RateServiceClient is the client API for RateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RateServiceClient interface {
	GetRates(ctx context.Context, in *GetRatesReq, opts ...grpc.CallOption) (*GetRatesResp, error)
	Healthcheck(ctx context.Context, in *HealthcheckReq, opts ...grpc.CallOption) (*HealthcheckResp, error)
}

type rateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRateServiceClient(cc grpc.ClientConnInterface) RateServiceClient {
	return &rateServiceClient{cc}
}

func (c *rateServiceClient) GetRates(ctx context.Context, in *GetRatesReq, opts ...grpc.CallOption) (*GetRatesResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRatesResp)
	err := c.cc.Invoke(ctx, RateService_GetRates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rateServiceClient) Healthcheck(ctx context.Context, in *HealthcheckReq, opts ...grpc.CallOption) (*HealthcheckResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthcheckResp)
	err := c.cc.Invoke(ctx, RateService_Healthcheck_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RateServiceServer is the server API for RateService service.
// All implementations must embed UnimplementedRateServiceServer
// for forward compatibility.
type RateServiceServer interface {
	GetRates(context.Context, *GetRatesReq) (*GetRatesResp, error)
	Healthcheck(context.Context, *HealthcheckReq) (*HealthcheckResp, error)
	mustEmbedUnimplementedRateServiceServer()
}

// UnimplementedRateServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRateServiceServer struct{}

func (UnimplementedRateServiceServer) GetRates(context.Context, *GetRatesReq) (*GetRatesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRates not implemented")
}
func (UnimplementedRateServiceServer) Healthcheck(context.Context, *HealthcheckReq) (*HealthcheckResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Healthcheck not implemented")
}
func (UnimplementedRateServiceServer) mustEmbedUnimplementedRateServiceServer() {}
func (UnimplementedRateServiceServer) testEmbeddedByValue()                     {}

// UnsafeRateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RateServiceServer will
// result in compilation errors.
type UnsafeRateServiceServer interface {
	mustEmbedUnimplementedRateServiceServer()
}

func RegisterRateServiceServer(s grpc.ServiceRegistrar, srv RateServiceServer) {
	// If the following call pancis, it indicates UnimplementedRateServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RateService_ServiceDesc, srv)
}

func _RateService_GetRates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRatesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateServiceServer).GetRates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateService_GetRates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateServiceServer).GetRates(ctx, req.(*GetRatesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RateService_Healthcheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthcheckReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateServiceServer).Healthcheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateService_Healthcheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateServiceServer).Healthcheck(ctx, req.(*HealthcheckReq))
	}
	return interceptor(ctx, in, info, handler)
}

// RateService_ServiceDesc is the grpc.ServiceDesc for RateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rateservice.v1.RateService",
	HandlerType: (*RateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRates",
			Handler:    _RateService_GetRates_Handler,
		},
		{
			MethodName: "Healthcheck",
			Handler:    _RateService_Healthcheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/rate-service.proto",
}
