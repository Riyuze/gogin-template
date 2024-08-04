package service

import (
	"context"
	"gogin-template/baselib/dto"
	"gogin-template/baselib/exception"
	"gogin-template/baselib/helper"
	"gogin-template/bootstrap"
	"gogin-template/internal/model"
	"gogin-template/internal/repository"
	"gogin-template/internal/viewmodel"
)

type SampleService interface {
	GetSamples(c context.Context, requestVM *viewmodel.SampleRqViewModel, dtoPage dto.PageRequest) (*[]viewmodel.SampleRsViewModel, *dto.PageInfo, error)
	GetSample(c context.Context, requestVM *viewmodel.SampleRqViewModel) (*viewmodel.SampleRsViewModel, error)
	GetSampleVersions(c context.Context, requestVM *viewmodel.SampleVersionRqViewModel) (*[]viewmodel.SampleVersionRsViewModel, error)
	GetSampleVersion(c context.Context, requestVM *viewmodel.SampleVersionRqViewModel) (*viewmodel.SampleVersionRsViewModel, error)
	SetSample(c context.Context, action string, requestVM *viewmodel.SampleRqViewModel) error
	SetSampleVersions(c context.Context, requestVM *[]viewmodel.SampleVersionRqViewModel) error
	SetSampleVersion(c context.Context, action string, requestVM *viewmodel.SampleVersionRqViewModel) error
}

type SampleServiceImpl struct {
	repository repository.SampleRepository
	cfg        *bootstrap.Container
}

func NewSampleService(repository repository.SampleRepository, cfg *bootstrap.Container) SampleService {
	return &SampleServiceImpl{repository: repository, cfg: cfg}
}

func (s *SampleServiceImpl) GetSamples(c context.Context, requestVM *viewmodel.SampleRqViewModel, dtoPage dto.PageRequest) (*[]viewmodel.SampleRsViewModel, *dto.PageInfo, error) {
	// Convert View Model to Model
	requestM := &model.SampleQueryModel{}

	s.cfg.CopyStruct(requestVM, requestM)

	// Process
	response, pageInfo, err := s.repository.GetSamples(c, requestM, dtoPage)
	if err != nil {
		return nil, nil, helper.CatchErr(err)
	}

	// Convert To View Model
	responseVM := &[]viewmodel.SampleRsViewModel{}
	s.cfg.CopyStruct(response, responseVM)

	return responseVM, pageInfo, nil
}

func (s *SampleServiceImpl) GetSample(c context.Context, requestVM *viewmodel.SampleRqViewModel) (*viewmodel.SampleRsViewModel, error) {
	// Convert View Model to Model
	requestM := &model.SampleQueryModel{}

	s.cfg.CopyStruct(requestVM, requestM)

	// Process
	response, err := s.repository.GetSample(c, requestM)
	if err != nil {
		return nil, helper.CatchErr(err)
	}

	if response == nil {
		return nil, exception.NotFoundException("404", "Not Found")
	}

	// Convert To View Model
	responseVM := &viewmodel.SampleRsViewModel{}
	s.cfg.CopyStruct(response, responseVM)

	return responseVM, nil
}

func (s *SampleServiceImpl) GetSampleVersions(c context.Context, requestVM *viewmodel.SampleVersionRqViewModel) (*[]viewmodel.SampleVersionRsViewModel, error) {
	// Convert View Model to Model
	requestM := &model.SampleVersionQueryModel{}

	s.cfg.CopyStruct(requestM, requestVM)

	// Process
	response, err := s.repository.GetSampleVersions(c, requestM)
	if err != nil {
		return nil, helper.CatchErr(err)
	}

	// Convert to View Model
	responseVM := &[]viewmodel.SampleVersionRsViewModel{}
	s.cfg.CopyStruct(response, responseVM)

	return responseVM, nil
}

func (s *SampleServiceImpl) GetSampleVersion(c context.Context, requestVM *viewmodel.SampleVersionRqViewModel) (*viewmodel.SampleVersionRsViewModel, error) {
	// Convert View Model to Model
	requestM := &model.SampleVersionQueryModel{}

	s.cfg.CopyStruct(requestM, requestVM)

	// Process
	response, err := s.repository.GetSampleVersion(c, requestM)
	if err != nil {
		return nil, helper.CatchErr(err)
	}

	// Convert to View Model
	responseVM := &viewmodel.SampleVersionRsViewModel{}
	s.cfg.CopyStruct(response, responseVM)

	return responseVM, nil
}

func (s *SampleServiceImpl) SetSample(c context.Context, action string, requestVM *viewmodel.SampleRqViewModel) error {
	// Convert View Model to Model
	requestM := &model.SampleModel{}

	s.cfg.CopyStruct(requestVM, requestM)

	// Process
	err := s.repository.SetSample(c, action, requestM)
	if err != nil {
		return helper.CatchErr(err)
	}

	return nil
}

func (s *SampleServiceImpl) SetSampleVersions(c context.Context, requestVM *[]viewmodel.SampleVersionRqViewModel) error {
	// Convert View Model to Model
	requestM := &[]model.SampleVersionModel{}

	s.cfg.CopyStruct(requestVM, requestM)

	// Process
	err := s.repository.SetSampleVersions(c, requestM)
	if err != nil {
		return helper.CatchErr(err)
	}

	return nil
}

func (s *SampleServiceImpl) SetSampleVersion(c context.Context, action string, requestVM *viewmodel.SampleVersionRqViewModel) error {
	// Convert View Model to Model
	requestM := &model.SampleVersionModel{}

	s.cfg.CopyStruct(requestVM, requestM)

	// Process
	err := s.repository.SetSampleVersion(c, action, requestM)
	if err != nil {
		return helper.CatchErr(err)
	}

	return nil
}
