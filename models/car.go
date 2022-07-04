package models

import (
	"errors"
	"fmt"
)

// 添加位置
func CarAdd(car *Car) error {

	err := GlobalDb.Create(car).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("该车号已存在")
	}

	return nil

}

//根据id获取货架
func CarGetById(id int) (Car, error) {

	car := Car{}
	err := GlobalDb.Where(&Car{ID: uint(id)}).First(&car).Error
	if err != nil {
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
		fmt.Println(err)
		return err
	}
	return nil
}

//删除车号
func CarDel(id int) error {
	car := Car{ID: uint(id)}

	//检测是否有物资在此位置，有的话则不能删除
	//m := Material{}
	//m.PlaceID = car.ID
	//var count int64 = -1
	//err := GlobalDb.Find(&m).Count(&count).Error
	//if err != nil {
	//	fmt.Println(err)
	//	return errors.New("删除失败")
	//}
	//
	//if count != 0 {
	//	return errors.New("此位置尚有物资，请先删除物资")
	//}

	err := GlobalDb.Delete(&car).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("删除失败")
	}
	return nil

}

//获取所有位置
func CarGetAll() ([]Car, error) {
	cars := []Car{}
	err := GlobalDb.Order("id DESC").Find(&cars).Error
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return nil, 0, errors.New("查询失败")
	}

	err = GlobalDb.Model(Car{}).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败")
	}

	return cars, count, nil
}
