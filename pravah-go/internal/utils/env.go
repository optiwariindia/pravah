package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

var requiredEnv = []string{
	"PORT",
	"MONGODB",
	"SUBSCRIPTION_JWT_SECRET",
	"ADMIN_API_KEY",
}

func CheckEnvironmentVariables() error {
	var missingEnv []string
	for _, env := range requiredEnv {
		if os.Getenv(env) == "" {
			missingEnv = append(missingEnv, env)
		}
	}

	if len(missingEnv) == 0 {
		if len(strings.TrimSpace(os.Getenv("SUBSCRIPTION_JWT_SECRET"))) < 32 {
			return errors.New("SUBSCRIPTION_JWT_SECRET must contain at least 32 characters")
		}
		if len(strings.TrimSpace(os.Getenv("ADMIN_API_KEY"))) < 32 {
			return errors.New("ADMIN_API_KEY must contain at least 32 characters")
		}
		for _, setting := range []struct {
			name string
			min  int
			max  int
		}{
			{"SSE_HEARTBEAT_INTERVAL_SECONDS", 5, 300},
			{"CLIENT_SIGNAL_BUFFER", 1, 10000},
			{"MAX_SUBSCRIPTION_STREAMS", 1, 128},
			{"MAX_REPLAY_EVENTS", 1, 100000},
			{"MAX_REQUEST_BODY_BYTES", 1024, 1048576},
			{"DATABASE_OPERATION_TIMEOUT_SECONDS", 1, 30},
		} {
			raw := os.Getenv(setting.name)
			if raw == "" {
				continue
			}
			value, err := strconv.Atoi(raw)
			if err != nil || value < setting.min || value > setting.max {
				return errors.New(setting.name + " is outside its supported range")
			}
		}
		return nil
	}

	msg := "Environment Variable"
	if len(missingEnv) > 1 {
		msg += "s " + joinWithAnd(missingEnv) + " are missing."
	} else {
		msg += " " + missingEnv[0] + " is missing."
	}

	return errors.New(msg)
}

func joinWithAnd(list []string) string {
	if len(list) == 0 {
		return ""
	}
	if len(list) == 1 {
		return list[0]
	}
	last := list[len(list)-1]
	rest := list[:len(list)-1]
	return strings.Join(rest, ", ") + " and " + last
}

func CORSAllowedOrigins() []string {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if raw == "" {
		return nil
	}

	origins := make([]string, 0)
	for _, origin := range strings.Split(raw, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func HeartbeatInterval() time.Duration {
	return time.Duration(positiveIntEnv("SSE_HEARTBEAT_INTERVAL_SECONDS", 30, 5, 300)) * time.Second
}

func ClientSignalBuffer() int {
	return positiveIntEnv("CLIENT_SIGNAL_BUFFER", 64, 1, 10000)
}

func MaxSubscriptionStreams() int {
	return positiveIntEnv("MAX_SUBSCRIPTION_STREAMS", 32, 1, 128)
}

func MaxReplayEvents() int {
	return positiveIntEnv("MAX_REPLAY_EVENTS", 1000, 1, 100000)
}

func MaxRequestBodyBytes() int64 {
	return int64(positiveIntEnv("MAX_REQUEST_BODY_BYTES", 16384, 1024, 1048576))
}

func DatabaseOperationTimeout() time.Duration {
	return time.Duration(positiveIntEnv("DATABASE_OPERATION_TIMEOUT_SECONDS", 3, 1, 30)) * time.Second
}

func positiveIntEnv(name string, defaultValue, min, max int) int {
	raw := os.Getenv(name)
	if raw == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < min || value > max {
		return defaultValue
	}
	return value
}
