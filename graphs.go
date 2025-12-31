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

// DailyIndexGraph returns intraday graph data points for any market index.
func (c *Client) DailyIndexGraph(ctx context.Context, indexType IndexType) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := c.apiRequest(ctx, c.indexEndpoint(indexType), &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// DailyNepseIndexGraph returns intraday graph data for the main NEPSE index.
func (c *Client) DailyNepseIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexNepse)
}

// DailySensitiveIndexGraph returns intraday graph data for the sensitive index.
func (c *Client) DailySensitiveIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexSensitive)
}

// DailyFloatIndexGraph returns intraday graph data for the float index.
func (c *Client) DailyFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexFloat)
}

// DailySensitiveFloatIndexGraph returns intraday graph data for the sensitive float index.
func (c *Client) DailySensitiveFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexSensitiveFloat)
}

// DailyBankSubindexGraph returns intraday graph data for the banking sector sub-index.
func (c *Client) DailyBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexBanking)
}

// DailyDevelopmentBankSubindexGraph returns intraday graph data for the development bank sector.
func (c *Client) DailyDevelopmentBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexDevBank)
}

// DailyFinanceSubindexGraph returns intraday graph data for the finance sector.
func (c *Client) DailyFinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexFinance)
}

// DailyHotelTourismSubindexGraph returns intraday graph data for the hotel & tourism sector.
func (c *Client) DailyHotelTourismSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexHotelTourism)
}

// DailyHydroSubindexGraph returns intraday graph data for the hydropower sector.
func (c *Client) DailyHydroSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexHydro)
}

// DailyInvestmentSubindexGraph returns intraday graph data for the investment sector.
func (c *Client) DailyInvestmentSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexInvestment)
}

// DailyLifeInsuranceSubindexGraph returns intraday graph data for the life insurance sector.
func (c *Client) DailyLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexLifeInsurance)
}

// DailyManufacturingSubindexGraph returns intraday graph data for the manufacturing sector.
func (c *Client) DailyManufacturingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexManufacturing)
}

// DailyMicrofinanceSubindexGraph returns intraday graph data for the microfinance sector.
func (c *Client) DailyMicrofinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexMicrofinance)
}

// DailyMutualfundSubindexGraph returns intraday graph data for the mutual fund sector.
func (c *Client) DailyMutualfundSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexMutualFund)
}

// DailyNonLifeInsuranceSubindexGraph returns intraday graph data for the non-life insurance sector.
func (c *Client) DailyNonLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexNonLifeInsurance)
}

// DailyOthersSubindexGraph returns intraday graph data for the others sector.
func (c *Client) DailyOthersSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexOthers)
}

// DailyTradingSubindexGraph returns intraday graph data for the trading sector.
func (c *Client) DailyTradingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexTrading)
}

// DailyScripGraph returns intraday price graph data for a specific security.
func (c *Client) DailyScripGraph(ctx context.Context, securityID int32) (*GraphResponse, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDailyGraph, securityID)
	var arr []GraphDataPoint
	if err := c.apiRequest(ctx, endpoint, &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// DailyScripGraphBySymbol returns intraday price graph data for a security by ticker symbol.
func (c *Client) DailyScripGraphBySymbol(ctx context.Context, symbol string) (*GraphResponse, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.DailyScripGraph(ctx, security.ID)
}
