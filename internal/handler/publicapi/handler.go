package publicapi

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	v1 "github.com/ch3nnn/webstack-go/api/v1"
	"github.com/ch3nnn/webstack-go/internal/dal/model"
	"github.com/ch3nnn/webstack-go/pkg/log"
	"github.com/ch3nnn/webstack-go/pkg/tools"
)

type Handler struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewHandler(db *gorm.DB, logger *log.Logger) *Handler {
	return &Handler{db: db, logger: logger}
}

func (h *Handler) SiteList(ctx *gin.Context) {
	var req v1.PublicSiteListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if req.PageSize > 100 {
		req.PageSize = 100
	}

	query := h.db.Model(&model.StSite{}).Where("is_used = ?", true).Where("deleted_at IS NULL")

	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.Search != "" {
		like := "%" + req.Search + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", like, like)
	}
	if req.TagID > 0 {
		query = query.Where("id IN (SELECT site_id FROM st_site_tag WHERE tag_id = ?)", req.TagID)
	}
	if req.ModelID > 0 {
		query = query.Where("id IN (SELECT site_id FROM st_site WHERE supported_models LIKE ?)", "%"+strconv.Itoa(req.ModelID)+"%")
	}

	var total int64
	query.Count(&total)

	orderClause := "sort ASC"
	switch req.SortBy {
	case "rating":
		orderClause = "avg_rating " + safeSortOrder(req.SortOrder)
	case "reviews":
		orderClause = "review_count " + safeSortOrder(req.SortOrder)
	case "newest":
		orderClause = "created_at " + safeSortOrder(req.SortOrder)
	}

	var sites []model.StSite
	query.Order(orderClause).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&sites)

	categoryMap := h.getCategoryMap()
	list := make([]v1.PublicSiteItem, 0, len(sites))
	for _, s := range sites {
		tags := h.getSiteTags(s.ID)
		list = append(list, v1.PublicSiteItem{
			ID:            s.ID,
			Title:         s.Title,
			Icon:          s.Icon,
			URL:           s.URL,
			Description:   s.Description,
			CategoryID:    s.CategoryID,
			CategoryName:  categoryMap[s.CategoryID],
			AvgRating:     math.Round(s.AvgRating*10) / 10,
			RatingCount:   s.RatingCount,
			ReviewCount:   s.ReviewCount,
			IsFeatured:    s.IsFeatured,
			IsRecommended: s.IsRecommended,
			Tags:          tags,
		})
	}

	v1.HandleSuccess(ctx, v1.PublicSiteListResp{
		List: list,
		Pagination: v1.SiteLisPagination{
			Total:        total,
			CurrentPage:  req.Page,
			PerPageCount: req.PageSize,
		},
	})
}

func (h *Handler) SiteDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	var site model.StSite
	if err := h.db.First(&site, id).Error; err != nil {
		v1.HandleError(ctx, http.StatusNotFound, v1.ErrNotFound, nil)
		return
	}

	h.db.Model(&site).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	categoryMap := h.getCategoryMap()

	var screenshots []string
	if site.Screenshots != "" {
		json.Unmarshal([]byte(site.Screenshots), &screenshots)
	}

	var models []v1.SiteModelInfo
	if site.SupportedModels != "" {
		var modelIDs []int
		json.Unmarshal([]byte(site.SupportedModels), &modelIDs)
		if len(modelIDs) > 0 {
			var dbModels []model.StModel
			h.db.Where("id IN ?", modelIDs).Find(&dbModels)
			for _, m := range dbModels {
				models = append(models, v1.SiteModelInfo{
					ID: m.ID, Name: m.Name, Provider: m.Provider,
				})
			}
		}
	}

	var pricing *v1.SitePricingInfo
	if site.PricingInfo != "" {
		pricing = &v1.SitePricingInfo{}
		json.Unmarshal([]byte(site.PricingInfo), pricing)
	}

	var health *v1.SiteHealthInfo
	if site.HealthStatus != "" {
		health = &v1.SiteHealthInfo{}
		json.Unmarshal([]byte(site.HealthStatus), health)
	}

	var tags []v1.SiteTagInfo
	var siteTags []model.StSiteTag
	h.db.Where("site_id = ?", site.ID).Find(&siteTags)
	if len(siteTags) > 0 {
		tagIDs := make([]int, len(siteTags))
		for i, st := range siteTags {
			tagIDs[i] = st.TagID
		}
		var dbTags []model.StTag
		h.db.Where("id IN ?", tagIDs).Find(&dbTags)
		for _, t := range dbTags {
			tags = append(tags, v1.SiteTagInfo{ID: t.ID, Name: t.Name, Color: t.Color})
		}
	}

	var payments []v1.SitePaymentInfo
	var sitePayments []model.StSitePayment
	h.db.Where("site_id = ?", site.ID).Find(&sitePayments)
	if len(sitePayments) > 0 {
		pmIDs := make([]int, len(sitePayments))
		for i, sp := range sitePayments {
			pmIDs[i] = sp.PaymentMethodID
		}
		var dbPMs []model.StPaymentMethod
		h.db.Where("id IN ?", pmIDs).Find(&dbPMs)
		for _, p := range dbPMs {
			payments = append(payments, v1.SitePaymentInfo{ID: p.ID, Name: p.Name, Icon: p.Icon})
		}
	}

	createdAt := ""
	if site.CreatedAt != nil {
		createdAt = site.CreatedAt.Format(time.DateTime)
	}

	v1.HandleSuccess(ctx, v1.SiteDetailResp{
		ID:              site.ID,
		Title:           site.Title,
		Icon:            site.Icon,
		URL:             site.URL,
		Description:     site.Description,
		Content:         site.Content,
		Screenshots:     screenshots,
		CategoryID:      site.CategoryID,
		CategoryName:    categoryMap[site.CategoryID],
		SupportedModels: models,
		PricingInfo:     pricing,
		HealthStatus:    health,
		Tags:            tags,
		PaymentMethods:  payments,
		AvgRating:       math.Round(site.AvgRating*10) / 10,
		RatingCount:     site.RatingCount,
		ReviewCount:     site.ReviewCount,
		ViewCount:       site.ViewCount + 1,
		IsFeatured:      site.IsFeatured,
		IsRecommended:   site.IsRecommended,
		CreatedAt:       createdAt,
	})
}

