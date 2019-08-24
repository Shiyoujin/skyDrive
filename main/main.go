package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"skyDrive/Dao"
	"skyDrive/UpDownload"
	sessions2 "src/github.com/gin-gonic/contrib/sessions"
	"strconv"
	"strings"
)

type User struct {
	UserID int
}

func main() {

	router := gin.Default()

	//这是设置session
	store := sessions2.NewCookieStore([]byte("secret"))
	router.Use(sessions2.Sessions("mysession", store))

	//注册接口
	router.GET("/register/:userID/:password", func(context *gin.Context) {

		userID := context.Param("userID")
		password := context.Param("password")

		int, _ := strconv.Atoi(userID)
		Exist := Dao.IfExistUserID(int)

		//为0代表user表没有这个userID即可以注册
		if Exist == 0 {

			//创建一个和userID一样的文件夹,当作用户的网盘空间
			os.Mkdir(userID, os.ModePerm)
			user, _ := strconv.Atoi(userID)

			//注册插入数据库
			Dao.Register(user, password)
		} else {

			context.String(http.StatusOK, "当前userID已被注册，请重新注册")
		}

	})

	//登录接口
	router.GET("/login/:userID/:password", func(context *gin.Context) {
		session := sessions2.Default(context)

		userID := context.Param("userID")
		password := context.Param("password")

		user, _ := strconv.Atoi(userID)
		num := Dao.Login(user, password)

		//查询到一条
		if num != 0 {
			session.Set("ID", userID)
			session.Save()
			//router.Static("/"+userID, "./"+userID)

			//创建一个可视化的用户网盘文件目录
			router.StaticFS("/"+userID, http.Dir(userID))
			context.String(http.StatusOK, "登录成功session中包含了{ID:userID}")

		} else {
			context.String(http.StatusOK, "你输入的密码有误")
		}
	})

	//这是个重定向的接口，主要是为了封装下载的接口，而不暴露给外部
	//用于加密的分享链接
	router.GET("/sharing/:afterParams/:code", UpDownload.Redirect)

	/*
	   加密分享接口
	   这里也需要用户的校验
	*/
	router.GET("/share/:file/:date", UpDownload.EnCrySharing)

	//上传单个文件接口
	router.POST("/upLoad", func(context *gin.Context) {

		session := sessions2.Default(context)
		userID := session.Get("ID")

		fmt.Println(userID)
		if userID != nil {
			router.MaxMultipartMemory = 100 << 20 // 设置最大上传大小为100M
			header, _ := context.FormFile("file")

			//上传到指定的用户目录
			path := "./" + userID.(string) + "/" + header.Filename // 上传存储到的地址
			fmt.Println(path)
			//保存到本地
			err := context.SaveUploadedFile(header, path)
			if err != nil {
				fmt.Println(err.Error())
			}

			//刚上传的文件默认是0，即公开，后面有接口可以修改文件权限
			Dao.AddFile(userID.(string), header.Filename, 0)

			log.Println(header.Filename, " ", header.Size, " ", header.Header)

			context.JSON(200, gin.H{
				"fileName": header.Filename,
			})
		} else {
			context.String(http.StatusOK, "你尚未登录不能上传文件")
		}

	})

	//上传多个文件接口
	router.POST("/upLoads", func(context *gin.Context) {

		session := sessions2.Default(context)
		userID := session.Get("ID")

		fmt.Println(userID)
		if userID.(string) != "" {
			form, _ := context.MultipartForm()

			files := form.File["file[]"]

			for _, file := range files {
				log.Println(file.Filename)
				fmt.Println(file.Filename)
				//保存到本地
				context.SaveUploadedFile(file, "./"+userID.(string)+"/"+file.Filename) // 文件夹需要先创建
			}
		} else {
			context.String(http.StatusOK, "你尚未登录，不能上传多个文件")
		}
	})

	//下载接口
	router.GET("/downLoad/:fileName", func(c *gin.Context) {

		//这是代表别人的ownerID
		ownerID := c.DefaultQuery("ownerID", "不包含ownerID")
		session := sessions2.Default(c)
		userID := session.Get("ID")

		fileName := c.Param("fileName")
		code := c.DefaultQuery("code", "无")

		fmt.Println(code)
		//校验登录状态
		if userID != nil {
			//如果带有code，代表这是一个加密的分享链接下载地址
			if code != "无" {

				result, ownerID := Dao.CheckCode(fileName, code)

				ownerIDS := strconv.Itoa(ownerID)
				//提取码校验成功
				//一次性上传的前提是要有验证码
				if result != 0 {

					//生成随机的6位数
					Ran := UpDownload.GetRandomString(6)
					router.StaticFile("/static.ico"+Ran, "./"+ownerIDS+"/"+fileName)

					response, err := http.Get("http://localhost:8083/static.ico" + Ran)
					if err != nil || response.StatusCode != http.StatusOK {
						c.Status(http.StatusServiceUnavailable)
						return
					}

					reader := response.Body
					contentLength := response.ContentLength
					contentType := response.Header.Get("Content-Type")

					extraHeaders := map[string]string{
						"Content-Disposition": `attachment; filename="` + fileName + "\"",
					}

					c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

				} else {
					c.String(http.StatusOK, "你的提取码有误，请重新尝试")
				}

				//这是下载自己的文件
			} else if ownerID == "不包含ownerID" {
				//这个代表从自己的网盘下载到任意位置
				//生成随机的6位数
				Ran := UpDownload.GetRandomString(6)
				fmt.Println(userID)
				router.StaticFile("/static.ico"+Ran, "./"+userID.(string)+"/"+fileName)
				fmt.Println(fileName)

				response, err := http.Get("http://localhost:8083/static.ico" + Ran)
				if err != nil || response.StatusCode != http.StatusOK {
					c.Status(http.StatusServiceUnavailable)
					return
				}
				reader := response.Body
				contentLength := response.ContentLength
				contentType := response.Header.Get("Content-Type")

				extraHeaders := map[string]string{
					"Content-Disposition": `attachment; filename="` + fileName + "\"",
				}
				c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

				//这是下载别人文件需要通过ownerID查询权限，接下来才能进行是否下载
			} else {

				ownerIDInt, _ := strconv.Atoi(ownerID)

				//0代表公开,1代表私密
				result := Dao.BoolPermit(ownerIDInt, fileName)

				if result != 1 {
					//这个代表从自己的网盘下载到任意位置
					//生成随机的6位数
					Ran := UpDownload.GetRandomString(6)
					router.StaticFile("/static.ico"+Ran, "./"+userID.(string)+"/"+fileName)

					response, err := http.Get("http://localhost:8083/static.ico" + Ran)
					if err != nil || response.StatusCode != http.StatusOK {
						c.Status(http.StatusServiceUnavailable)
						return
					}
					reader := response.Body
					contentLength := response.ContentLength
					contentType := response.Header.Get("Content-Type")

					extraHeaders := map[string]string{
						"Content-Disposition": `attachment; filename="` + fileName + "\"",
					}
					c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
				} else {
					c.String(http.StatusOK, "你没有权限下载该文件")
				}
			}
		}
	})

	//一次性快传接口
	router.POST("/oneupload", func(context *gin.Context) {

		session := sessions2.Default(context)
		userID := session.Get("ID")

		if userID != nil {
			router.MaxMultipartMemory = 100 << 20 // 设置最大上传大小为100M
			header, _ := context.FormFile("file")

			//上传到指定的用户目录
			path := "./" + userID.(string) + "/" + header.Filename // 上传存储到的地址

			fileName := header.Filename
			//在one表添加记录
			Dao.OneUp(userID.(string), header.Filename)

			//保存到本地
			err := context.SaveUploadedFile(header, path)
			if err != nil {
				fmt.Println(err.Error())
			}

			log.Println(header.Filename, " ", header.Size, " ", header.Header)

			fileName = strings.TrimSuffix(fileName, "/")
			context.String(http.StatusOK, "一次性快传成功，这是下载链接：http://localhost:8083/oneDown/"+userID.(string)+"/"+fileName)
		} else {
			context.String(http.StatusOK, "你尚未登录不能上传文件")
		}
	})

	//一次性快传的下载接口
	router.GET("/oneDown/:ownerID/:fileName", func(context *gin.Context) {

		ownerID := context.Param("ownerID")

		fileName := context.Param("fileName")

		ownerIDInt, _ := strconv.Atoi(ownerID)
		//返回查询到的id
		oneStatus := Dao.CheckOne(ownerIDInt, fileName)

		fmt.Println(oneStatus)

		if oneStatus != 0 {

			//下载后删除文件
			os.Remove(",/" + ownerID + "/" + fileName)

			//这个代表从别人发的链接的网盘上下载到任意位置
			//生成随机的6位数
			Ran := UpDownload.GetRandomString(6)
			router.StaticFile("/static.one"+Ran, "./"+ownerID+"/"+fileName)

			response, err := http.Get("http://localhost:8083/static.one" + Ran)
			if err != nil || response.StatusCode != http.StatusOK {
				context.Status(http.StatusServiceUnavailable)
				return
			}
			reader := response.Body
			contentLength := response.ContentLength
			contentType := response.Header.Get("Content-Type")

			extraHeaders := map[string]string{
				"Content-Disposition": `attachment; filename="` + fileName + "\"",
			}
			context.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

		} else {
			context.String(http.StatusOK, "该一次性快传下载连接已失效")
		}

		//删除数据库one表
		Dao.DeleteOne(ownerIDInt, fileName)
		//不是一次性快传的文件
		//if oneStatus ==0
	})

	//更改云盘文件权限的接口
	router.GET("/cFilePermit/:fileName/:public", func(context *gin.Context) {

		session := sessions2.Default(context)
		userID := session.Get("ID")
		fileName := context.Param("fileName")
		public := context.Param("public")
		str := userID.(string)
		user, _ := strconv.Atoi(str)
		pub, _ := strconv.Atoi(public)
		num := Dao.CFilePermit(user, fileName, pub)
		if num != 0 {
			context.String(http.StatusOK, "更改文件权限成功")
		} else {
			context.String(http.StatusOK, "更改文件权限失败")
		}

	})

	router.Run(":8083")
}
