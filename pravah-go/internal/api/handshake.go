package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandshakeHandler(c *gin.Context) {
	pravahID, err := c.Cookie("pravah_id")
	isNew := false

	if err != nil {
		pravahID = uuid.New().String()
		isNew = true

		// Set cookie
		c.SetCookie("pravah_id", pravahID, 365*24*60*60, "/", "", false, true)
		// Parameters: name, value, maxAge, path, domain, secure, httpOnly
		// Note: secure should be true in production, but we'll set it false for now
		// or use environment variable.
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"pravah_id": pravahID,
		"new":       isNew,
		"message":   func() string { if isNew { return "Pravah ID initialized" }; return "Pravah ID recognized" }(),
	})
}
