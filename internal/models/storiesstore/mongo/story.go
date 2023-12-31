package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/models/utils"
	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func storyFiltersToQuery(filters *storiesstore.StoryFilters) *mongoqb.QueryBuilder {
	qb := mongoqb.NewQueryBuilder()
	if len(filters.Ids) > 0 {
		qb.In("_id", filters.Ids)
	}
	if filters.CreatedBy != nil {
		qb.Eq("created_by", filters.CreatedBy)
	}
	if filters.User != nil {
		qb.Eq("user_id", filters.User.Id)
	}
	return qb
}
func (s *StoriesStore) CreateStory(ctx context.Context, story *storiesstore.Story, ownership *authorizer.Scope) (*storiesstore.Story, error) {
	now := time.Now()
	story.TimeCreated = now
	story.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	if dproc, ok := interface{}(story).(utils.DataProcessor); ok {
		if err := dproc.Process(ctx); err != nil {
			return nil, err
		}
	}
	story.Id = uuid.String()
	story.Ownership = ownership
	if _, err := s.getCollection(CollectionStories).InsertOne(ctx, story); err != nil {
		return nil, err
	}
	return story, nil
}

func (s *StoriesStore) GetStoryByID(ctx context.Context, id string) (*storiesstore.Story, error) {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	var story storiesstore.Story
	if err := s.getCollection(CollectionStories).FindOne(ctx, qb.Build()).Decode(&story); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &story, nil
}

func (s *StoriesStore) GetOneStory(ctx context.Context, filters *storiesstore.StoryFilters) (*storiesstore.Story, error) {
	qb := storyFiltersToQuery(filters)
	var story storiesstore.Story
	if err := s.getCollection(CollectionStories).FindOne(ctx, qb.Build()).Decode(&story); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &story, nil
}

func (s *StoriesStore) CountStories(ctx context.Context, filters *storiesstore.StoryFilters) (
	int64,
	error,
) {
	qb := storyFiltersToQuery(filters)

	count, err := s.getCollection(CollectionStories).CountDocuments(ctx, qb.Build())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *StoriesStore) GetStories(
	ctx context.Context,
	filters *storiesstore.StoryFilters,
	after *string,
	before *string,
	first *int64,
	last *int64,
	sortBy storiesstore.StorySortBy,
	sortOrder *string,
) (
	[]*storiesstore.Story,
	bool,
	bool,
	[]string,
	error,
) {
	qb := storyFiltersToQuery(filters)
	reqSortOrder := utils.GetSortOrderFromString(sortOrder)
	limit, paginationSortOrder, cursorStr, err := utils.GetLimitAndSortOrderAndCursor(first, last, after, before)
	if err != nil {
		return nil, false, false, nil, err
	}

	effectiveSortOrder := reqSortOrder * paginationSortOrder

	var c *cursor.Cursor
	if cursorStr != nil {
		c, err = cursor.FromString(*cursorStr)
		if err != nil {
			return nil, false, false, nil, err
		}
		if c != nil {
			if effectiveSortOrder == 1 {
				qb.Or(

					mongoqb.NewQueryBuilder().
						Eq(storiesstore.StorySortBy(c.SortBy).String(), c.OffsetValue).
						Gt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Gt(storiesstore.StorySortBy(c.SortBy).String(), c.OffsetValue),
				)
			} else {
				qb.Or(
					mongoqb.NewQueryBuilder().
						Eq(storiesstore.StorySortBy(c.SortBy).String(), c.OffsetValue).
						Lt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Lt(storiesstore.StorySortBy(c.SortBy).String(), c.OffsetValue),
				)
			}
		}
	}
	// incrementing limit by 2 to check if next, prev elements are present
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  utils.GetSortOrder(sortBy.String(), effectiveSortOrder),
	}

	var hasNextPage, hasPreviousPage bool

	var stories []*storiesstore.Story
	mongoCursor, err := s.getCollection(CollectionStories).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	err = mongoCursor.All(ctx, &stories)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	count := len(stories)
	if count == 0 {
		return stories, hasNextPage, hasPreviousPage, nil, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && stories[0].Id == c.ID {
		hasPreviousPage = true
		stories = stories[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		stories = stories[:limit-2]
		count = len(stories)
	}

	cursors := make([]string, count)
	for i, story := range stories {
		cursors[i] = cursor.NewCursor(story.Id, uint8(sortBy), story.Get(sortBy.String()), sortBy.CursorType()).String()
	}

	if paginationSortOrder < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		stories = utils.ReverseList(stories)
	}
	return stories, hasNextPage, hasPreviousPage, cursors, nil
}

func (s *StoriesStore) UpdateStory(ctx context.Context, id string, storyUpdate *storiesstore.StoryUpdate) error {

	if dproc, ok := interface{}(storyUpdate).(utils.DataProcessor); ok {
		if err := dproc.Process(ctx); err != nil {
			return err
		}
	}
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)

	now := time.Now()
	storyUpdate.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
		SetFields(storyUpdate)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(CollectionStories).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}

func (s *StoriesStore) DeleteStoryByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionStories).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}
