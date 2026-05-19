package repository

import (
	"context"

	"github.com/ch3nnn/webstack-go/internal/dal/model"
	"github.com/ch3nnn/webstack-go/internal/dal/query"
)

var DefaultSysMenuAdmin = []*model.SysMenu{
	{
		ID:     1,
		Pid:    0,
		Name:   "网站管理",
		Icon:   "users",
		Level:  1,
		Sort:   500,
		IsUsed: true,
	},
	{
		ID:     2,
		Pid:    1,
		Name:   "网站分类",
		Link:   "/admin/category",
		Level:  1,
		Sort:   501,
		IsUsed: true,
	},
	{
		ID:     3,
		Pid:    1,
		Name:   "网站信息",
		Link:   "/admin/site",
		Level:  1,
		Sort:   502,
		IsUsed: true,
	},
	{
		ID:     4,
		Pid:    0,
		Name:   "系统管理",
		Level:  1,
		Sort:   600,
		IsUsed: true,
	},
	{
		ID:     5,
		Pid:    4,
		Name:   "网站配置",
		Link:   "/admin/config",
		Level:  1,
		Sort:   601,
		IsUsed: true,
	},
	{
		ID:     6,
		Pid:    0,
		Name:   "评价管理",
		Icon:   "star",
		Level:  1,
		Sort:   550,
		IsUsed: true,
	},
	{
		ID:     7,
		Pid:    6,
		Name:   "评论审核",
		Link:   "/admin/reviews",
		Level:  1,
		Sort:   551,
		IsUsed: true,
	},
	{
		ID:     8,
		Pid:    0,
		Name:   "用户管理",
		Icon:   "user",
		Level:  1,
		Sort:   560,
		IsUsed: true,
	},
	{
		ID:     9,
		Pid:    8,
		Name:   "用户列表",
		Link:   "/admin/users",
		Level:  1,
		Sort:   561,
		IsUsed: true,
	},
	{
		ID:     10,
		Pid:    0,
		Name:   "字典管理",
		Icon:   "book",
		Level:  1,
		Sort:   570,
		IsUsed: true,
	},
	{
		ID:     11,
		Pid:    10,
		Name:   "模型管理",
		Link:   "/admin/models",
		Level:  1,
		Sort:   571,
		IsUsed: true,
	},
	{
		ID:     12,
		Pid:    10,
		Name:   "标签管理",
		Link:   "/admin/tags",
		Level:  1,
		Sort:   572,
		IsUsed: true,
	},
	{
		ID:     13,
		Pid:    10,
		Name:   "支付方式",
		Link:   "/admin/payments",
		Level:  1,
		Sort:   573,
		IsUsed: true,
	},
	{
		ID:     14,
		Pid:    6,
		Name:   "投稿管理",
		Link:   "/admin/submissions",
		Level:  1,
		Sort:   552,
		IsUsed: true,
	},
	{
		ID:     15,
		Pid:    4,
		Name:   "网站设置",
		Link:   "/admin/settings",
		Level:  1,
		Sort:   602,
		IsUsed: true,
	},
}

var _ iCustomGenSysMenuFunc = (*customSysMenuDao)(nil)

type (
	// ISysMenuDao not edit interface name
	ISysMenuDao interface {
		iWhereSysMenuFunc
		WithContext(ctx context.Context) iCustomGenSysMenuFunc

		// TODO Custom WhereFunc ....
		// ...
	}

	// not edit interface name
	iCustomGenSysMenuFunc interface {
		iGenSysMenuFunc

		// TODO Custom DaoFunc ....
		// ...
	}

	// not edit interface name
	customSysMenuDao struct {
		sysMenuDao
	}
)

func NewSysMenuDao() ISysMenuDao {
	return &customSysMenuDao{
		sysMenuDao{
			sysMenuDo: query.SysMenu.WithContext(context.Background()),
		},
	}
}

func (d *customSysMenuDao) WithContext(ctx context.Context) iCustomGenSysMenuFunc {
	d.sysMenuDo = d.sysMenuDo.WithContext(ctx)
	return d
}
