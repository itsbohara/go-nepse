package nepse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Index IDs used by NEPSE API.
const (
	nepseIndexID          = 58
	sensitiveIndexID      = 57
	floatIndexID          = 62
	sensitiveFloatIndexID = 63
)

// GetMarketSummary returns aggregate market statistics including turnover, volume, and capitalization.
func (c *Client) GetMarketSummary(ctx context.Context) (*MarketSummary, error) {
	var rawItems []MarketSummaryItem
	if err := c.apiRequest(ctx, c.config.Endpoints.MarketSummary, &rawItems); err != nil {
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

// GetMarketStatus returns whether the market is currently open or closed.
func (c *Client) GetMarketStatus(ctx context.Context) (*MarketStatus, error) {
	var status MarketStatus
	if err := c.apiRequest(ctx, c.config.Endpoints.MarketOpen, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// GetNepseIndex returns the main NEPSE index with current value, change, and 52-week range.
func (c *Client) GetNepseIndex(ctx context.Context) (*NepseIndex, error) {
	var rawIndices []NepseIndexRaw
	if err := c.apiRequest(ctx, c.config.Endpoints.NepseIndex, &rawIndices); err != nil {
		return nil, err
	}

	for i := range rawIndices {
		if rawIndices[i].ID == nepseIndexID {
			return &NepseIndex{
				IndexValue:       rawIndices[i].Close,
				PercentChange:    rawIndices[i].PerChange,
				PointChange:      rawIndices[i].Change,
				High:             rawIndices[i].High,
				Low:              rawIndices[i].Low,
				PreviousClose:    rawIndices[i].PreviousClose,
				FiftyTwoWeekHigh: rawIndices[i].FiftyTwoWeekHigh,
				FiftyTwoWeekLow:  rawIndices[i].FiftyTwoWeekLow,
				CurrentValue:     rawIndices[i].CurrentValue,
				GeneratedTime:    rawIndices[i].GeneratedTime,
			}, nil
		}
	}

	return nil, NewNotFoundError("NEPSE Index")
}

// GetNepseSubIndices returns all sector sub-indices excluding the main composite indices.
func (c *Client) GetNepseSubIndices(ctx context.Context) ([]SubIndex, error) {
	var rawIndices []NepseIndexRaw
	if err := c.apiRequest(ctx, c.config.Endpoints.NepseIndex, &rawIndices); err != nil {
		return nil, err
	}

	// Composite indices to exclude from sub-indices list.
	excluded := map[int32]bool{
		nepseIndexID:          true,
		sensitiveIndexID:      true,
		floatIndexID:          true,
		sensitiveFloatIndexID: true,
	}

	subIndices := make([]SubIndex, 0, len(rawIndices))
	for i := range rawIndices {
		if !excluded[rawIndices[i].ID] {
			subIndices = append(subIndices, SubIndex(rawIndices[i]))
		}
	}

	return subIndices, nil
}

// GetLiveMarket returns real-time price and volume data for all actively traded securities.
func (c *Client) GetLiveMarket(ctx context.Context) ([]LiveMarketEntry, error) {
	var liveMarket []LiveMarketEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.LiveMarket, &liveMarket); err != nil {
		return nil, err
	}
	return liveMarket, nil
}

// SupplyDemandData represents the combined supply and demand response.
type SupplyDemandData struct {
	SupplyList []SupplyDemandItem `json:"supplyList"`
	DemandList []SupplyDemandItem `json:"demandList"`
}

// SupplyDemandItem represents a single item in supply or demand list.
type SupplyDemandItem struct {
	SecurityID   int32  `json:"securityId"`
	Symbol       string `json:"symbol"`
	SecurityName string `json:"securityName"`
	TotalQuantity int64 `json:"totalQuantity"`
	TotalOrder   int32  `json:"totalOrder"`
}

// GetSupplyDemand returns aggregate supply and demand data.
func (c *Client) GetSupplyDemand(ctx context.Context) (*SupplyDemandData, error) {
	var data SupplyDemandData
	if err := c.apiRequest(ctx, c.config.Endpoints.SupplyDemand, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// GetTopGainers returns securities with the highest percentage gains for the trading day.
func (c *Client) GetTopGainers(ctx context.Context) ([]TopListEntry, error) {
	var topGainers []TopListEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopGainers, &topGainers); err != nil {
		return nil, err
	}
	return topGainers, nil
}

// GetTopLosers returns securities with the highest percentage losses for the trading day.
func (c *Client) GetTopLosers(ctx context.Context) ([]TopListEntry, error) {
	var topLosers []TopListEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopLosers, &topLosers); err != nil {
		return nil, err
	}
	return topLosers, nil
}

// GetTopTenTrade returns the ten securities with the highest traded share volume.
func (c *Client) GetTopTenTrade(ctx context.Context) ([]TopListEntry, error) {
	var topTrade []TopListEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopTrade, &topTrade); err != nil {
		return nil, err
	}
	return topTrade, nil
}

// GetTopTenTransaction returns the ten securities with the most transactions.
func (c *Client) GetTopTenTransaction(ctx context.Context) ([]TopListEntry, error) {
	var topTransaction []TopListEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopTransaction, &topTransaction); err != nil {
		return nil, err
	}
	return topTransaction, nil
}

// GetTopTenTurnover returns the ten securities with the highest trading turnover (value).
func (c *Client) GetTopTenTurnover(ctx context.Context) ([]TopListEntry, error) {
	var topTurnover []TopListEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopTurnover, &topTurnover); err != nil {
		return nil, err
	}
	return topTurnover, nil
}

