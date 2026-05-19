package v1

type (
	AdminReviewListReq struct {
		Page     int    `form:"page,default=1"`
		PageSize int    `form:"page_size,default=20"`
		Status   string `form:"status"`
		SiteID   int    `form:"site_id"`
	}

	AdminReviewItem struct {
		ID        int    `json:"id"`
		SiteID    int    `json:"site_id"`
		SiteTitle string `json:"site_title"`
		Nickname  string `json:"nickname"`
		Content   string `json:"content"`
		Status    string `json:"status"`
		AdminNote string `json:"admin_note"`
		CreatedAt string `json:"created_at"`
	}

	AdminReviewListResp struct {
		List       []AdminReviewItem `json:"list"`
		Pagination SiteLisPagination `json:"pagination"`
	}

	AdminReviewUpdateReq struct {
		Status    string `json:"status" binding:"required,oneof=approved rejected"`
		AdminNote string `json:"admin_note"`
	}
	AdminReviewUpdateResp struct{}
)

type (
	AdminRatingItem struct {
		ID        int     `json:"id"`
		SiteID    int     `json:"site_id"`
		SiteTitle string  `json:"site_title"`
		VisitorID string  `json:"visitor_id"`
		Score     int     `json:"score"`
		CreatedAt string  `json:"created_at"`
	}

	AdminRatingListResp struct {
		List       []AdminRatingItem `json:"list"`
		Pagination SiteLisPagination `json:"pagination"`
	}

	AdminRatingUpdateReq struct {
		Score int `json:"score" binding:"required,min=1,max=5"`
	}
	AdminRatingUpdateResp struct{}

	AdminRatingDeleteResp struct{}
)

type (
	DictItem struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Icon   string `json:"icon"`
		Sort   int    `json:"sort"`
		IsUsed bool   `json:"is_used"`
	}

	DictCreateReq struct {
		Name     string `json:"name" binding:"required"`
		Icon     string `json:"icon"`
		Sort     int    `json:"sort"`
		IsUsed   bool   `json:"is_used"`
		Provider string `json:"provider"`
		Color    string `json:"color"`
	}
	DictCreateResp struct {
		ID int `json:"id"`
	}

	DictUpdateReq struct {
		Name     string `json:"name"`
		Icon     string `json:"icon"`
		Sort     int    `json:"sort"`
		IsUsed   *bool  `json:"is_used"`
		Provider string `json:"provider"`
		Color    string `json:"color"`
	}
	DictUpdateResp struct{}

	DictDeleteResp struct{}

	ModelListResp struct {
		List []ModelItem `json:"list"`
	}

	ModelItem struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Provider string `json:"provider"`
		Icon     string `json:"icon"`
		Sort     int    `json:"sort"`
		IsUsed   bool   `json:"is_used"`
	}
)

type (
	HealthConfigResp struct {
		IntervalMinutes    int  `json:"interval_minutes"`
		MaxConsecutiveFails int `json:"max_consecutive_fails"`
		Enabled            bool `json:"enabled"`
	}

	HealthConfigUpdateReq struct {
		IntervalMinutes     *int  `json:"interval_minutes"`
		MaxConsecutiveFails *int  `json:"max_consecutive_fails"`
		Enabled             *bool `json:"enabled"`
	}
)

type (
	AdItem struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Slot        string `json:"slot"`
		ContentType string `json:"content_type"`
		Content     string `json:"content"`
		ImageURL    string `json:"image_url"`
		LinkURL     string `json:"link_url"`
		StartAt     string `json:"start_at"`
		EndAt       string `json:"end_at"`
		Sort        int    `json:"sort"`
		IsUsed      bool   `json:"is_used"`
		CreatedAt   string `json:"created_at"`
	}

	AdListReq struct {
		Page     int    `form:"page,default=1"`
		PageSize int    `form:"page_size,default=20"`
		Slot     string `form:"slot"`
	}

	AdListResp struct {
		List       []AdItem          `json:"list"`
		Pagination SiteLisPagination `json:"pagination"`
	}

	AdCreateReq struct {
		Name        string `json:"name" binding:"required"`
		Slot        string `json:"slot" binding:"required"`
		ContentType string `json:"content_type"`
		Content     string `json:"content"`
		ImageURL    string `json:"image_url"`
		LinkURL     string `json:"link_url"`
		StartAt     string `json:"start_at"`
		EndAt       string `json:"end_at"`
		Sort        int    `json:"sort"`
		IsUsed      bool   `json:"is_used"`
	}

	AdCreateResp struct {
		ID int `json:"id"`
	}

	AdUpdateReq struct {
		Name        string `json:"name"`
		Slot        string `json:"slot"`
		ContentType string `json:"content_type"`
		Content     string `json:"content"`
		ImageURL    string `json:"image_url"`
		LinkURL     string `json:"link_url"`
		StartAt     string `json:"start_at"`
		EndAt       string `json:"end_at"`
		Sort        int    `json:"sort"`
		IsUsed      *bool  `json:"is_used"`
	}

	AdUpdateResp struct{}

	AdDeleteResp struct{}
)

type (
	HomepageConfigResp struct {
		Sections         []HomepageSection `json:"sections"`
		DefaultSort      string            `json:"default_sort"`
		CardsPerRow      int               `json:"cards_per_row"`
		ShowFeatured     bool              `json:"show_featured"`
		ShowRecommended  bool              `json:"show_recommended"`
		FeaturedSiteIDs  []int             `json:"featured_site_ids"`
	}

	HomepageSection struct {
		Key      string `json:"key"`
		Title    string `json:"title"`
		Enabled  bool   `json:"enabled"`
		Collapsed bool  `json:"collapsed"`
	}

	HomepageConfigUpdateReq struct {
		Sections         []HomepageSection `json:"sections"`
		DefaultSort      *string           `json:"default_sort"`
		CardsPerRow      *int              `json:"cards_per_row"`
		ShowFeatured     *bool             `json:"show_featured"`
		ShowRecommended  *bool             `json:"show_recommended"`
		FeaturedSiteIDs  []int             `json:"featured_site_ids"`
	}
)
