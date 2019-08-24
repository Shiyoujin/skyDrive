package UpDownload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"src/github.com/gin-gonic/contrib/sessions"

	"math/rand"
	"net/http"
	"skyDrive/Dao"
	"strconv"
	"time"
)

//这个重定向的接口，主要是为了封装下载的接口，而不暴露给外部
func Redirect(context *gin.Context) {

	//从数据库中查询是否有这么一个参数
	afterParams := context.Param("afterParams")
	code := context.Param("code")

	//如果存在这么一个参数，则取出 fileName，并且重定向到shareUrl那里去

	fileName, ExtractedCode, dateNum, createdAt := Dao.GetFileName(afterParams)

	//dateNum是天数
	format := "2006-01-02 15:04:05"

	createdTime, _ := time.Parse(format, createdAt)

	now := time.Now()

	//过期的天数转化为 小时 h
	date := 24 * dateNum
	dateS := strconv.Itoa(date)
	//创建的时间加上 过期时间dateNum所需要的 add
	add, _ := time.ParseDuration(dateS + "h")
	fmt.Println(dateS)
	//share 分享链接 创建的时间加上 过期时间dateNum用来判断是否过期
	newTime := createdTime.Add(add)
	fmt.Println(createdTime)
	//再把 newTime转化格式方便进行 时间过期的比较
	fresh, _ := time.Parse(format, newTime.Format(format))

	fmt.Println(fresh.Format(format))
	if now.Before(fresh) {
		if ExtractedCode == code {
			downLoadUrl := "http://localhost:8083/downLoad/" + fileName + "?=" + ExtractedCode

			//重定向到下载的接口
			context.Redirect(http.StatusMovedPermanently, downLoadUrl)
		} else {
			context.String(http.StatusOK, "提取码错误，无法获得下载链接")
		}
	} else {
		context.String(http.StatusOK, "下载链接已过期")
	}

}

/*
加密分享接口
这里也需要用户的校验
*/
func EnCrySharing(c *gin.Context) {

	//这里获取用户的id进行判断，是否可以创建一个加密分享的连接
	//sesion
	session := sessions.Default(c)
	userID := session.Get("ID")

	file := c.Param("file")
	date := c.Param("date")

	if userID != nil {
		//生成四位由 大小字母+数字的 提取码
		ExtractedCode := GetRandomString(4)
		//对随机 生成的参数进行 MD5加密,并存入数据库
		afterParams := GetRandomString(11)

		//时间的格式 用于对过期时间的校验
		format := "2006-01-02 15:04:05"
		now := time.Now()
		createdAt := now.Format(format)

		//拼接好的链接以及参数都要存在数据库里面
		//这里应该是一个重定向的链接
		shareUrl := "http://localhost:8083/sharing/" + afterParams

		//string转int
		dateNum, _ := strconv.Atoi(date)

		photo := file

		Dao.AddShare(userID.(string), photo, afterParams, ExtractedCode, dateNum, createdAt)

		c.String(http.StatusOK, "成功创建分享链接，请复制粘贴发送给好友："+shareUrl+"  提取码为："+ExtractedCode)
	} else {
		c.String(http.StatusOK, "你尚未登录")
	}

}

// 随机生成指定位数的大写字母和数字的组合
func GetRandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
