package models

type Group struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Feed struct {
	ID              int    `json:"id"`
	FaviconID       int    `json:"favicon_id"`
	Title           string `json:"title"`
	URL             string `json:"url"`
	SiteURL         string `json:"site_url"`
	IsSpark         int    `json:"is_spark"`
	LastUpdatedTime int64  `json:"last_updated_on_time"`
}

type FeedsGroup struct {
	GroupID int    `json:"group_id"`
	FeedIDs string `json:"feed_ids"`
}

type Item struct {
	ID            string `json:"id"`
	FeedID        int    `json:"feed_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	HTML          string `json:"html"`
	URL           string `json:"url"`
	IsSaved       int    `json:"is_saved"`
	IsRead        int    `json:"is_read"`
	CreatedOnTime int64  `json:"created_on_time"`
}

type GroupsResponse struct {
	APIVersion          int          `json:"api_version"`
	Auth                int          `json:"auth"`
	LastRefreshedOnTime int64        `json:"last_refreshed_on_time"`
	Groups              []Group      `json:"groups"`
	FeedsGroups         []FeedsGroup `json:"feeds_groups"`
}

type FeedsResponse struct {
	APIVersion          int          `json:"api_version"`
	Auth                int          `json:"auth"`
	LastRefreshedOnTime int64        `json:"last_refreshed_on_time"`
	Feeds               []Feed       `json:"feeds"`
	FeedsGroups         []FeedsGroup `json:"feeds_groups"`
}

type UnreadItemIDsResponse struct {
	APIVersion          int    `json:"api_version"`
	Auth                int    `json:"auth"`
	LastRefreshedOnTime int64  `json:"last_refreshed_on_time"`
	UnreadItemIDs       string `json:"unread_item_ids"`
}

type ItemsResponse struct {
	APIVersion          int    `json:"api_version"`
	Auth                int    `json:"auth"`
	LastRefreshedOnTime int64  `json:"last_refreshed_on_time"`
	Items               []Item `json:"items"`
	TotalItems          int    `json:"total_items"`
}

type BaseResponse struct {
	APIVersion          int   `json:"api_version"`
	Auth                int   `json:"auth"`
	LastRefreshedOnTime int64 `json:"last_refreshed_on_time"`
}
