package models

import (
	"errors"
	"fmt"
)

// 添加位置
func PlaceAdd(place *Place) error {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return err
	//}
	//defer GLOBAL_DB.Close()

	err := GlobalDb.Create(place).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("该位置已存在")
	}

	return nil

}

//根据id获取货架
func PlaceGetById(id int) (Place, error) {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return Place{}, err
	//}
	//defer GLOBAL_DB.Close()

	place := Place{}
	err := GlobalDb.Where(&Place{ID: uint(id)}).First(&place).Error
	if err != nil {
		return Place{}, err
	}

	return place, nil
}

//修改位置
func PlaceEditByID(place Place) error {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return err
	//}
	//defer GLOBAL_DB.Close()

	err := GlobalDb.Model(&Place{ID: place.ID}).Updates(
		map[string]interface{}{
			"Position": place.Position,
			"Remarks":  place.Remarks,
		}).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//删除位置
func PlaceDel(id int) error {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return err
	//}
	//defer GLOBAL_DB.Close()

	place := Place{ID: uint(id)}

	//检测是否有物资在此位置，有的话则不能删除
	m := Material{}
	m.PlaceID = place.ID
	var count int64 = -1
	err := GlobalDb.Find(&m).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("删除失败")
	}

	if count != 0 {
		return errors.New("此位置尚有物资，请先删除物资")
	}

	err = GlobalDb.Delete(&place).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("删除失败")
	}
	return nil

}

//获取所有位置
func PlaceAll() ([]Place, error) {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return nil, err
	//}
	//defer GLOBAL_DB.Close()

	places := []Place{}
	err := GlobalDb.Order("id DESC").Find(&places).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("查询失败")
	}

	return places, nil
}

//获取分页位置
func PlaceAllGetPlaceByPage(page, perPage int) ([]Place, int64, error) {
	//GLOBAL_DB, err := getConn()
	//if err != nil {
	//	return nil, 0, err
	//}
	//defer GLOBAL_DB.Close()

	places := []Place{}
	var count int64 = 0
	err := GlobalDb.Offset((page - 1) * perPage).Limit(perPage).Order("id DESC").Find(&places).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败")
	}

	err = GlobalDb.Model(Place{}).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败")
	}

	return places, count, nil
}
