package service

import (
	"context"
	"sync"

	courierv1 "github.com/moondoggy/courier/proto/courier/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	courierv1.UnimplementedCourierServiceServer

	mu          sync.Mutex
	data        map[int32]*courierv1.Order
	lastOrderId int32
}

func NewServer() *Server {
	return &Server{
		mu:   sync.Mutex{},
		data: make(map[int32]*courierv1.Order, 256),
	}
}

func (s *Server) CreateOrder(_ context.Context, request *courierv1.CreateOrderRequest) (*courierv1.Order, error) {
	if request.SellerAddress == "" ||
		request.ReceiverAddress == "" ||
		request.WeightGrams <= 0 {
		return nil, status.Error(codes.InvalidArgument, "missing required fields or its invalid")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order := courierv1.Order{
		SellerAddress:   request.SellerAddress,
		ReceiverAddress: request.ReceiverAddress,
		WeightGrams:     request.WeightGrams,
		Id:              s.lastOrderId + 1,
		Status:          courierv1.OrderStatus_ORDER_STATUS_CREATED,
		CreatedAt:       timestamppb.Now(),
	}

	s.data[order.Id] = &order
	s.lastOrderId++
	return &order, nil
}

func (s *Server) GetOrder(ctx context.Context, request *courierv1.GetOrderRequest) (*courierv1.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.data[request.GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	return order, nil
}
