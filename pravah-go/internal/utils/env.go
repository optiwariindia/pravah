package utils

import (
	"errors"
	"os"
	"strings"
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
