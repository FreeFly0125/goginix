// Code generated by github.com/gokultp/go-envparser. DO NOT EDIT.
package configs

import (
	"fmt"
	"os"

	commonconfigs "github.com/firstcontributions/backend/internal/configs"
	"github.com/gokultp/go-envparser/pkg/envdecoder"
)

func (t *Config) DecodeEnv() error {
	_recLog := commonconfigs.LogConfig{}
	if err := envdecoder.Decode(&_recLog); err != nil {
		return fmt.Errorf("type commonconfigs.LogConfignot implemts env Decoder interface, %w", envdecoder.ErrDecoderNotImplemented)
	}
	t.Log = &_recLog
	if _recPortStr := os.Getenv("PROFILE_PORT"); _recPortStr != "" {
		_recPort := _recPortStr
		t.Port = &_recPort
	}
	if _recMongourlStr := os.Getenv("MONGO_URL"); _recMongourlStr != "" {
		_recMongourl := _recMongourlStr
		t.MongoURL = &_recMongourl
	}
	return nil
}
