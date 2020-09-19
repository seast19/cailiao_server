package models

import "fmt"

//添加材料
func MaterialAdd(m *Material) error {
	db, err := getConn()
	if err != nil {
		return err
	}
	defer db.Close()
	//db.Model(m).Related(&profile)
	err=db.Create(m).Error
	if err!=nil{
		fmt.Println(err)
		return err
	}

	return nil
}
