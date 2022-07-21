package models

import (
	"errors"
	"github.com/beego/beego/v2/core/logs"
)

// 添加位置
func PlaceAdd(place *Place) error {
	err := GlobalDb.Create(place).Error
	if err != nil {
		logs.Error(err)
		return errors.New("该位置已存在")
	}
	return nil
}

//根据id获取货架
func PlaceGetById(id int) (Place, error) {
	place := Place{}
	err := GlobalDb.Where(&Place{ID: uint(id)}).First(&place).Error
	if err != nil {
		logs.Error(err)
		return Place{}, err
	}

	return place, nil
}

//修改位置
func PlaceEditByID(place Place) error {
	err := GlobalDb.Model(&Place{ID: place.ID}).Updates(
		map[string]interface{}{
			"CarID":    place.CarID,
			"Position": place.Position,
			"Remarks":  place.Remarks,
		}).Error
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

// PlaceDel 删除位置
func PlaceDel(id int) error {
	place := Place{ID: uint(id)}

	//检测是否有物资在此位置，有的话则不能删除
	var count int64 = -1
	err := GlobalDb.
		Model(&Material{}).
		Where("place_id = ?", id).
		Count(&count).Error
	if err != nil {
		logs.Error(err)
		return errors.New("删除失败#1")
	}

	if count != 0 {
		return errors.New("此位置尚有物资，请先删除物资")
	}

	err = GlobalDb.Delete(&place).Error
	if err != nil {
		logs.Error(err)
		return errors.New("删除失败#2")
	}
	return nil

}

//获取所有位置
func PlaceAll(carId uint) ([]Place, error) {

	places := []Place{}
	err := GlobalDb.
		Preload("Car").
		Where("car_id = ?", carId).
		Order("id").
		Find(&places).Error
	if err != nil {
		logs.Error(err)
		return nil, errors.New("查询失败")
	}

	return places, nil
}

//获取分页位置
func PlaceAllGetPlaceByPage(carId uint, page, perPage int) ([]Place, int64, error) {
	places := []Place{}
	var count int64 = 0
	err := GlobalDb.
		Preload("Car").
		Where("car_id = ? or 0 = ?", carId, carId).
		Order("id DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&places).Error
	if err != nil {
		logs.Error(err)
		return nil, 0, errors.New("查询失败")
	}

	err = GlobalDb.
		Model(Place{}).
		Where("car_id = ? or 0 = ?", carId, carId).
		Count(&count).Error
	if err != nil {
		logs.Error(err)
		return nil, 0, errors.New("查询失败")
	}

	return places, count, nil
}
