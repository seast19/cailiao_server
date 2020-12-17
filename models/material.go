package models

import (
	"cailiao_server/utils"
	"errors"
	"fmt"
)

//删除材料
func MaterialDel(id int) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()

	material := Material{ID: uint(id)}

	tx := db.Begin()

	//出料记录存在则不能删除该物资
	rc := 0
	err = tx.Model(&Record{}).Where(&Record{MaterialID: material.ID}).Count(&rc).Error
	if err != nil {
		utils.Mlogger.Errorf("#1 | %s | %s", "models/MaterialDel", err.Error())
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
		utils.Mlogger.Errorf("#2 | %s | %s", "models/MaterialDel", err.Error())
		return errors.New("#2 删除失败")
	}

	tx.Commit()
	return nil

}

//添加材料
func MaterialAdd(m *Material) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()
	//db.Model(m).Related(&profile)
	err = db.Create(m).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

//分页获取材料
func MaterialGetByPage(page, perPage, placeID int) ([]Material, int, error) {
	db, err := getConn()
	if err != nil {
		return []Material{}, 0, err
	}
	defer db.Close()

	//placeID=36

	material := []Material{}
	count := 0
	err = db.Preload("Place").Where(&Material{PlaceID: uint(placeID)}).Offset((page - 1) * perPage).Limit(perPage).Order("id DESC").Find(&material).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败1")
	}
	//fmt.Println(material[0])
	//m := material[0]
	//m := Material{}

	//fmt.Println(m)
	//	db.Last(&m)
	//	fmt.Println(m)
	//	err=db.Model(&m).Related(&m.Place,"PlaceID").Error
	//	//res:=db.First(&material[0].Place).Error
	//	fmt.Println(m)
	//	if err != nil {
	//		fmt.Println(err)
	//		return nil, 0, errors.New("查询失败2")
	//	}

	err = db.Model(Material{}).Where(&Material{PlaceID: uint(placeID)}).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败3")
	}

	return material, count, nil
}

//根据id获取材料
func MaterialGetById(id int) (Material, error) {
	db, err := getConn()
	if err != nil {
		return Material{}, err
	}
	defer db.Close()

	material := Material{}
	err = db.Where(&Material{ID: uint(id)}).First(&material).Error
	if err != nil {
		return Material{}, err
	}

	return material, nil
}

//更新材料
func MaterialEditByID(material Material) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println(material)
	err = db.Model(&Material{ID: material.ID}).Update(
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

//搜索材料
func MaterialSearchByKey(key string, page, perPage int) ([]Material, int, error) {
	db, err := getConn()
	if err != nil {
		return []Material{}, 0, err
	}
	defer db.Close()

	materials := []Material{}

	regKey := "%" + key + "%"
	//fmt.Println(regKey)
	err = db.Preload("Place").Where("name LIKE ? OR model LIKE ?", regKey, regKey).Offset((page - 1) * perPage).Limit(perPage).Order("id DESC").Find(&materials).Error
	if err != nil {
		fmt.Println(err)
		return []Material{}, 0, err
	}

	count := 0
	err = db.Model(Material{}).Where("name LIKE ? OR model LIKE ?", regKey, regKey).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return []Material{}, 0, err
	}

	//fmt.Println(materials)
	//fmt.Println(count)

	return materials, count, nil

}
