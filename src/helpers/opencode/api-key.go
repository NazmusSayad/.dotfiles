package opencode

import "os"

func ResolveApiKey(providerId string, modelsDotDevProvider ModelsDotDevProvider, authConfig AuthConfig) string {
	if auth, ok := authConfig[providerId]; ok {
		if auth.Type == "api" && auth.Key != "" {
			return auth.Key
		}
	}

	for _, env := range modelsDotDevProvider.Env {
		envVar := os.Getenv(env)
		if envVar != "" {
			return envVar
		}
	}

	return ""
}
