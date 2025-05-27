package server

import (
	"context"
	apiCtx "github.com/EduardMikhrin/territorial-fishermans/internal/api/ctx"
	apitypes "github.com/EduardMikhrin/territorial-fishermans/internal/api/types"
	"github.com/EduardMikhrin/territorial-fishermans/internal/data"
	"github.com/EduardMikhrin/territorial-fishermans/internal/data/pg"
	"github.com/EduardMikhrin/territorial-fishermans/internal/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g grpcImpl) SubmitClient(ctx2 context.Context, request *types.SubmitClientRequest) (*types.SubmitClientResponse, error) {

	var (
		log = apiCtx.Logger(ctx2)
		d   = apiCtx.DB(ctx2)
	)
	err := validation.Errors{
		"name":    validation.Validate(request.Client.Name, validation.Required),
		"surname": validation.Validate(request.Client.Surname, validation.Required),
		"contact": validation.Validate(request.Client.Contact, validation.Required),
		"photo":   validation.Validate(request.Client.Photo, validation.Required),
	}.Filter()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var id int64

	err = d.Transaction(func(dt interface{}) error {
		id, err = d.ClientQ().Insert(&data.Client{
			FirstName: request.Client.Name,
			LastName:  request.Client.Surname,
			Contact:   request.Client.Contact,
			Photo:     request.Client.Photo})

		return errors.Wrap(err, "failed to insert client")
	}, request.Client)

	if err != nil {
		switch {
		case errors.Is(err, pg.InvalidForeignKeyError):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, pg.ErrUnique):
			return nil, status.Error(codes.AlreadyExists, "client already exists")
		default:
			log.WithError(err).Error("insert application failed")
			return nil, apitypes.ErrInternal
		}
	}

	return &types.SubmitClientResponse{ClientId: id}, nil
}

func (g grpcImpl) UpdateClient(ctx2 context.Context, request *types.UpdateClientRequest) (*types.UpdateClientResponse, error) {

	var (
		log = apiCtx.Logger(ctx2)
		d   = apiCtx.DB(ctx2)
	)

	err := validation.Errors{
		"id":      validation.Validate(request.Client.Id, validation.Required, validatePositiveID()),
		"name":    validation.Validate(request.Client.Name, validation.Required),
		"surname": validation.Validate(request.Client.Surname, validation.Required),
		"contact": validation.Validate(request.Client.Contact, validation.Required),
		"photo":   validation.Validate(request.Client.Photo, validation.Required),
	}.Filter()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clt, err := apiCtx.DB(ctx2).ClientQ().FilterByID(int64(request.Client.Id)).Get()
	if err != nil {
		log.WithError(err).Error("failed to get client to update client")
		return nil, apitypes.ErrInternal
	}
	if clt == nil {
		return nil, status.Error(codes.NotFound, "client not found")
	}

	clt.FirstName = request.Client.Name
	clt.LastName = request.Client.Surname
	clt.Contact = request.Client.Contact
	clt.Photo = request.Client.Photo

	if err = d.Transaction(func(data interface{}) error {
		return errors.Wrap(d.ClientQ().FilterByID(int64(request.Client.Id)).Update(clt), "failed to update client")
	}, clt); err != nil {
		switch {
		case errors.Is(err, pg.InvalidForeignKeyError):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, pg.ErrUnique):
			return nil, status.Error(codes.AlreadyExists, "client already exists")
		default:
			log.WithError(err).Error("insert application failed")
			return nil, apitypes.ErrInternal
		}
	}

	return &types.UpdateClientResponse{}, nil
}

func (g grpcImpl) DeleteClient(ctx2 context.Context, request *types.DeleteClientRequest) (*types.DeleteClientResponse, error) {

	var (
		log = apiCtx.Logger(ctx2)
		d   = apiCtx.DB(ctx2)
	)
	err := validation.Errors{
		"clientId": validation.Validate(request.ClientId, validation.Required, validatePositiveID()),
	}.Filter()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clt, err := apiCtx.DB(ctx2).ClientQ().FilterByID(int64(request.ClientId)).Get()
	if err != nil {
		log.WithError(err).Error("failed to get client to delete client")
		return nil, apitypes.ErrInternal
	}

	if clt == nil {
		return nil, status.Error(codes.NotFound, "client not found")
	}

	if err = d.Transaction(func(data interface{}) error {
		return errors.Wrap(d.ClientQ().Delete(int64(request.ClientId)), "failed to delete client")
	}, request.ClientId); err != nil {
		log.WithError(err).Error("failed to delete client")
		return nil, apitypes.ErrInternal
	}

	return &types.DeleteClientResponse{}, nil
}

func (g grpcImpl) GetAllClients(ctx2 context.Context, _ *types.GetAllClientsRequest) (*types.GetAllClientsResponse, error) {
	var (
		log = apiCtx.Logger(ctx2)
		d   = apiCtx.DB(ctx2).ClientQ()
	)

	clients, err := d.GetAll()
	if err != nil {
		log.WithError(err).Error("failed to get clients")
		return nil, apitypes.ErrInternal
	}

	var clts []*types.Client
	for _, client := range clients {
		clts = append(clts, data.ToClient(&client))
	}

	return &types.GetAllClientsResponse{Clients: clts}, nil
}

func (g grpcImpl) GetClientById(ctx2 context.Context, request *types.GetClientByIdRequest) (*types.GetClientByIdResponse, error) {

	var (
		log = apiCtx.Logger(ctx2)
		d   = apiCtx.DB(ctx2).ClientQ().FilterByID(int64(request.ClientId))
	)
	err := validation.Errors{
		"clientId": validation.Validate(request.ClientId, validation.Required, validatePositiveID()),
	}.Filter()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	client, err := d.Get()
	if err != nil {
		log.WithError(err).Error("failed to get client to get client")
		return nil, apitypes.ErrInternal
	}

	if client == nil {
		return nil, status.Error(codes.NotFound, "client not found")
	}

	return &types.GetClientByIdResponse{Client: data.ToClient(client)}, nil
}
