package ravandlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	ContextKeyUserUUID    = "user_uuid"
	ContextKeyRequestUUID = "request_uuid"
	ContextKeyFunction    = "function"
	ContextKeyFile        = "file"
)

func (l *Logger) extractFieldsFromContext(ctx context.Context) logrus.Fields {
	fields := make(map[string]interface{})

	userUUIDValue := ctx.Value(ContextKeyUserUUID)
	if userUUIDValue != nil {
		fields[ContextKeyUserUUID] = userUUIDValue.(string)
	}
	requestUUIDValue := ctx.Value(ContextKeyRequestUUID)
	if requestUUIDValue != nil {
		fields[ContextKeyRequestUUID] = requestUUIDValue.(string)
	}
	functionNameValue := ctx.Value(ContextKeyFunction)
	if functionNameValue != nil {
		fields[ContextKeyFunction] = functionNameValue.(string)
	}
	if l.config.IsReportCaller {
		fileValue := ctx.Value(ContextKeyFile)
		if fileValue != nil {
			fields[ContextKeyFile] = fileValue.(string)
		}
	}

	return fields
}
