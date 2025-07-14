package service

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type TagQueryService struct {
	tagRepository port.TagRepository
}

func NewTagQueryService(
	tagRepository port.TagRepository,
) *TagQueryService {
	if tagRepository == nil {
		panic("nil TagRepository")
	}
	return &TagQueryService{
		tagRepository: tagRepository,
	}
}

type TagQueryOutput struct {
	Id   uuid.UUID
	Name string
}

// タグを一覧取得
func (s *TagQueryService) FindAllTags() ([]TagQueryOutput, error) {
	tags, err := s.tagRepository.FindAll()
	if err != nil {
		return nil, err
	}
	outputs := []TagQueryOutput{}
	for _, tag := range tags {
		outputs = append(outputs, TagQueryOutput{
			Id:   tag.Id(),
			Name: tag.Name(),
		})
	}
	return outputs, nil
}
