package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mickaelroger/mcp-freshrss/models"
)

type FeverClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewFeverClient(baseURL, apiKey string) *FeverClient {
	return &FeverClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *FeverClient) makeRequest(ctx context.Context, endpoint string, extraParams map[string]string) ([]byte, error) {
	apiURL := fmt.Sprintf("%s?api", c.baseURL)
	if endpoint != "" {
		apiURL = fmt.Sprintf("%s&%s", apiURL, endpoint)
	}

	formData := url.Values{}
	formData.Set("api_key", c.apiKey)

	for key, value := range extraParams {
		formData.Set(key, value)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

func (c *FeverClient) GetGroups(ctx context.Context) (*models.GroupsResponse, error) {
	body, err := c.makeRequest(ctx, "groups", nil)
	if err != nil {
		return nil, err
	}

	var response models.GroupsResponse
	if err := parseJSON(body, &response); err != nil {
		return nil, err
	}

	if response.Auth != 1 {
		return nil, fmt.Errorf("authentication failed")
	}

	return &response, nil
}

func (c *FeverClient) GetFeeds(ctx context.Context) (*models.FeedsResponse, error) {
	body, err := c.makeRequest(ctx, "feeds", nil)
	if err != nil {
		return nil, err
	}

	var response models.FeedsResponse
	if err := parseJSON(body, &response); err != nil {
		return nil, err
	}

	if response.Auth != 1 {
		return nil, fmt.Errorf("authentication failed")
	}

	return &response, nil
}

func (c *FeverClient) GetUnreadItemIDs(ctx context.Context) ([]int, error) {
	body, err := c.makeRequest(ctx, "unread_item_ids", nil)
	if err != nil {
		return nil, err
	}

	var response models.UnreadItemIDsResponse
	if err := parseJSON(body, &response); err != nil {
		return nil, err
	}

	if response.Auth != 1 {
		return nil, fmt.Errorf("authentication failed")
	}

	return parseIDs(response.UnreadItemIDs), nil
}

func (c *FeverClient) GetUnreadItemStringIDs(ctx context.Context) ([]string, error) {
	body, err := c.makeRequest(ctx, "unread_item_ids", nil)
	if err != nil {
		return nil, err
	}

	var response models.UnreadItemIDsResponse
	if err := parseJSON(body, &response); err != nil {
		return nil, err
	}

	if response.Auth != 1 {
		return nil, fmt.Errorf("authentication failed")
	}

	return parseStringIDs(response.UnreadItemIDs), nil
}

func (c *FeverClient) GetItem(ctx context.Context, itemID int) (*models.Item, error) {
	body, err := c.makeRequest(ctx, "items", map[string]string{
		"with_ids": strconv.Itoa(itemID),
	})
	if err != nil {
		return nil, err
	}

	var response models.ItemsResponse
	if err := parseJSON(body, &response); err != nil {
		return nil, err
	}

	if response.Auth != 1 {
		return nil, fmt.Errorf("authentication failed")
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("item %d not found", itemID)
	}

	return &response.Items[0], nil
}

func (c *FeverClient) GetItems(ctx context.Context) ([]models.Item, error) {
	body, err := c.makeRequest(ctx, "items", nil)
	if err != nil {
		return nil, err
	}

	var response models.ItemsResponse
	if err := parseJSON(body, &response); err != nil {
		return nil, err
	}

	if response.Auth != 1 {
		return nil, fmt.Errorf("authentication failed")
	}

	return response.Items, nil
}

func (c *FeverClient) GetItemsByIDs(ctx context.Context, ids []string) ([]models.Item, error) {
	if len(ids) == 0 {
		return []models.Item{}, nil
	}

	idStr := strings.Join(ids, ",")
	body, err := c.makeRequest(ctx, "items", map[string]string{
		"with_ids": idStr,
	})
	if err != nil {
		return nil, err
	}

	var response models.ItemsResponse
	if err := parseJSON(body, &response); err != nil {
		return nil, err
	}

	if response.Auth != 1 {
		return nil, fmt.Errorf("authentication failed")
	}

	return response.Items, nil
}

func (c *FeverClient) MarkItemRead(ctx context.Context, itemID int) error {
	return c.markItem(ctx, itemID, "read")
}

func (c *FeverClient) MarkItemUnread(ctx context.Context, itemID int) error {
	return c.markItem(ctx, itemID, "unread")
}

func (c *FeverClient) markItem(ctx context.Context, itemID int, status string) error {
	body, err := c.makeRequest(ctx, "", map[string]string{
		"mark": "item",
		"as":   status,
		"id":   strconv.Itoa(itemID),
	})
	if err != nil {
		return err
	}

	var response models.BaseResponse
	if err := parseJSON(body, &response); err != nil {
		return err
	}

	if response.Auth != 1 {
		return fmt.Errorf("authentication failed")
	}

	return nil
}

func parseJSON(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}

func parseIDs(idsStr string) []int {
	if idsStr == "" {
		return []int{}
	}

	parts := strings.Split(idsStr, ",")
	ids := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err == nil {
			ids = append(ids, int(id))
		}
	}

	return ids
}

func parseStringIDs(idsStr string) []string {
	if idsStr == "" {
		return []string{}
	}

	parts := strings.Split(idsStr, ",")
	ids := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			ids = append(ids, part)
		}
	}

	return ids
}
