/**
 * @Author: chentong
 * @Date: 2024/05/26 上午1:52
 */

package index

import (
	"context"
	"encoding/json"
	"sort"

	"golang.org/x/sync/errgroup"

	v1 "github.com/ch3nnn/webstack-go/api/v1"
	"github.com/ch3nnn/webstack-go/internal/dal/model"
	"github.com/ch3nnn/webstack-go/internal/dal/query"
	"github.com/ch3nnn/webstack-go/internal/i18n"
)

// buildTree 构建树形结构
func buildTree(nodes []*v1.TreeNode, pid int) []*v1.TreeNode {
	var treeNodes []*v1.TreeNode
	for _, node := range nodes {
		if node.Pid == pid {
			node.Child = buildTree(nodes, node.Id)
			treeNodes = append(treeNodes, node)
		}
	}
	return treeNodes
}

// categoryTree 对树形结构按 Sort 字段排序
func categoryTree(nodes []*v1.TreeNode) []*v1.TreeNode {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Sort < nodes[j].Sort
	})

	for _, node := range nodes {
		if len(node.Child) > 0 {
			categoryTree(node.Child)
		}
	}
	return nodes
}

// categorySites 将站点数据归类到分类站点中
func categorySites(sites []*model.StSite, treeNodes []*v1.TreeNode) (data []*v1.CategorySite) {
	for _, node := range treeNodes {
		categorySite := &v1.CategorySite{
			Category: node.Name,
			SiteList: []model.StSite{},
		}

		for _, site := range sites {
			if site.CategoryID == node.Id {
				categorySite.SiteList = append(categorySite.SiteList, *site)
			}
		}
		//  Sort 字段进行升序排序
		sort.Slice(categorySite.SiteList, func(i, j int) bool {
			return categorySite.SiteList[i].Sort < categorySite.SiteList[j].Sort
		})

		if len(categorySite.SiteList) > 0 {
			data = append(data, categorySite)
		}

		if len(node.Child) > 0 {
			childCategorySites := categorySites(sites, node.Child)
			data = append(data, childCategorySites...)
		}
	}

	return data
}

// Index 获取首页数据
func (s *service) Index(ctx context.Context) (*v1.IndexResp, error) {
	var (
		g          errgroup.Group
		sysConfig  *model.SysConfig
		sites      []*model.StSite
		categories []*model.StCategory
	)

	g.Go(func() (err error) {
		categories, err = s.categoryRepo.WithContext(ctx).FindAllOrderBySort(query.StCategory.Sort.Abs(), s.categoryRepo.WhereByIsUsed(true))
		return err
	})

	g.Go(func() (err error) {
		sites, err = s.siteRepo.WithContext(ctx).FindAll(s.siteRepo.WhereByIsUsed(true))
		return err
	})

	g.Go(func() (err error) {
		sysConfig, err = s.configRepo.WithContext(ctx).FindOne()
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	nodes := make([]*v1.TreeNode, len(categories))
	for i, category := range categories {
		nodes[i] = &v1.TreeNode{
			Id:   category.ID,
			Pid:  category.ParentID,
			Name: category.Title,
			Icon: category.Icon,
			Sort: category.Sort,
		}
	}

	categoryTree := categoryTree(buildTree(nodes, 0))
	categorySites := categorySites(sites, categoryTree)

	// 解析友情链接 JSON
	var friendLinks []v1.FriendLink
	if sysConfig.FriendLinks != "" {
		_ = json.Unmarshal([]byte(sysConfig.FriendLinks), &friendLinks)
	}
	// 若数据库为空则提供默认友情链接
	if len(friendLinks) == 0 {
		friendLinks = []v1.FriendLink{
			{Name: "幻城网安科技", URL: "https://hcnsec.cn", Desc: "新疆幻城网安科技有限责任公司官网"},
			{Name: "公益大模型 API 分发", URL: "https://api.iamhc.cn", Desc: "公益大模型 API 分发平台"},
		}
	}

	return &v1.IndexResp{
		ConfigSite: &v1.ConfigSite{
			SiteTitle:             sysConfig.SiteTitle,
			SiteKeyword:           sysConfig.SiteKeyword,
			SiteDesc:              sysConfig.SiteDesc,
			SiteRecord:            sysConfig.SiteRecord,
			SiteURL:               sysConfig.SiteURL,
			SiteLogo:              sysConfig.SiteLogo,
			SiteFavicon:           sysConfig.SiteFavicon,
			SiteLang:              sysConfig.SiteLang,
			AllowAnonymousRating:  sysConfig.AllowAnonymousRating,
			AllowAnonymousReview:  sysConfig.AllowAnonymousReview,
			AllowSubmission:       sysConfig.AllowSubmission,
			SeoAuthor:             sysConfig.SeoAuthor,
			SeoOgImage:            sysConfig.SeoOgImage,
			SeoTwitterCard:        sysConfig.SeoTwitterCard,
			SeoCanonicalURL:       sysConfig.SeoCanonicalURL,
			SeoRobots:             sysConfig.SeoRobots,
			FriendLinks:           friendLinks,
			SiteCopyright:         sysConfig.SiteCopyright,
		},
		About: &v1.About{
			AboutSite:   sysConfig.AboutSite,
			AboutAuthor: sysConfig.AboutAuthor,
			IsAbout:     sysConfig.IsAbout,
		},
		CategoryTree:  categoryTree,
		CategorySites: categorySites,
		Lang:          buildLangMap(sysConfig.SiteLang),
	}, nil
}

func buildLangMap(lang string) map[string]string {
	if lang == "" {
		lang = "zh"
	}
	m := make(map[string]string)
	for key := range i18n.Translations {
		m[key] = i18n.T(key, lang)
	}
	return m
}
