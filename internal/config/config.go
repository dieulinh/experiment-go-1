package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// load priority (later files override earlier ones)
	// same pattern as Next.js / Laravel
	files := []string{
		".env",                   // base — always loaded
		".env." + env,            // .env.development or .env.production
		".env.local",             // local overrides — always gitignored
		".env." + env + ".local", // .env.development.local
	}

	for _, file := range files {
		if err := godotenv.Overload(file); err == nil {
			log.Printf("loaded config from %s", file)
		}
	}

	log.Printf("running in %s mode", env)
}
