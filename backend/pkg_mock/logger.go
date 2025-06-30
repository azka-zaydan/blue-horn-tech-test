package pkgmock

import (
	"mini-evv-logger-backend/utils"
	"os"

	"github.com/rs/zerolog"
)

func InitMockLogger() zerolog.Logger {
	utils.InitLogger()
	mockLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return mockLogger
}
