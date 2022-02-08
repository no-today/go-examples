package response

import (
	"cathub.me/go-gin-examples/pkg/errors"
	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context, entity *ResponseEntity) {
	for k, v := range entity.Headers {
		c.Header(k, v)
	}

	if entity.Body == nil {
		c.JSON(entity.Status, gin.H{"msg": entity.Msg})
	} else {
		c.JSON(entity.Status, entity.Body)
	}
}

func Fail(c *gin.Context, error error) {
	problem := errors.ErrorHandler(error)
	problem.Path = c.Request.RequestURI

	c.JSON(problem.Status, problem)
}
