package errorsx

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/hayrat/go-template2/backend/pkg/log"
	"github.com/hayrat/go-template2/backend/pkg/viewmodel"
)

const (
	ErrorMesasgeNotFound       = "not found"
	ErrorMesasgeDatabase       = "database error"
	ErrorMesasgeInternalServer = "internal server error"
)

// ErrorType is the type of error
type ErrorType int

const (
	ErrorTypeNoType       ErrorType = 0
	ErrorTypeBadRequest   ErrorType = 400
	ErrorTypeUnauthorized ErrorType = 401
	ErrorTypeForbidden    ErrorType = 403
	ErrorTypeNotFound     ErrorType = 404
	ErrorTypeInternal     ErrorType = 500
)

// ErrorCode is the identifier for frontend error management
type ErrorCode int

const (
// ErrorCode ErrorCode = 1
)

type Errorx interface {
	Error() string
}

type errorx struct {
	errorType     ErrorType
	errorCode     ErrorCode
	message       string
	originalError error
	callerPath    string
}

// Error returns the mssage of a errorx
func (e errorx) Error() string {
	return e.message
}

func NewErrorx(errorType ErrorType, msg string) Errorx {
	return errorx{errorType: errorType, originalError: errors.New(msg)}
}

func New(msg string) error {
	return errors.New(msg)
}

func Is(err1, err2 error) bool {
	return errors.Is(err1, err2)
}

func BadRequestError(msg string) Errorx {
	return errorx{errorType: ErrorTypeBadRequest, message: msg}
}

func ValidationError(errs []error) Errorx {
	msg := ""
	for _, err := range errs {
		msg = strings.Join([]string{msg, err.Error()}, "\n") // kim uğraşacak son elemana \n filan koymayla
	}
	return errorx{errorType: ErrorTypeBadRequest, message: msg}
}

func UnauthorizedError(msg string) Errorx {
	return errorx{errorType: ErrorTypeUnauthorized, message: msg}
}

func ForbiddenError(msg ...string) Errorx {
	m := "permission error"
	if len(msg) > 0 {
		m = msg[0]
	}
	return errorx{errorType: ErrorTypeForbidden, message: m}
}

func NotFoundError(msg string) Errorx {
	return errorx{errorType: ErrorTypeNotFound, message: msg}
}

func Database(err error, recordNotFoundMsg ...string) Errorx {
	if err == nil {
		return nil
	}
	if IsDBNotFoundError(err) {
		msg := ErrorMesasgeNotFound
		if len(recordNotFoundMsg) > 0 {
			msg = recordNotFoundMsg[0]
		}
		return NotFoundError(msg)
	}

	return InternalError(err, ErrorMesasgeDatabase)
}

func IsDBNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func InternalError(err error, message ...string) Errorx {
	if err == nil {
		return nil
	}
	m := ErrorMesasgeInternalServer
	if len(message) > 0 {
		m = message[0]
	}

	skip := 1
	if m == ErrorMesasgeDatabase {
		skip = 2
	}
	callerPath := ""
	_, fullPath, lineNo, ok := runtime.Caller(skip)
	if ok {
		baseDir, _ := os.Getwd()
		relPath, _ := filepath.Rel(baseDir, fullPath)

		callerPath = fmt.Sprintf("%s:%d", relPath, lineNo)
	}

	return errorx{errorType: ErrorTypeInternal, originalError: err, message: m, callerPath: callerPath}
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(errorx); ok {
		return customErr.errorType
	}

	return ErrorTypeNoType
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	respModel := &viewmodel.ResponseModel{}

	logger, ok := c.Locals("logger").(*log.Logger)
	if !ok {
		l := log.GetLogger("error_handler")
		logger = &l
	}

	errType := GetType(err)
	if errType == ErrorTypeNoType {
		if !strings.HasPrefix(err.Error(), "Cannot ") {
			logger.Error("TypeNoType", zap.Error(err))
		}
		code := http.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		respModel.ErrorMessage = ""
		return c.Status(code).JSON(respModel)
	}

	errx, _ := err.(errorx)
	respModel.ErrorMessage = errx.Error()

	switch errType {
	case ErrorTypeBadRequest:
		return c.Status(http.StatusBadRequest).JSON(respModel)
	case ErrorTypeNotFound:
		return c.Status(http.StatusNotFound).JSON(respModel)
	case ErrorTypeUnauthorized:
		return c.Status(http.StatusUnauthorized).JSON(respModel)
	case ErrorTypeForbidden:
		return c.Status(http.StatusForbidden).JSON(respModel)
	case ErrorTypeInternal:
		logger.Error(errx.message, zap.Error(errx.originalError), zap.Int("customCode", int(errx.errorCode)), zap.String("callerPath", errx.callerPath))
		return c.Status(http.StatusInternalServerError).JSON(respModel)
	default:
		return err
	}
}
