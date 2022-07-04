package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
)

//删除材料
func MaterialDel(id int) error {
	//GLOBAL_DB, err := GLOBAL_DB
	//if err != nil {
	//	return err
	//}
	//defer GLOBAL_DB.Close()

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

//添加材料
func MaterialAdd(m *Material) error {
	err := GlobalDb.Create(m).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//批量添加材料
func MaterialAddAll(m []Material) error {
	err := GlobalDb.Create(&m).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//分页获取材料
func MaterialGetByPage(page, perPage, placeID int) ([]Material, int64, error) {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return []Material{}, 0, err
	//}
	//defer GLOBAL_DB.Close()

	//placeID=36

	material := []Material{}
	var count int64 = 0
	err := GlobalDb.Preload("Place").Where(&Material{PlaceID: uint(placeID)}).Offset((page - 1) * perPage).Limit(perPage).Order("id DESC").Find(&material).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败1")
	}
	//fmt.Println(material[0])
	//m := material[0]
	//m := Material{}

	//fmt.Println(m)
	//	GLOBAL_DB.Last(&m)
	//	fmt.Println(m)
	//	err=GLOBAL_DB.Model(&m).Related(&m.Place,"PlaceID").Error
	//	//res:=GLOBAL_DB.First(&material[0].Place).Error
	//	fmt.Println(m)
	//	if err != nil {
	//		fmt.Println(err)
	//		return nil, 0, errors.New("查询失败2")
	//	}

	err = GlobalDb.Model(Material{}).Where(&Material{PlaceID: uint(placeID)}).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败3")
	}

	return material, count, nil
}

//根据id获取材料
func MaterialGetById(id int) (Material, error) {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return Material{}, err
	//}
	//defer GLOBAL_DB.Close()

	material := Material{}
	err := GlobalDb.Where(&Material{ID: uint(id)}).First(&material).Error
	if err != nil {
		return Material{}, err
	}

	return material, nil
}

//更新材料
func MaterialEditByID(material Material) error {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return err
	//}
	//defer GLOBAL_DB.Close()

	fmt.Println(material)
	err := GlobalDb.Model(&Material{ID: material.ID}).Updates(
		map[string]interface{}{
			"Name":     material.Name,
			"Model":    material.Model,
			"Unit":     material.Unit,
			"NickName": material.NickName,

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
		fmt.Println(err)
		return err
	}
	return nil
}

// MaterialDownloadByKey 下载材料清单
//下载时必须选择车号
func MaterialDownloadByKey(key string, car, place, page, perPage int) ([]Material, error) {
	var materials []Material
	users := []User{}
	userId := []uint{}

	if car <= 0 {
		return nil, errors.New("下载材料清单必须选择车号")
	}

	//存在car参数则查找car对应的user ID
	err := GlobalDb.Preload("Car").Where("car_id = ? or 0 = ?", car, car).Find(&users).Error
	if err != nil {
		fmt.Println(err)
		return []Material{}, err
	}
	for _, user := range users {
		userId = append(userId, user.ID)
	}

	//查询符合条件的材料
	regKey := "%" + key + "%"
	err = GlobalDb.Preload("Place").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Car").Omit("password")
	}).
		Where("place_id = ? or 0 = ?", place, place).
		Where("user_id in ? or 0 = ?", userId, car).
		Where("name LIKE ? OR model LIKE ?", regKey, regKey).
		Order("place_id").
		Order("floor").
		Order("location").
		Find(&materials).Error
	if err != nil {
		fmt.Println(err)
		return []Material{}, err
	}

	return materials, nil
}

//搜索材料
func MaterialSearchByKey(key string, car, place, page, perPage int) ([]Material, int64, error) {

	var materials []Material
	users := []User{}
	userId := []uint{}
	var count int64 = 0

	//存在car参数则查找car对应的user ID
	if car > 0 {
		err := GlobalDb.Preload("Car").Where("car_id = ? or 0 = ?", car, car).Find(&users).Error
		if err != nil {
			fmt.Println(err)
			return []Material{}, 0, err
		}
		for _, user := range users {
			userId = append(userId, user.ID)
		}
	}
	//查询符合条件的材料
	regKey := "%" + key + "%"
	err := GlobalDb.Preload("Place").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Car").Omit("password")
	}).
		Where("place_id = ? or 0 = ?", place, place).
		Where("user_id in ? or 0 = ?", userId, car).
		Where("name LIKE ? OR model LIKE ?", regKey, regKey).
		Offset((page - 1) * perPage).Limit(perPage).
		Order("id DESC").
		Find(&materials).Error
	if err != nil {
		fmt.Println(err)
		return []Material{}, 0, err
	}

	//查询数量
	err = GlobalDb.Model(&Material{}).Preload("Place").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Car").Omit("password")
	}).
		Where("place_id = ? or 0 = ?", place, place).
		Where("user_id in ? or 0 = ?", userId, car).
		Where("name LIKE ? OR model LIKE ?", regKey, regKey).
		Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return []Material{}, 0, err
	}

	return materials, count, nil
}