// GetTodaysPrices returns price data for all securities on a given business date.
// If businessDate is empty, returns data for the current trading day.
func (c *Client) GetTodaysPrices(ctx context.Context, businessDate string) ([]TodayPrice, error) {
	endpoint := c.config.Endpoints.TodaysPrice
	if businessDate != "" {
		params := url.Values{}
		params.Set("businessDate", businessDate)
		params.Set("size", "500")
		endpoint += "?" + params.Encode()
	}

	var todayPrices []TodayPrice
	if err := c.apiRequest(ctx, endpoint, &todayPrices); err != nil {
		return nil, err
	}
	return todayPrices, nil
}

// GetPriceVolumeHistory returns historical OHLCV data for a security within a date range.
func (c *Client) GetPriceVolumeHistory(ctx context.Context, securityID int32, startDate, endDate string) ([]PriceHistory, error) {
	params := url.Values{}
	params.Set("size", "500")
	params.Set("startDate", startDate)
	params.Set("endDate", endDate)
	endpoint := fmt.Sprintf("%s/%d?%s", c.config.Endpoints.CompanyPriceHistory, securityID, params.Encode())

	var response struct {
		Content []PriceHistory `json:"content"`
	}

	if err := c.apiRequest(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	return response.Content, nil
}

// GetPriceVolumeHistoryBySymbol returns historical OHLCV data for a security by symbol.
func (c *Client) GetPriceVolumeHistoryBySymbol(ctx context.Context, symbol string, startDate, endDate string) ([]PriceHistory, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.GetPriceVolumeHistory(ctx, security.ID, startDate, endDate)
}

// GetMarketDepth returns the order book (bid/ask levels) for a security.
func (c *Client) GetMarketDepth(ctx context.Context, securityID int32) (*MarketDepth, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.MarketDepth, securityID)

	var raw MarketDepthRaw
	if err := c.apiRequest(ctx, endpoint, &raw); err != nil {
		return nil, err
	}

	return &MarketDepth{
		TotalBuyQty:  raw.TotalBuyQty,
		TotalSellQty: raw.TotalSellQty,
		BuyDepth:     raw.MarketDepth.BuyList,
		SellDepth:    raw.MarketDepth.SellList,
	}, nil
}

// GetMarketDepthBySymbol returns the order book for a security by ticker symbol.
func (c *Client) GetMarketDepthBySymbol(ctx context.Context, symbol string) (*MarketDepth, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.GetMarketDepth(ctx, security.ID)
}

// GetSecurityList returns all tradable securities on the exchange.
func (c *Client) GetSecurityList(ctx context.Context) ([]Security, error) {
	var securities []Security
	if err := c.apiRequest(ctx, c.config.Endpoints.SecurityList, &securities); err != nil {
		return nil, err
	}
	return securities, nil
}

// GetCompanyList returns all listed companies on the exchange.
func (c *Client) GetCompanyList(ctx context.Context) ([]Company, error) {
	var companies []Company
	if err := c.apiRequest(ctx, c.config.Endpoints.CompanyList, &companies); err != nil {
		return nil, err
	}
	return companies, nil
}

