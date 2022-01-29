package gateway

import (
	"context"
	"log"

	"github.com/firstcontributions/backend/internal/profile/proto"
)

func (s *Server) UpdateProfileReputation(profile *proto.Profile) {
	log.Println("from update profile, profile.Reputation", profile.Reputation)
	profile.Reputation++
	s.ProfileManager.UpdateProfile(context.Background(), profile)
}
