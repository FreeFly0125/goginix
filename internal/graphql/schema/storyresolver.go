package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/firstcontributions/backend/pkg/cursor"
	graphql "github.com/graph-gophers/graphql-go"
)

type Story struct {
	ref             *storiesstore.Story
	AbstractContent string
	ContentJson     string
	createdBy       string
	Id              string
	Thumbnail       string
	TimeCreated     graphql.Time
	TimeUpdated     graphql.Time
	Title           string
	UrlSuffix       string
}

func NewStory(m *storiesstore.Story) *Story {
	if m == nil {
		return nil
	}
	return &Story{
		ref:             m,
		AbstractContent: m.AbstractContent,
		ContentJson:     m.ContentJson,
		createdBy:       m.CreatedBy,
		Id:              m.Id,
		Thumbnail:       m.Thumbnail,
		TimeCreated:     graphql.Time{Time: m.TimeCreated},
		TimeUpdated:     graphql.Time{Time: m.TimeUpdated},
		Title:           m.Title,
		UrlSuffix:       m.UrlSuffix,
	}
}
func (n *Story) CreatedBy(ctx context.Context) (*User, error) {

	data, err := storemanager.FromContext(ctx).UsersStore.GetUserByID(ctx, n.createdBy)
	if err != nil {
		return nil, err
	}
	return NewUser(data), nil
}

type CreateStoryInput struct {
	AbstractContent string
	ContentJson     string
	Thumbnail       string
	Title           string
	UrlSuffix       string
	UserID          graphql.ID
}

func (n *CreateStoryInput) ToModel() (*storiesstore.Story, error) {
	if n == nil {
		return nil, nil
	}

	return &storiesstore.Story{
		AbstractContent: n.AbstractContent,
		ContentJson:     n.ContentJson,
		Thumbnail:       n.Thumbnail,
		Title:           n.Title,
		UrlSuffix:       n.UrlSuffix,
	}, nil
}

type UpdateStoryInput struct {
	ID graphql.ID
}

func (n *UpdateStoryInput) ToModel() *storiesstore.StoryUpdate {
	if n == nil {
		return nil
	}
	return &storiesstore.StoryUpdate{}
}
func (n *Story) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("story", n.Id).
		ToGraphqlID()
}

type StoriesConnection struct {
	Edges    []*StoryEdge
	PageInfo *PageInfo
}

func NewStoriesConnection(
	data []*storiesstore.Story,
	hasNextPage bool,
	hasPreviousPage bool,
	firstCursor *string,
	lastCursor *string,
) *StoriesConnection {
	edges := []*StoryEdge{}
	for _, d := range data {
		node := NewStory(d)

		edges = append(edges, &StoryEdge{
			Node:   node,
			Cursor: cursor.NewCursor(d.Id, d.TimeCreated).String(),
		})
	}
	return &StoriesConnection{
		Edges: edges,
		PageInfo: &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     firstCursor,
			EndCursor:       lastCursor,
		},
	}
}

type StoryEdge struct {
	Node   *Story
	Cursor string
}
