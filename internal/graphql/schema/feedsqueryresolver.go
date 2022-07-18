package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/storemanager"
)

type FeedsInput struct {
	First     *int32
	Last      *int32
	After     *string
	Before    *string
	SortBy    *string
	SortOrder *string
}

func (r *Resolver) Feeds(ctx context.Context, in *FeedsInput) (*StoriesConnection, error) {
	var first, last *int64
	if in.First != nil {
		tmp := int64(*in.First)
		first = &tmp
	}
	if in.Last != nil {
		tmp := int64(*in.Last)
		last = &tmp
	}
	store := storemanager.FromContext(ctx)

	filters := &storiesstore.StoryFilters{}
	sortByStr := ""
	if in.SortBy != nil {
		sortByStr = *in.SortBy
	}
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := store.StoriesStore.GetStories(
		ctx,
		filters,
		in.After,
		in.Before,
		first,
		last,
		storiesstore.GetStorySortByFromString(sortByStr),
		in.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	return NewStoriesConnection(filters, data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}
