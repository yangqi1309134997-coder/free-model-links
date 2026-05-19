package adminext

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
	"github.com/ch3nnn/webstack-go/internal/healthcheck"
	"github.com/ch3nnn/webstack-go/pkg/log"
)

type Handler struct {
	db       *gorm.DB
	logger   *log.Logger
	checker  *healthcheck.Checker
}

func NewHandler(db *gorm.DB, logger *log.Logger) *Handler {
	return &Handler{db: db, logger: logger}
}

func (h *Handler) SetChecker(c *healthcheck.Checker) {
	h.checker = c
}

// Review management

func (h *Handler) ReviewList(ctx *gin.Context) {
	var req v1.AdminReviewListReq
	ctx.ShouldBindQuery(&req)

	query := h.db.Model(&model.StReview{}).Where("deleted_at IS NULL")
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.SiteID > 0 {
		query = query.Where("site_id = ?", req.SiteID)
	}

	var total int64
	query.Count(&total)

	var reviews []model.StReview
	query.Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&reviews)

	list := make([]v1.AdminReviewItem, 0, len(reviews))
	for _, r := range reviews {
		var site model.StSite
		h.db.Select("title").First(&site, r.SiteID)
		nickname := r.Nickname
		if nickname == "" {
			nickname = "Anonymous"
		}
		createdAt := ""
		if r.CreatedAt != nil {
			createdAt = r.CreatedAt.Format(time.DateTime)
		}
		list = append(list, v1.AdminReviewItem{
			ID: r.ID, SiteID: r.SiteID, SiteTitle: site.Title,
			Nickname: nickname, Content: r.Content, Status: r.Status,
			AdminNote: r.AdminNote, CreatedAt: createdAt,
		})
	}

	v1.HandleSuccess(ctx, v1.AdminReviewListResp{
		List: list,
		Pagination: v1.SiteLisPagination{Total: total, CurrentPage: req.Page, PerPageCount: req.PageSize},
	})
}

func (h *Handler) ReviewUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req v1.AdminReviewUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	h.db.Model(&model.StReview{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     req.Status,
		"admin_note": req.AdminNote,
	})

	if req.Status == "approved" {
		var review model.StReview
		h.db.First(&review, id)
		var count int64
		h.db.Model(&model.StReview{}).Where("site_id = ? AND status = ?", review.SiteID, "approved").Count(&count)
		h.db.Model(&model.StSite{}).Where("id = ?", review.SiteID).Update("review_count", count)
	} else if req.Status == "rejected" {
		var review model.StReview
		h.db.First(&review, id)
		var count int64
		h.db.Model(&model.StReview{}).Where("site_id = ? AND status = ?", review.SiteID, "approved").Count(&count)
		h.db.Model(&model.StSite{}).Where("id = ?", review.SiteID).Update("review_count", count)
	}

	v1.HandleSuccess(ctx, v1.AdminReviewUpdateResp{})
}

// Rating management

func (h *Handler) RatingList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	siteID, _ := strconv.Atoi(ctx.DefaultQuery("site_id", "0"))

	query := h.db.Model(&model.StRating{})
	if siteID > 0 {
		query = query.Where("site_id = ?", siteID)
	}

	var total int64
	query.Count(&total)

	var ratings []model.StRating
	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&ratings)

	siteMap := h.getSiteMap()

	list := make([]v1.AdminRatingItem, 0, len(ratings))
	for _, r := range ratings {
		createdAt := ""
		if r.CreatedAt != nil {
			createdAt = r.CreatedAt.Format(time.DateTime)
		}
		list = append(list, v1.AdminRatingItem{
			ID: r.ID, SiteID: r.SiteID, SiteTitle: siteMap[r.SiteID],
			VisitorID: r.VisitorID, Score: r.Score, CreatedAt: createdAt,
		})
	}

	v1.HandleSuccess(ctx, v1.AdminRatingListResp{
		List: list,
		Pagination: v1.SiteLisPagination{Total: total, CurrentPage: page, PerPageCount: pageSize},
	})
}

