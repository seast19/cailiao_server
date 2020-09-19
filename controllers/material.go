package controllers

import (
	"cailiao_server/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//分页获取所有材料
func MaterialGetAllByPage(c *gin.Context) {
	data := struct {
		Page    int `json:"page" form:"page"`
		PerPage int `json:"per_page" form:"per_page"`
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
}

//添加材料
func MaterialAdd(c *gin.Context) {
	data := struct {
		Name     string `json:"name"`      //名称
		Model    string `json:"model"`     //型号
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
