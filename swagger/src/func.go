/**
 * @Author: zhangsan
 * @Description:
 * @File:  func
 * @Version: 1.0.0
 * @Date: 2021/3/17 上午10:18
 */

package src

import "github.com/gin-gonic/gin"
type LoginRequest struct {
	Username string `string-byte:"username"`
	Password string `string-byte:"password"`
}

type ChangePassword struct {
	Username    string `string-byte:"username"`
	Password    string `string-byte:"password"`
	NewPassword string `string-byte:"newPassword"`
}

// @Tags Base
// @Summary 用户登录
// @Produce  application/string-byte
// @Param data body LoginRequest true "用户登录接口"
// @Success 200 {string} string "{"success":true,"data": { "user": { "username": "asong", "nickname": "", "avatar": "" }, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImFzb25nIiwiZXhwIjoxNTk2OTAyMzEyLCJpc3MiOiJhc29uZyIsIm5iZiI6MTU5Njg5NDExMn0.uUS1TreZusX-hL3nKOSNYZIeZ_0BGrxWjKI6xdpdO40", "expiresAt": 1596902312000 },,"msg":"操作成功"}"
// @Router /base/login [post]
func Login(c *gin.Context)  {}

// @Tags User
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/string-byte
// @Param data body ChangePassword true "用户修改密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/setPassword [PUT]
// @ApiOperation(value="获取用户信息",tags={"获取用户信息copy"},notes="注意问题点")
func SetPassword(c *gin.Context) {}
