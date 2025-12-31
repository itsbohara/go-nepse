package nepse

import (
	"context"
	"fmt"
)

// IndexType represents the type of market index for graph data retrieval.
type IndexType int

const (
	IndexNepse IndexType = iota
	IndexSensitive
	IndexFloat
	IndexSensitiveFloat
	IndexBanking
	IndexDevBank
	IndexFinance
	IndexHotelTourism
	IndexHydro
	IndexInvestment
	IndexLifeInsurance
	IndexManufacturing
	IndexMicrofinance
	IndexMutualFund
	IndexNonLifeInsurance
	IndexOthers
	IndexTrading
)

// indexEndpoint returns the endpoint for a given index type.
func (c *Client) indexEndpoint(indexType IndexType) string {
	endpoints := c.config.Endpoints
	switch indexType {
	case IndexNepse:
		return endpoints.GraphNepseIndex
	case IndexSensitive:
		return endpoints.GraphSensitiveIndex
	case IndexFloat:
		return endpoints.GraphFloatIndex
	case IndexSensitiveFloat:
		return endpoints.GraphSensitiveFloatIndex
	case IndexBanking:
		return endpoints.GraphBankingSubindex
	case IndexDevBank:
		return endpoints.GraphDevBankSubindex
	case IndexFinance:
		return endpoints.GraphFinanceSubindex
	case IndexHotelTourism:
		return endpoints.GraphHotelSubindex
	case IndexHydro:
		return endpoints.GraphHydroSubindex
	case IndexInvestment:
		return endpoints.GraphInvestmentSubindex
	case IndexLifeInsurance:
		return endpoints.GraphLifeInsSubindex
	case IndexManufacturing:
		return endpoints.GraphManufacturingSubindex
	case IndexMicrofinance:
		return endpoints.GraphMicrofinanceSubindex
	case IndexMutualFund:
		return endpoints.GraphMutualFundSubindex
	case IndexNonLifeInsurance:
		return endpoints.GraphNonLifeInsSubindex
	case IndexOthers:
		return endpoints.GraphOthersSubindex
	case IndexTrading:
		return endpoints.GraphTradingSubindex
	default:
		return endpoints.GraphNepseIndex
	}
}

// GetDailyIndexGraph returns intraday graph data points for any market index.
func (c *Client) GetDailyIndexGraph(ctx context.Context, indexType IndexType) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := c.apiRequest(ctx, c.indexEndpoint(indexType), &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// GetDailyNepseIndexGraph returns intraday graph data for the main NEPSE index.
func (c *Client) GetDailyNepseIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexNepse)
}

// GetDailySensitiveIndexGraph returns intraday graph data for the sensitive index.
func (c *Client) GetDailySensitiveIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexSensitive)
}

// GetDailyFloatIndexGraph returns intraday graph data for the float index.
func (c *Client) GetDailyFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexFloat)
}

// GetDailySensitiveFloatIndexGraph returns intraday graph data for the sensitive float index.
func (c *Client) GetDailySensitiveFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexSensitiveFloat)
}

// GetDailyBankSubindexGraph returns intraday graph data for the banking sector sub-index.
func (c *Client) GetDailyBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexBanking)
}

// GetDailyDevelopmentBankSubindexGraph returns intraday graph data for the development bank sector.
func (c *Client) GetDailyDevelopmentBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexDevBank)
}

// GetDailyFinanceSubindexGraph returns intraday graph data for the finance sector.
func (c *Client) GetDailyFinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexFinance)
}

// GetDailyHotelTourismSubindexGraph returns intraday graph data for the hotel & tourism sector.
func (c *Client) GetDailyHotelTourismSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexHotelTourism)
}

// GetDailyHydroSubindexGraph returns intraday graph data for the hydropower sector.
func (c *Client) GetDailyHydroSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexHydro)
}

// GetDailyInvestmentSubindexGraph returns intraday graph data for the investment sector.
func (c *Client) GetDailyInvestmentSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexInvestment)
}

// GetDailyLifeInsuranceSubindexGraph returns intraday graph data for the life insurance sector.
func (c *Client) GetDailyLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexLifeInsurance)
}

// GetDailyManufacturingSubindexGraph returns intraday graph data for the manufacturing sector.
func (c *Client) GetDailyManufacturingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexManufacturing)
}

// GetDailyMicrofinanceSubindexGraph returns intraday graph data for the microfinance sector.
func (c *Client) GetDailyMicrofinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexMicrofinance)
}

// GetDailyMutualfundSubindexGraph returns intraday graph data for the mutual fund sector.
func (c *Client) GetDailyMutualfundSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexMutualFund)
}

// GetDailyNonLifeInsuranceSubindexGraph returns intraday graph data for the non-life insurance sector.
func (c *Client) GetDailyNonLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexNonLifeInsurance)
}

// GetDailyOthersSubindexGraph returns intraday graph data for the others sector.
func (c *Client) GetDailyOthersSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexOthers)
}

// GetDailyTradingSubindexGraph returns intraday graph data for the trading sector.
func (c *Client) GetDailyTradingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.GetDailyIndexGraph(ctx, IndexTrading)
}

// GetDailyScripPriceGraph returns intraday price graph data for a specific security.
func (c *Client) GetDailyScripPriceGraph(ctx context.Context, securityID int32) (*GraphResponse, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDailyGraph, securityID)
	var arr []GraphDataPoint
	if err := c.apiRequest(ctx, endpoint, &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// GetDailyScripPriceGraphBySymbol returns intraday price graph data for a security by ticker symbol.
func (c *Client) GetDailyScripPriceGraphBySymbol(ctx context.Context, symbol string) (*GraphResponse, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.GetDailyScripPriceGraph(ctx, security.ID)
}
