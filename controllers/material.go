package controllers

import (
	"cailiao_server/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//分页获取所有材料
func MaterialGetAllByPage(c *gin.Context) {
	data := struct {
		Page    int `json:"page" form:"page"`
		PerPage int `json:"per_page" form:"per_page"`
		PlaceID int `json:"place_id" form:"place_id"`
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

	//	查询
	materials, count, err := models.MaterialGetByPage(data.Page, data.PerPage, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 5000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "ok",
		"data": gin.H{
			"materials": materials,
			"count":     count,
			"page":      data.Page,
		},
	})
}

//添加材料
func MaterialAdd(c *gin.Context) {
	data := struct {
		Name     string `json:"name"`  //名称
		Model    string `json:"model"` //型号
		Unit     string `json:"unit"`
		NickName string `json:"nick_name"` //俗称

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
	uid, err := models.UserGetIDByJwt(jwt)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}
	//fmt.Println(data)
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
		UserID: uid,

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

//搜索材料
func MaterialSearch(c *gin.Context) {
	data := struct {
		Page    int    `json:"page" form:"page"`
		PerPage int    `json:"per_page" form:"per_page"`
		Key     string `json:"key" form:"key"`
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
	materials, count, err := models.MaterialSearchByKey(data.Key, data.Page, data.PerPage)
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
	uid, err := models.UserGetIDByJwt(jwt)
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
		Unit: data.Unit,
		NickName: data.NickName,

		PlaceID:  data.PlaceID,
		Floor:    data.Floor,
		Location: data.Location,

		Count:        data.Count,
		PrepareCount: data.PrepareCount,
		WarnCount:    data.WarnCount,
		Marks:        data.Marks,

		UserID: uid,
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
