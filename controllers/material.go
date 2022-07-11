package controllers

import (
	"cailiao_server/models"
	"cailiao_server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"time"
)

//添加材料
func MaterialAdd(c *gin.Context) {
	data := struct {
		Name     string `json:"name"`  //名称
		Model    string `json:"model"` //型号
		Unit     string `json:"unit"`
		NickName string `json:"nick_name"` //俗称

		CarID    uint `json:"car_id"`
		PlaceID  uint `json:"place_id"`
		Floor    int  `json:"floor"`    //层
		Location int  `json:"location"` //位

		Count        int `json:"count"`         //数量
		PrepareCount int `json:"prepare_count"` //常备数量
		WarnCount    int `json:"warn_count"`    //警报数量

		Marks string `json:"marks"` //备注
	}{}

	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	//获取该用户id
	jwt := c.GetHeader("jwt")

	user, err := models.UserGetByJwt(jwt)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}
	//	构建待添加的数据
	material := models.Material{
		Name:     data.Name,
		Model:    data.Model,
		Unit:     data.Unit,
		NickName: data.NickName,
		//Place:        models.Place{ID: data.PlaceID},
		PlaceID:      data.PlaceID,
		Floor:        data.Floor,
		Location:     data.Location,
		Count:        data.Count,
		PrepareCount: data.PrepareCount,
		WarnCount:    data.WarnCount,
		Marks:        data.Marks,
		//User:         models.User{ID: uid},
		UserID: user.ID,
		CarID:  user.Car.ID,

		CreateAt: time.Now().UnixNano(),
	}

	err = models.MaterialAdd(&material)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 2000,
		"msg":  "添加成功",
	})
}

//通过xlsx批量添加材料信息
func MaterialAddAll(c *gin.Context) {
	//	获取xlsx文件
	// 单文件
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusCreated, gin.H{
			"code": 4000,
			"msg":  "上传失败，获取文件失败",
		})
		return
	}
	fmt.Println(file.Filename)
	dst := "./uploads/" + file.Filename
	// 上传文件至指定的完整文件路径
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"code": 4000,
			"msg":  "上传失败，保存文件失败",
		})
		return
	}

	//解析xlsx文件
	f, err := excelize.OpenFile(dst)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"code": 4000,
			"msg":  "上传失败，文件格式错误",
		})
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	//获取该用户id
	jwt := c.GetHeader("jwt")
	user, err := models.UserGetByJwt(jwt)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	//构建材料数据
	ms := []models.Material{}
	for rowIndex, _ := range rows {
		if rowIndex <= 2 {
			continue
		}
		//构建axis
		axis1 := fmt.Sprintf("B%d", rowIndex+1)
		name, _ := f.GetCellValue("Sheet1", axis1)

		axis2 := fmt.Sprintf("C%d", rowIndex+1)
		model, _ := f.GetCellValue("Sheet1", axis2)

		axis3 := fmt.Sprintf("D%d", rowIndex+1)
		nickname, _ := f.GetCellValue("Sheet1", axis3)

		axis4 := fmt.Sprintf("E%d", rowIndex+1)
		unit, _ := f.GetCellValue("Sheet1", axis4)

		axis5 := fmt.Sprintf("F%d", rowIndex+1)
		placeID, _ := f.GetCellValue("Sheet1", axis5)

		axis6 := fmt.Sprintf("G%d", rowIndex+1)
		floor, _ := f.GetCellValue("Sheet1", axis6)

		axis7 := fmt.Sprintf("H%d", rowIndex+1)
		location, _ := f.GetCellValue("Sheet1", axis7)

		axis8 := fmt.Sprintf("I%d", rowIndex+1)
		count, _ := f.GetCellValue("Sheet1", axis8)

		axis9 := fmt.Sprintf("J%d", rowIndex+1)
		prepareCount, _ := f.GetCellValue("Sheet1", axis9)

		axis10 := fmt.Sprintf("K%d", rowIndex+1)
		warnCount, _ := f.GetCellValue("Sheet1", axis10)

		axis11 := fmt.Sprintf("L%d", rowIndex+1)
		marks, _ := f.GetCellValue("Sheet1", axis11)

		axis12 := fmt.Sprintf("M%d", rowIndex+1)
		carID, _ := f.GetCellValue("Sheet1", axis12)

		ms = append(ms, models.Material{
			Name:         name,
			Model:        model,
			NickName:     nickname,
			Unit:         unit,
			PlaceID:      uint(utils.StringToInt(placeID)),
			CarID:        uint(utils.StringToInt(carID)),
			Floor:        utils.StringToInt(floor),
			Location:     utils.StringToInt(location),
			Count:        utils.StringToInt(count),
			PrepareCount: utils.StringToInt(prepareCount),
			WarnCount:    utils.StringToInt(warnCount),
			Marks:        marks,
			UserID:       user.ID,
			CreateAt:     time.Now().UnixNano(),
		})
	}

	//批量导入
	//fmt.Println(ms)
	err = models.MaterialAddAll(ms)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "批量添加失败",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 2000,
		"msg":  "上传成功",
	})
}

