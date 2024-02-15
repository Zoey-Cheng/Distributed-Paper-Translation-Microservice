// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: paper.proto

package v1

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for PaperService service

func NewPaperServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for PaperService service

type PaperService interface {
	Create(ctx context.Context, in *CreatePaper, opts ...client.CallOption) (*Paper, error)
	Fetch(ctx context.Context, in *PaperID, opts ...client.CallOption) (*Paper, error)
	Delete(ctx context.Context, in *PaperID, opts ...client.CallOption) (*DeletePaper, error)
	Fetchs(ctx context.Context, in *ReqFetchs, opts ...client.CallOption) (*RespFetchs, error)
}

type paperService struct {
	c    client.Client
	name string
}

func NewPaperService(name string, c client.Client) PaperService {
	return &paperService{
		c:    c,
		name: name,
	}
}

func (c *paperService) Create(ctx context.Context, in *CreatePaper, opts ...client.CallOption) (*Paper, error) {
	req := c.c.NewRequest(c.name, "PaperService.Create", in)
	out := new(Paper)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paperService) Fetch(ctx context.Context, in *PaperID, opts ...client.CallOption) (*Paper, error) {
	req := c.c.NewRequest(c.name, "PaperService.Fetch", in)
	out := new(Paper)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paperService) Delete(ctx context.Context, in *PaperID, opts ...client.CallOption) (*DeletePaper, error) {
	req := c.c.NewRequest(c.name, "PaperService.Delete", in)
	out := new(DeletePaper)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paperService) Fetchs(ctx context.Context, in *ReqFetchs, opts ...client.CallOption) (*RespFetchs, error) {
	req := c.c.NewRequest(c.name, "PaperService.Fetchs", in)
	out := new(RespFetchs)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PaperService service

type PaperServiceHandler interface {
	Create(context.Context, *CreatePaper, *Paper) error
	Fetch(context.Context, *PaperID, *Paper) error
	Delete(context.Context, *PaperID, *DeletePaper) error
	Fetchs(context.Context, *ReqFetchs, *RespFetchs) error
}

func RegisterPaperServiceHandler(s server.Server, hdlr PaperServiceHandler, opts ...server.HandlerOption) error {
	type paperService interface {
		Create(ctx context.Context, in *CreatePaper, out *Paper) error
		Fetch(ctx context.Context, in *PaperID, out *Paper) error
		Delete(ctx context.Context, in *PaperID, out *DeletePaper) error
		Fetchs(ctx context.Context, in *ReqFetchs, out *RespFetchs) error
	}
	type PaperService struct {
		paperService
	}
	h := &paperServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&PaperService{h}, opts...))
}

type paperServiceHandler struct {
	PaperServiceHandler
}

func (h *paperServiceHandler) Create(ctx context.Context, in *CreatePaper, out *Paper) error {
	return h.PaperServiceHandler.Create(ctx, in, out)
}

func (h *paperServiceHandler) Fetch(ctx context.Context, in *PaperID, out *Paper) error {
	return h.PaperServiceHandler.Fetch(ctx, in, out)
}

func (h *paperServiceHandler) Delete(ctx context.Context, in *PaperID, out *DeletePaper) error {
	return h.PaperServiceHandler.Delete(ctx, in, out)
}

func (h *paperServiceHandler) Fetchs(ctx context.Context, in *ReqFetchs, out *RespFetchs) error {
	return h.PaperServiceHandler.Fetchs(ctx, in, out)
}