/**
 * @Author: chentong
 * @Date: 2024/05/26 上午1:50
 */

package v1

import "github.com/ch3nnn/webstack-go/internal/dal/model"

type TreeNode struct {
	Id    int         // 节点ID
	Pid   int         // 父节点ID
	Name  string      // 节点名称
	Icon  string      // 图标
	Sort  int         // 排序
	Child []*TreeNode // 获取子节点切片
}

type CategorySite struct {
	Category string         // 分类
	SiteList []model.StSite // 站点列表
}

type About struct {
	AboutSite   string `json:"about_site"`   // 关于站点
	AboutAuthor string `json:"about_author"` // 关于作者
	IsAbout     bool   `json:"is_about"`     // 是否开启关于
}

type FriendLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Desc string `json:"desc,omitempty"`
}

type ConfigSite struct {
	SiteTitle   string `json:"site_title"`
	SiteKeyword string `json:"site_keyword"`
	SiteDesc    string `json:"site_desc"`
	SiteRecord  string `json:"site_record"`
	SiteURL     string `json:"site_url"`
	SiteLogo    string `json:"site_logo"`
	SiteFavicon string `json:"site_favicon"`
	SiteLang    string `json:"site_lang"`
	AllowAnonymousRating  bool `json:"allow_anonymous_rating"`
	AllowAnonymousReview  bool `json:"allow_anonymous_review"`
	AllowSubmission       bool `json:"allow_submission"`
	// 首页布局
	CardsPerRow     int  `json:"cards_per_row"`
	ShowFeatured    bool `json:"show_featured"`
	ShowRecommended bool `json:"show_recommended"`
	ShowSearch      bool `json:"show_search"`
	ShowHero        bool `json:"show_hero"`
	ShowCategoryNav bool `json:"show_category_nav"`
	ShowFooter      bool `json:"show_footer"`
	// SEO 扩展
	SeoAuthor       string       `json:"seo_author"`
	SeoOgImage      string       `json:"seo_og_image"`
	SeoTwitterCard  string       `json:"seo_twitter_card"`
	SeoCanonicalURL string       `json:"seo_canonical_url"`
	SeoRobots       string       `json:"seo_robots"`
	FriendLinks     []FriendLink `json:"friend_links"`
	SiteCopyright   string       `json:"site_copyright"`
}

type IndexResp struct {
	About            *About
	ConfigSite       *ConfigSite
	CategoryTree     []*TreeNode
	CategorySites    []*CategorySite
	FeaturedSites    []model.StSite
	RecommendedSites []model.StSite
	Lang             map[string]string
}

type AboutResp struct {
	About
}
