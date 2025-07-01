package utils

import (
	"fmt"
	"log/slog"

	"github.com/davecgh/go-spew/spew"
	"github.com/localrivet/gomcp/server"
)

// HandlerFunc represents a generic handler function signature
type HandlerFunc[T any, R any] func(*server.Context, T) (R, error)

// CallHandlerDirectly is a generic utility function that can call any handler directly
// with proper logging context. This replaces the package-specific helper functions.
func CallHandlerDirectly[T any, R any](logger *slog.Logger, functionName string, args T, handler HandlerFunc[T, R]) error {
	serverContext := &server.Context{
		Logger: logger,
	}

	logger.Info("Calling handler directly", "function", functionName)

	result, err := handler(serverContext, args)
	if err != nil {
		logger.Error("Handler call failed", "function", functionName, "error", err)
		return fmt.Errorf("failed to call %s: %w", functionName, err)
	}

	logger.Info("Handler call successful", "function", functionName)
	spew.Dump(result)

	return nil
}

// CreateServerContext creates a new server context with the provided logger
// This is useful for testing or direct handler calls
func CreateServerContext(logger *slog.Logger) *server.Context {
	return &server.Context{
		Logger: logger,
	}
}