func (h *Handler) RatingUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req v1.AdminRatingUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	h.db.Model(&model.StRating{}).Where("id = ?", id).Update("score", req.Score)

	// Recalculate site avg
	var rating model.StRating
	h.db.First(&rating, id)
	h.recalcSiteRating(rating.SiteID)

	v1.HandleSuccess(ctx, v1.AdminRatingUpdateResp{})
}

func (h *Handler) RatingDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	var rating model.StRating
	h.db.First(&rating, id)
	siteID := rating.SiteID

	h.db.Delete(&model.StRating{}, id)
	h.recalcSiteRating(siteID)

	v1.HandleSuccess(ctx, v1.AdminRatingDeleteResp{})
}

func (h *Handler) recalcSiteRating(siteID int) {
	var avgRating float64
	var ratingCount int64
	h.db.Model(&model.StRating{}).Where("site_id = ?", siteID).Count(&ratingCount)
	h.db.Model(&model.StRating{}).Where("site_id = ?", siteID).Select("COALESCE(AVG(score), 0)").Scan(&avgRating)
	h.db.Model(&model.StSite{}).Where("id = ?", siteID).Updates(map[string]interface{}{
		"avg_rating":   math.Round(avgRating*10) / 10,
		"rating_count": ratingCount,
	})
}

// User management (admin-only users)

func (h *Handler) UserList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	query := h.db.Model(&model.SysUser{}).Where("deleted_at IS NULL")

	var total int64
	query.Count(&total)

	var users []model.SysUser
	query.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users)

	type userItem struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}

	list := make([]userItem, 0, len(users))
	for _, u := range users {
		list = append(list, userItem{ID: u.ID, Username: u.Username, Role: u.Role, Status: u.Status})
	}

	v1.HandleSuccess(ctx, map[string]interface{}{
		"list":       list,
		"pagination": v1.SiteLisPagination{Total: total, CurrentPage: page, PerPageCount: pageSize},
	})
}

