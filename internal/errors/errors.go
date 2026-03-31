package errors

import (
	"net/http"
)

// ErrorCode - типизированный код ошибки
type ErrorCode string

const (
	// Общие ошибки
	INTERNAL_ERROR    ErrorCode = "INTERNAL_ERROR"
	INVALID_INPUT     ErrorCode = "INVALID_INPUT"
	UNAUTHORIZED      ErrorCode = "UNAUTHORIZED"
	FORBIDDEN         ErrorCode = "FORBIDDEN"
	NOT_FOUND         ErrorCode = "NOT_FOUND"
	
	// Ошибки аутентификации
	INVALID_CREDENTIALS ErrorCode = "INVALID_CREDENTIALS"
	TOKEN_REQUIRED      ErrorCode = "TOKEN_REQUIRED"
	TOKEN_INVALID       ErrorCode = "TOKEN_INVALID"
	
	// Ошибки матча
	MATCH_NOT_FOUND   ErrorCode = "MATCH_NOT_FOUND"
	SAME_TEAMS        ErrorCode = "SAME_TEAMS"
	INVALID_DATE      ErrorCode = "INVALID_DATE"
	
	// Ошибки пользователя
	USER_NOT_FOUND    ErrorCode = "USER_NOT_FOUND"
	USER_EXISTS       ErrorCode = "USER_EXISTS"
	INVALID_ROLE      ErrorCode = "INVALID_ROLE"
	
	// Ошибки прогноза
	PREDICTION_ERROR  ErrorCode = "PREDICTION_ERROR"
)

// AppError - структурированная ошибка приложения
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Status  int       `json:"status"`
	Details string    `json:"details,omitempty"`
}

// Error - реализация интерфейса error
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError - создание новой ошибки приложения
func NewAppError(code ErrorCode, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// NewAppErrorWithDetails - создание ошибки с деталями
func NewAppErrorWithDetails(code ErrorCode, message string, status int, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Details: details,
	}
}

// Предопределённые ошибки для удобства
var (
	ErrInternalError    = NewAppError(INTERNAL_ERROR, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	ErrInvalidInput     = NewAppError(INVALID_INPUT, "Некорректные входные данные", http.StatusBadRequest)
	ErrUnauthorized     = NewAppError(UNAUTHORIZED, "Требуется авторизация", http.StatusUnauthorized)
	ErrForbidden        = NewAppError(FORBIDDEN, "Недостаточно прав", http.StatusForbidden)
	ErrNotFound         = NewAppError(NOT_FOUND, "Ресурс не найден", http.StatusNotFound)
	ErrInvalidCredentials = NewAppError(INVALID_CREDENTIALS, "Неверный email или пароль", http.StatusUnauthorized)
	ErrTokenRequired    = NewAppError(TOKEN_REQUIRED, "Требуется JWT токен", http.StatusUnauthorized)
	ErrTokenInvalid     = NewAppError(TOKEN_INVALID, "Неверный или истёкший токен", http.StatusUnauthorized)
	ErrMatchNotFound    = NewAppError(MATCH_NOT_FOUND, "Матч не найден", http.StatusNotFound)
	ErrSameTeams        = NewAppError(SAME_TEAMS, "Команды не могут совпадать", http.StatusBadRequest)
	ErrInvalidDate      = NewAppError(INVALID_DATE, "Некорректная дата матча", http.StatusBadRequest)
	ErrUserNotFound     = NewAppError(USER_NOT_FOUND, "Пользователь не найден", http.StatusNotFound)
	ErrUserExists       = NewAppError(USER_EXISTS, "Пользователь уже существует", http.StatusConflict)
	ErrInvalidRole      = NewAppError(INVALID_ROLE, "Неверная роль", http.StatusBadRequest)
)
