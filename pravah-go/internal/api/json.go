package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/utils"
)

func decodeStrictJSON(c *gin.Context, target interface{}) error {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, utils.MaxRequestBodyBytes())
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("request body must contain one JSON object")
	}
	return nil
}

func badRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Bad Request",
		"message": message,
	})
}
