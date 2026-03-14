package system

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

type SysErrorService struct{}

func (sysErrorService *SysErrorService) CreateSysError(ctx context.Context, sysError *system.SysError) error {
	if global.GVA_DB == nil {
		return nil
	}
	return global.GVA_DB.Create(sysError).Error
}

func (sysErrorService *SysErrorService) DeleteSysError(ctx context.Context, id string) error {
	return global.GVA_DB.Delete(&system.SysError{}, "id = ?", id).Error
}

func (sysErrorService *SysErrorService) DeleteSysErrorByIds(ctx context.Context, ids []string) error {
	return global.GVA_DB.Delete(&[]system.SysError{}, "id in ?", ids).Error
}

func (sysErrorService *SysErrorService) UpdateSysError(ctx context.Context, sysError system.SysError) error {
	return global.GVA_DB.Model(&system.SysError{}).Where("id = ?", sysError.ID).Updates(&sysError).Error
}

func (sysErrorService *SysErrorService) GetSysError(ctx context.Context, id string) (system.SysError, error) {
	var sysError system.SysError
	err := global.GVA_DB.Where("id = ?", id).First(&sysError).Error
	return sysError, err
}

func (sysErrorService *SysErrorService) GetSysErrorInfoList(ctx context.Context, info systemReq.SysErrorSearch) ([]system.SysError, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&system.SysError{}).Order("created_at desc")
	var sysErrors []system.SysError
	var total int64

	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if info.Form != nil && *info.Form != "" {
		db = db.Where("form = ?", *info.Form)
	}
	if info.Info != nil && *info.Info != "" {
		db = db.Where("info LIKE ?", "%"+*info.Info+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err := db.Find(&sysErrors).Error
	return sysErrors, total, err
}

func (sysErrorService *SysErrorService) GetSysErrorSolution(ctx context.Context, id string) error {
	err := global.GVA_DB.WithContext(ctx).Model(&system.SysError{}).Where("id = ?", id).Update("status", "处理中").Error
	if err != nil {
		return err
	}

	go func(targetID string) {
		var se system.SysError
		_ = global.GVA_DB.Model(&system.SysError{}).Where("id = ?", targetID).First(&se).Error

		solution := "当前脚手架已移除 AI 自动分析能力，请根据错误日志手动排查。"
		if se.Info != nil && *se.Info != "" {
			solution = fmt.Sprintf("当前脚手架已移除 AI 自动分析能力，请根据以下错误日志手动排查：%s", *se.Info)
		}

		_ = global.GVA_DB.Model(&system.SysError{}).Where("id = ?", targetID).Updates(map[string]interface{}{
			"status":   "处理完成",
			"solution": solution,
		}).Error
	}(id)

	return nil
}