func (h *Handler) UserUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		Status string `json:"status" binding:"required,oneof=active banned"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	h.db.Model(&model.SysUser{}).Where("id = ?", id).Update("status", req.Status)
	v1.HandleSuccess(ctx, map[string]interface{}{})
}

// Dictionary CRUD (models, tags, payment methods)

func (h *Handler) ModelList(ctx *gin.Context) {
	var models []model.StModel
	h.db.Where("deleted_at IS NULL").Order("sort ASC").Find(&models)
	list := make([]v1.ModelItem, 0, len(models))
	for _, m := range models {
		list = append(list, v1.ModelItem{
			ID: m.ID, Name: m.Name, Provider: m.Provider, Icon: m.Icon, Sort: m.Sort, IsUsed: m.IsUsed,
		})
	}
	v1.HandleSuccess(ctx, v1.ModelListResp{List: list})
}

func (h *Handler) ModelCreate(ctx *gin.Context) {
	var req v1.DictCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	now := time.Now()
	m := model.StModel{Name: req.Name, Provider: req.Provider, Icon: req.Icon, Sort: req.Sort, IsUsed: req.IsUsed, CreatedAt: &now, UpdatedAt: &now}
	h.db.Create(&m)
	v1.HandleSuccess(ctx, v1.DictCreateResp{ID: m.ID})
}

func (h *Handler) ModelUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req v1.DictUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	updates := map[string]interface{}{"name": req.Name, "provider": req.Provider, "icon": req.Icon, "sort": req.Sort}
	if req.IsUsed != nil {
		updates["is_used"] = *req.IsUsed
	}
	h.db.Model(&model.StModel{}).Where("id = ?", id).Updates(updates)
	v1.HandleSuccess(ctx, v1.DictUpdateResp{})
}

func (h *Handler) ModelDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	now := time.Now()
	h.db.Model(&model.StModel{}).Where("id = ?", id).Update("deleted_at", &now)
	v1.HandleSuccess(ctx, v1.DictDeleteResp{})
}

func (h *Handler) TagList(ctx *gin.Context) {
	var tags []model.StTag
	h.db.Where("deleted_at IS NULL").Order("sort ASC").Find(&tags)
	list := make([]v1.SiteTagInfo, 0, len(tags))
	for _, t := range tags {
		list = append(list, v1.SiteTagInfo{ID: t.ID, Name: t.Name, Color: t.Color})
	}
	v1.HandleSuccess(ctx, map[string]interface{}{"list": list})
}

func (h *Handler) TagCreate(ctx *gin.Context) {
	var req v1.DictCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	now := time.Now()
	t := model.StTag{Name: req.Name, Color: req.Color, Sort: req.Sort, IsUsed: req.IsUsed, CreatedAt: &now, UpdatedAt: &now}
	h.db.Create(&t)
	v1.HandleSuccess(ctx, v1.DictCreateResp{ID: t.ID})
}

func (h *Handler) TagUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req v1.DictUpdateReq
	ctx.ShouldBindJSON(&req)
	updates := map[string]interface{}{"name": req.Name, "color": req.Color, "sort": req.Sort}
	if req.IsUsed != nil {
		updates["is_used"] = *req.IsUsed
	}
	h.db.Model(&model.StTag{}).Where("id = ?", id).Updates(updates)
	v1.HandleSuccess(ctx, v1.DictUpdateResp{})
}

func (h *Handler) TagDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	now := time.Now()
	h.db.Model(&model.StTag{}).Where("id = ?", id).Update("deleted_at", &now)
	h.db.Where("tag_id = ?", id).Delete(&model.StSiteTag{})
	v1.HandleSuccess(ctx, v1.DictDeleteResp{})
}

func (h *Handler) PaymentMethodList(ctx *gin.Context) {
	var pms []model.StPaymentMethod
	h.db.Where("deleted_at IS NULL").Order("sort ASC").Find(&pms)
	list := make([]v1.DictItem, 0, len(pms))
	for _, p := range pms {
		list = append(list, v1.DictItem{ID: p.ID, Name: p.Name, Icon: p.Icon, Sort: p.Sort, IsUsed: p.IsUsed})
	}
	v1.HandleSuccess(ctx, map[string]interface{}{"list": list})
}

func (h *Handler) PaymentMethodCreate(ctx *gin.Context) {
	var req v1.DictCreateReq
	ctx.ShouldBindJSON(&req)
	now := time.Now()
	p := model.StPaymentMethod{Name: req.Name, Icon: req.Icon, Sort: req.Sort, IsUsed: req.IsUsed, CreatedAt: &now, UpdatedAt: &now}
	h.db.Create(&p)
	v1.HandleSuccess(ctx, v1.DictCreateResp{ID: p.ID})
}

func (h *Handler) PaymentMethodUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req v1.DictUpdateReq
	ctx.ShouldBindJSON(&req)
	updates := map[string]interface{}{"name": req.Name, "icon": req.Icon, "sort": req.Sort}
	if req.IsUsed != nil {
		updates["is_used"] = *req.IsUsed
	}
	h.db.Model(&model.StPaymentMethod{}).Where("id = ?", id).Updates(updates)
	v1.HandleSuccess(ctx, v1.DictUpdateResp{})
}

func (h *Handler) PaymentMethodDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	now := time.Now()
	h.db.Model(&model.StPaymentMethod{}).Where("id = ?", id).Update("deleted_at", &now)
	h.db.Where("payment_method_id = ?", id).Delete(&model.StSitePayment{})
	v1.HandleSuccess(ctx, v1.DictDeleteResp{})
}

// Submission management

func (h *Handler) SubmissionList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	status := ctx.Query("status")

	query := h.db.Model(&model.StSubmission{}).Where("deleted_at IS NULL")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var subs []model.StSubmission
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&subs)

	type subItem struct {
		ID          int    `json:"id"`
		URL         string `json:"url"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Contact     string `json:"contact"`
		Status      string `json:"status"`
		AdminNote   string `json:"admin_note"`
		CreatedAt   string `json:"created_at"`
	}

	list := make([]subItem, 0, len(subs))
	for _, s := range subs {
		ca := ""
		if s.CreatedAt != nil {
			ca = s.CreatedAt.Format(time.DateTime)
		}
		list = append(list, subItem{
			ID: s.ID, URL: s.URL, Title: s.Title, Description: s.Description,
			Contact: s.Contact, Status: s.Status, AdminNote: s.AdminNote, CreatedAt: ca,
		})
	}

	v1.HandleSuccess(ctx, map[string]interface{}{
		"list":       list,
		"pagination": v1.SiteLisPagination{Total: total, CurrentPage: page, PerPageCount: pageSize},
	})
}

