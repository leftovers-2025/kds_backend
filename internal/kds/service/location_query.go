package service

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type LocationQueryService struct {
	locationReository port.LocationRepository
}

func NewLocationQueryService(
	locationReository port.LocationRepository,
) *LocationQueryService {
	if locationReository == nil {
		panic("nil LocationRepository")
	}
	return &LocationQueryService{
		locationReository: locationReository,
	}
}

type LocationQueryOutput struct {
	Id   uuid.UUID
	Name string
}

// ロケーションを一覧取得
func (s *LocationQueryService) GetAllLocations() ([]LocationQueryOutput, error) {
	locations, err := s.locationReository.FindAll()
	if err != nil {
		return nil, err
	}
	outputs := []LocationQueryOutput{}
	for _, location := range locations {
		outputs = append(outputs, LocationQueryOutput{
			Id:   location.Id(),
			Name: location.Name(),
		})
	}
	return outputs, nil
}
