package models

import (
	"errors"
	"github.com/beego/beego/v2/core/logs"
)

// 添加位置
func CarAdd(car *Car) error {

	err := GlobalDb.Create(car).Error
	if err != nil {
		logs.Error(err)
		return errors.New("该车号已存在")
	}

	return nil

}

//根据id获取货架
func CarGetById(id int) (Car, error) {

	car := Car{}
	err := GlobalDb.Where(&Car{ID: uint(id)}).First(&car).Error
	if err != nil {
		logs.Error(err)
		return Car{}, err
	}

	return car, nil
}

//修改车号
func CarEditByID(car Car) error {

	err := GlobalDb.Model(&Car{ID: car.ID}).Updates(
		map[string]interface{}{
			"Car":     car.Car,
			"Remarks": car.Remarks,
		}).Error
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

//删除车号
func CarDel(id int) error {
	car := Car{ID: uint(id)}

	err := GlobalDb.Delete(&car).Error
	if err != nil {
		logs.Error(err)
		return errors.New("删除失败")
	}
	return nil

}

//获取所有位置
func CarGetAll() ([]Car, error) {
	cars := []Car{}
	err := GlobalDb.
		Order("car").
		Find(&cars).Error
	if err != nil {
		logs.Error(err)
		return nil, errors.New("查询失败")
	}

	return cars, nil
}

//获取分页位置
func CarGetAllPlaceByPage(page, perPage int) ([]Car, int64, error) {

	cars := []Car{}
	var count int64 = 0
	err := GlobalDb.Offset((page - 1) * perPage).Limit(perPage).Order("id DESC").Find(&cars).Error
	if err != nil {
		logs.Error(err)
		return nil, 0, errors.New("查询失败")
	}

	err = GlobalDb.Model(Car{}).Count(&count).Error
	if err != nil {
		logs.Error(err)
		return nil, 0, errors.New("查询失败")
	}

	return cars, count, nil
}

//随机获取一个车号
func CarRandomGetOne() (Car, error) {
	car := Car{}
	if err := GlobalDb.First(&car).Error; err != nil {
		logs.Error(err)
		return Car{}, errors.New("获取随机车号失败")
	}
	return car, nil
}
