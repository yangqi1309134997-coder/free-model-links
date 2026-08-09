package v1

type (
	PublicSiteListReq struct {
		Page       int    `form:"page,default=1"`
		PageSize   int    `form:"page_size,default=20"`
		CategoryID int    `form:"category_id"`
		ModelID    int    `form:"model_id"`
		TagID      int    `form:"tag_id"`
		Search     string `form:"search"`
		SortBy     string `form:"sort_by,default=sort"` // sort | rating | reviews | newest
		SortOrder  string `form:"sort_order,default=asc"`
	}

	PublicSiteItem struct {
		ID            int      `json:"id"`
		Title         string   `json:"title"`
		Icon          string   `json:"icon"`
		URL           string   `json:"url"`
		Description   string   `json:"description"`
		CategoryID    int      `json:"category_id"`
		CategoryName  string   `json:"category_name"`
		AvgRating     float64  `json:"avg_rating"`
		RatingCount   int      `json:"rating_count"`
		ReviewCount   int      `json:"review_count"`
		IsFeatured    bool     `json:"is_featured"`
		IsRecommended bool     `json:"is_recommended"`
		Tags          []string `json:"tags"`
	}

	PublicSiteListResp struct {
		List       []PublicSiteItem  `json:"list"`
		Pagination SiteLisPagination `json:"pagination"`
	}
)

type (
	SiteDetailResp struct {
		ID              int                `json:"id"`
		Title           string             `json:"title"`
		Icon            string             `json:"icon"`
		URL             string             `json:"url"`
		Description     string             `json:"description"`
		Content         string             `json:"content"`
		Screenshots     []string           `json:"screenshots"`
		CategoryID      int                `json:"category_id"`
		CategoryName    string             `json:"category_name"`
		SupportedModels []SiteModelInfo    `json:"supported_models"`
		PricingInfo     *SitePricingInfo   `json:"pricing_info"`
		HealthStatus    *SiteHealthInfo    `json:"health_status"`
		Tags            []SiteTagInfo      `json:"tags"`
		PaymentMethods  []SitePaymentInfo  `json:"payment_methods"`
		AvgRating       float64            `json:"avg_rating"`
		RatingCount     int                `json:"rating_count"`
		ReviewCount     int                `json:"review_count"`
		ViewCount       int                `json:"view_count"`
		IsFeatured      bool               `json:"is_featured"`
		IsRecommended   bool               `json:"is_recommended"`
		CreatedAt       string             `json:"created_at"`
	}

	SiteModelInfo struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Provider string `json:"provider"`
	}

	SitePricingInfo struct {
		MinCharge      float64  `json:"min_charge"`
		RateMultiplier float64  `json:"rate_multiplier"`
		PaymentMethods []string `json:"payment_methods"`
		PricingNote    string   `json:"pricing_note"`
	}

	SiteHealthInfo struct {
		Uptime7d     float64 `json:"uptime_7d"`
		AvgLatencyMs int     `json:"avg_latency_ms"`
		LastCheck    string  `json:"last_check"`
		IsOnline     bool    `json:"is_online"`
	}

	SiteTagInfo struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
		Sort  int    `json:"sort"`
		IsUsed bool  `json:"is_used"`
	}

	SitePaymentInfo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Icon string `json:"icon"`
	}
)

type (
	SubmitRatingReq struct {
		Score int `json:"score" binding:"required,min=1,max=5"`
	}
	SubmitRatingResp struct {
		AvgRating   float64 `json:"avg_rating"`
		RatingCount int     `json:"rating_count"`
	}
)

type (
	SubmitReviewReq struct {
		Nickname string `json:"nickname"`
		Content  string `json:"content" binding:"required,min=10,max=2000"`
	}
	SubmitReviewResp struct {
		ID      int    `json:"id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
)

type (
	ReviewListReq struct {
		Page     int `form:"page,default=1"`
		PageSize int `form:"page_size,default=10"`
	}

	ReviewItem struct {
		ID        int    `json:"id"`
		Nickname  string `json:"nickname"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}

	ReviewListResp struct {
		List       []ReviewItem      `json:"list"`
		Pagination SiteLisPagination `json:"pagination"`
	}
)
