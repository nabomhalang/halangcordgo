package config

func Init() {
	logger := NewLogger("config")
	if err := initializeEnvironment(); err != nil {
		logger.Fatalf("failed to initialize environment variables: %v", err)
	}
}