// MaterialSearch 搜索材料
func MaterialSearch(c *gin.Context) {
	data := struct {
		Page    int    `json:"page" form:"page"`         //页码
		PerPage int    `json:"per_page" form:"per_page"` //每页数量
		Key     string `json:"key" form:"key"`           // 搜索关键字
		Car     int    `json:"car" form:"car"`           // 选择车号
		Place   int    `json:"place" form:"place"`       // 选择位置

	}{}

	err := c.BindQuery(&data)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	//fmt.Println(data)
	materials, count, err := models.MaterialSearchByKey(data.Key, data.Car, data.Place, data.Page, data.PerPage)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  2000,
		"msg":   "ok",
		"data":  materials,
		"count": count,
		"page":  data.Page,
	})
}

// MaterialWarn  获取达到常备数量以下的材料
func MaterialWarn(c *gin.Context) {
	data := struct {
		Page    int `json:"page" form:"page"`         //页码
		PerPage int `json:"per_page" form:"per_page"` //每页数量
		//Key     string `json:"key" form:"key"`           // 搜索关键字
		Car   int `json:"car" form:"car"`     // 选择车号
		Place int `json:"place" form:"place"` // 选择位置

	}{}

	err := c.BindQuery(&data)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	//fmt.Println(data)
	materials, count, err := models.MaterialWarnByCar(data.Car, data.Place, data.Page, data.PerPage)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  2000,
		"msg":   "ok",
		"data":  materials,
		"count": count,
		"page":  data.Page,
	})
}

// 删除单个材料
func MaterialDelOneByID(c *gin.Context) {
	id := c.Param("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	err = models.MaterialDel(idNum)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "删除成功",
	})

}

//获取单个材料
func MaterialGetOneById(c *gin.Context) {
	id := c.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	material, err := models.MaterialGetById(idNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"data": material,
		"msg":  "ok",
	})
}

//更新某个材料
func MaterialUpdateOneById(c *gin.Context) {
	idStr := c.Param("id")

	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误1",
		})
		return
	}

	data := struct {
		Name     string `json:"name"`  //名称
		Model    string `json:"model"` //型号
		Unit     string `json:"unit"`
		NickName string `json:"nick_name"` //俗称

		CarID    uint `json:"car_id"`
		PlaceID  uint `json:"place_id"`
		Floor    int  `json:"floor"`    //层
		Location int  `json:"location"` //位

		Count        int `json:"count"`         //数量
		PrepareCount int `json:"prepare_count"` //常备数量
		WarnCount    int `json:"warn_count"`    //警报数量

		Marks string `json:"marks"` //备注
	}{}

	err = c.BindJSON(&data)

	//fmt.Println(data)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误2",
		})
		return
	}

	jwt := c.GetHeader("jwt")
	user, err := models.UserGetByJwt(jwt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误3",
		})
		return
	}

	material := models.Material{
		ID: uint(idNum),

		Name:     data.Name,
		Model:    data.Model,
		Unit:     data.Unit,
		NickName: data.NickName,

		CarID:    data.CarID,
		PlaceID:  data.PlaceID,
		Floor:    data.Floor,
		Location: data.Location,

		Count:        data.Count,
		PrepareCount: data.PrepareCount,
		WarnCount:    data.WarnCount,
		Marks:        data.Marks,

		UserID: user.ID,
	}

	err = models.MaterialEditByID(material)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,

		"msg": "ok",
	})

}