func (h *Handler) SubmissionUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		Status    string `json:"status" binding:"required"`
		AdminNote string `json:"admin_note"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	var sub model.StSubmission
	if err := h.db.First(&sub, id).Error; err != nil {
		v1.HandleError(ctx, http.StatusNotFound, v1.ErrNotFound, nil)
		return
	}

	h.db.Model(&model.StSubmission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     req.Status,
		"admin_note": req.AdminNote,
	})

	// 审核通过后自动创建网站记录
	if req.Status == "approved" {
		now := time.Now()
		title := sub.Title
		if title == "" {
			title = sub.URL
		}
		site := model.StSite{
			Title:       title,
			URL:         sub.URL,
			Description: sub.Description,
			CategoryID:  sub.CategoryID,
			IsUsed:      true,
			Sort:        0,
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}
		h.db.Create(&site)
	}
	v1.HandleSuccess(ctx, map[string]interface{}{})
}

// Site settings (lang, anonymous, submission toggle)

func (h *Handler) UpdateSiteSettings(ctx *gin.Context) {
	var req struct {
		SiteLang             string `json:"site_lang"`
		DefaultSort          string `json:"default_sort"`
		AllowAnonymousRating *bool  `json:"allow_anonymous_rating"`
		AllowAnonymousReview *bool  `json:"allow_anonymous_review"`
		AllowSubmission      *bool  `json:"allow_submission"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	updates := map[string]interface{}{}
	if req.SiteLang != "" {
		updates["site_lang"] = req.SiteLang
	}
	if req.DefaultSort != "" {
		updates["default_sort"] = req.DefaultSort
	}
	if req.AllowAnonymousRating != nil {
		updates["allow_anonymous_rating"] = *req.AllowAnonymousRating
	}
	if req.AllowAnonymousReview != nil {
		updates["allow_anonymous_review"] = *req.AllowAnonymousReview
	}
	if req.AllowSubmission != nil {
		updates["allow_submission"] = *req.AllowSubmission
	}

	if len(updates) > 0 {
		h.db.Model(&model.SysConfig{}).Where("id = 1").Updates(updates)
	}
	v1.HandleSuccess(ctx, map[string]interface{}{})
}

// Health check config

func (h *Handler) GetHealthConfig(ctx *gin.Context) {
	var config model.SysConfig
	h.db.First(&config)

	interval := 30
	maxFails := 10
	enabled := true

	if config.HealthCheckInterval > 0 {
		interval = config.HealthCheckInterval
	}
	if config.HealthCheckMaxFails > 0 {
		maxFails = config.HealthCheckMaxFails
	}

	v1.HandleSuccess(ctx, v1.HealthConfigResp{
		IntervalMinutes:     interval,
		MaxConsecutiveFails: maxFails,
		Enabled:             enabled,
	})
}

func (h *Handler) UpdateHealthConfig(ctx *gin.Context) {
	var req v1.HealthConfigUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	updates := map[string]interface{}{}
	if req.IntervalMinutes != nil {
		updates["health_check_interval"] = *req.IntervalMinutes
	}
	if req.MaxConsecutiveFails != nil {
		updates["health_check_max_fails"] = *req.MaxConsecutiveFails
	}
	if req.Enabled != nil {
		updates["health_check_enabled"] = *req.Enabled
	}

	if len(updates) > 0 {
		h.db.Model(&model.SysConfig{}).Where("id = 1").Updates(updates)
	}
	v1.HandleSuccess(ctx, map[string]interface{}{})
}

