package api

import "go.uber.org/zap"

func zapError(err error) zap.Field {
	return zap.Error(err)
}
