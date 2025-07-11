package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	POST_DESC_MIN = 1
	POST_DESC_MAX = 128
)

var (
	ErrPostIdRequired               = errors.New("post id is required")
	ErrPostUserIdRequired           = errors.New("post userId is required")
	ErrPostDescriptionRequired      = errors.New("post userId is required")
	ErrPostDescriptionInvalidLength = errors.New("post length is invalid")
	ErrPostTagInvalid               = errors.New("post tags is invalid")
	ErrPostImagesInvalid            = errors.New("post images is invalid")
	ErrPostTimeZero                 = errors.New("post time is zero")
)

type Post struct {
	id          uuid.UUID
	userId      uuid.UUID
	location    Location
	description string
	tags        []Tag
	images      []Image
	createdAt   time.Time
	updatedAt   time.Time
}

func NewPost(
	id, userId uuid.UUID,
	location Location,
	description string,
	tags []Tag,
	images []Image,
	createdAt,
	updatedAt time.Time,
) (*Post, error) {
	if id == uuid.Nil {
		return nil, ErrPostIdRequired
	}
	if userId == uuid.Nil {
		return nil, ErrPostUserIdRequired
	}
	if description == "" {
		return nil, ErrPostDescriptionRequired
	}
	if len(description) < POST_DESC_MIN || len(description) > POST_DESC_MAX {
		return nil, ErrPostDescriptionInvalidLength
	}
	if tags == nil {
		return nil, ErrPostTagInvalid
	}
	if images == nil {
		return nil, ErrPostImagesInvalid
	}
	if createdAt.IsZero() || updatedAt.IsZero() {
		return nil, ErrPostTimeZero
	}
	return &Post{
		id:          id,
		userId:      userId,
		location:    location,
		description: description,
		tags:        tags,
		images:      images,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

func (p *Post) Id() uuid.UUID {
	return p.id
}

func (p *Post) UserId() uuid.UUID {
	return p.userId
}

func (p *Post) Location() Location {
	return p.location
}

func (p *Post) Description() string {
	return p.description
}

func (p *Post) Tags() []Tag {
	return p.tags
}

func (p *Post) Images() []Image {
	return p.images
}

func (p *Post) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Post) UpdatedAt() time.Time {
	return p.updatedAt
}
