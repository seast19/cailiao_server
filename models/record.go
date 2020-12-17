package models

import (
	"errors"
	"fmt"
	"time"
)

//添加记录
func RecordAdd(r *Record) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()

	// 构建完整数据
	t := time.Now().UnixNano() / 1e6
	r.CreateAt = t
	r.UpdateAt = t

	tx := db.Begin()

	//分发料和领料操作
	if r.Type == "send" {
		m := Material{}
		m.ID = r.MaterialID
		//查找该物资
		err = tx.First(&m).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return errors.New("无此物资id")
		}
		//fmt.Println(m)
		//库存为0
		if m.Count <= 0 {
			tx.Rollback()
			return errors.New("此物资库存数为0")
		}

		//库存数比发料数少
		if m.Count < r.CountChange {
			tx.Rollback()
			return errors.New("此物资库存数不足以发料")
		}

		//创建record记录
		r.BeforeCount = m.Count
		r.AfterCount = m.Count - r.CountChange
		err = tx.Create(r).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return errors.New("添加记录失败")
		}

		//修改material库存数
		err = tx.Model(&m).Select("count").Updates(map[string]interface{}{"count": m.Count - r.CountChange}).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}

		tx.Commit()
		return nil

	} else if r.Type == "receive" {
		m := Material{}
		m.ID = r.MaterialID
		//查找该物资
		err = tx.First(&m).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return errors.New("无此物资id")
		}
		//fmt.Println(m)
		//库存为0
		//if m.Count <= 0 {
		//	tx.Rollback()
		//	return errors.New("此物资库存数为0")
		//}

		//库存数比发料数少
		//if m.Count < r.CountChange {
		//	tx.Rollback()
		//	return errors.New("此物资库存数不足以发料")
		//}

		//创建record记录
		r.BeforeCount = m.Count
		r.AfterCount = m.Count + r.CountChange
		err = tx.Create(r).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return errors.New("添加记录失败")
		}

		//修改material库存数
		err = tx.Model(&m).Select("count").Updates(map[string]interface{}{"count": m.Count + r.CountChange}).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}

		tx.Commit()
		return nil
	}
	tx.Rollback()
	return errors.New("操作类型无效")

}

//分页&搜索 获取记录列表
func RecordGetAllByPageAndSearch(key string, page, perPage int) ([]Record, int, error) {
	db, err := getConn()
	if err != nil {
		return []Record{}, 0, err
	}
	defer db.Close()

	//构建参数
	records := []Record{}
	regKey := "%" + key + "%"
	err = db.Preload("Place").Preload("User.RealName").Where("name LIKE ? OR model LIKE ?", regKey, regKey).Offset((page - 1) * perPage).Limit(perPage).Order("id DESC").Find(&records).Error
	if err != nil {
		fmt.Println(err)
		return []Record{}, 0, err
	}

	count := 0
	err = db.Model(Record{}).Where("name LIKE ? OR model LIKE ?", regKey, regKey).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return []Record{}, 0, err
	}

	return records, count, nil

}

//分页获取记录列表
func RecordGetAllByPage(page, perPage int, t string, startT, stopT int) ([]Record, int, error) {
	db, err := getConn()
	if err != nil {
		return []Record{}, 0, err
	}
	defer db.Close()

	//fmt.Println(t)
	//构建参数
	records := []Record{}
	count := 0
	typeQuery := ""
	timeQuery := ""

	//构建操作参数
	if t != "" {
		typeQuery = fmt.Sprintf("Type = '%s'", t)
	}

	if startT != 0 {
		timeQuery = fmt.Sprintf("create_at >= %d AND create_at <= %d", startT, stopT)
	}

	err = db.Preload("Material").Preload("Material.Place").Preload("User").Where(typeQuery).Where(timeQuery).Offset((page - 1) * perPage).Limit(perPage).Order("id DESC").Find(&records).Error
	if err != nil {
		fmt.Println(err)
		return []Record{}, 0, err
	}

	err = db.Model(Record{}).Where(typeQuery).Where(timeQuery).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return []Record{}, 0, err
	}

	//过滤掉密码
	records2 := []Record{}
	for _, item := range records {
		item.User.Password = ""
		records2 = append(records2, item)
	}

	return records2, count, nil

}

// 根据id删除出入库记录
func RecordDelById(id int) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()

	r := Record{ID: uint(id)}

	tx := db.Begin()

	//找到该条记录
	err = tx.First(&r).Error
	if err != nil {
		fmt.Println("#1",err)
		tx.Rollback()
		return err
	}

	// 查询该物资
	m := Material{}
	m.ID = r.MaterialID
	err = tx.First(&m).Error
	if err != nil {
		fmt.Println("#2",err)
		tx.Rollback()
		return err
	}

	//按操作类型 更新材料表
	switch r.Type {
	case "send":
		//修改material库存数
		err = tx.Model(&m).Select("count").Updates(map[string]interface{}{"count": m.Count + r.CountChange}).Error
		if err != nil {
			fmt.Println("#3",err)
			tx.Rollback()
			return err
		}
	case "receive":
		//修改material库存数
		err = tx.Model(&m).Select("count").Updates(map[string]interface{}{"count": m.Count - r.CountChange}).Error
		if err != nil {
			fmt.Println("#4",err)
			tx.Rollback()
			return err
		}
	default:
		tx.Rollback()
		return errors.New("操作类型错误")

	}

	//删除该记录
	err = tx.Delete(&r).Error
	if err != nil {
		fmt.Println("#5",err)
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}
