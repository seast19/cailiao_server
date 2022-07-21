package controllers

import (
	"cailiao_server/models"
	"github.com/beego/beego/v2/core/logs"
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
		logs.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}
	//fmt.Println(data)
	//获取该用户id
	jwt := c.GetHeader("jwt")
	user, err := models.UserGetByJwt(jwt)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "参数错误",
		})
		return
	}

	//fmt.Println(uid)

	record := models.Record{}
	record.UserID = user.ID
	record.MaterialID = uint(data.Id)
	record.Marks = data.Marks
	//record.CountChange = data.ChangeCount
	record.Type = data.Action

	//出库
	if data.Action == "send" {
		record.SendCount = data.ChangeCount
		err = models.RecordAddSend(&record)
		if err != nil {
			logs.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 4000,
				"msg":  err.Error(),
			})
			return
		}
	} else if data.Action == "receive" {
		record.ReceiveCount = data.ChangeCount
		err = models.RecordAddReceive(&record)
		if err != nil {
			logs.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 4000,
				"msg":  err.Error(),
			})
			return
		}
	}
	//入库

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "添加成功",
	})
}

//聚合获取材料领用记录
func RecordGetPolyWithCardByPage(c *gin.Context) {
	data := struct {
		Page    int `json:"page" form:"page"`
		PerPage int `json:"per_page" form:"per_page"`
		//Type       string `json:"type" form:"type"`
		StartTime int `json:"start_time" form:"start_time"`
		StopTime  int `json:"stop_time" form:"stop_time"` //ms
		//MaterialID uint   `json:"id" form:"id"`               //待id则只查看该材料id的记录
		CarID uint `json:"car_id" form:"car_id"` //含carid则查看该车id
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

	recordPoly, count, err := models.RecordGetPolyWithCardByPage(data.Page, data.PerPage, data.StartTime, data.StopTime, data.CarID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 5000,
			"msg":  err.Error(),
		})
		return
	}

	//logs.Info(recordPoly)

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "ok",
		"data": gin.H{
			"records": recordPoly,
			"count":   count,
			"page":    data.Page,
		},
	})

}

//分页带搜索获取记录列表
func RecordGetAllByPageAndSearch(c *gin.Context) {
	data := struct {
		Page       int    `json:"page" form:"page"`
		PerPage    int    `json:"per_page" form:"per_page"`
		Type       string `json:"type" form:"type"`
		StartTime  int    `json:"start_time" form:"start_time"`
		StopTime   int    `json:"stop_time" form:"stop_time"` //ms
		MaterialID uint   `json:"mid" form:"mid"`             //待id则只查看该材料id的记录
		CarID      uint   `json:"car_id" form:"car_id"`       //含carid则查看该车id
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

	//	查询
	//fmt.Println(data)
	records := []models.Record{}
	var count int64
	if data.MaterialID == 0 {
		records, count, err = models.RecordGetAllWithCarByPage(data.Page, data.PerPage, data.Type, data.StartTime, data.StopTime, data.CarID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 5000,
				"msg":  err.Error(),
			})
			return
		}
	} else {
		records, count, err = models.RecordGetAllWithMIdByPage(data.Page, data.PerPage, data.Type, data.StartTime, data.StopTime, data.MaterialID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 5000,
				"msg":  err.Error(),
			})
			return
		}
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
		logs.Error(err)
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
