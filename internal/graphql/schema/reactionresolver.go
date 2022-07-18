package schema

import (
	"context"
	"errors"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/firstcontributions/backend/pkg/cursor"
	graphql "github.com/graph-gophers/graphql-go"
)

type Reaction struct {
	ref         *storiesstore.Reaction
	createdBy   string
	Id          string
	TimeCreated graphql.Time
	TimeUpdated graphql.Time
}

func NewReaction(m *storiesstore.Reaction) *Reaction {
	if m == nil {
		return nil
	}
	return &Reaction{
		ref:         m,
		createdBy:   m.CreatedBy,
		Id:          m.Id,
		TimeCreated: graphql.Time{Time: m.TimeCreated},
		TimeUpdated: graphql.Time{Time: m.TimeUpdated},
	}
}
func (n *Reaction) CreatedBy(ctx context.Context) (*User, error) {

	data, err := storemanager.FromContext(ctx).UsersStore.GetUserByID(ctx, n.createdBy)
	if err != nil {
		return nil, err
	}
	return NewUser(data), nil
}

type CreateReactionInput struct {
	StoryID graphql.ID
}

func (n *CreateReactionInput) ToModel() (*storiesstore.Reaction, error) {
	if n == nil {
		return nil, nil
	}
	storyID, err := ParseGraphqlID(n.StoryID)
	if err != nil {
		return nil, err
	}

	return &storiesstore.Reaction{
		StoryID: storyID.ID,
	}, nil
}

type UpdateReactionInput struct {
	ID graphql.ID
}

func (n *UpdateReactionInput) ToModel() *storiesstore.ReactionUpdate {
	if n == nil {
		return nil
	}
	return &storiesstore.ReactionUpdate{}
}
func (n *Reaction) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("reaction", n.Id).
		ToGraphqlID()
}

type ReactionsConnection struct {
	Edges    []*ReactionEdge
	PageInfo *PageInfo
	filters  *storiesstore.ReactionFilters
}

func NewReactionsConnection(
	filters *storiesstore.ReactionFilters,
	data []*storiesstore.Reaction,
	hasNextPage bool,
	hasPreviousPage bool,
	firstCursor *string,
	lastCursor *string,
) *ReactionsConnection {
	edges := []*ReactionEdge{}
	for _, d := range data {
		node := NewReaction(d)

		edges = append(edges, &ReactionEdge{
			Node:   node,
			Cursor: cursor.NewCursor(d.Id, "time_created", d.TimeCreated).String(),
		})
	}
	return &ReactionsConnection{
		filters: filters,
		Edges:   edges,
		PageInfo: &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     firstCursor,
			EndCursor:       lastCursor,
		},
	}
}

func (c ReactionsConnection) TotalCount(ctx context.Context) (int32, error) {
	count, err := storemanager.FromContext(ctx).StoriesStore.CountReactions(ctx, c.filters)
	return int32(count), err
}
func (c ReactionsConnection) HasViewerAssociation(ctx context.Context) (bool, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return false, errors.New("Unauthorized")
	}
	userID := session.UserID()

	newFilter := *c.filters
	newFilter.CreatedBy = &userID

	data, err := storemanager.FromContext(ctx).StoriesStore.GetOneReaction(ctx, c.filters)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

type ReactionEdge struct {
	Node   *Reaction
	Cursor string
}
