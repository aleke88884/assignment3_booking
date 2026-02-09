package logger

import (
	"log"
	"os"
	"time"
)

// Logger provides structured logging for the application
type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = New()
}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs informational messages
func Info(format string, v ...interface{}) {
	defaultLogger.infoLog.Printf(format, v...)
}

// Error logs error messages
func Error(format string, v ...interface{}) {
	defaultLogger.errorLog.Printf(format, v...)
}

// Debug logs debug messages
func Debug(format string, v ...interface{}) {
	defaultLogger.debugLog.Printf(format, v...)
}

// LogRequest logs HTTP request information
func LogRequest(method, path string, statusCode int, duration time.Duration) {
	Info("HTTP %s %s - Status: %d - Duration: %v", method, path, statusCode, duration)
}

// LogDatabaseQuery logs database queries
func LogDatabaseQuery(query string, duration time.Duration, err error) {
	if err != nil {
		Error("DB Query failed: %s - Duration: %v - Error: %v", query, duration, err)
	} else {
		Debug("DB Query: %s - Duration: %v", query, duration)
	}
}

// LogServiceCall logs service layer operations
func LogServiceCall(service, method string, err error) {
	if err != nil {
		Error("Service %s.%s failed: %v", service, method, err)
	} else {
		Info("Service %s.%s completed successfully", service, method)
	}
}

// LogAuth logs authentication events
func LogAuth(event, userEmail string, success bool) {
	if success {
		Info("Auth event: %s - User: %s - Success", event, userEmail)
	} else {
		Error("Auth event: %s - User: %s - Failed", event, userEmail)
	}
}

// LogResourceOperation logs resource-related operations
func LogResourceOperation(operation string, resourceID int64, userID int64, err error) {
	if err != nil {
		Error("Resource %s - ResourceID: %d - UserID: %d - Error: %v", operation, resourceID, userID, err)
	} else {
		Info("Resource %s - ResourceID: %d - UserID: %d - Success", operation, resourceID, userID)
	}
}

// LogBookingOperation logs booking-related operations
func LogBookingOperation(operation string, bookingID int64, userID int64, resourceID int64, err error) {
	if err != nil {
		Error("Booking %s - BookingID: %d - UserID: %d - ResourceID: %d - Error: %v",
			operation, bookingID, userID, resourceID, err)
	} else {
		Info("Booking %s - BookingID: %d - UserID: %d - ResourceID: %d - Success",
			operation, bookingID, userID, resourceID)
	}
}
