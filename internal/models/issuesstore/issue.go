package issuesstore

import (
	"time"

	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type IssueSortBy uint8

const (
	IssueSortByDefault = iota
	IssueSortByTimeCreated
)

type Issue struct {
	UserID              string    `bson:"user_id"`
	Body                string    `bson:"body,omitempty"`
	CommentCount        int64     `bson:"comment_count,omitempty"`
	Id                  string    `bson:"_id"`
	IssueType           string    `bson:"issue_type,omitempty"`
	Labels              []*string `bson:"labels,omitempty"`
	Repository          string    `bson:"repository,omitempty"`
	RepositoryAvatar    string    `bson:"repository_avatar,omitempty"`
	RepositoryUpdatedAt time.Time `bson:"repository_updated_at,omitempty"`
	Title               string    `bson:"title,omitempty"`
	Url                 string    `bson:"url,omitempty"`
	Cursor              string
}

func NewIssue() *Issue {
	return &Issue{}
}

type IssueFilters struct {
	Ids       []string
	IssueType *string
	User      *usersstore.User
}

func (s IssueSortBy) String() string {
	switch s {
	case IssueSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetIssueSortByFromString(s string) IssueSortBy {
	switch s {
	case "time_created":
		return IssueSortByTimeCreated
	default:
		return IssueSortByDefault
	}
}
