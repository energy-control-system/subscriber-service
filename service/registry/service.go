package registry

import (
	"fmt"
	"mime/multipart"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/golog"
	"github.com/xuri/excelize/v2"
)

type Service struct {
	subscriberRepository SubscriberRepository
	objectRepository     ObjectRepository
	contractRepository   ContractRepository
}

func NewService(subscriberRepository SubscriberRepository, objectRepository ObjectRepository, contractRepository ContractRepository) *Service {
	return &Service{
		subscriberRepository: subscriberRepository,
		objectRepository:     objectRepository,
		contractRepository:   contractRepository,
	}
}

func (s *Service) Parse(ctx goctx.Context, log golog.Logger, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("open file %s: %w", fileHeader.Filename, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Errorf("close file %s: %v", fileHeader.Filename, closeErr)
		}
	}()

	registry, err := excelize.OpenReader(file)
	if err != nil {
		return fmt.Errorf("open registry %s: %w", fileHeader.Filename, err)
	}
	defer func() {
		if closeErr := registry.Close(); closeErr != nil {
			log.Errorf("close registry %s: %v", fileHeader.Filename, closeErr)
		}
	}()

	subs, err := parseSubscribers(registry)
	if err != nil {
		return fmt.Errorf("parse subscribers: %w", err)
	}

	objects, err := parseObjects(registry)
	if err != nil {
		return fmt.Errorf("parse objects: %w", err)
	}

	devices, err := parseDevices(registry)
	if err != nil {
		return fmt.Errorf("parse devices: %w", err)
	}

	seals, err := parseSeals(registry)
	if err != nil {
		return fmt.Errorf("parse seals: %w", err)
	}

	contracts, err := parseContracts(registry)
	if err != nil {
		return fmt.Errorf("parse contracts: %w", err)
	}

	if err = s.subscriberRepository.UpsertSubscribers(ctx, subs); err != nil {
		return fmt.Errorf("upsert subscribers: %w", err)
	}

	if err = s.objectRepository.UpsertObjects(ctx, objects); err != nil {
		return fmt.Errorf("upsert objects: %w", err)
	}

	if err = s.objectRepository.UpsertDevices(ctx, devices); err != nil {
		return fmt.Errorf("upsert devices: %w", err)
	}

	if err = s.objectRepository.UpsertSeals(ctx, seals); err != nil {
		return fmt.Errorf("upsert seals: %w", err)
	}

	if err = s.contractRepository.UpsertContracts(ctx, contracts); err != nil {
		return fmt.Errorf("upsert contracts: %w", err)
	}

	return nil
}
