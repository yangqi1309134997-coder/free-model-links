package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/ch3nnn/webstack-go/internal/dal/model"
	"github.com/ch3nnn/webstack-go/internal/dal/query"
	"github.com/ch3nnn/webstack-go/pkg/gormx"
	"github.com/ch3nnn/webstack-go/pkg/log"
	"github.com/ch3nnn/webstack-go/pkg/zapgorm2"
)

type Repository struct {
	logger *log.Logger
	db     *gorm.DB
}

func NewRepository(logger *log.Logger, db *gorm.DB) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	var (
		db        *gorm.DB
		err       error
		dialector gorm.Dialector
	)

	dsn := conf.GetString("data.db.user.dsn")
	switch conf.GetString("data.db.user.driver") {
	case gormx.MYSQL:
		dialector = mysql.Open(dsn)
	case gormx.POSTGRES:
		dialector = postgres.Open(dsn)
	case gormx.SQLITE:
		dialector = sqlite.Open(dsn)
	default:
		panic("unknown db driver")
	}
	l.Info("db driver", zap.String("driver", dialector.Name()))

	db, err = gorm.Open(dialector, &gorm.Config{
		QueryFields:    true,
		TranslateError: true, // 方言转换错误
		Logger:         zapgorm2.New(l.Logger),
	})
	if err != nil {
		panic(err)
	}

	cachePlugin, _ := cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         config.CacheStorageMemory,
		InvalidateWhenUpdate: true,  // when you create/update/delete objects, invalidate cache
		CacheTTL:             30000, // 30 s
		CacheSize:            10000, // max items
	})
	if err := db.Use(cachePlugin); err != nil {
		panic(err)
	}

	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// set gorm-gen
	query.SetDefault(db)

	// 迁移和初始化操作
	autoMigrateAndInitialize(db)

	return db.Debug()
}

func autoMigrateAndInitialize(db *gorm.DB) {
	ctx := context.Background()

	err := db.AutoMigrate(
		&model.SysConfig{},
		&model.SysUser{},
		&model.SysUserMenu{},
		&model.SysMenu{},
		&model.StCategory{},
		&model.StSite{},
		&model.StModel{},
		&model.StTag{},
		&model.StSiteTag{},
		&model.StPaymentMethod{},
		&model.StSitePayment{},
		&model.StRating{},
		&model.StReview{},
		&model.StSubmission{},
		&model.StSubmissionTag{},
		&model.StSubmissionPayment{},
		&model.StAd{},
	)
	if err != nil {
		fmt.Println("user migrate error")
		os.Exit(0)
	}

	_, err = query.SysUser.WithContext(ctx).
		Where(
			query.SysUser.ID.Eq(1),
		).
		Attrs(
			field.Attrs(&model.SysUser{
				Username: DefaultUname,
				Password: cryptor.Md5String(DefaultUPassword),
				Role:     "admin",
				Status:   "active",
			}),
		).
		FirstOrCreate()
	if err != nil {
		fmt.Println("user migrate error")
		os.Exit(0)
	}

	if _, err = query.SysMenu.WithContext(ctx).
		Where(query.SysMenu.ID.In(menuIDs(DefaultSysMenuAdmin)...)).
		Attrs(
			field.Attrs(&model.SysMenu{
				Pid:    0,
				Name:   "菜单",
				IsUsed: true,
			}),
		).
		FirstOrCreate(); err != nil {
		fmt.Println("menu migrate error")
		os.Exit(0)
	}
	for _, menu := range DefaultSysMenuAdmin {
		count, err := query.SysMenu.WithContext(ctx).Where(query.SysMenu.ID.Eq(menu.ID)).Count()
		if err != nil {
			fmt.Println("menu migrate error")
			os.Exit(0)
		}
		if count == 0 {
			if err := query.SysMenu.WithContext(ctx).Create(menu); err != nil {
				fmt.Println("menu migrate error")
				os.Exit(0)
			}
		}
	}

	if _, err = query.SysUserMenu.WithContext(ctx).
		Where(query.SysUserMenu.MenuID.In(menuIDs(DefaultSysMenuAdmin)...), query.SysUserMenu.UserID.Eq(1)).
		Attrs(
			field.Attrs(&model.SysUserMenu{
				UserID: 1,
			}),
		).
		FirstOrCreate(); err != nil {
		fmt.Println("user menu migrate error")
		os.Exit(0)
	}
	existingUserMenus, err := query.SysUserMenu.WithContext(ctx).Where(query.SysUserMenu.UserID.Eq(1)).Find()
	if err != nil {
		fmt.Println("user menu migrate error")
		os.Exit(0)
	}
	have := make(map[int]bool, len(existingUserMenus))
	for _, um := range existingUserMenus {
		have[um.MenuID] = true
	}
	for _, menu := range DefaultSysMenuAdmin {
		if !have[menu.ID] {
			if err := query.SysUserMenu.WithContext(ctx).Create(&model.SysUserMenu{UserID: 1, MenuID: menu.ID}); err != nil {
				fmt.Println("user menu migrate error")
				os.Exit(0)
			}
		}
	}

	_, err = query.SysConfig.WithContext(ctx).
		Where(
			query.SysConfig.ID.Eq(1),
		).
		Attrs(
			field.Attrs(&model.SysConfig{
				AboutSite:     DefaultAboutSite,
				AboutAuthor:   DefaultAuthor,
				SiteTitle:     DefaultSiteTitle,
				SiteKeyword:   DefaultSiteKeyword,
				SiteDesc:      DefaultSiteDesc,
				SiteFavicon:   DefaultFaviconBase64,
				SiteLogo:      DefaultLogoBase64,
				SiteCopyright: DefaultSiteCopyright,
			}),
		).
		FirstOrCreate()
	if err != nil {
		fmt.Println("config migrate error")
		os.Exit(0)
	}

	// 兼容旧库：若此前从未有过评分/评论数据且匿名开关被默认关闭，
	// 则自动开启匿名评分与匿名评论，保证前台功能开箱即用
	migrateAnonymousFlags(db)

	// 内置一套完整的 AI 大模型导航默认分类
	seedDefaultCategories(db)
	// 内置覆盖主流大模型的模型字典
	seedDefaultModels(db)

	fmt.Println("success initialize")
}

