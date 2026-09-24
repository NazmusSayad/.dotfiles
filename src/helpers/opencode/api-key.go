package opencode

import "os"

func ResolveApiKey(modelsDotDevProvider ModelsDotDevProvider) string {
	for _, env := range modelsDotDevProvider.Env {
		envVar := os.Getenv(env)
		if envVar != "" {
			return envVar
		}
	}

	return ""
}
