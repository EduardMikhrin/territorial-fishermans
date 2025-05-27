package server

import (
	"context"
	"database/sql"
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

// getApplicationByID returns (nil, nil) if the row is not found.
func getApplicationByID(q data.ApplicationQ) (*data.Application, error) {
	app, err := q.Get()
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "query application by id")
	}
	return app, nil
}

func (g grpcImpl) SubmitApplication(ctx context.Context, req *types.SubmitApplicationRequest) (*types.SubmitApplicationResponse, error) {
	err := validation.Errors{
		"application":  validation.Validate(req.Application, validation.Required),
		"client_id":    validation.Validate(req.Application.ClientId, validation.Required, validatePositiveID()),
		"fishing_date": validation.Validate(req.Application.FishingDate, validation.Required),
		"location_id":  validation.Validate(req.Application.LocationId, validation.Required, validatePositiveID()),
		"status_id":    validation.Validate(req.Application.StatusId, validation.Required, validatePositiveID()),
	}.Filter()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	repo := apiCtx.DB(ctx)
	log := apiCtx.Logger(ctx)

	app := data.Application{
		FishingDate: req.Application.FishingDate.AsTime(),
		ClientID:    req.Application.ClientId,
		LocationId:  req.Application.LocationId,
		StatusId:    req.Application.StatusId,
	}

	var id int64
	err = repo.Transaction(func(data interface{}) error {
		id, err = repo.ApplicationQ().Insert(&app)
		return errors.Wrap(err, "insert failed")

	}, app)

	if err != nil {
		switch {
		case errors.Is(err, pg.InvalidForeignKeyError):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, pg.ErrUnique):
			return nil, status.Error(codes.AlreadyExists, "application already exists")
		default:
			log.WithError(err).Error("insert application failed")
			return nil, apitypes.ErrInternal
		}
	}

	return &types.SubmitApplicationResponse{ApplicationId: id}, nil
}

// UpdateApplication updates mutable fields of an existing application.
func (g grpcImpl) UpdateApplication(ctx context.Context, req *types.UpdateApplicationRequest) (*types.UpdateApplicationResponse, error) {
	err := validation.Errors{
		"id":           validation.Validate(req.Application.Id, validation.Required, validatePositiveID()),
		"application":  validation.Validate(req.Application, validation.Required),
		"client_id":    validation.Validate(req.Application.ClientId, validation.Required, validatePositiveID()),
		"fishing_date": validation.Validate(req.Application.FishingDate, validation.Required),
		"location_id":  validation.Validate(req.Application.LocationId, validation.Required, validatePositiveID()),
		"status_id":    validation.Validate(req.Application.StatusId, validation.Required, validatePositiveID()),
	}.Filter()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	repo := apiCtx.DB(ctx)
	log := apiCtx.Logger(ctx)

	app, err := getApplicationByID(apiCtx.DB(ctx).ApplicationQ().FilterByID(int64(req.Application.Id)))
	if err != nil {
		log.WithError(err).Errorf("query application %d", req.Application.Id)
		return nil, apitypes.ErrInternal
	}
	if app == nil {
		return nil, status.Error(codes.NotFound, "application not found")
	}

	app.FishingDate = req.Application.FishingDate.AsTime()
	app.ClientID = req.Application.ClientId
	app.LocationId = req.Application.LocationId
	app.StatusId = req.Application.StatusId

	err = repo.Transaction(func(data interface{}) error {
		return errors.Wrap(repo.ApplicationQ().FilterByID(int64(req.Application.Id)).Update(app), "update application failed")
	}, app)

	if err != nil {
		if errors.Is(err, pg.InvalidForeignKeyError) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		switch {
		case errors.Is(err, pg.InvalidForeignKeyError):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, pg.ErrUnique):
			return nil, status.Error(codes.AlreadyExists, "application already exists")
		default:
			log.WithError(err).Error("insert application failed")
			return nil, apitypes.ErrInternal
		}
	}

	return &types.UpdateApplicationResponse{}, nil
}

// DeleteApplication removes a row by primary key.
func (g grpcImpl) DeleteApplication(ctx context.Context, req *types.DeleteApplicationRequest) (*types.DeleteApplicationResponse, error) {
	err := validation.Errors{
		"applicationId": validation.Validate(req.ApplicationId, validation.Required, validatePositiveID()),
	}.Filter()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	repo := apiCtx.DB(ctx)
	log := apiCtx.Logger(ctx)

	if application, err := getApplicationByID(apiCtx.DB(ctx).ApplicationQ().FilterByID(int64(req.ApplicationId))); err != nil {
		log.WithError(err).Errorf("query application %d", req.ApplicationId)
		return nil, apitypes.ErrInternal
	} else if application == nil {
		return nil, status.Error(codes.NotFound, "application not found")
	}

	err = repo.Transaction(func(data interface{}) error {
		return errors.Wrap(repo.ApplicationQ().Delete(req.ApplicationId), "update application failed")
	}, req.ApplicationId)

	if err != nil {
		log.WithError(err).Error("delete application failed")
		return nil, apitypes.ErrInternal
	}

	return &types.DeleteApplicationResponse{}, nil
}

// GetAllApplications returns all rows.
func (g grpcImpl) GetAllApplications(ctx context.Context, _ *types.ListApplicationsRequest) (*types.ListApplicationsResponse, error) {
	repo := apiCtx.DB(ctx).ApplicationQ()
	log := apiCtx.Logger(ctx)

	list, err := repo.GetAll()
	if err != nil {
		log.WithError(err).Error("get all applications failed")
		return nil, apitypes.ErrInternal
	}

	apps := make([]*types.Application, 0, len(list))
	for i := range list {
		apps = append(apps, data.ToApplication(&list[i]))
	}

	return &types.ListApplicationsResponse{Applications: apps}, nil
}

// GetApplicationById returns a single application.
func (g grpcImpl) GetApplicationById(ctx context.Context, req *types.GetApplicationByIdRequest) (*types.GetApplicationByIdResponse, error) {
	err := validation.Errors{
		"applicationId": validation.Validate(req.ApplicationId, validation.Required, validatePositiveID()),
	}.Filter()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log := apiCtx.Logger(ctx)

	app, err := getApplicationByID(apiCtx.DB(ctx).ApplicationQ().FilterByID(int64(req.ApplicationId)))
	if err != nil {
		log.WithError(err).Error("get application failed")
		return nil, apitypes.ErrInternal
	}
	if app == nil {
		return nil, status.Error(codes.NotFound, "application not found")
	}

	return &types.GetApplicationByIdResponse{Application: data.ToApplication(app)}, nil
}
