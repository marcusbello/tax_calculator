package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// TaxData represents the data structure for storing tax information.
type TaxData struct {
    ID             string
	AnnualIncome string
	Rent string
	Investments string
	TaxAmount uint64
}

var (
    templates = template.Must(template.ParseGlob("templates/*.html"))
    store     = make(map[string]TaxData)
    mu        sync.Mutex
	logger = log.New(os.Stdout, "tax-calculator: ", log.LstdFlags|log.Lshortfile)
)

func main() {
    rand.Seed(time.Now().UnixNano())

	mux := http.NewServeMux()

    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    mux.HandleFunc("/", formHandler)
    mux.HandleFunc("/tax-calculator", calculateTaxHandler)
    mux.HandleFunc("/tax/", func(w http.ResponseWriter, r *http.Request) {
       	taxHandler(w, r)
    })
	// Setting up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	go func() {
		logger.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP:
			logger.Println("Received SIGHUP (reload requested)")
		case syscall.SIGINT, syscall.SIGTERM:
			logger.Printf("Received %s — shutting down...", sig)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				logger.Fatalf("Graceful shutdown failed: %v", err)
			}
			logger.Println("Server stopped cleanly.")
			return
		default:
			logger.Printf("Unhandled signal: %v", sig)
		}
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "form.html", nil)
}

func calculateTaxHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    annualEarning := r.FormValue("annualIncome")
	rentAmount := r.FormValue("rentAmount")
	businessExpense := r.FormValue("businessExpense")

    id := fmt.Sprintf("%d", rand.Intn(1000000))

	taxAmount, err := taxCalculator(annualEarning, rentAmount, businessExpense)
	if err != nil {
		logger.Print("taxCalculator: ", err)
		return
	}
	logger.Printf("Tax amount calculated: %d", taxAmount)
    data := TaxData{
        ID:          id,
		AnnualIncome: annualEarning,
		Rent: rentAmount,
		Investments: businessExpense,
		TaxAmount: taxAmount,
    }

    mu.Lock()
    store[id] = data
    mu.Unlock()

	http.Redirect(w, r, "/tax/"+id, http.StatusSeeOther)
}

func taxHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/tax/"):]
    mu.Lock()
    data, ok := store[id]
    mu.Unlock()

    if !ok {
        http.NotFound(w, r)
        return
    }

    templates.ExecuteTemplate(w, "tax.html", data)
}

func parseOrZero(s string) uint64 {
	if s == "0" {
		return 0
	}
	if s == "" {
		return 0
	}
	if decimal := strings.HasSuffix(s, ".00"); s != "" && decimal {
		s = strings.TrimSuffix(s, ".00")
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logger.Printf("parseOrZero: error parsing '%s': %v", s, err)
		return 0 // or handle error differently
	}
	return uint64(n)
}

func percentageOf(percent, amount int64) uint64 {
	return uint64((percent * amount) / 100)
}

func taxCalculator(annualEarnings, rentAmount, businessExpenses string) (uint64, error) {
	logger.Printf("annualEarning: %s; rentAmount: %s; businessExpenses: %s;", annualEarnings, rentAmount, businessExpenses)
	annualIncome := parseOrZero(annualEarnings)
	// base case if annual income is less than 800k
	baseValue := uint64(800_000)
	if annualIncome < baseValue {
		return 0, nil
	}
	// Rate is the tax brackets
	type Rate struct {
		Amount uint64
		Percentage string
		Payment uint64 // This feild  is for storing the payment amount i.e 15% of 2,200,000
	}
	rates := [5]Rate{
		{
			Amount: 2_200_000,
			Percentage: "15%",
			Payment: percentageOf(15, 2_200_000),
		},
		{
			Amount: 9_000_000,
			Percentage: "18%",
			Payment: percentageOf(18, 9_000_000),
		},
		{
			Amount: 13_000_000,
			Percentage: "21%",
			Payment: percentageOf(21, 13_000_000),
		},
		{
			Amount: 25_000_000,
			Percentage: "23%",
			Payment: percentageOf(23, 25_000_000),
		},
		{
			Amount: 50_000_000,
			Percentage: "25%",
			Payment: percentageOf(25, 50_000_000),
		},
	}
	// const
	var taxAmount uint64
	var lastRate uint64
	// rate calculator
	for i, rate := range rates {
		if annualIncome > rate.Amount {
			taxAmount += rate.Payment
			annualIncome -= rate.Amount
			if i+1 < len(rates) {
				lastRate = rates[i+1].Payment
			}
		}
	}
	// If there is any remaining annual income, apply the last rate
	if annualIncome > 0 {
		if lastRate == 0 {
			lastRate =  rates[0].Payment // Default to the first rate if no last rate is set
		}
		taxAmount += lastRate
		logger.Printf("Applying last rate: %v; taxAmount: %d", lastRate, taxAmount)
	}
	// rent and investments
	rent := parseOrZero(rentAmount)
	investments := parseOrZero(businessExpenses)
	// rentRefund
	if rent > 0 {
		rentRefund := percentageOf(20, int64(rent))
		logger.Printf("rentRefund: %d", rentRefund)
		if rentRefund > 500_000 {
			taxAmount -= 500_000
		} else {
			taxAmount -= rentRefund
		}
	}
	// investment refund
	if investments > 0 {
		taxAmount -= investments
	}

	return taxAmount, nil
}

