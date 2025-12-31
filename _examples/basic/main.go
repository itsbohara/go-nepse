package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/voidarchive/go-nepse"
)

func main() {
	withGraphs := flag.Bool("with-graphs", false, "include graph endpoints in the run")
	withFloor := flag.Bool("with-floorsheet", false, "include floorsheet endpoints in the run")
	symbolFlag := flag.String("symbol", "NABIL", "symbol to use for symbol-based calls")
	bizDateFlag := flag.String("business-date", "", "business date (YYYY-MM-DD) for today's prices and floorsheet; defaults to last weekday")
	flag.Parse()

	fmt.Println("NEPSE Go Library - Full API Example")
	fmt.Println("====================================")

	opts := nepse.DefaultOptions()
	opts.TLSVerification = false // For development only

	client, err := nepse.NewClient(opts)
	if err != nil {
		log.Fatalf("Failed to create NEPSE client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Close client: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	symbol := *symbolFlag
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	start := startDate.Format("2006-01-02")
	end := startDate.Format("2006-01-02")
	today := now.Format("2006-01-02")
	userBizDate := *bizDateFlag

	// 1) Market data
	fmt.Println("\n[Market] Summary")
	if s, err := client.GetMarketSummary(ctx); err != nil {
		log.Printf("Market summary: %v", err)
	} else {
		fmt.Printf("- Turnover: %.2f, Trades: %.0f\n", s.TotalTurnover, s.TotalTransactions)
	}

	fmt.Println("[Market] Status")
	marketOpen := false
	if st, err := client.GetMarketStatus(ctx); err != nil {
		log.Printf("Market status: %v", err)
	} else {
		fmt.Printf("- Open: %v (%s)\n", st.IsMarketOpen(), st.IsOpen)
		marketOpen = st.IsMarketOpen()
	}

	fmt.Println("[Market] NEPSE Index + Sub-Indices")
	if idx, err := client.GetNepseIndex(ctx); err != nil {
		log.Printf("NEPSE index: %v", err)
	} else {
		fmt.Printf("- NEPSE: %.2f (%.2f%%)\n", idx.IndexValue, idx.PercentChange)
	}
	if subs, err := client.GetNepseSubIndices(ctx); err != nil {
		log.Printf("Sub-indices: %v", err)
	} else {
		fmt.Printf("- Sub-Indices: %d\n", len(subs))
	}

	fmt.Println("[Market] Live + Supply/Demand")
	if live, err := client.GetLiveMarket(ctx); err != nil {
		log.Printf("Live market: %v", err)
	} else {
		fmt.Printf("- Live entries: %d\n", len(live))
	}
	if sd, err := client.GetSupplyDemand(ctx); err != nil {
		log.Printf("Supply/Demand: %v", err)
	} else {
		fmt.Printf("- Supply: %d entries, Demand: %d entries\n", len(sd.SupplyList), len(sd.DemandList))
	}

	// 2) Securities & Companies
	fmt.Println("\n[Securities] Lists & Details")
	var nabilID int32
	if secs, err := client.GetSecurityList(ctx); err != nil {
		log.Printf("Security list: %v", err)
	} else {
		fmt.Printf("- Securities: %d\n", len(secs))
	}
	if cos, err := client.GetCompanyList(ctx); err != nil {
		log.Printf("Company list: %v", err)
	} else {
		fmt.Printf("- Companies: %d\n", len(cos))
	}
	if sec, err := client.FindSecurityBySymbol(ctx, symbol); err != nil {
		log.Printf("Find security %s: %v", symbol, err)
	} else {
		nabilID = sec.ID
		fmt.Printf("- %s ID: %d\n", symbol, nabilID)
		fmt.Printf("- Company: %s (%s)\n", sec.SecurityName, sec.SectorName)
	}
	if nabilID != 0 {
		if det, err := client.GetCompanyDetails(ctx, nabilID); err != nil {
			log.Printf("Company details %d: %v", nabilID, err)
		} else {
			fmt.Printf("- %s Close: %.2f, LTP: %.2f\n", det.Symbol, det.ClosePrice, det.LastTradedPrice)
		}
	}
	if ss, err := client.GetSectorScrips(ctx); err != nil {
		log.Printf("Sector scrips: %v", err)
	} else {
		fmt.Printf("- Sectors: %d\n", len(ss))
	}

	// 3) Price & Trading
	fmt.Println("\n[Trading] Today, History, Depth, Floor")
	effBizDate := userBizDate
	if effBizDate == "" {
		if marketOpen {
			effBizDate = today
		} else {
			effBizDate = lastWeekday(now).Format("2006-01-02")
		}
	}
	if todays, err := client.GetTodaysPrices(ctx, effBizDate); err != nil {
		log.Printf("Today's prices (%s): %v", effBizDate, err)
	} else {
		fmt.Printf("- Today prices (%s): %d\n", effBizDate, len(todays))
	}
	if nabilID != 0 {
		if hist, err := client.GetPriceVolumeHistory(ctx, nabilID, start, end); err != nil {
			log.Printf("History %d: %v", nabilID, err)
		} else {
			fmt.Printf("- History records: %d\n", len(hist))
		}
		if md, err := client.GetMarketDepthBySymbol(ctx, symbol); err != nil {
			fmt.Printf("- Market depth(%s): unavailable (%v)\n", symbol, err)
		} else {
			fmt.Printf("- Depth(%s) buy/sell qty: %d/%d, levels: %d/%d\n", symbol, md.TotalBuyQty, md.TotalSellQty, len(md.BuyDepth), len(md.SellDepth))
		}
		if *withFloor {
			if fs, err := client.GetFloorSheetOf(ctx, nabilID, effBizDate); err != nil {
				fmt.Printf("- Floorsheet(%s): error (%v)\n", effBizDate, err)
			} else {
				fmt.Printf("- Floorsheet(%s): %d\n", effBizDate, len(fs))
			}
		}
	}
	if *withFloor {
		if fsAll, err := client.GetFloorSheet(ctx); err != nil {
			log.Printf("Floorsheet(all): %v", err)
		} else {
			fmt.Printf("- Floorsheet(all): %d\n", len(fsAll))
		}
	}

	// 4) Top Lists
	fmt.Println("\n[Top] Gainers/Losers/Trade/Transaction/Turnover")
	if v, err := client.GetTopGainers(ctx); err != nil {
		log.Printf("Top gainers: %v", err)
	} else {
		fmt.Printf("- Gainers: %d\n", len(v))
	}
	if v, err := client.GetTopLosers(ctx); err != nil {
		log.Printf("Top losers: %v", err)
	} else {
		fmt.Printf("- Losers: %d\n", len(v))
	}
	if v, err := client.GetTopTenTrade(ctx); err != nil {
		log.Printf("Top trade: %v", err)
	} else {
		fmt.Printf("- Top trade: %d\n", len(v))
	}
	if v, err := client.GetTopTenTransaction(ctx); err != nil {
		log.Printf("Top transaction: %v", err)
	} else {
		fmt.Printf("- Top transaction: %d\n", len(v))
	}
	if v, err := client.GetTopTenTurnover(ctx); err != nil {
		log.Printf("Top turnover: %v", err)
	} else {
		fmt.Printf("- Top turnover: %d\n", len(v))
	}

	// 5) Graphs
	if *withGraphs {
		fmt.Println("\n[Graphs] Main indices")
		if g, err := client.GetDailyNepseIndexGraph(ctx); err != nil {
			log.Printf("NEPSE graph: %v", err)
		} else {
			fmt.Printf("- NEPSE pts: %d\n", len(g.Data))
		}
		if g, err := client.GetDailySensitiveIndexGraph(ctx); err != nil {
			log.Printf("Sensitive graph: %v", err)
		} else {
			fmt.Printf("- Sensitive pts: %d\n", len(g.Data))
		}
		if g, err := client.GetDailyFloatIndexGraph(ctx); err != nil {
			log.Printf("Float graph: %v", err)
		} else {
			fmt.Printf("- Float pts: %d\n", len(g.Data))
		}
		if g, err := client.GetDailySensitiveFloatIndexGraph(ctx); err != nil {
			log.Printf("Sensitive Float graph: %v", err)
		} else {
			fmt.Printf("- Sensitive Float pts: %d\n", len(g.Data))
		}

		fmt.Println("[Graphs] Sub-indices")
		if g, err := client.GetDailyBankSubindexGraph(ctx); err != nil {
			log.Printf("Banking graph: %v", err)
		} else {
			fmt.Printf("- Banking pts: %d\n", len(g.Data))
		}

		fmt.Println("[Graphs] Company")
		if g, err := client.GetDailyScripPriceGraphBySymbol(ctx, symbol); err != nil {
			log.Printf("Company graph %s: %v", symbol, err)
		} else {
			fmt.Printf("- %s graph pts: %d\n", symbol, len(g.Data))
		}
	}

	fmt.Println("\nFinished exercising all public APIs.")
}

func lastWeekday(t time.Time) time.Time {
	switch t.Weekday() {
	case time.Saturday:
		return t.AddDate(0, 0, -1)
	case time.Sunday:
		return t.AddDate(0, 0, -2)
	default:
		return t
	}
}
