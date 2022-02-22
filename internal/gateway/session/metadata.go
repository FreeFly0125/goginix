package session

import (
	"context"
	"encoding/json"

	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type CxtKey int

const (
	CxtKeySession CxtKey = iota
)

// MetaData encapsulates session info
type MetaData struct {
	usersstore.User
}

func NewMetaData(user *usersstore.User) MetaData {
	return MetaData{
		User: *user,
	}
}

func FromContext(ctx context.Context) MetaData {
	return ctx.Value(CxtKeySession).(MetaData)
}
func (m MetaData) SetHandle(h string) MetaData {
	m.User.Handle = h
	return m
}

func (m MetaData) SetUserID(uid string) MetaData {
	m.User.Id = uid
	return m
}

func (m MetaData) Handle() string {
	return m.User.Handle
}

func (m MetaData) UserID() string {
	return m.User.Id
}

func (m MetaData) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MetaData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}