// GetCompanyDetails returns comprehensive information including price data for a security.
func (c *Client) GetCompanyDetails(ctx context.Context, securityID int32) (*CompanyDetails, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDetails, securityID)

	var rawDetails CompanyDetailsRaw
	if err := c.apiRequest(ctx, endpoint, &rawDetails); err != nil {
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

// GetCompanyDetailsBySymbol returns comprehensive information for a security by ticker symbol.
func (c *Client) GetCompanyDetailsBySymbol(ctx context.Context, symbol string) (*CompanyDetails, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.GetCompanyDetails(ctx, security.ID)
}

// GetSectorScrips returns a map of sector names to their constituent security symbols.
func (c *Client) GetSectorScrips(ctx context.Context) (SectorScrips, error) {
	// Use company list which includes sector information
	companies, err := c.GetCompanyList(ctx)
	if err != nil {
		return nil, err
	}

	sectorScrips := make(SectorScrips)

	for _, company := range companies {
		sectorName := company.SectorName
		if strings.HasSuffix(company.Symbol, "P") {
			sectorName = "Promoter Share"
		} else if sectorName == "" {
			sectorName = "Others"
		}

		sectorScrips[sectorName] = append(sectorScrips[sectorName], company.Symbol)
	}

	return sectorScrips, nil
}

// FindSecurity returns the security with the given ID.
func (c *Client) FindSecurity(ctx context.Context, securityID int32) (*Security, error) {
	return c.findSecurityByID(ctx, securityID)
}

// FindSecurityBySymbol returns the security with the given ticker symbol.
func (c *Client) FindSecurityBySymbol(ctx context.Context, symbol string) (*Security, error) {
	return c.findSecurityBySymbol(ctx, symbol)
}

func (c *Client) findSecurityByID(ctx context.Context, id int32) (*Security, error) {
	if id <= 0 {
		return nil, NewInvalidClientRequestError("security ID must be positive")
	}

	securities, err := c.GetSecurityList(ctx)
	if err != nil {
		return nil, err
	}

	for i := range securities {
		if securities[i].ID == id {
			return &securities[i], nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("security with ID %d", id))
}

func (c *Client) findSecurityBySymbol(ctx context.Context, symbol string) (*Security, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, NewInvalidClientRequestError("symbol cannot be empty")
	}

	securities, err := c.GetSecurityList(ctx)
	if err != nil {
		return nil, err
	}

	for i := range securities {
		if securities[i].Symbol == symbol {
			return &securities[i], nil
		}
	}

	return nil, NewNotFoundError("security with symbol " + symbol)
}

// GetFloorSheet returns all trades executed on the exchange for the current trading day.
// Handles both array and paginated response formats.
// Note: Returns empty slice if no trades have occurred yet.
func (c *Client) GetFloorSheet(ctx context.Context) ([]FloorSheetEntry, error) {
	params := url.Values{}
	params.Set("size", "500")
	params.Set("sort", "contractId,desc")
	endpoint := c.config.Endpoints.FloorSheet + "?" + params.Encode()

	data, err := c.apiRequestRaw(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	// Try direct array format (may be empty during market hours before trades occur).
	var floorSheetArray []FloorSheetEntry
	if err := json.Unmarshal(data, &floorSheetArray); err == nil {
		return floorSheetArray, nil
	}

	// Try paginated format.
	var firstPage FloorSheetResponse
	if err := json.Unmarshal(data, &firstPage); err != nil {
		return nil, NewInvalidServerResponseError("unrecognized floor sheet response format")
	}

	all := firstPage.FloorSheets.Content
	total := firstPage.FloorSheets.TotalPages
	for p := int32(1); p < total; p++ {
		pageEndpoint := fmt.Sprintf("%s&page=%d", endpoint, p)
		var page FloorSheetResponse
		if err := c.apiRequest(ctx, pageEndpoint, &page); err != nil {
			return nil, err
		}
		all = append(all, page.FloorSheets.Content...)
	}
	return all, nil
}

// GetFloorSheetOf returns all trades for a specific security on a given business date.
// Note: This endpoint may return 403 Forbidden if NEPSE has restricted access.
func (c *Client) GetFloorSheetOf(ctx context.Context, securityID int32, businessDate string) ([]FloorSheetEntry, error) {
	params := url.Values{}
	params.Set("businessDate", businessDate)
	params.Set("size", "500")
	params.Set("sort", "contractid,desc")
	endpoint := fmt.Sprintf("%s/%d?%s", c.config.Endpoints.CompanyFloorsheet, securityID, params.Encode())

	var firstPage FloorSheetResponse
	if err := c.apiRequest(ctx, endpoint, &firstPage); err != nil {
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
		if err := c.apiRequest(ctx, pageEndpoint, &pageResponse); err != nil {
			return nil, err
		}

		allEntries = append(allEntries, pageResponse.FloorSheets.Content...)
	}

	return allEntries, nil
}

// GetFloorSheetBySymbol returns all trades for a specific security by symbol on a given date.
// Note: This endpoint may return 403 Forbidden if NEPSE has restricted access.
func (c *Client) GetFloorSheetBySymbol(ctx context.Context, symbol string, businessDate string) ([]FloorSheetEntry, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.GetFloorSheetOf(ctx, security.ID, businessDate)
}