func (h *Handler) TriggerHealthCheck(ctx *gin.Context) {
	if h.checker != nil {
		go h.checker.RunOnce()
	}
	v1.HandleSuccess(ctx, map[string]interface{}{"message": "Health check triggered"})
}

// Homepage layout config

func (h *Handler) GetHomepageConfig(ctx *gin.Context) {
	var config model.SysConfig
	h.db.First(&config)

	var sections []v1.HomepageSection
	if config.HomepageSections != "" {
		json.Unmarshal([]byte(config.HomepageSections), &sections)
	}
	if len(sections) == 0 {
		sections = defaultHomepageSections()
	}

	v1.HandleSuccess(ctx, v1.HomepageConfigResp{
		Sections:        sections,
		DefaultSort:     config.DefaultSort,
		CardsPerRow:     config.CardsPerRow,
		ShowFeatured:    config.ShowFeatured,
		ShowRecommended: config.ShowRecommended,
		FeaturedSiteIDs: config.FeaturedSiteIDs,
	})
}

func (h *Handler) UpdateHomepageConfig(ctx *gin.Context) {
	var req v1.HomepageConfigUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	updates := map[string]interface{}{}
	if req.Sections != nil {
		data, _ := json.Marshal(req.Sections)
		updates["homepage_sections"] = string(data)
	}
	if req.DefaultSort != nil {
		updates["default_sort"] = *req.DefaultSort
	}
	if req.CardsPerRow != nil {
		updates["cards_per_row"] = *req.CardsPerRow
	}
	if req.ShowFeatured != nil {
		updates["show_featured"] = *req.ShowFeatured
	}
	if req.ShowRecommended != nil {
		updates["show_recommended"] = *req.ShowRecommended
	}
	if req.FeaturedSiteIDs != nil {
		data, _ := json.Marshal(req.FeaturedSiteIDs)
		updates["featured_site_ids"] = string(data)
	}

	if len(updates) > 0 {
		h.db.Model(&model.SysConfig{}).Where("id = 1").Updates(updates)
	}
	v1.HandleSuccess(ctx, map[string]interface{}{})
}

func defaultHomepageSections() []v1.HomepageSection {
	return []v1.HomepageSection{
		{Key: "featured", Title: "精选推荐", Enabled: true, Collapsed: false},
		{Key: "recommended", Title: "站长推荐", Enabled: true, Collapsed: false},
		{Key: "categories", Title: "全部分类", Enabled: true, Collapsed: false},
	}
}

// helpers

func (h *Handler) getSiteMap() map[int]string {
	var sites []model.StSite
	h.db.Select("id, title").Find(&sites)
	m := make(map[int]string)
	for _, s := range sites {
		m[s.ID] = s.Title
	}
	return m
}

// SEO and Friend links config

func (h *Handler) GetSEOConfig(ctx *gin.Context) {
	var config model.SysConfig
	h.db.First(&config)
	v1.HandleSuccess(ctx, map[string]interface{}{
		"seo_author":        config.SeoAuthor,
		"seo_og_image":      config.SeoOgImage,
		"seo_twitter_card":  config.SeoTwitterCard,
		"seo_canonical_url": config.SeoCanonicalURL,
		"seo_robots":        config.SeoRobots,
		"friend_links":      config.FriendLinks,
	})
}

