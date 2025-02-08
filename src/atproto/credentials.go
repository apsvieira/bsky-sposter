package atproto

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Credentials struct {
	Handle string
	AppKey string
}

// GetCredentials reads the .env file and returns BlueSky credentials.
func GetCredentials() (*Credentials, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading .env file: %w", err)
	}

	handle, ok := myEnv["BSKY_HANDLE"]
	if !ok {
		return nil, fmt.Errorf("BSKY_HANDLE not found in .env file: %w", err)
	}
	appkey, ok := myEnv["BSKY_APPKEY"]
	if !ok {
		return nil, fmt.Errorf("BSKY_APPKEY not found in .env file: %w", err)
	}

	return &Credentials{Handle: handle, AppKey: appkey}, nil
}
