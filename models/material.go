package models

import (
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
)

// MaterialAdd 添加材料
func MaterialAdd(m *Material) error {
	err := GlobalDb.Create(m).Error
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

//批量添加材料
func MaterialAddAll(m []Material) error {
	err := GlobalDb.Create(&m).Error
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

//删除材料
func MaterialDel(id int) error {
	material := Material{ID: uint(id)}

	tx := GlobalDb.Begin()

	//出料记录存在则不能删除该物资
	var rc int64 = 0
	err := tx.Model(&Record{}).Where(&Record{MaterialID: material.ID}).Count(&rc).Error
	if err != nil {
		logs.Error("#1 | %s | %s", "models/MaterialDel", err.Error())
		//utils.Mlogger.Errorf("#1 | %s | %s", "models/MaterialDel", err.Error())
		tx.Rollback()
		return errors.New("#1 删除失败")
	}
	if rc != 0 {
		tx.Rollback()
		return errors.New("该物资存在出入库记录，暂无法删除")
	}

	err = tx.Delete(&material).Error
	if err != nil {
		tx.Rollback()
		logs.Error("#2 | %s | %s", "models/MaterialDel", err.Error())
		//utils.Mlogger.Errorf("#2 | %s | %s", "models/MaterialDel", err.Error())
		return errors.New("#2 删除失败")
	}

	tx.Commit()
	return nil

}

//更新材料
func MaterialEditByID(material Material) error {

	//fmt.Println(material)
	err := GlobalDb.Model(&Material{ID: material.ID}).Updates(
		map[string]interface{}{
			"Name":     material.Name,
			"Model":    material.Model,
			"Unit":     material.Unit,
			"NickName": material.NickName,

			//""
			"CarID":    material.CarID,
			"PlaceID":  material.PlaceID,
			"Floor":    material.Floor,
			"Location": material.Location,

			"Count":        material.Count,
			"PrepareCount": material.PrepareCount,
			"WarnCount":    material.WarnCount,
			"Marks":        material.Marks,

			"UserID": material.UserID,
		}).Error
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

//根据id获取材料
func MaterialGetById(id int) (Material, error) {

	material := Material{}
	err := GlobalDb.
		Preload("Car").
		Preload("Place").
		Where(&Material{ID: uint(id)}).First(&material).Error
	if err != nil {
		return Material{}, err
	}

	return material, nil
}

//搜索材料
func MaterialSearchByKey(key string, car, place, page, perPage int) ([]Material, int64, error) {
	var materials []Material
	var count int64 = 0

	//查询符合条件的材料
	regKey := "%" + key + "%"
	err := GlobalDb.
		Preload("Place").
		Preload("Car").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Car").Omit("password")
		}).
		Where("materials.place_id = ? or 0 = ?", place, place).
		Where("materials.car_id = ? or 0 = ?", car, car).
		Where("name LIKE ? OR model LIKE ?", regKey, regKey).
		Joins("join places on places.id = materials.place_id").
		Order("places.position").
		Order("floor").
		Order("location").
		Order("id").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&materials).Error
	if err != nil {
		logs.Error(err)
		return []Material{}, 0, err
	}

	//查询数量
	err = GlobalDb.Model(&Material{}).Preload("Place").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Car").Omit("password")
	}).
		Where("place_id = ? or 0 = ?", place, place).
		Where("car_id = ? or 0 = ?", car, car).
		Where("name LIKE ? OR model LIKE ?", regKey, regKey).
		Count(&count).Error
	if err != nil {
		logs.Error(err)
		return []Material{}, 0, err
	}

	return materials, count, nil
}

//搜索材料
func MaterialWarnByCar(car, place, page, perPage int) ([]Material, int64, error) {
	var materials []Material
	var count int64 = 0

	//查询符合条件的材料
	//regKey := "%" + key + "%"
	err := GlobalDb.
		Preload("Place").
		Preload("Car").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Car").Omit("password")
		}).
		Where("place_id = ? or 0 = ?", place, place).
		Where("materials.car_id = ? or 0 = ?", car, car).
		Where("count < prepare_count").
		//Where("name LIKE ? OR model LIKE ?", regKey, regKey).
		Joins("join places on places.id = materials.place_id").
		Order("places.position").
		Order("floor").
		Order("location").
		Offset((page - 1) * perPage).Limit(perPage).
		Find(&materials).Error
	if err != nil {
		logs.Error(err)
		return []Material{}, 0, err
	}

	//查询数量
	err = GlobalDb.Model(&Material{}).Preload("Place").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Car").Omit("password")
	}).
		Where("place_id = ? or 0 = ?", place, place).
		Where("car_id = ? or 0 = ?", car, car).
		Where("count < prepare_count").
		Count(&count).Error
	if err != nil {
		logs.Error(err)
		return []Material{}, 0, err
	}

	return materials, count, nil
}

// MaterialDownloadByKey 下载材料清单
//下载时必须选择车号
func MaterialDownloadByKey(key string, carID, placeID int) ([]Material, error) {
	var materials []Material
	if carID <= 0 {
		return nil, errors.New("下载材料清单必须选择车号")
	}

	//查询符合条件的材料
	regKey := "%" + key + "%"
	err := GlobalDb.
		Preload("Place").
		Preload("Car").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Car").Omit("password")
		}).
		Where("place_id = ? or 0 = ?", placeID, placeID).
		Where("car_id = ? or 0 = ?", carID, carID).
		Where("name LIKE ? OR model LIKE ?", regKey, regKey).
		Order("place_id").
		Order("floor").
		Order("location").
		Order("id").
		Find(&materials).Error
	if err != nil {
		logs.Error(err)
		return []Material{}, err
	}

	return materials, nil
}

// MaterialDownloadWarnByCar  下载备料材料清单
//下载时必须选择车号
func MaterialDownloadWarnByCar(carID, placeID int) ([]Material, error) {
	var materials []Material
	if carID <= 0 {
		return nil, errors.New("下载材料清单必须选择车号")
	}

	//查询符合条件的材料
	//regKey := "%" + key + "%"
	err := GlobalDb.
		Preload("Place").
		Preload("Car").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Car").Omit("password")
		}).
		Where("place_id = ? or 0 = ?", placeID, placeID).
		Where("car_id = ? or 0 = ?", carID, carID).
		Where("count < prepare_count").
		Order("place_id").
		Order("floor").
		Order("location").
		Order("id").
		Find(&materials).Error
	if err != nil {
		logs.Error(err)
		return []Material{}, err
	}

	return materials, nil
}
