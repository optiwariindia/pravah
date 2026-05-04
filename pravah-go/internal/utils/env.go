package utils

import (
	"fmt"
	"os"
	"strings"
)

var requiredEnv = []string{
	"PORT",
	"MONGODB",
}

func CheckEnvironmentVariables() error {
	var missingEnv []string
	for _, env := range requiredEnv {
		if os.Getenv(env) == "" {
			missingEnv = append(missingEnv, env)
		}
	}

	if len(missingEnv) == 0 {
		return nil
	}

	msg := "Environment Variable"
	if len(missingEnv) > 1 {
		msg += "s " + joinWithAnd(missingEnv) + " are missing."
	} else {
		msg += " " + missingEnv[0] + " is missing."
	}

	return fmt.Errorf(msg)
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
