package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/models/utils"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *StoriesStore) CreateReaction(ctx context.Context, reaction *storiesstore.Reaction) (*storiesstore.Reaction, error) {
	now := time.Now()
	reaction.TimeCreated = now
	reaction.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	reaction.Id = uuid.String()
	if _, err := s.getCollection(CollectionReactions).InsertOne(ctx, reaction); err != nil {
		return nil, err
	}
	return reaction, nil
}

func (s *StoriesStore) GetReactionByID(ctx context.Context, id string) (*storiesstore.Reaction, error) {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	var reaction storiesstore.Reaction
	if err := s.getCollection(CollectionReactions).FindOne(ctx, qb.Build()).Decode(&reaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &reaction, nil
}

func (s *StoriesStore) GetReactions(
	ctx context.Context,
	ids []string,
	comment *storiesstore.Comment,
	story *storiesstore.Story,
	after *string,
	before *string,
	first *int64,
	last *int64,
) (
	[]*storiesstore.Reaction,
	bool,
	bool,
	string,
	string,
	error,
) {
	qb := mongoqb.NewQueryBuilder()
	if len(ids) > 0 {
		qb.In("_id", ids)
	}
	if comment != nil {
		qb.Eq("comment_id", comment.Id)
	}
	if story != nil {
		qb.Eq("story_id", story.Id)
	}

	limit, order, cursorStr := utils.GetLimitAndSortOrderAndCursor(first, last, after, before)
	var c *cursor.Cursor
	if cursorStr != nil {
		c = cursor.FromString(*cursorStr)
		if c != nil {
			if order == 1 {
				qb.Lte("time_created", c.TimeStamp)
				qb.Lte("_id", c.ID)
			} else {
				qb.Gte("time_created", c.TimeStamp)
				qb.Gte("_id", c.ID)
			}
		}
	}
	sortOrder := utils.GetSortOrder(order)
	// incrementing limit by 2 to check if next, prev elements are present
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  sortOrder,
	}

	var firstCursor, lastCursor string
	var hasNextPage, hasPreviousPage bool

	var reactions []*storiesstore.Reaction
	mongoCursor, err := s.getCollection(CollectionReactions).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	err = mongoCursor.All(ctx, &reactions)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	count := len(reactions)
	if count == 0 {
		return reactions, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && reactions[0].Id == c.ID {
		hasPreviousPage = true
		reactions = reactions[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		reactions = reactions[:limit-2]
		count = len(reactions)
	}

	if count > 0 {
		firstCursor = cursor.NewCursor(reactions[0].Id, reactions[0].TimeCreated).String()
		lastCursor = cursor.NewCursor(reactions[count-1].Id, reactions[count-1].TimeCreated).String()
	}
	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		firstCursor, lastCursor = lastCursor, firstCursor
	}
	return reactions, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
}
func (s *StoriesStore) UpdateReaction(ctx context.Context, id string, reactionUpdate *storiesstore.ReactionUpdate) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)

	now := time.Now()
	reactionUpdate.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
		SetFields(reactionUpdate)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(CollectionReactions).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}

func (s *StoriesStore) DeleteReactionByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionReactions).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}
