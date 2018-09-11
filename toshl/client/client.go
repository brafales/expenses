package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is a client to talk to Toshl's REST API
type Client struct {
	AuthToken   string
	AccountID   string
	HTTPClient  http.Client
	CategoryMap map[string]string
}

// Expense defines a Toshl Entry
type Expense struct {
	Amount      int       `json:"amount"`
	Created     time.Time `json:"created"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
}

type toshlCurrency struct {
	Code string `json:"code"`
}

type toshlEntry struct {
	Amount   float64       `json:"amount"`
	Currency toshlCurrency `json:"currency"`
	Date     string        `json:"date"`
	Desc     string        `json:"desc"`
	Account  string        `json:"account"`
	Category string        `json:"category"`
}

type toshlCategories map[string]string

// CreateEntry entry will create a new entry in Toshl based on the Expense data
func (c *Client) CreateEntry(expense *Expense) error {
	client := http.Client{}
	data, err := c.newToshlEntry(expense)
	if err != nil {
		return err
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "ttps://api.toshl.com/entries", bytes.NewBuffer(dataJSON))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.AuthToken, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respCode := resp.StatusCode
	if respCode != 200 {
		return fmt.Errorf("Something went wrong, response code was %d", respCode)
	}

	return nil
}

func (c *Client) newToshlEntry(expense *Expense) (toshlEntry, error) {
	category, err := c.categoryID(expense.Category)
	if err != nil {
		return toshlEntry{}, err
	}

	entry := toshlEntry{
		Amount: float64(expense.Amount) / 100,
		Currency: toshlCurrency{
			Code: expense.Currency,
		},
		Date:     expense.Created.Format("2000-12-31"),
		Desc:     expense.Description,
		Account:  c.AccountID,
		Category: category,
	}

	return entry, nil
}

func (c *Client) categoryID(categoryName string) (string, error) {
	mappedCategoryName := c.CategoryMap[categoryName]
	if mappedCategoryName == "" {
		return "", fmt.Errorf("Category %s not in the mapping", categoryName)
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.toshl.com/categories", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(c.AuthToken, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var categories toshlCategories
	err = json.NewDecoder(resp.Body).Decode(categories)
	if err != nil {
		return "", err
	}
	category := categories[mappedCategoryName]
	if category == "" {
		return "", fmt.Errorf("Category not found in Toshl for name %s", mappedCategoryName)
	}
	return category, nil
}
