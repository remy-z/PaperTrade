package api

import (
	"context"
	"fmt"
)

func (s *APIServer) finnhubQuote(symbol string) (interface{}, error) {
	res, _, err := s.finnhubClient.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve data: %v", err)
	}
	return res, nil
}
