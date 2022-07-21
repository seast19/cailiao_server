package models

import (
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"time"
)

// RecordAddSend RecordAdd 添加出库记录
func RecordAddSend(r *Record) error {
	// 构建完整数据
	t := time.Now().UnixNano() / 1e6
	r.CreateAt = t
	r.UpdateAt = t

	//启用事物
	tx := GlobalDb.Begin()

	//查找该物资
	m := Material{}
	m.ID = r.MaterialID
	err := tx.First(&m).Error
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		return errors.New("无此物资id")
	}
	//库存为0
	if m.Count <= 0 {
		tx.Rollback()
		return errors.New("此物资库存数为0")
	}

	//库存数比发料数少
	if m.Count < r.SendCount {
		tx.Rollback()
		return errors.New("此物资库存数不足")
	}

	//创建record记录
	err = tx.Create(r).Error
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		return errors.New("添加记录失败")
	}

	//修改material库存数
	err = tx.Model(&m).
		Select("count").
		Updates(map[string]interface{}{"count": m.Count - r.SendCount}).Error
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		return errors.New("修改材料数量失败")
	}

	tx.Commit()
	return nil
}

// RecordAddReceive RecordAdd 添加入库记录
func RecordAddReceive(r *Record) error {
	// 构建完整数据
	t := time.Now().UnixNano() / 1e6
	r.CreateAt = t
	r.UpdateAt = t
	m := Material{}
	m.ID = r.MaterialID

	//开启事物
	tx := GlobalDb.Begin()

	//查找该物资
	err := tx.First(&m).Error
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		return errors.New("无此物资id")
	}

	//创建record记录
	err = tx.Create(r).Error
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		return errors.New("添加记录失败")
	}

	//修改material库存数
	err = tx.Model(&m).
		Select("count").
		Updates(map[string]interface{}{"count": m.Count + r.ReceiveCount}).Error
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		return errors.New("添加记录失败")
	}

	tx.Commit()
	return nil

}

// RecordDelById 根据id删除出入库记录
func RecordDelById(id int) error {
	r := Record{ID: uint(id)}

	tx := GlobalDb.Begin()

	//找到该条记录
	err := tx.First(&r).Error
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		return errors.New("找不到记录id")
	}

	// 查询该物资
	m := Material{}
	m.ID = r.MaterialID
	err = tx.First(&m).Error
	if err != nil {
		logs.Error(err)
		tx.Rollback()
		return errors.New("找不到物资id")
	}

	//按操作类型 更新材料表
	switch r.Type {
	case "send":
		//修改material库存数
		err = tx.Model(&m).Select("count").Updates(map[string]interface{}{"count": m.Count + r.SendCount}).Error
		if err != nil {
			logs.Error("#3", err)
			tx.Rollback()
			return errors.New("操作失败")
		}
	case "receive":
		//修改material库存数
		err = tx.Model(&m).Select("count").Updates(map[string]interface{}{"count": m.Count - r.ReceiveCount}).Error
		if err != nil {
			logs.Error("#4", err)
			tx.Rollback()
			return errors.New("操作失败")
		}
	default:
		tx.Rollback()
		return errors.New("操作类型错误")

	}

	//删除该记录
	err = tx.Delete(&r).Error
	if err != nil {
		logs.Error("#5", err)
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

// RecordGetAllWithCarByPage 根据车号获取消耗记录
func RecordGetAllWithCarByPage(page, perPage int, t string, startT, stopT int, cid uint) ([]Record, int64, error) {
	//构建参数
	records := []Record{}
	materials := []Material{}
	mids := []uint{}
	var count int64 = 0

	//查找符合carid 的所有material
	if cid != 0 {
		err := GlobalDb.
			Preload("Car").
			Where("car_id = ? or 0 = ?", cid, cid).
			Find(&materials).Error
		if err != nil {
			logs.Error(err)
			return []Record{}, 0, err
		}
		for _, material := range materials {
			mids = append(mids, material.ID)
		}
	}

	err := GlobalDb.
		Preload("Material.Place").
		Preload("Material.Car").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Car").Omit("password")
		}).
		Where("type = ? or '' = ?", t, t).
		Where("(create_at >= ? AND create_at <= ?) or (0 = ?)", startT, stopT, startT).
		Where("material_id in ? or 0 = ?", mids, cid).
		Order("id DESC").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&records).Error
	if err != nil {
		logs.Error(err)
		return []Record{}, 0, err
	}

	//查数量
	err = GlobalDb.
		Model(&Record{}).
		Preload("Material.Place").
		Preload("Material.Car").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Car").Omit("password")
		}).
		Where("type = ? or '' = ?", t, t).
		Where("(create_at >= ? AND create_at <= ?) or (0 = ?)", startT, stopT, startT).
		Where("material_id in ? or 0 = ?", mids, cid).
		//Offset((page - 1) * perPage).Limit(perPage).
		//Order("id DESC").
		Count(&count).Error
	if err != nil {
		logs.Error(err)
		return []Record{}, 0, err
	}

	return records, count, nil

}

