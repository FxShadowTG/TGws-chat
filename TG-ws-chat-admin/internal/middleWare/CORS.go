/**
 * @Author: FxShadow
 * @Description:
 * @Date: 2022/12/26 22:51
 */

package MiddleWare

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CorsHandler 跨域请求处理
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 这是允许访问所有域
		// c.Header().Set("Access-Control-Allow-Origin", "")           // 跨域请求是否需要带cookie信息 默认设置为true
		c.Header("Access-Control-Allow-Credentials", "true") // 跨域请求是否需要带cookie信息 默认设置为true
		//  header的类型
		c.Header("Access-Control-Allow-Origin", "*") // 这是允许访问所有域
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control,XMLHttpRequest, X-Requested-With")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, custom-header, Cache-Control,XMLHttpRequest, X-Requested-With")
		//c.Header("Access-Control-Allow-Headers", "*")
		//服务器支持的所有跨域请求的方法
		// c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		// c.Header("Access-Control-Max-Age", "21600") //可以缓存预检请求结果的时间（以秒为单位）
		// c.Set("content-type", "application/json")   // 设置返回格式是json
		if c.Request.Method == "OPTIONS" {
			// 	// c.AbortWithStatus(204)
			// 	c.AbortWithStatus(http.StatusNoContent)
			// 	return
			// 	放行所有OPTIONS方法，本项目直接返回204
			c.JSON(200, "Options Request!")
		}

		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Cors() \n")
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
