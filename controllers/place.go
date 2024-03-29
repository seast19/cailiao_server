package controllers

import (
	"cailiao_server/models"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//添加货架位置
func PlaceAdd(c *gin.Context) {
	data := struct {
		Position string `json:"position"`
		Remarks  string `json:"remarks"`
		CarId    uint   `json:"car_id"`
	}{}

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	place := models.Place{
		Position: data.Position,
		Remarks:  data.Remarks,
		CarID:    data.CarId,
	}

	err = models.PlaceAdd(&place)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "添加成功",
	})
}

// 分页获取货架位置
func PlaceGetPlaceByPage(c *gin.Context) {
	data := struct {
		Page    int  `json:"page" form:"page"`
		PerPage int  `json:"per_page" form:"per_page"`
		CarId   uint `json:"car_id" form:"car_id"`
	}{}
	err := c.BindQuery(&data)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	//fmt.Println(data)
	//fmt.Println(data)

	places, count, err := models.PlaceAllGetPlaceByPage(data.CarId, data.Page, data.PerPage)
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

			"places": places,
			"count":  count,
			"page":   data.Page,
		},
	})

}

//获取所有货架
func PlaceGetAll(c *gin.Context) {
	data := struct {
		CarId uint `json:"car_id" form:"car_id"`
	}{}
	err := c.BindQuery(&data)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	places, err := models.PlaceAll(data.CarId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"data": places,
		"msg":  "ok",
	})
}

// 删除某个货架
func PlaceDelById(c *gin.Context) {
	id := c.Param("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}
	err = models.PlaceDel(idNum)
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

//根据id获取货架内容
func PlaceGetOneById(c *gin.Context) {
	id := c.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	place, err := models.PlaceGetById(idNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"data": place,
		"msg":  "ok",
	})

}

//更新某个货架
func PlaceUpdateById(c *gin.Context) {
	idStr := c.Param("id")

	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	data := struct {
		Position string `json:"position"`
		Remarks  string `json:"remarks"`
		CarId    uint   `json:"car_id"`
	}{}

	err = c.BindJSON(&data)

	//fmt.Println(data)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	place := models.Place{
		ID:       uint(idNum),
		Position: data.Position,
		CarID:    data.CarId,
		Remarks:  data.Remarks,
	}

	err = models.PlaceEditByID(place)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "修改成功",
	})

}