// RecordGetAllWithMIdByPage 根据材料id获取消耗记录
func RecordGetAllWithMIdByPage(page, perPage int, t string, startT, stopT int, mid uint) ([]Record, int64, error) {
	//构建参数
	records := []Record{}
	var count int64 = 0

	err := GlobalDb.
		Preload("Material.Place").
		Preload("Material.Car").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Car").Omit("password")
		}).
		Where("type = ? or '' = ?", t, t).
		Where("(create_at >= ? AND create_at <= ?) or (0 = ?)", startT, stopT, startT).
		Where("material_id = ? or 0 = ?", mid, mid).
		Order("id DESC").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&records).Error
	if err != nil {
		logs.Error(err)
		return []Record{}, 0, err
	}

	//查数量
	err = GlobalDb.Model(&Record{}).
		Preload("Material.Place").
		Preload("Material.Car").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Car").Omit("password")
		}).
		Where("type = ? or '' = ?", t, t).
		Where("(create_at >= ? AND create_at <= ?) or (0 = ?)", startT, stopT, startT).
		Where("material_id = ? or 0 = ?", mid, mid).
		//Offset((page - 1) * perPage).Limit(perPage).
		//Order("id DESC").
		Count(&count).Error
	if err != nil {
		logs.Error(err)
		return []Record{}, 0, err
	}

	return records, count, nil
}

//聚合车号获取时间范围内的材料入库、出库记录
func RecordGetPolyWithCardByPage(page, perPage int, startT, stopT int, cid uint) ([]RecordPoly, int64, error) {
	recordPoly := []RecordPoly{}
	materials := []Material{}
	mids := []uint{}
	var count int64

	//查找符合carid 的所有material
	if cid != 0 {
		err := GlobalDb.
			Preload("Car").
			Where("car_id = ? or 0 = ?", cid, cid).
			Find(&materials).Error
		if err != nil {
			return []RecordPoly{}, 0, errors.New("查询失败")
		}
		for _, material := range materials {
			mids = append(mids, material.ID)
		}
	}

	err := GlobalDb.Model(&Record{}).
		Preload("Material.Car").
		Preload("Material.Place").
		Where("(records.create_at >= ? AND records.create_at <= ? ) or (0 = ?)", startT, stopT, startT).
		Where("records.material_id in ? or 0 = ?", mids, cid).
		Select("material_id, SUM( send_count ) as send ,SUM( receive_count ) as receive").
		Joins("JOIN materials ON materials.id = material_id ").
		Group("material_id").
		Find(&recordPoly).Error
	if err != nil {
		logs.Error(err)
		return nil, 0, errors.New("查询失败")
	}
	//查数量
	err = GlobalDb.Model(&Record{}).
		Preload("Material.Car").
		Preload("Material.Place").
		Where("(records.create_at >= ? AND records.create_at <= ? ) or (0 = ?)", startT, stopT, startT).
		Where("records.material_id in ? or 0 = ?", mids, cid).
		Select("material_id, SUM( send_count ) as send ,SUM( receive_count ) as receive").
		Joins("JOIN materials ON materials.id = material_id ").
		Group("material_id").
		Count(&count).Error
	if err != nil {
		logs.Error(err)
		return nil, 0, errors.New("查询失败")
	}

	return recordPoly, count, nil
}
