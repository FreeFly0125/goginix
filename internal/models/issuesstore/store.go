package issuesstore

import "context"

type Store interface {
	// issue methods
	GetIssueByID(ctx context.Context, id string) (*Issue, error)
	GetOneIssue(ctx context.Context, filters *IssueFilters) (*Issue, error)
	GetIssues(ctx context.Context, filters *IssueFilters, after *string, before *string, first *int64, last *int64, sortBy IssueSortBy, sortOrder *string) ([]*Issue, bool, bool, []string, error)
	CountIssues(ctx context.Context, filters *IssueFilters) (int64, error)
}
