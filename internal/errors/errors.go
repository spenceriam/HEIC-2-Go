package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrorCode represents different types of errors that can occur in the application
type ErrorCode int

const (
	// File operations
	ErrFileNotFound ErrorCode = iota + 1000
	ErrFileRead
	ErrFileWrite
	ErrFileExists
	ErrFileCreate
	
	// Directory operations
	ErrDirCreate
	ErrDirRead
	
	// Conversion errors
	ErrInvalidImage
	ErrDecodeFailed
	ErrEncodeFailed
	ErrMetadataPreservation
	
	// Permission errors
	ErrPermissionDenied
	ErrAdminRequired
	
	// Input validation
	ErrInvalidInput
	ErrInvalidFormat
	
	// System errors
	ErrSystem
	ErrNotSupported
)

// AppError represents an application error with a code and message
type AppError struct {
	Code    ErrorCode
	Message string
	Details string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	msg := fmt.Sprintf("[%d] %s", e.Code, e.Message)
	if e.Details != "" {
		msg = fmt.Sprintf("%s: %s", msg, e.Details)
	}
	if e.Err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Err)
	}
	return msg
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new application error
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithError wraps another error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// Error utilities

// Wrap wraps an error with additional context
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is checks if the error is of a specific type
func Is(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// Common error constructors

// FileNotFound creates a new file not found error
func FileNotFound(path string) *AppError {
	return New(ErrFileNotFound, "file not found").WithDetails(path)
}

// FileExists creates a new file exists error
func FileExists(path string) *AppError {
	return New(ErrFileExists, "file already exists").WithDetails(path)
}

// PermissionDenied creates a new permission denied error
func PermissionDenied(path string) *AppError {
	return New(ErrPermissionDenied, "permission denied").WithDetails(path)
}

// InvalidInput creates a new invalid input error
func InvalidInput(field string, value interface{}) *AppError {
	details := fmt.Sprintf("field: %s, value: %v", field, value)
	return New(ErrInvalidInput, "invalid input").WithDetails(details)
}

// HandleError handles an error by logging it and returning a user-friendly message
func HandleError(err error) string {
	if err == nil {
		return ""
	}

	switch e := err.(type) {
	case *AppError:
		// Format the error message based on the error code
		switch e.Code {
		case ErrFileNotFound:
			return fmt.Sprintf("File not found: %s", e.Details)
		case ErrFileExists:
			return fmt.Sprintf("File already exists: %s", e.Details)
		case ErrPermissionDenied:
			return fmt.Sprintf("Permission denied: %s", e.Details)
		case ErrAdminRequired:
			return "Administrator/root privileges are required to perform this action"
		case ErrInvalidImage:
			return "The file is not a valid image"
		case ErrDecodeFailed:
			return "Failed to decode the image"
		case ErrEncodeFailed:
			return "Failed to encode the image"
		default:
			return fmt.Sprintf("Error: %s", e.Message)
		}
	default:
		// For non-AppError errors, return a generic message
		return fmt.Sprintf("An unexpected error occurred: %v", err)
	}
}

// IsFileNotFound checks if the error is a file not found error
func IsFileNotFound(err error) bool {
	return Is(err, ErrFileNotFound)
}

// IsPermissionDenied checks if the error is a permission denied error
func IsPermissionDenied(err error) bool {
	return Is(err, ErrPermissionDenied)
}

// HandleFileError handles file-related errors
func HandleFileError(err error, path string) error {
	if os.IsNotExist(err) {
		return FileNotFound(path)
	}
	if os.IsPermission(err) {
		return PermissionDenied(path)
	}
	return Wrap(err, ErrSystem, "file operation failed").WithDetails(path)
}
