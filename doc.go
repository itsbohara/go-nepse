// Package nepse provides a modern, type-safe Go client for the NEPSE (Nepal Stock Exchange) API.
//
// This package offers comprehensive access to NEPSE market data including:
//   - Market summaries and status information
//   - Security and company listings
//   - Real-time and historical price data
//   - Trading volume and floor sheet information
//   - Market indices and sub-indices
//   - Top gainers, losers, and trading statistics
//   - Supply and demand data
//   - Market depth information
//   - Graph data for various indices and securities
//
// The client is built with clean architecture principles, proper error handling,
// and type safety throughout. It provides both high-level convenience methods
// and low-level access to the underlying API endpoints.
//
// Example usage:
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"log"
//
//		"github.com/voidarchive/go-nepse"
//	)
//
//	func main() {
//		// Create a new NEPSE client with default options
//		opts := nepse.DefaultOptions()
//		opts.TLSVerification = false // Required due to NEPSE server TLS issues
//
//		client, err := nepse.NewClient(opts)
//		if err != nil {
//			log.Fatalf("Failed to create NEPSE client: %v", err)
//		}
//		defer client.Close()
//
//		ctx := context.Background()
//
//		// Get market summary
//		summary, err := client.MarketSummary(ctx)
//		if err != nil {
//			log.Fatalf("Failed to get market summary: %v", err)
//		}
//		fmt.Printf("Total Turnover: Rs. %.2f\n", summary.TotalTurnover)
//
//		// Get company details directly by symbol
//		details, err := client.CompanyBySymbol(ctx, "NABIL")
//		if err != nil {
//			log.Fatalf("Failed to get company details: %v", err)
//		}
//		fmt.Printf("Company: %s, LTP: Rs. %.2f\n", details.SecurityName, details.LastTradedPrice)
//	}
package nepse

// Version information
const (
	// Version is the current version of the nepse package.
	Version = "0.1.2"

	// UserAgent is the default user agent string used by the client.
	UserAgent = "go-nepse/" + Version
)

// Common business date formats used by the NEPSE API.
const (
	// DateFormat is the standard date format used by NEPSE API (YYYY-MM-DD).
	DateFormat = "2006-01-02"

	// DateTimeFormat is the standard datetime format used by NEPSE API.
	DateTimeFormat = "2006-01-02 15:04:05"
)

// Sector names commonly used in the NEPSE market.
const (
	SectorBanking          = "Banking"
	SectorDevelopmentBank  = "Development Bank"
	SectorFinance          = "Finance"
	SectorHotelTourism     = "Hotel Tourism"
	SectorHydro            = "Hydro"
	SectorInvestment       = "Investment"
	SectorLifeInsurance    = "Life Insurance"
	SectorManufacturing    = "Manufacturing"
	SectorMicrofinance     = "Microfinance"
	SectorMutualFund       = "Mutual Fund"
	SectorNonLifeInsurance = "Non Life Insurance"
	SectorOthers           = "Others"
	SectorTrading          = "Trading"
	SectorPromoterShare    = "Promoter Share"
)
