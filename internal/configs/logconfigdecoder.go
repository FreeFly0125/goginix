package configs

// Code generated by github.com/gokultp/go-envparser. DO NOT EDIT.
import (
	"os"
)

func (t *LogConfig) DecodeEnv() error {
	if _recLevelStr := os.Getenv("LOG_LEVEL"); _recLevelStr != "" {
		_recLevel := _recLevelStr
		t.Level = &_recLevel
	}
	if _recPathStr := os.Getenv("LOG_PATH"); _recPathStr != "" {
		_recPath := _recPathStr
		t.Path = &_recPath
	}
	return nil
}
