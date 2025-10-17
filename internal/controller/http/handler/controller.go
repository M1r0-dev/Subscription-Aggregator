package handler

import (
	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/mapper"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/parser"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/usecase"
	"github.com/M1r0-dev/Subscription-Aggregator/pkg/logger"
)

type SubscriptionHandler struct {
	usecase usecase.SubscriptionUsecase
	logger logger.Interface
	parser *parser.SubscriptionParser
	mapper *mapper.SubscriptionMapper
}

func New(usecase usecase.SubscriptionUsecase, logger logger.Interface, parser *parser.SubscriptionParser, mapper *mapper.SubscriptionMapper) *SubscriptionHandler {
	return &SubscriptionHandler{
		usecase: usecase,
		logger: logger,
		parser: parser,
		mapper: mapper,
	}
}