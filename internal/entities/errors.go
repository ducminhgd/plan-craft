package entities

import "errors"

var (
	ErrRecordNotFound                = errors.New("record not found")
	ErrInvalidTransaction            = errors.New("invalid transaction")
	ErrNotImplemented                = errors.New("not implemented")
	ErrMissingWhereClause            = errors.New("missing where clause")
	ErrUnsupportedRelation           = errors.New("unsupported relation")
	ErrPrimaryKeyRequired            = errors.New("primary key required")
	ErrModelValueRequired            = errors.New("model value required")
	ErrModelAccessibleFieldsRequired = errors.New("model accessible fields required")
	ErrSubQueryRequired              = errors.New("sub query required")
	ErrInvalidData                   = errors.New("invalid data")
	ErrUnsupportedDriver             = errors.New("unsupported driver")
	ErrRegistered                    = errors.New("registered")
	ErrInvalidField                  = errors.New("invalid field")
	ErrEmptySlice                    = errors.New("empty slice")
	ErrDryRunModeUnsupported         = errors.New("dry run mode unsupported")
	ErrInvalidDB                     = errors.New("invalid db")
	ErrInvalidValue                  = errors.New("invalid value")
	ErrInvalidValueOfLength          = errors.New("invalid value of length")
	ErrPreloadNotAllowed             = errors.New("preload not allowed")
	ErrDuplicatedKey                 = errors.New("duplicated key")
	ErrForeignKeyViolated            = errors.New("foreign key violated")
	ErrCheckConstraintViolated       = errors.New("check constraint violated")
)
