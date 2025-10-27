package object

import (
	"fmt"

	"github.com/sunshineOfficial/golib/goctx"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddObject(ctx goctx.Context, request AddObjectRequest) (Object, error) {
	obj, err := s.repository.AddObject(ctx, request)
	if err != nil {
		return Object{}, fmt.Errorf("add object to repository: %w", err)
	}

	return obj, nil
}

func (s *Service) GetObjectByID(ctx goctx.Context, id int) (Object, error) {
	obj, err := s.repository.GetObjectByID(ctx, id)
	if err != nil {
		return Object{}, fmt.Errorf("get object from repository: %w", err)
	}

	return obj, nil
}

func (s *Service) GetObjectByDeviceID(ctx goctx.Context, deviceID int) (Object, error) {
	obj, err := s.repository.GetObjectByDeviceID(ctx, deviceID)
	if err != nil {
		return Object{}, fmt.Errorf("get object from repository: %w", err)
	}

	return obj, nil
}

func (s *Service) GetObjectBySealID(ctx goctx.Context, sealID int) (Object, error) {
	obj, err := s.repository.GetObjectBySealID(ctx, sealID)
	if err != nil {
		return Object{}, fmt.Errorf("get object from repository: %w", err)
	}

	return obj, nil
}
