package nepse

import (
	"context"
	"fmt"
)

// getIndexGraph is a helper for fetching index graph data.
func (h *nepseClient) getIndexGraph(ctx context.Context, endpointKey string) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints[endpointKey], &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// GetDailyNepseIndexGraph retrieves the daily NEPSE index graph.
func (h *nepseClient) GetDailyNepseIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "nepse_index_daily_graph")
}

// GetDailySensitiveIndexGraph retrieves the daily sensitive index graph.
func (h *nepseClient) GetDailySensitiveIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "sensitive_index_daily_graph")
}

// GetDailyFloatIndexGraph retrieves the daily float index graph.
func (h *nepseClient) GetDailyFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "float_index_daily_graph")
}

// GetDailySensitiveFloatIndexGraph retrieves the daily sensitive float index graph.
func (h *nepseClient) GetDailySensitiveFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "sensitive_float_index_daily_graph")
}

// GetDailyBankSubindexGraph retrieves the daily banking sub-index graph.
func (h *nepseClient) GetDailyBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "banking_sub_index_graph")
}

// GetDailyDevelopmentBankSubindexGraph retrieves the daily development bank sub-index graph.
func (h *nepseClient) GetDailyDevelopmentBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "development_bank_sub_index_graph")
}

// GetDailyFinanceSubindexGraph retrieves the daily finance sub-index graph.
func (h *nepseClient) GetDailyFinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "finance_sub_index_graph")
}

// GetDailyHotelTourismSubindexGraph retrieves the daily hotel & tourism sub-index graph.
func (h *nepseClient) GetDailyHotelTourismSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "hotel_tourism_sub_index_graph")
}

// GetDailyHydroSubindexGraph retrieves the daily hydro sub-index graph.
func (h *nepseClient) GetDailyHydroSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "hydro_sub_index_graph")
}

// GetDailyInvestmentSubindexGraph retrieves the daily investment sub-index graph.
func (h *nepseClient) GetDailyInvestmentSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "investment_sub_index_graph")
}

// GetDailyLifeInsuranceSubindexGraph retrieves the daily life insurance sub-index graph.
func (h *nepseClient) GetDailyLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "life_insurance_sub_index_graph")
}

// GetDailyManufacturingSubindexGraph retrieves the daily manufacturing sub-index graph.
func (h *nepseClient) GetDailyManufacturingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "manufacturing_sub_index_graph")
}

// GetDailyMicrofinanceSubindexGraph retrieves the daily microfinance sub-index graph.
func (h *nepseClient) GetDailyMicrofinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "microfinance_sub_index_graph")
}

// GetDailyMutualfundSubindexGraph retrieves the daily mutual fund sub-index graph.
func (h *nepseClient) GetDailyMutualfundSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "mutual_fund_sub_index_graph")
}

// GetDailyNonLifeInsuranceSubindexGraph retrieves the daily non-life insurance sub-index graph.
func (h *nepseClient) GetDailyNonLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "non_life_insurance_sub_index_graph")
}

// GetDailyOthersSubindexGraph retrieves the daily others sub-index graph.
func (h *nepseClient) GetDailyOthersSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "others_sub_index_graph")
}

// GetDailyTradingSubindexGraph retrieves the daily trading sub-index graph.
func (h *nepseClient) GetDailyTradingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return h.getIndexGraph(ctx, "trading_sub_index_graph")
}

// GetDailyScripPriceGraph retrieves the daily scrip price graph for a security by ID.
func (h *nepseClient) GetDailyScripPriceGraph(ctx context.Context, securityID int32) (*GraphResponse, error) {
	endpoint := fmt.Sprintf("%s%d", h.config.APIEndpoints["company_daily_graph"], securityID)
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, endpoint, &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// GetDailyScripPriceGraphBySymbol retrieves the daily scrip price graph for a security by symbol.
func (h *nepseClient) GetDailyScripPriceGraphBySymbol(ctx context.Context, symbol string) (*GraphResponse, error) {
	security, err := h.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return h.GetDailyScripPriceGraph(ctx, security.ID)
}
