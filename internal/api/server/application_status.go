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

func (g grpcImpl) SubmitApplicationStatus(ctx2 context.Context, request *types.SubmitApplicationStatusRequest) (*types.SubmitApplicationStatusResponse, error) {
	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2).ApplicationStatusQ()
	)

	if err := validateStatus(request.Name); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := db.Insert(&data.ApplicationStatus{
		Name: request.Name,
	})

	if err != nil {
		switch {
		case errors.Is(err, pg.InvalidForeignKeyError):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, pg.ErrUnique):
			return nil, status.Error(codes.AlreadyExists, "application_status already exists")
		default:
			log.WithError(err).Error("insert application_status failed")
			return nil, apitypes.ErrInternal
		}
	}

	return &types.SubmitApplicationStatusResponse{Id: id}, nil
}

func (g grpcImpl) UpdateApplicationStatus(ctx2 context.Context, request *types.UpdateApplicationStatusRequest) (*types.UpdateApplicationStatusResponse, error) {
	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2).ApplicationStatusQ().FilterByID(int64(request.ApplicationStatus.Id))
	)

	if err := validateStatus(request.ApplicationStatus.Name); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appStatus, err := db.Get()
	if err != nil {
		log.WithError(err).Error("get application_status failed")
		return nil, apitypes.ErrInternal
	}

	if appStatus == nil {
		return nil, status.Error(codes.NotFound, "application_status not found")
	}

	appStatus.Name = request.ApplicationStatus.Name

	if err := db.Update(appStatus); err != nil {
		log.WithError(err).Error("update application_status failed")
		switch {
		case errors.Is(err, pg.InvalidForeignKeyError):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, pg.ErrUnique):
			return nil, status.Error(codes.AlreadyExists, "application_status already exists")
		default:
			log.WithError(err).Error("insert application_status failed")
			return nil, apitypes.ErrInternal
		}
	}

	return &types.UpdateApplicationStatusResponse{}, nil
}

func (g grpcImpl) DeleteApplicationStatus(ctx2 context.Context, request *types.DeleteApplicationStatusRequest) (*types.DeleteApplicationStatusResponse, error) {
	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2).ApplicationStatusQ().FilterByID(int64(request.Id))
	)

	err := validation.Errors{
		"id": validation.Validate(request.Id, validation.Required, validatePositiveID()),
	}.Filter()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appStatus, err := db.Get()
	if err != nil {
		log.WithError(err).Error("get application_status failed")
		return nil, apitypes.ErrInternal
	}

	if appStatus == nil {
		return nil, status.Error(codes.NotFound, "application_status not found")
	}

	if err = db.Delete(int64(request.Id)); err != nil {
		log.WithError(err).Error("delete application_status failed")
		return nil, apitypes.ErrInternal
	}

	return &types.DeleteApplicationStatusResponse{}, nil
}

func (g grpcImpl) GetApplicationStatus(ctx2 context.Context, request *types.GetApplicationStatusRequest) (*types.GetApplicationStatusResponse, error) {
	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2).ApplicationStatusQ().FilterByID(int64(request.Id))
	)

	err := validation.Errors{
		"id": validation.Validate(request.Id, validation.Required, validatePositiveID()),
	}.Filter()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	appStatus, err := db.Get()
	if err != nil {
		log.WithError(err).Error("get application_status failed")
		return nil, apitypes.ErrInternal
	}
	if appStatus == nil {
		return nil, status.Error(codes.NotFound, "application_status not found")
	}

	return &types.GetApplicationStatusResponse{ApplicationStatus: data.ToApplicationStatus(*appStatus)}, nil
}

func (g grpcImpl) GetAllApplicationStatuses(ctx2 context.Context, request *types.GetAllApplicationsStatusRequest) (*types.GetAllApplicationsStatusResponse, error) {
	var (
		log = apiCtx.Logger(ctx2)
		db  = apiCtx.DB(ctx2).ApplicationStatusQ()
	)

	appStatuses, err := db.GetAll()
	if err != nil {
		log.WithError(err).Error("get application_status failed")
		return nil, apitypes.ErrInternal
	}

	var statuses []*types.ApplicationStatus
	for _, appStatus := range appStatuses {
		statuses = append(statuses, data.ToApplicationStatus(appStatus))
	}

	return &types.GetAllApplicationsStatusResponse{ApplicationStatuses: statuses}, nil
}

func validateStatus(status string) error {
	return validation.Errors{
		"name": validation.Validate(status, validation.Required),
	}.Filter()
}
