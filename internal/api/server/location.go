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

func (g grpcImpl) SubmitLocation(ctx2 context.Context, request *types.SubmitLocationRequest) (*types.SubmitLocationResponse, error) {

	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2)
	)

	err := validation.Errors{
		"location":    validation.Validate(request.Location, validation.Required),
		"name":        validation.Validate(request.Location.Name, validation.Required),
		"description": validation.Validate(request.Location.Description, validation.Required),
	}.Filter()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var id int64
	err = db.Transaction(func(dt interface{}) error {

		id, err = db.LocationQ().Insert(&data.Location{
			Name:        request.Location.Name,
			Description: request.Location.Description,
			Photo:       request.Location.Photo,
		})

		return errors.Wrap(err, "failed to insert location")
	}, id)

	if err != nil {
		switch {
		case errors.Is(err, pg.InvalidForeignKeyError):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, pg.ErrUnique):
			return nil, status.Error(codes.AlreadyExists, "location already exists")
		default:
			log.WithError(err).Error("insert location failed")
			return nil, apitypes.ErrInternal
		}
	}

	return &types.SubmitLocationResponse{LocationId: id}, nil

}

func (g grpcImpl) UpdateLocation(ctx2 context.Context, request *types.UpdateLocationRequest) (*types.UpdateLocationResponse, error) {

	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2)
	)

	err := validation.Errors{
		"location":    validation.Validate(request.Location, validation.Required),
		"name":        validation.Validate(request.Location.Name, validation.Required),
		"description": validation.Validate(request.Location.Description, validation.Required),
		"id":          validation.Validate(request.Location.Id, validation.Required, validatePositiveID()),
	}.Filter()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	location, err := apiCtx.DB(ctx2).LocationQ().FilterByID(int64(request.Location.Id)).Get()
	if err != nil {
		log.WithError(err).Error("get location failed")
		return nil, apitypes.ErrInternal
	}

	if location == nil {
		return nil, status.Error(codes.NotFound, "location not found")
	}

	location.Name = request.Location.Name
	location.Description = request.Location.Description
	location.Photo = request.Location.Photo

	err = db.Transaction(func(data interface{}) error {
		return errors.Wrap(db.LocationQ().FilterByID(int64(request.Location.Id)).Update(location), "failed to update location")
	}, location)
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

	return &types.UpdateLocationResponse{}, nil
}

func (g grpcImpl) DeleteLocation(ctx2 context.Context, request *types.DeleteLocationRequest) (*types.DeleteLocationResponse, error) {
	log := apiCtx.Logger(ctx2)

	err := validation.Errors{
		"locationId": validation.Validate(request.LocationId, validation.Required, validatePositiveID()),
	}.Filter()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if location, err := apiCtx.DB(ctx2).LocationQ().FilterByID(int64(request.LocationId)).Get(); err != nil {
		log.WithError(err).Error("get location failed")
		return nil, apitypes.ErrInternal
	} else if location == nil {
		return nil, status.Error(codes.NotFound, "location not found")
	}

	err = apiCtx.DB(ctx2).Transaction(func(data interface{}) error {
		return errors.Wrap(apiCtx.DB(ctx2).LocationQ().Delete(int64(request.LocationId)), "failed to delete location")
	}, request)

	if err != nil {
		log.WithError(err).Error("delete location failed")
		return nil, apitypes.ErrInternal
	}

	return &types.DeleteLocationResponse{}, nil
}

func (g grpcImpl) GetAllLocations(ctx2 context.Context, _ *types.GetAllLocationsRequest) (*types.GetAllLocationsResponse, error) {
	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2).LocationQ()
	)

	var locations []*types.Location

	locs, err := db.GetAll()
	if err != nil {
		log.WithError(err).Error("get locations failed")
		return nil, apitypes.ErrInternal
	}

	for _, loc := range locs {
		locations = append(locations, data.ToLocation(&loc))
	}

	return &types.GetAllLocationsResponse{Locations: locations}, nil
}

func (g grpcImpl) GetLocationById(ctx2 context.Context, request *types.GetLocationByIdRequest) (*types.GetLocationByIdResponse, error) {

	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2).LocationQ()
	)

	err := validation.Errors{
		"locationId": validation.Validate(request.LocationId, validation.Required, validatePositiveID()),
	}.Filter()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	location, err := db.FilterByID(int64(request.LocationId)).Get()
	if err != nil {
		log.WithError(err).Error("get location failed")
		return nil, apitypes.ErrInternal
	}

	if location == nil {
		return nil, status.Error(codes.NotFound, "location not found")
	}

	return &types.GetLocationByIdResponse{Location: data.ToLocation(location)}, nil
}
