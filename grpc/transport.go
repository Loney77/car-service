package grpc

import (
	"context"
	"errors"

	kittransport "github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"example.com/car-service/car"
	"example.com/car-service/pb"
	"example.com/car-service/transport"
)

type server struct {
	pb.UnimplementedCarServiceServer
	create   kitgrpc.Handler
	getByID  kitgrpc.Handler
	markSold kitgrpc.Handler
}

func NewService(endpoints transport.Endpoints, logger log.Logger) pb.CarServiceServer {
	opt := kitgrpc.ServerErrorHandler(kittransport.NewLogErrorHandler(logger))
	return &server{
		create:   kitgrpc.NewServer(endpoints.Create, decodeCreateRequest, encodeCreateResponse, opt),
		getByID:  kitgrpc.NewServer(endpoints.GetByID, decodeGetByIDRequest, encodeGetByIDResponse, opt),
		markSold: kitgrpc.NewServer(endpoints.MarkSold, decodeMarkSoldRequest, encodeMarkSoldResponse, opt),
	}
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, response, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return response.(*pb.CreateResponse), nil
}

func (s *server) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	_, response, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return response.(*pb.GetByIDResponse), nil
}

func (s *server) MarkSold(ctx context.Context, req *pb.MarkSoldRequest) (*pb.MarkSoldResponse, error) {
	_, response, err := s.markSold.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return response.(*pb.MarkSoldResponse), nil
}

func decodeCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	c := request.(*pb.CreateRequest).GetCar()
	if c == nil {
		return nil, status.Error(codes.InvalidArgument, "car is required")
	}
	return transport.CreateRequest{Car: car.Car{
		Brand: c.Brand,
		Model: c.Model,
		Year:  int(c.Year),
		Price: c.Price,
	}}, nil
}

func decodeGetByIDRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetByIDRequest)
	return transport.GetByIDRequest{ID: req.Id}, nil
}

func decodeMarkSoldRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.MarkSoldRequest)
	return transport.MarkSoldRequest{ID: req.Id}, nil
}

func encodeCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(transport.CreateResponse)
	return &pb.CreateResponse{Id: res.ID}, nil
}

func encodeGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	c := response.(transport.GetByIDResponse).Car
	return &pb.GetByIDResponse{Car: &pb.Car{
		Id:    c.ID,
		Brand: c.Brand,
		Model: c.Model,
		Year:  int64(c.Year),
		Price: c.Price,
		Sold:  c.Sold,
	}}, nil
}

func encodeMarkSoldResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &pb.MarkSoldResponse{}, nil
}

func encodeError(err error) error {
	if errors.Is(err, car.ErrCarNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	return status.Error(codes.Internal, err.Error())
}
