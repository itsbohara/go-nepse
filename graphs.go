package nepse

import (
	"context"
	"fmt"
)

// getIndexGraph is a helper for fetching index graph data.
func (c *Client) getIndexGraph(ctx context.Context, endpoint string) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := c.apiRequest(ctx, endpoint, &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// GetDailyNepseIndexGraph retrieves the daily NEPSE index graph.
func (c *Client) GetDailyNepseIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphNepseIndex)
}

// GetDailySensitiveIndexGraph retrieves the daily sensitive index graph.
func (c *Client) GetDailySensitiveIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphSensitiveIndex)
}

// GetDailyFloatIndexGraph retrieves the daily float index graph.
func (c *Client) GetDailyFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphFloatIndex)
}

// GetDailySensitiveFloatIndexGraph retrieves the daily sensitive float index graph.
func (c *Client) GetDailySensitiveFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphSensitiveFloatIndex)
}

// GetDailyBankSubindexGraph retrieves the daily banking sub-index graph.
func (c *Client) GetDailyBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphBankingSubindex)
}

// GetDailyDevelopmentBankSubindexGraph retrieves the daily development bank sub-index graph.
func (c *Client) GetDailyDevelopmentBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphDevBankSubindex)
}

// GetDailyFinanceSubindexGraph retrieves the daily finance sub-index graph.
func (c *Client) GetDailyFinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphFinanceSubindex)
}

// GetDailyHotelTourismSubindexGraph retrieves the daily hotel & tourism sub-index graph.
func (c *Client) GetDailyHotelTourismSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphHotelSubindex)
}

// GetDailyHydroSubindexGraph retrieves the daily hydro sub-index graph.
func (c *Client) GetDailyHydroSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphHydroSubindex)
}

// GetDailyInvestmentSubindexGraph retrieves the daily investment sub-index graph.
func (c *Client) GetDailyInvestmentSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphInvestmentSubindex)
}

// GetDailyLifeInsuranceSubindexGraph retrieves the daily life insurance sub-index graph.
func (c *Client) GetDailyLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphLifeInsSubindex)
}

// GetDailyManufacturingSubindexGraph retrieves the daily manufacturing sub-index graph.
func (c *Client) GetDailyManufacturingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphManufacturingSubindex)
}

// GetDailyMicrofinanceSubindexGraph retrieves the daily microfinance sub-index graph.
func (c *Client) GetDailyMicrofinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphMicrofinanceSubindex)
}

// GetDailyMutualfundSubindexGraph retrieves the daily mutual fund sub-index graph.
func (c *Client) GetDailyMutualfundSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphMutualFundSubindex)
}

// GetDailyNonLifeInsuranceSubindexGraph retrieves the daily non-life insurance sub-index graph.
func (c *Client) GetDailyNonLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphNonLifeInsSubindex)
}

// GetDailyOthersSubindexGraph retrieves the daily others sub-index graph.
func (c *Client) GetDailyOthersSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphOthersSubindex)
}

// GetDailyTradingSubindexGraph retrieves the daily trading sub-index graph.
func (c *Client) GetDailyTradingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.getIndexGraph(ctx, c.config.Endpoints.GraphTradingSubindex)
}

// GetDailyScripPriceGraph retrieves the daily scrip price graph for a security by ID.
func (c *Client) GetDailyScripPriceGraph(ctx context.Context, securityID int32) (*GraphResponse, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDailyGraph, securityID)
	var arr []GraphDataPoint
	if err := c.apiRequest(ctx, endpoint, &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// GetDailyScripPriceGraphBySymbol retrieves the daily scrip price graph for a security by symbol.
func (c *Client) GetDailyScripPriceGraphBySymbol(ctx context.Context, symbol string) (*GraphResponse, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.GetDailyScripPriceGraph(ctx, security.ID)
}
