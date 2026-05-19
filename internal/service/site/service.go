/**
 * @Author: chentong
 * @Date: 2024/05/27 下午5:56
 */

package site

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/mat/besticon/besticon"
	"github.com/pkg/errors"

	v1 "github.com/ch3nnn/webstack-go/api/v1"
	"github.com/ch3nnn/webstack-go/internal/dal/repository"
	s "github.com/ch3nnn/webstack-go/internal/service"
	"github.com/ch3nnn/webstack-go/pkg/tools"
)

// normalizeURL ensures the URL has a scheme; defaults to https if missing.
func normalizeURL(u string) string {
	u = strings.TrimSpace(u)
	if u == "" {
		return u
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return "https://" + u
	}
	return u
}

// Service 层
var _ Service = (*service)(nil)

type Service interface {
	i()

	// List 站点列表
	List(ctx context.Context, req *v1.SiteListReq) (resp *v1.SiteListResp, err error)
	// Delete 删除站点
	Delete(ctx context.Context, req *v1.SiteDeleteReq) (resp *v1.SiteDeleteResp, err error)
	// Update 更新站点
	Update(ctx *gin.Context, req *v1.SiteUpdateReq) (resp *v1.SiteUpdateResp, err error)
	// BatchCreate 批量创建站点
	BatchCreate(ctx context.Context, req *v1.SiteCreateReq) (resp *v1.SiteCreateResp, err error)
	// Sync 同步站点信息
	Sync(ctx *gin.Context, req *v1.SiteSyncReq) (resp *v1.SiteSyncResp, err error)
	// Export 导出站点信息
	Export(ctx *gin.Context, req *v1.SiteExportReq) (resp *v1.SiteExportResp, err error)
}

type service struct {
	*s.Service
	siteRepository     repository.IStSiteDao
	categoryRepository repository.IStCategoryDao
}

func NewService(
	s *s.Service,
	siteRepository repository.IStSiteDao,
	categoryRepository repository.IStCategoryDao,
) Service {
	return &service{
		Service:            s,
		siteRepository:     siteRepository,
		categoryRepository: categoryRepository,
	}
}

func (s *service) i() {}

// -----------------------------------------------------------------------------------------------------------------------------------------

func getWebLogoIconBase64(url string) (string, error) {
	url = normalizeURL(url)
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: tools.NewHTTPTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}),
	}

	b := besticon.New(
		besticon.WithHTTPClient(client),
		besticon.WithLogger(besticon.NewDefaultLogger(io.Discard)), // disable verbose logging
	)

	finder := b.NewIconFinder()
	if _, err := finder.FetchIcons(url); err != nil {
		return repository.DefaultFaviconBase64, err
	}

	icon := finder.IconInSizeRange(besticon.SizeRange{Min: 16, Perfect: 64, Max: 512})
	if icon == nil {
		return repository.DefaultFaviconBase64, errors.New("failed to fetch icons")
	}

	return base64.StdEncoding.EncodeToString(icon.ImageData), nil
}

func getWebTitle(url string) (title string, err error) {
	url = normalizeURL(url)
	c := tools.NewColly()
	c.OnHTML("title", func(e *colly.HTMLElement) {
		if title == "" {
			title = strings.TrimSpace(e.Text)
		}
	})
	// fallback to og:title
	c.OnHTML("meta[property='og:title']", func(e *colly.HTMLElement) {
		if title == "" {
			title = strings.TrimSpace(e.Attr("content"))
		}
	})
	if err := c.Visit(url); err != nil {
		return "", err
	}

	return
}

func getWebDescription(url string) (doc string, err error) {
	url = normalizeURL(url)
	c := tools.NewColly()
	// 兼容多种 description meta 写法
	c.OnHTML("meta[name='description']", func(e *colly.HTMLElement) {
		if doc == "" {
			doc = strings.TrimSpace(e.Attr("content"))
		}
	})
	c.OnHTML("meta[name='Description']", func(e *colly.HTMLElement) {
		if doc == "" {
			doc = strings.TrimSpace(e.Attr("content"))
		}
	})
	c.OnHTML("meta[property='og:description']", func(e *colly.HTMLElement) {
		if doc == "" {
			doc = strings.TrimSpace(e.Attr("content"))
		}
	})
	c.OnHTML("meta[name='twitter:description']", func(e *colly.HTMLElement) {
		if doc == "" {
			doc = strings.TrimSpace(e.Attr("content"))
		}
	})
	if err := c.Visit(url); err != nil {
		return "", err
	}

	return
}
