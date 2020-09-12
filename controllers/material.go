package controllers

import (
	"cailiao_server/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//添加货架位置
func PlaceAdd(c *gin.Context) {
	data := struct {
		Position string `json:"position"`
		Remarks  string `json:"remarks"`
	}{}

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	place := models.Place{
		Position: data.Position,
		Remarks:  data.Remarks,
	}

	err = models.PlaceAdd(&place)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg": "添加成功",
	})
}

// 获取货架位置
func PlaceGetPlaceByPage(c *gin.Context) {
	data := struct {
		Page    int `json:"page" form:"page"`
		PerPage int `json:"per_page" form:"per_page"`
	}{}

	err := c.BindQuery(&data)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"参数错误",
		})
		return
	}

	//fmt.Println(data)


	places,count, err := models.PlaceAllGetPlaceByPage(data.Page,data.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"places": places,
		"count":count,
		"page":data.Page,
	})

}

// 删除某个货架
func PlaceDelById(c *gin.Context) {
	id := c.Param("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"参数错误",
		})
		return
	}
	err = models.PlaceDel(idNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"msg": "删除成功",
	})
}
