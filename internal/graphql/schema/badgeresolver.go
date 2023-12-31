package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/storemanager"
	graphql "github.com/graph-gophers/graphql-go"
)

type Badge struct {
	ref                           *usersstore.Badge
	CurrentLevel                  int32
	DisplayName                   string
	Id                            string
	LinesOfCodeToNextLevel        int32
	Points                        int32
	ProgressPercentageToNextLevel int32
	TimeCreated                   graphql.Time
	TimeUpdated                   graphql.Time
}

func NewBadge(m *usersstore.Badge) *Badge {
	if m == nil {
		return nil
	}
	return &Badge{
		ref:                           m,
		CurrentLevel:                  int32(m.CurrentLevel),
		DisplayName:                   m.DisplayName,
		Id:                            m.Id,
		LinesOfCodeToNextLevel:        int32(m.LinesOfCodeToNextLevel),
		Points:                        int32(m.Points),
		ProgressPercentageToNextLevel: int32(m.ProgressPercentageToNextLevel),
		TimeCreated:                   graphql.Time{Time: m.TimeCreated},
		TimeUpdated:                   graphql.Time{Time: m.TimeUpdated},
	}
}

type CreateBadgeInput struct {
	CurrentLevel                  int32
	DisplayName                   string
	LinesOfCodeToNextLevel        int32
	Points                        int32
	ProgressPercentageToNextLevel int32
	UserID                        graphql.ID
}

func (n *CreateBadgeInput) ToModel() (*usersstore.Badge, error) {
	if n == nil {
		return nil, nil
	}

	return &usersstore.Badge{
		CurrentLevel:                  int64(n.CurrentLevel),
		DisplayName:                   n.DisplayName,
		LinesOfCodeToNextLevel:        int64(n.LinesOfCodeToNextLevel),
		Points:                        int64(n.Points),
		ProgressPercentageToNextLevel: int64(n.ProgressPercentageToNextLevel),
	}, nil
}
func (n *Badge) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller(NodeTypeBadge, n.Id, true).
		ToGraphqlID()
}

type BadgesConnection struct {
	Edges    []*BadgeEdge
	PageInfo *PageInfo
	filters  *usersstore.BadgeFilters
}

func NewBadgesConnection(
	filters *usersstore.BadgeFilters,
	data []*usersstore.Badge,
	hasNextPage bool,
	hasPreviousPage bool,
	cursors []string,
) *BadgesConnection {
	edges := []*BadgeEdge{}
	for i, d := range data {
		node := NewBadge(d)

		edges = append(edges, &BadgeEdge{
			Node:   node,
			Cursor: cursors[i],
		})
	}
	var startCursor, endCursor *string
	if len(cursors) > 0 {
		startCursor = &cursors[0]
		endCursor = &cursors[len(cursors)-1]
	}
	return &BadgesConnection{
		filters: filters,
		Edges:   edges,
		PageInfo: &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}
}

func (c BadgesConnection) TotalCount(ctx context.Context) (int32, error) {
	count, err := storemanager.FromContext(ctx).UsersStore.CountBadges(ctx, c.filters)
	return int32(count), err
}

type BadgeEdge struct {
	Node   *Badge
	Cursor string
}
