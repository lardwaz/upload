package core

const (
	// NoLimit define no limits
	NoLimit = -1

	// EnvironmentPROD defines production environment
	EnvironmentPROD = "PROD"

	// EnvironmentDEV defines development environment
	EnvironmentDEV = "DEV"
)

var (
	// Env defines current environment
	Env = EnvironmentDEV
)

// SetEnv sets the environment gocipe-upload operates in
func SetEnv(env string) {
	switch env {
	case EnvironmentDEV, EnvironmentPROD:
		// We are good :)
	default:
		// Invalid environment
		return
	}
	Env = env
}
