package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/storemanager"
)

type UserInput struct {
	Handle *string
}

func (r *Resolver) User(ctx context.Context, in *UserInput) (*User, error) {

	filters := &usersstore.UserFilters{
		Handle: in.Handle,
	}
	store := storemanager.FromContext(ctx)
	user, err := store.UsersStore.GetOneUser(ctx, filters)
	if err != nil {
		return nil, err
	}
	return NewUser(user), nil
}
