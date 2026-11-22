package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sampledge-sse/pravah-go/internal/utils"
)

func HandshakeHandler(c *gin.Context) {
	pravahID, err := c.Cookie("pravah_id")
	isNew := false

	if err != nil {
		pravahID = uuid.New().String()
		isNew = true

		switch utils.CookieSameSite() {
		case "strict":
			c.SetSameSite(http.SameSiteStrictMode)
		case "none":
			c.SetSameSite(http.SameSiteNoneMode)
		default:
			c.SetSameSite(http.SameSiteLaxMode)
		}
		c.SetCookie(
			"pravah_id",
			pravahID,
			utils.CookieMaxAgeSeconds(),
			"/",
			utils.CookieDomain(),
			utils.CookieSecure(),
			true,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"pravah_id": pravahID,
		"new":       isNew,
		"message": func() string {
			if isNew {
				return "Pravah ID initialized"
			}
			return "Pravah ID recognized"
		}(),
	})
}