func (h *Handler) SubmitRating(ctx *gin.Context) {
	siteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	var config model.SysConfig
	h.db.First(&config)

	if !config.AllowAnonymousRating {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	var req v1.SubmitRatingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// Use visitor IP as identifier for anonymous rating dedup
	visitorID := ctx.ClientIP()

	now := time.Now()
	var existing model.StRating
	result := h.db.Where("site_id = ? AND visitor_id = ?", siteID, visitorID).First(&existing)
	if result.Error == nil {
		h.db.Model(&model.StRating{}).Where("id = ?", existing.ID).
			Updates(map[string]interface{}{"score": req.Score, "updated_at": &now})
	} else {
		h.db.Create(&model.StRating{SiteID: siteID, UserID: 0, VisitorID: visitorID, Score: req.Score, CreatedAt: &now, UpdatedAt: &now})
	}

	var avgRating float64
	var ratingCount int64
	h.db.Model(&model.StRating{}).Where("site_id = ?", siteID).Count(&ratingCount)
	h.db.Model(&model.StRating{}).Where("site_id = ?", siteID).Select("COALESCE(AVG(score), 0)").Scan(&avgRating)

	h.db.Model(&model.StSite{}).Where("id = ?", siteID).Updates(map[string]interface{}{
		"avg_rating":   avgRating,
		"rating_count": ratingCount,
	})

	v1.HandleSuccess(ctx, v1.SubmitRatingResp{
		AvgRating:   math.Round(avgRating*10) / 10,
		RatingCount: int(ratingCount),
	})
}

func (h *Handler) SubmitReview(ctx *gin.Context) {
	siteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	var config model.SysConfig
	h.db.First(&config)

	if !config.AllowAnonymousReview {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	var req v1.SubmitReviewReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	now := time.Now()
	review := model.StReview{
		SiteID:    siteID,
		UserID:    0,
		Nickname:  req.Nickname,
		Content:   req.Content,
		Status:    "pending",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := h.db.Create(&review).Error; err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServer, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.SubmitReviewResp{
		ID:      review.ID,
		Status:  "pending",
		Message: "Review submitted. It will be visible after approval.",
	})
}

func (h *Handler) ReviewList(ctx *gin.Context) {
	siteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	var req v1.ReviewListReq
	ctx.ShouldBindQuery(&req)
	if req.PageSize > 50 {
		req.PageSize = 50
	}

	query := h.db.Model(&model.StReview{}).Where("site_id = ? AND status = ?", siteID, "approved")

	var total int64
	query.Count(&total)

	var reviews []model.StReview
	query.Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&reviews)

	list := make([]v1.ReviewItem, 0, len(reviews))
	for _, r := range reviews {
		createdAt := ""
		if r.CreatedAt != nil {
			createdAt = r.CreatedAt.Format(time.DateTime)
		}
		nickname := r.Nickname
		if nickname == "" {
			nickname = "Anonymous"
		}
		list = append(list, v1.ReviewItem{
			ID:        r.ID,
			Nickname:  nickname,
			Content:   r.Content,
			CreatedAt: createdAt,
		})
	}

	v1.HandleSuccess(ctx, v1.ReviewListResp{
		List: list,
		Pagination: v1.SiteLisPagination{
			Total:        total,
			CurrentPage:  req.Page,
			PerPageCount: req.PageSize,
		},
	})
}

func (h *Handler) ModelList(ctx *gin.Context) {
	var models []model.StModel
	h.db.Where("is_used = ? AND deleted_at IS NULL", true).Order("sort ASC").Find(&models)

	list := make([]v1.ModelItem, 0, len(models))
	for _, m := range models {
		list = append(list, v1.ModelItem{
			ID: m.ID, Name: m.Name, Provider: m.Provider, Icon: m.Icon,
		})
	}
	v1.HandleSuccess(ctx, v1.ModelListResp{List: list})
}

func (h *Handler) TagList(ctx *gin.Context) {
	var tags []model.StTag
	h.db.Where("is_used = ? AND deleted_at IS NULL", true).Order("sort ASC").Find(&tags)

	list := make([]v1.SiteTagInfo, 0, len(tags))
	for _, t := range tags {
		list = append(list, v1.SiteTagInfo{ID: t.ID, Name: t.Name, Color: t.Color})
	}
	v1.HandleSuccess(ctx, map[string]interface{}{"list": list})
}

func (h *Handler) SubmitSite(ctx *gin.Context) {
	var config model.SysConfig
	h.db.First(&config)
	if !config.AllowSubmission {
		v1.HandleError(ctx, http.StatusForbidden, v1.ErrBadRequest, nil)
		return
	}

	var req struct {
		URL         string `json:"url" binding:"required"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CategoryID  int    `json:"category_id"`
		Contact     string `json:"contact"`
		Icon        string `json:"icon"`
		ModelIDs    string `json:"model_ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	now := time.Now()
	sub := model.StSubmission{
		UserID:      0,
		URL:         req.URL,
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Contact:     req.Contact,
		ModelIDs:    req.ModelIDs,
		Status:      "pending",
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	if req.Icon != "" {
		if len(req.Icon) > 20 && req.Icon[:20] == "data:image/png;base64," {
			sub.Icon = req.Icon[22:]
		} else {
			icon, err := tools.ResizeURLImgToBase64(req.Icon, 64, 64)
			if err == nil && icon != "" {
				sub.Icon = icon
			}
		}
	}

	if err := h.db.Create(&sub).Error; err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServer, nil)
		return
	}

	v1.HandleSuccess(ctx, map[string]interface{}{"id": sub.ID, "status": "pending"})
}

func (h *Handler) GetSiteConfig(ctx *gin.Context) {
	var config model.SysConfig
	h.db.First(&config)
	v1.HandleSuccess(ctx, map[string]interface{}{
		"site_lang":              config.SiteLang,
		"allow_anonymous_rating": config.AllowAnonymousRating,
		"allow_anonymous_review": config.AllowAnonymousReview,
		"allow_submission":       config.AllowSubmission,
		"site_copyright":         config.SiteCopyright,
	})
}

// helpers

func (h *Handler) getCategoryMap() map[int]string {
	var cats []model.StCategory
	h.db.Where("deleted_at IS NULL").Find(&cats)
	m := make(map[int]string)
	for _, c := range cats {
		m[c.ID] = c.Title
	}
	return m
}

func (h *Handler) getSiteTags(siteID int) []string {
	var siteTags []model.StSiteTag
	h.db.Where("site_id = ?", siteID).Find(&siteTags)
	if len(siteTags) == 0 {
		return []string{}
	}
	tagIDs := make([]int, len(siteTags))
	for i, st := range siteTags {
		tagIDs[i] = st.TagID
	}
	var tags []model.StTag
	h.db.Where("id IN ?", tagIDs).Find(&tags)
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	return names
}

func (h *Handler) AdList(ctx *gin.Context) {
	slot := ctx.Query("slot")
	if slot == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	now := time.Now()
	var ads []model.StAd
	h.db.Model(&model.StAd{}).
		Where("slot = ? AND is_used = ? AND deleted_at IS NULL", slot, true).
		Where("(start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)", now, now).
		Order("sort ASC, id DESC").
		Find(&ads)

	type adResp struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ContentType string `json:"content_type"`
		Content     string `json:"content"`
		ImageURL    string `json:"image_url"`
		LinkURL     string `json:"link_url"`
	}

	list := make([]adResp, 0, len(ads))
	for _, a := range ads {
		list = append(list, adResp{
			ID:          a.ID,
			Name:        a.Name,
			ContentType: a.ContentType,
			Content:     a.Content,
			ImageURL:    a.ImageURL,
			LinkURL:     a.LinkURL,
		})
	}

	v1.HandleSuccess(ctx, map[string]interface{}{"list": list})
}

func safeSortOrder(order string) string {
	if order == "desc" || order == "DESC" {
		return "DESC"
	}
	return "ASC"
}
