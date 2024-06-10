package utils

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/greenplum-db/gp-common-go-libs/gplog"
)

type HelpErr struct {
	err      error
	helpText string
}

func NewHelpErr(err error, helpText string) HelpErr {
	return HelpErr{
		err:      err,
		helpText: helpText,
	}
}

func (h HelpErr) Error() string {
	return h.err.Error()
}

func (h HelpErr) Help() {
	fmt.Println(h.helpText)
}

// FormatGrpcError formats the given error according to gRPC conventions.
// If the error is nil, it returns nil. If the error is a gRPC error, it extracts
// the error message and returns a new error with the extracted message. Otherwise,
// it returns the original error.
func FormatGrpcError(err error) error {
	if err == nil {
		return nil
	}

	grpcErr, ok := status.FromError(err)
	if ok {
		errorDescription := grpcErr.Message()
		return fmt.Errorf(errorDescription)
	}

	return err
}

// IsGrpcServerUnavailableErr checks if the given error is a gRPC server unavailable error.
// It returns true if the error is a gRPC server unavailable error, and false otherwise.
func IsGrpcServerUnavailableErr(err error) bool {
	if err == nil {
		return false
	}

	grpcErr, ok := status.FromError(err)
	if ok {
		errorCode := grpcErr.Code()
		return errorCode == codes.Unavailable
	}

	return false
}

// LogAndReturnError logs the given error and returns it.
func LogAndReturnError(err error) error {
	gplog.Error(err.Error())
	return err
}

// LogErrorAndExit logs the given error and exits the program with the specified exit code.
func LogErrorAndExit(err error, exitCode int) {
	gplog.Error(err.Error())

	var helpErr HelpErr
	if errors.As(err, &helpErr) {
		fmt.Println()
		helpErr.Help()
	}

	System.OSExit(exitCode)
}
