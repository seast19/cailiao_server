package models

import (
	"errors"
	"fmt"
)

// 添加位置
func PlaceAdd(place *Place) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Create(place).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("该位置已存在")
	}

	return nil

}

//根据id获取货架
func PlaceGetById(id int) (Place, error) {
	db, err := getConn()
	if err != nil {
		return Place{}, err
	}
	defer db.Close()

	place := Place{}
	err = db.Where(&Place{ID: uint(id)}).First(&place).Error
	if err != nil {
		return Place{}, err
	}

	return place, nil
}

//修改位置
func PlaceEditByID(place Place) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Model(&Place{ID: place.ID}).Update(
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
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()

	place := Place{ID: uint(id)}
	err = db.Delete(&place).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("删除失败")
	}
	return nil

}

//获取所有位置
func PlaceAll() ([]Place, error) {
	db, err := getConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	places := []Place{}
	err = db.Order("id DESC").Find(&places).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("查询失败")
	}

	return places, nil
}

//获取分页位置
func PlaceAllGetPlaceByPage(page, perPage int) ([]Place, int, error) {
	db, err := getConn()
	if err != nil {
		return nil, 0, err
	}
	defer db.Close()

	places := []Place{}
	count := 0
	err = db.Offset((page - 1) * perPage).Limit(perPage).Order("id DESC").Find(&places).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败")
	}

	err = db.Model(Place{}).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return nil, 0, errors.New("查询失败")
	}

	return places, count, nil
}
