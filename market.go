package nepse

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetMarketSummary retrieves the overall market summary.
func (h *nepseClient) GetMarketSummary(ctx context.Context) (*MarketSummary, error) {
	var rawItems []MarketSummaryItem
	if err := h.apiRequest(ctx, h.config.APIEndpoints["market_summary"], &rawItems); err != nil {
		return nil, err
	}

	summary := &MarketSummary{}
	for _, item := range rawItems {
		switch item.Detail {
		case "Total Turnover Rs:":
			summary.TotalTurnover = item.Value
		case "Total Traded Shares":
			summary.TotalTradedShares = item.Value
		case "Total Transactions":
			summary.TotalTransactions = item.Value
		case "Total Scrips Traded":
			summary.TotalScripsTraded = item.Value
		case "Total Market Capitalization Rs:":
			summary.TotalMarketCapitalization = item.Value
		case "Total Float Market Capitalization Rs:":
			summary.TotalFloatMarketCap = item.Value
		}
	}

	return summary, nil
}

// GetMarketStatus retrieves the current market status.
func (h *nepseClient) GetMarketStatus(ctx context.Context) (*MarketStatus, error) {
	var status MarketStatus
	if err := h.apiRequest(ctx, h.config.APIEndpoints["market_open"], &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// GetNepseIndex retrieves the NEPSE index information.
func (h *nepseClient) GetNepseIndex(ctx context.Context) (*NepseIndex, error) {
	var rawIndices []NepseIndexRaw
	if err := h.apiRequest(ctx, h.config.APIEndpoints["nepse_index"], &rawIndices); err != nil {
		return nil, err
	}

	for _, rawIndex := range rawIndices {
		if rawIndex.ID == 58 && rawIndex.Index == "NEPSE Index" {
			return &NepseIndex{
				IndexValue:       rawIndex.Close,
				PercentChange:    rawIndex.PerChange,
				PointChange:      rawIndex.Change,
				High:             rawIndex.High,
				Low:              rawIndex.Low,
				PreviousClose:    rawIndex.PreviousClose,
				FiftyTwoWeekHigh: rawIndex.FiftyTwoWeekHigh,
				FiftyTwoWeekLow:  rawIndex.FiftyTwoWeekLow,
				CurrentValue:     rawIndex.CurrentValue,
				GeneratedTime:    rawIndex.GeneratedTime,
			}, nil
		}
	}

	return nil, NewNotFoundError("NEPSE Index")
}

// GetNepseSubIndices retrieves all NEPSE sub-indices.
func (h *nepseClient) GetNepseSubIndices(ctx context.Context) ([]SubIndex, error) {
	var rawIndices []NepseIndexRaw
	if err := h.apiRequest(ctx, h.config.APIEndpoints["nepse_index"], &rawIndices); err != nil {
		return nil, err
	}

	var subIndices []SubIndex
	for _, rawIndex := range rawIndices {
		if rawIndex.ID != 58 && rawIndex.ID != 57 && rawIndex.ID != 62 && rawIndex.ID != 63 {
			subIndices = append(subIndices, SubIndex(rawIndex))
		}
	}

	if len(subIndices) == 0 {
		for _, rawIndex := range rawIndices {
			if rawIndex.Index != "NEPSE Index" {
				subIndices = append(subIndices, SubIndex(rawIndex))
			}
		}
	}

	return subIndices, nil
}

// GetLiveMarket retrieves live market data.
func (h *nepseClient) GetLiveMarket(ctx context.Context) ([]LiveMarketEntry, error) {
	var liveMarket []LiveMarketEntry
	if err := h.apiRequest(ctx, h.config.APIEndpoints["live_market"], &liveMarket); err != nil {
		return nil, err
	}
	return liveMarket, nil
}

// GetSupplyDemand retrieves supply and demand data.
func (h *nepseClient) GetSupplyDemand(ctx context.Context) ([]SupplyDemandEntry, error) {
	endpoint := h.config.APIEndpoints["supply_demand"]

	var arr []SupplyDemandEntry
	if err := h.apiRequest(ctx, endpoint, &arr); err == nil {
		return arr, nil
	}

	data, err := h.apiRequestRaw(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &arr); err == nil {
		return arr, nil
	}

	var pagRoot struct {
		Content []SupplyDemandEntry `json:"content"`
	}
	if err := json.Unmarshal(data, &pagRoot); err == nil && len(pagRoot.Content) > 0 {
		return pagRoot.Content, nil
	}

	var nested map[string]json.RawMessage
	if err := json.Unmarshal(data, &nested); err == nil {
		for _, v := range nested {
			var maybe struct {
				Content []SupplyDemandEntry `json:"content"`
			}
			if json.Unmarshal(v, &maybe) == nil && len(maybe.Content) > 0 {
				return maybe.Content, nil
			}
		}
	}

	return nil, NewInvalidServerResponseError("unrecognized supply/demand response shape")
}

// GetTopGainers retrieves the top gainers list.
func (h *nepseClient) GetTopGainers(ctx context.Context) ([]TopListEntry, error) {
	var topGainers []TopListEntry
	if err := h.apiRequest(ctx, h.config.APIEndpoints["top_gainers"], &topGainers); err != nil {
		return nil, err
	}
	return topGainers, nil
}

// GetTopLosers retrieves the top losers list.
func (h *nepseClient) GetTopLosers(ctx context.Context) ([]TopListEntry, error) {
	var topLosers []TopListEntry
	if err := h.apiRequest(ctx, h.config.APIEndpoints["top_losers"], &topLosers); err != nil {
		return nil, err
	}
	return topLosers, nil
}

// GetTopTenTrade retrieves the top ten trade list.
func (h *nepseClient) GetTopTenTrade(ctx context.Context) ([]TopListEntry, error) {
	var topTrade []TopListEntry
	if err := h.apiRequest(ctx, h.config.APIEndpoints["top_ten_trade"], &topTrade); err != nil {
		return nil, err
	}
	return topTrade, nil
}

// GetTopTenTransaction retrieves the top ten transaction list.
func (h *nepseClient) GetTopTenTransaction(ctx context.Context) ([]TopListEntry, error) {
	var topTransaction []TopListEntry
	if err := h.apiRequest(ctx, h.config.APIEndpoints["top_ten_transaction"], &topTransaction); err != nil {
		return nil, err
	}
	return topTransaction, nil
}

// GetTopTenTurnover retrieves the top ten turnover list.
func (h *nepseClient) GetTopTenTurnover(ctx context.Context) ([]TopListEntry, error) {
	var topTurnover []TopListEntry
	if err := h.apiRequest(ctx, h.config.APIEndpoints["top_ten_turnover"], &topTurnover); err != nil {
		return nil, err
	}
	return topTurnover, nil
}

// GetTodaysPrices retrieves today's price data.
func (h *nepseClient) GetTodaysPrices(ctx context.Context, businessDate string) ([]TodayPrice, error) {
	endpoint := h.config.APIEndpoints["todays_price"]
	if businessDate != "" {
		endpoint += "?businessDate=" + businessDate + "&size=500"
	}

	var todayPrices []TodayPrice
	if err := h.apiRequest(ctx, endpoint, &todayPrices); err != nil {
		return nil, err
	}
	return todayPrices, nil
}

// GetPriceVolumeHistory retrieves price volume history for a security by ID.
func (h *nepseClient) GetPriceVolumeHistory(ctx context.Context, securityID int32, startDate, endDate string) ([]PriceHistory, error) {
	endpoint := fmt.Sprintf("%s%d?size=500&startDate=%s&endDate=%s",
		h.config.APIEndpoints["company_price_volume_history"], securityID, startDate, endDate)

	var response struct {
		Content []PriceHistory `json:"content"`
	}

	if err := h.apiRequest(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	return response.Content, nil
}

// GetPriceVolumeHistoryBySymbol retrieves price volume history for a security by symbol.
func (h *nepseClient) GetPriceVolumeHistoryBySymbol(ctx context.Context, symbol string, startDate, endDate string) ([]PriceHistory, error) {
	security, err := h.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return h.GetPriceVolumeHistory(ctx, security.ID, startDate, endDate)
}

// GetMarketDepth retrieves market depth information for a security by ID.
func (h *nepseClient) GetMarketDepth(ctx context.Context, securityID int32) (*MarketDepth, error) {
	endpoint := fmt.Sprintf("%s%d/", h.config.APIEndpoints["market_depth"], securityID)

	var marketDepth MarketDepth
	if err := h.apiRequest(ctx, endpoint, &marketDepth); err != nil {
		return nil, err
	}
	return &marketDepth, nil
}

// GetMarketDepthBySymbol retrieves market depth information for a security by symbol.
func (h *nepseClient) GetMarketDepthBySymbol(ctx context.Context, symbol string) (*MarketDepth, error) {
	security, err := h.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return h.GetMarketDepth(ctx, security.ID)
}

// GetSecurityList retrieves the list of all securities.
func (h *nepseClient) GetSecurityList(ctx context.Context) ([]Security, error) {
	var securities []Security
	if err := h.apiRequest(ctx, h.config.APIEndpoints["security_list"], &securities); err != nil {
		return nil, err
	}
	return securities, nil
}

// GetCompanyList retrieves the list of all companies.
func (h *nepseClient) GetCompanyList(ctx context.Context) ([]Company, error) {
	var companies []Company
	if err := h.apiRequest(ctx, h.config.APIEndpoints["company_list"], &companies); err != nil {
		return nil, err
	}
	return companies, nil
}

// GetCompanyDetails retrieves detailed information about a specific company/security by ID.
func (h *nepseClient) GetCompanyDetails(ctx context.Context, securityID int32) (*CompanyDetails, error) {
	endpoint := fmt.Sprintf("%s%d", h.config.APIEndpoints["company_details"], securityID)

	var rawDetails CompanyDetailsRaw
	if err := h.apiRequest(ctx, endpoint, &rawDetails); err != nil {
		return nil, err
	}

	details := &CompanyDetails{
		ID:               rawDetails.SecurityData.ID,
		Symbol:           rawDetails.SecurityData.Symbol,
		SecurityName:     rawDetails.SecurityData.SecurityName,
		SectorName:       rawDetails.SecurityData.Sector,
		Email:            rawDetails.SecurityData.Email,
		ActiveStatus:     rawDetails.SecurityData.ActiveStatus,
		PermittedToTrade: rawDetails.SecurityData.PermittedToTrade,

		OpenPrice:           rawDetails.SecurityMcsData.OpenPrice,
		HighPrice:           rawDetails.SecurityMcsData.HighPrice,
		LowPrice:            rawDetails.SecurityMcsData.LowPrice,
		ClosePrice:          rawDetails.SecurityMcsData.ClosePrice,
		LastTradedPrice:     rawDetails.SecurityMcsData.LastTradedPrice,
		PreviousClose:       rawDetails.SecurityMcsData.PreviousClose,
		TotalTradeQuantity:  rawDetails.SecurityMcsData.TotalTradeQuantity,
		TotalTrades:         rawDetails.SecurityMcsData.TotalTrades,
		FiftyTwoWeekHigh:    rawDetails.SecurityMcsData.FiftyTwoWeekHigh,
		FiftyTwoWeekLow:     rawDetails.SecurityMcsData.FiftyTwoWeekLow,
		BusinessDate:        rawDetails.SecurityMcsData.BusinessDate,
		LastUpdatedDateTime: rawDetails.SecurityMcsData.LastUpdatedDateTime,
	}

	return details, nil
}

// GetCompanyDetailsBySymbol retrieves detailed information about a specific company/security by symbol.
func (h *nepseClient) GetCompanyDetailsBySymbol(ctx context.Context, symbol string) (*CompanyDetails, error) {
	security, err := h.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return h.GetCompanyDetails(ctx, security.ID)
}

// GetSectorScrips groups securities by their sector.
func (h *nepseClient) GetSectorScrips(ctx context.Context) (SectorScrips, error) {
	securities, err := h.GetSecurityList(ctx)
	if err != nil {
		return nil, err
	}

	sectorScrips := make(SectorScrips)

	for _, security := range securities {
		if security.IsSuspended {
			continue
		}

		var sectorName string
		if strings.Contains(security.Symbol, "P") && strings.HasSuffix(security.Symbol, "P") {
			sectorName = "Promoter Share"
		} else {
			sectorName = security.SectorName
			if sectorName == "" {
				sectorName = "Others"
			}
		}

		if sectorScrips[sectorName] == nil {
			sectorScrips[sectorName] = make([]string, 0)
		}
		sectorScrips[sectorName] = append(sectorScrips[sectorName], security.Symbol)
	}

	return sectorScrips, nil
}

// FindSecurity finds a security by ID.
func (h *nepseClient) FindSecurity(ctx context.Context, securityID int32) (*Security, error) {
	return h.findSecurityByID(ctx, securityID)
}

// FindSecurityBySymbol finds a security by symbol.
func (h *nepseClient) FindSecurityBySymbol(ctx context.Context, symbol string) (*Security, error) {
	return h.findSecurityBySymbol(ctx, symbol)
}

func (h *nepseClient) findSecurityByID(ctx context.Context, id int32) (*Security, error) {
	if id <= 0 {
		return nil, NewInvalidClientRequestError("security ID must be positive")
	}

	securities, err := h.GetSecurityList(ctx)
	if err != nil {
		return nil, err
	}

	for _, security := range securities {
		if security.ID == id {
			return &security, nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("security with ID %d", id))
}

func (h *nepseClient) findSecurityBySymbol(ctx context.Context, symbol string) (*Security, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, NewInvalidClientRequestError("symbol cannot be empty")
	}

	securities, err := h.GetSecurityList(ctx)
	if err != nil {
		return nil, err
	}

	for _, security := range securities {
		if security.Symbol == symbol {
			return &security, nil
		}
	}

	return nil, NewNotFoundError("security with symbol " + symbol)
}

// GetFloorSheet retrieves the complete floor sheet data.
func (h *nepseClient) GetFloorSheet(ctx context.Context) ([]FloorSheetEntry, error) {
	endpoint := fmt.Sprintf("%s?size=500&sort=contractId,desc", h.config.APIEndpoints["floor_sheet"])

	var floorSheetArray []FloorSheetEntry
	if err := h.apiRequest(ctx, endpoint, &floorSheetArray); err == nil {
		return floorSheetArray, nil
	}

	var firstPage FloorSheetResponse
	if err := h.apiRequest(ctx, endpoint, &firstPage); err != nil {
		return nil, err
	}

	all := firstPage.FloorSheets.Content
	total := firstPage.FloorSheets.TotalPages
	for p := int32(1); p < total; p++ {
		pageEndpoint := fmt.Sprintf("%s&page=%d", endpoint, p)
		var page FloorSheetResponse
		if err := h.apiRequest(ctx, pageEndpoint, &page); err != nil {
			return nil, err
		}
		all = append(all, page.FloorSheets.Content...)
	}
	return all, nil
}

// GetFloorSheetOf retrieves floor sheet data for a specific security on a specific business date by ID.
func (h *nepseClient) GetFloorSheetOf(ctx context.Context, securityID int32, businessDate string) ([]FloorSheetEntry, error) {
	endpoint := fmt.Sprintf("%s%d?businessDate=%s&size=500&sort=contractid,desc",
		h.config.APIEndpoints["company_floorsheet"], securityID, businessDate)

	var firstPage FloorSheetResponse
	if err := h.apiRequest(ctx, endpoint, &firstPage); err != nil {
		return nil, err
	}

	if len(firstPage.FloorSheets.Content) == 0 {
		return []FloorSheetEntry{}, nil
	}

	allEntries := firstPage.FloorSheets.Content
	totalPages := firstPage.FloorSheets.TotalPages

	for page := int32(1); page < totalPages; page++ {
		pageEndpoint := fmt.Sprintf("%s&page=%d", endpoint, page)

		var pageResponse FloorSheetResponse
		if err := h.apiRequest(ctx, pageEndpoint, &pageResponse); err != nil {
			return nil, err
		}

		allEntries = append(allEntries, pageResponse.FloorSheets.Content...)
	}

	return allEntries, nil
}

// GetFloorSheetBySymbol retrieves floor sheet data for a specific security by symbol.
func (h *nepseClient) GetFloorSheetBySymbol(ctx context.Context, symbol string, businessDate string) ([]FloorSheetEntry, error) {
	security, err := h.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return h.GetFloorSheetOf(ctx, security.ID, businessDate)
}