func (h *Handler) UpdateSEOConfig(ctx *gin.Context) {
	var req struct {
		SeoAuthor       string `json:"seo_author"`
		SeoOgImage      string `json:"seo_og_image"`
		SeoTwitterCard  string `json:"seo_twitter_card"`
		SeoCanonicalURL string `json:"seo_canonical_url"`
		SeoRobots       string `json:"seo_robots"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	h.db.Model(&model.SysConfig{}).Where("id = ?", 1).Updates(map[string]interface{}{
		"seo_author":        req.SeoAuthor,
		"seo_og_image":      req.SeoOgImage,
		"seo_twitter_card":  req.SeoTwitterCard,
		"seo_canonical_url": req.SeoCanonicalURL,
		"seo_robots":        req.SeoRobots,
	})
	v1.HandleSuccess(ctx, map[string]interface{}{})
}

func (h *Handler) UpdateFriendLinks(ctx *gin.Context) {
	var req struct {
		FriendLinks string `json:"friend_links"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	h.db.Model(&model.SysConfig{}).Where("id = ?", 1).Update("friend_links", req.FriendLinks)
	v1.HandleSuccess(ctx, map[string]interface{}{})
}

// Ad management

func (h *Handler) AdList(ctx *gin.Context) {
	var req v1.AdListReq
	ctx.ShouldBindQuery(&req)
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	query := h.db.Model(&model.StAd{}).Where("deleted_at IS NULL")
	if req.Slot != "" {
		query = query.Where("slot = ?", req.Slot)
	}

	var total int64
	query.Count(&total)

	var ads []model.StAd
	query.Order("sort ASC, id DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&ads)

	list := make([]v1.AdItem, 0, len(ads))
	for _, a := range ads {
		list = append(list, adToItem(a))
	}

	v1.HandleSuccess(ctx, v1.AdListResp{
		List: list,
		Pagination: v1.SiteLisPagination{Total: total, CurrentPage: req.Page, PerPageCount: req.PageSize},
	})
}

func (h *Handler) AdCreate(ctx *gin.Context) {
	var req v1.AdCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	if req.ContentType == "" {
		req.ContentType = "html"
	}

	now := time.Now()
	ad := model.StAd{
		Name:        req.Name,
		Slot:        req.Slot,
		ContentType: req.ContentType,
		Content:     req.Content,
		ImageURL:    req.ImageURL,
		LinkURL:     req.LinkURL,
		Sort:        req.Sort,
		IsUsed:      req.IsUsed,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	if req.StartAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.StartAt); err == nil {
			ad.StartAt = &t
		}
	}
	if req.EndAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.EndAt); err == nil {
			ad.EndAt = &t
		}
	}

	h.db.Create(&ad)
	v1.HandleSuccess(ctx, v1.AdCreateResp{ID: ad.ID})
}

func (h *Handler) AdUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req v1.AdUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	updates := map[string]interface{}{
		"name":         req.Name,
		"slot":         req.Slot,
		"content_type": req.ContentType,
		"content":      req.Content,
		"image_url":    req.ImageURL,
		"link_url":     req.LinkURL,
		"sort":         req.Sort,
		"updated_at":   time.Now(),
	}
	if req.IsUsed != nil {
		updates["is_used"] = *req.IsUsed
	}
	if req.StartAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.StartAt); err == nil {
			updates["start_at"] = &t
		} else {
			updates["start_at"] = nil
		}
	} else {
		updates["start_at"] = nil
	}
	if req.EndAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.EndAt); err == nil {
			updates["end_at"] = &t
		} else {
			updates["end_at"] = nil
		}
	} else {
		updates["end_at"] = nil
	}

	h.db.Model(&model.StAd{}).Where("id = ?", id).Updates(updates)
	v1.HandleSuccess(ctx, v1.AdUpdateResp{})
}

func (h *Handler) AdDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	now := time.Now()
	h.db.Model(&model.StAd{}).Where("id = ?", id).Update("deleted_at", &now)
	v1.HandleSuccess(ctx, v1.AdDeleteResp{})
}

func adToItem(a model.StAd) v1.AdItem {
	item := v1.AdItem{
		ID:          a.ID,
		Name:        a.Name,
		Slot:        a.Slot,
		ContentType: a.ContentType,
		Content:     a.Content,
		ImageURL:    a.ImageURL,
		LinkURL:     a.LinkURL,
		Sort:        a.Sort,
		IsUsed:      a.IsUsed,
	}
	if a.StartAt != nil {
		item.StartAt = a.StartAt.Format(time.DateTime)
	}
	if a.EndAt != nil {
		item.EndAt = a.EndAt.Format(time.DateTime)
	}
	if a.CreatedAt != nil {
		item.CreatedAt = a.CreatedAt.Format(time.DateTime)
	}
	return item
}
