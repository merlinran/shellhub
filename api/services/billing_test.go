package services

import (
	"errors"
	"testing"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/api/store/mocks"
	"github.com/shellhub-io/shellhub/pkg/cache"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
)

func mockBillingEvaluation(canAccept, canConnect bool) *models.BillingEvaluation {
	return &models.BillingEvaluation{
		CanAccept:  canAccept,
		CanConnect: canConnect,
	}
}

func TestBillingEvaluate(t *testing.T) {
    mock := new(mocks.Store)

	type Expected struct {
		canAccept bool
		err       error
	}

	cases := []struct {
		description   string
		tenant        string
		requiredMocks func()
		expected      Expected
	}{
		{
			description: "succeeds when \"client.BillingEvaluate\" err is nil",
			tenant:      "00000000-0000-0000-0000-000000000000",
			requiredMocks: func() {
				clientMock.On("BillingEvaluate", "00000000-0000-0000-0000-000000000000").Return(mockBillingEvaluation(true, true), 0, nil).Once()
			},
			expected: Expected{canAccept: true, err: nil},
		},
		{
			description: "fails when \"client.BillingEvaluate\" err is different than nil",
			tenant:      "00000000-0000-0000-0000-000000000000",
			requiredMocks: func() {
				clientMock.On("BillingEvaluate", "00000000-0000-0000-0000-000000000000").Return(mockBillingEvaluation(false, false), 0, ErrEvaluate).Once()
			},
			expected: Expected{canAccept: false, err: ErrEvaluate},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			tc.requiredMocks()

            s := NewService(store.Store(mock), privateKey, publicKey, cache.NewNullCache(), clientMock, nil)
			canAccept, err := s.billingEvaluate(clientMock, tc.tenant)
			assert.Equal(t, tc.expected, Expected{canAccept: canAccept, err: err})
		})
	}

	mock.AssertExpectations(t)
}

func TestBillingReport(t *testing.T) {
    mock := new(mocks.Store)

	cases := []struct {
		description   string
		tenant        string
		action        string
		requiredMocks func()
		expected      error
	}{
		{
			description: "succeeds when \"client.BillingReport\" status is 200",
			tenant:      "00000000-0000-0000-0000-000000000000",
			action:      "device_accept",
			requiredMocks: func() {
				clientMock.On("BillingReport", "00000000-0000-0000-0000-000000000000", "device_accept").Return(200, nil).Once()
			},
			expected: nil,
		},
		{
			description: "fails when \"client.BillingReport\" status is 402",
			tenant:      "00000000-0000-0000-0000-000000000000",
			action:      "device_accept",
			requiredMocks: func() {
				clientMock.On("BillingReport", "00000000-0000-0000-0000-000000000000", "device_accept").Return(402, nil).Once()
			},
			expected: ErrPaymentRequired,
		},
		{
			description: "fails when \"client.BillingReport\" status is other than 200 or 402",
			tenant:      "00000000-0000-0000-0000-000000000000",
			action:      "device_accept",
			requiredMocks: func() {
				clientMock.On("BillingReport", "00000000-0000-0000-0000-000000000000", "device_accept").Return(500, nil).Once()
			},
			expected: ErrReport,
		},
		{
			description: "fails when \"client.BillingReport\" returns an error",
			tenant:      "00000000-0000-0000-0000-000000000000",
			action:      "device_accept",
			requiredMocks: func() {
				clientMock.On("BillingReport", "00000000-0000-0000-0000-000000000000", "device_accept").Return(0, errors.New("error")).Once()
			},
			expected: errors.New("error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			tc.requiredMocks()

            s := NewService(store.Store(mock), privateKey, publicKey, cache.NewNullCache(), clientMock, nil)
			err := s.billingReport(clientMock, tc.tenant, tc.action)
			assert.Equal(t, tc.expected, err)
		})
	}

	mock.AssertExpectations(t)
}