// migrateAnonymousFlags 处理历史数据库：默认关闭的匿名开关在无任何历史评分/评论数据时自动打开
func migrateAnonymousFlags(db *gorm.DB) {
	var config model.SysConfig
	if err := db.First(&config, 1).Error; err != nil {
		return
	}

	if config.AllowAnonymousRating && config.AllowAnonymousReview {
		return
	}

	var ratingCount int64
	db.Model(&model.StRating{}).Count(&ratingCount)
	var reviewCount int64
	db.Model(&model.StReview{}).Count(&reviewCount)
	if ratingCount == 0 && reviewCount == 0 {
		updates := map[string]interface{}{}
		if !config.AllowAnonymousRating {
			updates["allow_anonymous_rating"] = true
		}
		if !config.AllowAnonymousReview {
			updates["allow_anonymous_review"] = true
		}
		if len(updates) > 0 {
			db.Model(&model.SysConfig{}).Where("id = ?", 1).Updates(updates)
		}
	}
}

// menuIDs 提取默认菜单 ID 列表
func menuIDs(menus []*model.SysMenu) []int {
	ids := make([]int, 0, len(menus))
	for _, m := range menus {
		ids = append(ids, m.ID)
	}
	return ids
}

// seedDefaultCategories 仅在分类表为空时，内置一套完整的 AI 大模型导航默认分类
func seedDefaultCategories(db *gorm.DB) {
	var count int64
	db.Model(&model.StCategory{}).Count(&count)
	if count > 0 {
		return
	}

	now := time.Now()
	create := func(id, parentID int, title, icon string, sort int, level int32) model.StCategory {
		return model.StCategory{
			ID:        id,
			ParentID:  parentID,
			Title:     title,
			Icon:      icon,
			Sort:      sort,
			Level:     level,
			IsUsed:    true,
			CreatedAt: &now,
			UpdatedAt: &now,
		}
	}

	categories := []model.StCategory{
		// 一级分类：官方模型平台
		create(1, 0, "官方模型平台", "linecons-star", 1, 1),
		create(2, 1, "国产大模型官方", "linecons-lightbulb", 10, 2),
		create(3, 1, "国外大模型官方", "linecons-globe", 20, 2),
		// 一级分类：中转 API 服务商
		create(4, 0, "中转 API 服务商", "linecons-cloud", 2, 1),
		create(5, 4, "大厂中转站", "linecons-diamond", 10, 2),
		create(6, 4, "普通中转站", "linecons-shop", 20, 2),
		create(7, 4, "国内外全系中转", "linecons-cup", 30, 2),
		// 一级分类：公益与免费
		create(8, 0, "公益与免费", "linecons-heart", 3, 1),
		create(9, 8, "纯公益网站", "linecons-fire", 10, 2),
		create(10, 8, "免费额度网站", "linecons-money", 20, 2),
		// 一级分类：模型工具
		create(11, 0, "模型工具", "linecons-cog", 4, 1),
		create(12, 11, "模型聚合平台", "linecons-database", 10, 2),
		create(13, 11, "模型评测与对比", "linecons-search", 20, 2),
		create(14, 11, "开发工具与社区", "linecons-comment", 30, 2),
	}

	if err := db.Create(&categories).Error; err != nil {
		fmt.Println("seed default categories error")
		return
	}
	fmt.Println("success seed default categories")
}
