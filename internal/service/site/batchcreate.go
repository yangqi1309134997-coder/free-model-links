/**
 * @Author: chentong
 * @Date: 2024/06/04 下午4:33
 */

package site

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/validator"
	"golang.org/x/sync/errgroup"

	v1 "github.com/ch3nnn/webstack-go/api/v1"
	"github.com/ch3nnn/webstack-go/internal/dal/model"
	"github.com/ch3nnn/webstack-go/internal/dal/query"
	"github.com/ch3nnn/webstack-go/internal/dal/repository"
	"github.com/ch3nnn/webstack-go/pkg/tools"
)

func (s *service) parseURL(u string) (urls []string) {
	for _, u := range strings.Split(u, "\n") {
		if validator.IsUrl(u) || validator.IsIp(u) || validator.IsIpPort(u) {
			urls = append(urls, u)
		}
	}

	return
}

func (s *service) BatchCreate(ctx context.Context, req *v1.SiteCreateReq) (*v1.SiteCreateResp, error) {
	workerPool := tools.NewWorkerPool(5, 20)
	workerPool.Start()

	var successCnt int
	var failURLs []string
	for _, u := range s.parseURL(req.Url) {
		workerPool.AddJob(func() {
			var (
				g                 errgroup.Group
				title, icon, desc string
			)

			u = strings.TrimSpace(u)
			u = normalizeURL(u)

			// 如果手动指定了字段则直接使用，跳过对应的抓取
			if req.ManualTitle != "" {
				title = req.ManualTitle
			} else {
				g.Go(func() (err error) {
					title, err = getWebTitle(u)
					return
				})
			}

			if req.ManualIcon != "" {
				if b64, ierr := tools.ResizeURLImgToBase64(req.ManualIcon, FaviconWidth, FaviconHeight); ierr == nil {
					icon = b64
				} else {
					icon = repository.DefaultFaviconBase64
				}
			} else {
				g.Go(func() (err error) {
					icon, err = getWebLogoIconBase64(u)
					return
				})
			}

			if req.ManualDescription != "" {
				desc = req.ManualDescription
			} else {
				g.Go(func() (err error) {
					desc, err = getWebDescription(u)
					return
				})
			}

			if err := g.Wait(); err != nil {
				if !req.FailSwitch {
					failURLs = append(failURLs, u)
					return
				}
			}

			site := &model.StSite{
				Title:       condition.Ternary(title != "", title, u),
				Icon:        icon,
				Description: desc,
				URL:         u,
				CategoryID:  req.CategoryID,
				IsUsed:      req.IsUsed,
				Sort:        0,
			}
			created, err := s.siteRepository.WithContext(ctx).Create(site)
			if err != nil {
				failURLs = append(failURLs, u)
				return
			}

			db := query.StSite.WithContext(ctx).UnderlyingDB()

			// Save tag associations
			if req.TagIDs != "" {
				for _, idStr := range strings.Split(req.TagIDs, ",") {
					tagID, _ := strconv.Atoi(strings.TrimSpace(idStr))
					if tagID > 0 {
						db.Table("st_site_tag").Create(map[string]interface{}{"site_id": created.ID, "tag_id": tagID})
					}
				}
			}

			// Save model associations (JSON in supported_models)
			if req.ModelIDs != "" {
				var modelIDs []int
				for _, idStr := range strings.Split(req.ModelIDs, ",") {
					mid, _ := strconv.Atoi(strings.TrimSpace(idStr))
					if mid > 0 {
						modelIDs = append(modelIDs, mid)
					}
				}
				if len(modelIDs) > 0 {
					jsonData, _ := json.Marshal(modelIDs)
					db.Table("st_site").Where("id = ?", created.ID).Update("supported_models", string(jsonData))
				}
			}

			// Save payment method associations
			if req.PaymentMethodIDs != "" {
				for _, idStr := range strings.Split(req.PaymentMethodIDs, ",") {
					pmID, _ := strconv.Atoi(strings.TrimSpace(idStr))
					if pmID > 0 {
						db.Table("st_site_payment").Create(map[string]interface{}{"site_id": created.ID, "payment_method_id": pmID})
					}
				}
			}

			successCnt++
		})
	}

	workerPool.Wait()

	return &v1.SiteCreateResp{
		FailCount:    len(failURLs),
		SuccessCount: successCnt,
		FailURLs:     failURLs,
	}, nil
}
