package storiesstore

import (
	"time"

	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type ReactionSortBy uint8

const (
	ReactionSortByDefault ReactionSortBy = iota
	ReactionSortByTimeCreated
)

type Reaction struct {
	StoryID     string            `bson:"story_id"`
	CreatedBy   string            `bson:"created_by,omitempty"`
	Id          string            `bson:"_id"`
	TimeCreated time.Time         `bson:"time_created,omitempty"`
	TimeUpdated time.Time         `bson:"time_updated,omitempty"`
	Ownership   *authorizer.Scope `bson:"ownership,omitempty"`
}

func NewReaction() *Reaction {
	return &Reaction{}
}
func (reaction *Reaction) Get(field string) interface{} {
	switch field {
	case "story_id":
		return reaction.StoryID
	case "created_by":
		return reaction.CreatedBy
	case "_id":
		return reaction.Id
	case "time_created":
		return reaction.TimeCreated
	case "time_updated":
		return reaction.TimeUpdated
	default:
		return nil
	}
}

type ReactionUpdate struct {
	TimeUpdated *time.Time `bson:"time_updated,omitempty"`
}

type ReactionFilters struct {
	Ids       []string
	CreatedBy *string
	Story     *Story
}

func (s ReactionSortBy) String() string {
	switch s {
	case ReactionSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetReactionSortByFromString(s string) ReactionSortBy {
	switch s {
	case "time_created":
		return ReactionSortByTimeCreated
	default:
		return ReactionSortByDefault
	}
}

func (s ReactionSortBy) CursorType() cursor.ValueType {
	switch s {
	case ReactionSortByTimeCreated:
		return cursor.ValueTypeTime
	default:
		return cursor.ValueTypeTime
	}
}
