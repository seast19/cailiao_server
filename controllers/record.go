package controllers

import (
	"cailiao_server/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//添加出入库记录
func RecordAdd(c *gin.Context) {
	data := struct {
		Id          int    `json:"id"`          //物资id
		ChangeCount int    `json:"changeCount"` //变动数量
		Marks       string `json:"marks"`       //俗称
		Action      string `json:"action"`      //操作动作 send receive
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
	//fmt.Println(data)
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

	//fmt.Println(uid)

	record := models.Record{}
	record.UserID = uid
	record.MaterialID = uint(data.Id)
	record.Marks = data.Marks
	record.CountChange = data.ChangeCount
	record.Type = data.Action
	err = models.RecordAdd(&record)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "添加成功",
	})
}

//分页带搜索获取记录列表
func RecordGetAllByPageAndSearch(c *gin.Context) {
	data := struct {
		Page      int    `json:"page" form:"page"`
		PerPage   int    `json:"per_page" form:"per_page"`
		Type      string `json:"type" form:"type"`
		StartTime int    `json:"start_time" form:"start_time"`
		StopTime  int    `json:"stop_time" form:"stop_time"` //ms
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
	//	查询
	records, count, err := models.RecordGetAllByPage(data.Page, data.PerPage, data.Type, data.StartTime, data.StopTime)
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
			"records": records,
			"count":   count,
			"page":    data.Page,
		},
	})
}

// 删除出入库记录
func RecordDelById(c *gin.Context) {
	id := c.Param("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	err = models.RecordDelById(idNum)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "删除失败",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "ok",
	})
}
