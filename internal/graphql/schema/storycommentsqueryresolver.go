package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/storemanager"
)

type StoryCommentsInput struct {
	First     *int32
	Last      *int32
	After     *string
	Before    *string
	SortOrder *string
	SortBy    *string
}

func (n *Story) Comments(ctx context.Context, in *StoryCommentsInput) (*CommentsConnection, error) {
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

	filters := &storiesstore.CommentFilters{
		Story: n.ref,
	}
	sortByStr := ""
	if in.SortBy != nil {
		sortByStr = *in.SortBy
	}
	data, hasNextPage, hasPreviousPage, cursors, err := store.StoriesStore.GetComments(
		ctx,
		filters,
		in.After,
		in.Before,
		first,
		last,
		storiesstore.GetCommentSortByFromString(sortByStr),
		in.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	return NewCommentsConnection(filters, data, hasNextPage, hasPreviousPage, cursors), nil
}
