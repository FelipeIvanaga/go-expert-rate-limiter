package mocks

import (
	"context"
	"net/http"

	"github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/ratelimiter/strategies"
	"github.com/stretchr/testify/mock"
)

type RateLimiterMock struct {
	mock.Mock
}

func (m *RateLimiterMock) Check(ctx context.Context, r *http.Request) (*strategies.RateLimiterResult, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*strategies.RateLimiterResult), args.Error(1)
}