// MaterialDownload 下载材料清单
func MaterialDownload(c *gin.Context) {
	data := struct {
		Page    int    `json:"page" form:"page"`         //页码
		PerPage int    `json:"per_page" form:"per_page"` //每页数量
		Key     string `json:"key" form:"key"`           // 搜索关键字
		Car     int    `json:"car" form:"car"`           // 选择车号
		Place   int    `json:"place" form:"place"`       // 选择位置

	}{}

	err := c.BindQuery(&data)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	materials, err := models.MaterialDownloadByKey(data.Key, data.Car, data.Place, data.Page, data.PerPage)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  err.Error(),
		})
		return
	}

	//	构建xlsx文件
	f, err := excelize.OpenFile("assets/下载材料清单模板.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	//fmt.Println(materials)
	//循环填入数据
	for i, material := range materials {
		//	按行赋值
		axis := fmt.Sprintf("A%d", i+4)
		//fmt.Println(axis)
		err := f.SetSheetRow("Sheet1", axis, &[]interface{}{
			material.ID,
			material.Name,
			material.Model,
			material.NickName,
			material.Unit,
			material.Place.Position,
			material.Floor,
			material.Location,
			material.Count,
			material.PrepareCount,
			material.WarnCount,
			material.User.RealName,
			material.Car.Car,
			material.Marks,
		})
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	filename := fmt.Sprintf("statics/%s材料清单-%d.xlsx", materials[0].Car.Car, time.Now().Unix())
	err = f.SaveAs(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"url":  filename,
		"msg":  "下载成功",
	})

}

// MaterialDownloadWarn 下载材料清单
func MaterialDownloadWarn(c *gin.Context) {
	data := struct {
		Page    int `json:"page" form:"page"`         //页码
		PerPage int `json:"per_page" form:"per_page"` //每页数量
		//Key     string `json:"key" form:"key"`           // 搜索关键字
		Car   int `json:"car" form:"car"`     // 选择车号
		Place int `json:"place" form:"place"` // 选择位置

	}{}

	err := c.BindQuery(&data)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	materials, err := models.MaterialDownloadWarnByCar(data.Car, data.Place, data.Page, data.PerPage)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  err.Error(),
		})
		return
	}

	//	构建xlsx文件
	f, err := excelize.OpenFile("assets/下载材料清单模板.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	//fmt.Println(materials)
	//循环填入数据
	for i, material := range materials {
		//	按行赋值
		axis := fmt.Sprintf("A%d", i+4)
		//fmt.Println(axis)
		err := f.SetSheetRow("Sheet1", axis, &[]interface{}{
			material.ID,
			material.Name,
			material.Model,
			material.NickName,
			material.Unit,
			material.Place.Position,
			material.Floor,
			material.Location,
			material.Count,
			material.PrepareCount,
			material.WarnCount,
			material.User.RealName,
			material.Car.Car,
			material.Marks,
		})
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	filename := fmt.Sprintf("statics/%s备料清单-%d.xlsx", materials[0].Car.Car, time.Now().Unix())
	err = f.SaveAs(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"url":  filename,
		"msg":  "下载成功",
	})

}
