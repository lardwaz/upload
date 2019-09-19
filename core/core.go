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
