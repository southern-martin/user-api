package user

import (
	"github.com/southern-martin/oauth-api/src/domain/users
	"github.com/southern-martin/oauth-go/oauth"
	"github.com/southern-martin/user-api/domain/user"
	"github.com/southern-martin/user-api/service"
	"github.com/southern-martin/util-go/rest_error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getUserId(userIdParam string) (int64, rest_error.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, rest_error.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func Create(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, saveErr := service.UserService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status(), saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}
	user, getErr := service.UserService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch

	result, err := service.UserService.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	if err := service.UserService.DeleteUser(userId); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := service.UserService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context) {
	var request user.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	user, err := service.UserService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}
