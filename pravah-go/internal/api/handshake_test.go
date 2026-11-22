package api

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandshakeSetsSecureCookie(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("COOKIE_SAMESITE", "strict")
	t.Setenv("COOKIE_DOMAIN", "events.example.com")
	t.Setenv("COOKIE_MAX_AGE_SECONDS", "7200")

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest("GET", "/pravah/init", nil)

	HandshakeHandler(context)

	cookie := recorder.Header().Get("Set-Cookie")
	for _, required := range []string{
		"pravah_id=",
		"Domain=events.example.com",
		"Max-Age=7200",
		"HttpOnly",
		"Secure",
		"SameSite=Strict",
	} {
		if !strings.Contains(cookie, required) {
			t.Fatalf("Set-Cookie header %q is missing %q", cookie, required)
		}
	}
}
