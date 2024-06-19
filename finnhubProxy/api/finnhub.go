package api

import (
	"context"
	"strings"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/remy-z/PaperTrade/finnhubProxy/common"
)

func (s *APIServer) initFinnhubData() error {
	s.tree = &common.TST{}

	res, err := s.finnhubSymbols("US")
	if err != nil {
		return err
	}

	for _, symbol := range res {
		displaySymbol := strings.ToUpper(symbol.GetDisplaySymbol())
		description := strings.ToUpper(symbol.GetDescription())
		s.tree.Put(displaySymbol, description)
		s.descriptionToSymbol[description] = displaySymbol
		s.companyInfo[displaySymbol] = NewCompanyProfile(displaySymbol, description)
	}
	return nil
}

func (s *APIServer) finnhubQuote(symbol string) (*finnhub.Quote, error) {
	res, _, err := s.finnhubClient.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *APIServer) finnhubFinancial(symbol string) (*finnhub.BasicFinancials, error) {
	res, _, err := s.finnhubClient.CompanyBasicFinancials(context.Background()).Symbol(symbol).Metric("all").Execute()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *APIServer) finnhubSymbols(market string) ([]finnhub.StockSymbol, error) {
	res, _, err := s.finnhubClient.StockSymbols(context.Background()).Exchange(market).Execute()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *APIServer) finnhubProfile(symbol string) (*finnhub.CompanyProfile2, error) {
	res, _, err := s.finnhubClient.CompanyProfile2(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return nil, err
	}
	return &res, nil
}
