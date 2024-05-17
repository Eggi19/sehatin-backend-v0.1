package custom_errors

import (
	"errors"
	"fmt"
	"net/http"

	constants "github.com/tsanaativa/sehatin-backend-v0.1/constants"
)

var (
	ErrForbidden           = errors.New(constants.ResponseMsgForbidden)
	ErrInvalidEmail        = errors.New(constants.InvalidEmailErrMsg)
	ErrInvalidAuthToken    = errors.New(constants.InvalidAuthTokenErrMsg)
	ErrExpiredResetPwdCode = errors.New(constants.ExpiredResetPwdCodeErrMsg)
	ErrSendEmail           = errors.New(constants.SendEmailErrMsg)
	ErrEmailNotFound       = errors.New(constants.ResponseMsgEmailNotFound)
	ErrInvalidRole         = errors.New(constants.InvalidRoleErrMsg)
	ErrVerifiedEmail       = errors.New(constants.VerifiedEmailErrMsg)
	ErrUploadFile          = errors.New(constants.UploadFileErrMsg)
	ErrFileRequired        = errors.New(constants.FileNotFoundErrMsg)
	ErrFileTooLarge        = errors.New(constants.FileSizeErrMsg)
	ErrFileNotPdf          = errors.New(constants.FileIsNotPdfErrMsg)
	ErrFileNotPng          = errors.New(constants.FileIsNotPngErrMsg)
	ErrNotVerified         = errors.New(constants.AccountNotVerifiedErrMsg)
	ErrDataNotFound        = errors.New(constants.ResponseMsgErrorNotFound)
	ErrTokenExpired        = errors.New(constants.ExpiredTokenErrMsg)
	ErrRoomAlreadyExists   = errors.New(constants.RoomNotUniqueErrMsg)
	ErrDoctorIsNotVerified = errors.New(constants.DoctorIsNotVerifiedErrMsg)
	ErrContextNotFound     = errors.New(constants.ContextNotFoundErrMsg)
	ErrNoGoogleAuthCode    = errors.New(constants.NoGoogleAuthCodeErrMsg)
	ErrNonNumberCoordinate = errors.New(constants.NonNumberCoordinateErrMsg)
	ErrFileNotImage        = errors.New(constants.FileIsNotImageErrMsg)
	ErrNotEnoughStock      = errors.New(constants.StockIsNotEnoughErrMsg)
)

type AppError struct {
	Code    int
	Message string
	err     error
}

func (e AppError) Error() string {
	return fmt.Sprint(e.Message)
}

func BadRequest(err error, message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		err:     err,
	}
}

func InternalServerError(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: constants.ResponseMsgErrorInternalServer,
		err:     err,
	}
}

func NotFound(err error) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: constants.ResponseMsgErrorNotFound,
		err:     err,
	}
}

func Forbidden() *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: constants.ResponseMsgForbidden,
		err:     ErrForbidden,
	}
}

func Unauthorized(err error, message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		err:     err,
	}
}

func InvalidAuthToken() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: constants.InvalidAuthTokenErrMsg,
		err:     ErrInvalidAuthToken,
	}
}

func InvalidEmail() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: constants.InvalidEmailErrMsg,
		err:     ErrInvalidEmail,
	}
}

func ExpiredResetPwdCode() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: constants.ExpiredResetPwdCodeErrMsg,
		err:     ErrExpiredResetPwdCode,
	}
}

func SendEmail() *AppError {
	return &AppError{
		Code:    http.StatusFailedDependency,
		Message: constants.SendEmailErrMsg,
		err:     ErrSendEmail,
	}
}

func EmailNotFound() *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: constants.ResponseMsgEmailNotFound,
		err:     ErrEmailNotFound,
	}
}

func InvalidRole() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: constants.InvalidRoleErrMsg,
		err:     ErrInvalidRole,
	}
}

func VerifiedEmail() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: constants.VerifiedEmailErrMsg,
		err:     ErrVerifiedEmail,
	}
}

func UploadFile() *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: constants.UploadFileErrMsg,
		err:     ErrUploadFile,
	}
}

func FileRequired() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: constants.FileNotFoundErrMsg,
		err:     ErrFileRequired,
	}
}

func FileTooLarge() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: constants.FileSizeErrMsg,
		err:     ErrFileTooLarge,
	}
}

func FileNotPdf() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: constants.FileIsNotPdfErrMsg,
		err:     ErrFileNotPdf,
	}
}

func FileNotPng() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: constants.FileIsNotPngErrMsg,
		err:     ErrFileNotPng,
	}
}

func NotVerified() *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: constants.AccountNotVerifiedErrMsg,
		err:     ErrNotVerified,
	}
}

func TokenExpired() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: constants.ExpiredTokenErrMsg,
		err:     ErrTokenExpired,
	}
}

func DoctorIsNotVerified() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: constants.DoctorIsNotVerifiedErrMsg,
		err:     ErrDoctorIsNotVerified,
	}
}

func ContextNotFound() *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: constants.ContextNotFoundErrMsg,
		err:     ErrContextNotFound,
	}
}

func FileNotImage() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: constants.FileIsNotImageErrMsg,
		err:     ErrFileNotImage,
	}
}

func NotEnoughStock() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: constants.StockIsNotEnoughErrMsg,
		err:     ErrNotEnoughStock,
	}
}
