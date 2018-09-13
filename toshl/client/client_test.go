package client_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/brafales/expenses/toshl/client"
)

const token = "token"
const accountID = "1234"
const amount = 100.0
const currency = "GBP"
const description = "description"
const category = "category"
const categoriesJSON = `[
	{
	  "id": "53459076",
	  "name": "Category",
	  "name_override": true,
	  "modified": "2016-12-29 12:25:07.010",
	  "type": "expense",
	  "deleted": false
	}
]`

var categories = map[string]string{"category": "Category"}
var created = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func TestHappyPath(t *testing.T) {
	// Mock calls to Toshl
	// /categories to return a list of fake categories
	// /entries to return a `created` status code
	// we will make the test fail if calls are made with the wrong parameters
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/categories" {
			if r.Method != "GET" {
				t.Errorf("Expected GET request, got ‘%s’", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, categoriesJSON)
			return
		}
		if r.URL.String() == "/entries" {
			if r.Method != "POST" {
				t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
			}
			reqJson, err := simplejson.NewFromReader(r.Body)
			if err != nil {
				t.Errorf("Error while reading request JSON: %s", err)
			}
			reqAmount := reqJson.GetPath("amount").MustFloat64()
			if reqAmount != amount/100 {
				t.Errorf("Expected request JSON to have amount = 100.0, got %f", reqAmount)
			}
			reqCurrency := reqJson.GetPath("currency", "code").MustString()
			if reqCurrency != currency {
				t.Errorf("Expected request JSON to have currency = %s, got %s", currency, reqCurrency)
			}
			reqDate := reqJson.GetPath("date").MustString()
			if reqDate != created.Format("2006-01-02") {
				t.Errorf("Expected request JSON to have date = %s, got %s", created.Format("2006-01-02"), reqDate)
			}
			reqDescription := reqJson.GetPath("desc").MustString()
			if reqDescription != description {
				t.Errorf("Expected request JSON to have desc = %s, got %s", description, reqDescription)
			}
			reqAccount := reqJson.GetPath("account").MustString()
			if reqAccount != accountID {
				t.Errorf("Expected request JSON to have account = %s, got %s", accountID, reqAccount)
			}
			reqCategory := reqJson.GetPath("category").MustString()
			if reqCategory != "53459076" {
				t.Errorf("Expected request JSON to have category = %s, got %s", "53459076", reqCategory)
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
		t.Errorf("Unknown request made to Toshl")
	}))
	defer ts.Close()

	toshlClient := client.Client{
		AuthToken:    token,
		AccountID:    accountID,
		HTTPClient:   *ts.Client(),
		CategoryMap:  categories,
		ToshlBaseURL: ts.URL,
	}

	expense := client.Expense{
		Category:    "category",
		Amount:      amount,
		Created:     created,
		Currency:    currency,
		Description: description,
	}

	err := toshlClient.CreateEntry(&expense)
	if err != nil {
		t.Error(err)
	}
}
