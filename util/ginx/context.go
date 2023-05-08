package ginx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *Context)

func WrapHandler(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&Context{
			c,
		})
	}
}

type Context struct {
	*gin.Context
}

func (c *Context) Success(code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func (c *Context) Failure(code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func (c *Context) BindJSONEx(obj interface{}) error {
	if c.Request == nil || c.Request.Body == nil {
		return fmt.Errorf("invalid request")
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(removeComments(data), obj)
}
