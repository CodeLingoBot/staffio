package web

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/settings"
)

var (
	CookieName = "_user"
)

const (
	kAuthUser = "user"
)

func AuthUserMiddleware(redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := UserFromRequest(c.Request)
		if err != nil {
			log.Printf("user from request ERR %s", err)
			if redirect {
				markReferer(c)
				c.Redirect(302, UrlFor("login"))
				c.Abort()
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
				// apiError(c, ERROR_INTERNAL, err)
			}
			return
		}
		// log.Printf("got user %q", user.Uid)
		c.Set(kAuthUser, user)
		c.Next()
		user.Refresh()
		user.toResponse(c.Writer)
	}
}

func (s *server) authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exist := c.Get(kAuthUser)
		if !exist {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Abort()
			return
		}
		user := v.(*User)
		if !s.IsKeeper(user.Uid) {
			c.AbortWithStatus(http.StatusForbidden)
			c.Abort()
			return
		}
	}
}

func UserWithContext(c *gin.Context) (user *User) {
	if v, ok := c.Get(kAuthUser); ok {
		user = v.(*User)
	}
	if user == nil {
		panic("user not found in request")
	}

	return
}

func UserFromRequest(r *http.Request) (user *User, err error) {
	var cookie *http.Cookie
	cookie, err = r.Cookie(CookieName)
	if err != nil {
		log.Printf("cookie %q ERR %s", CookieName, err)
		return
	}
	var b []byte
	b, err = base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Printf("base64decode %q ERR %s", cookie.Value, err)
		return
	}
	// log.Printf("got encrypted %s", b)
	user = new(User)
	err = user.Decode(b)
	if err != nil {
		log.Printf("decode msgpack ERR %s", err)
	}
	if user.IsExpired(settings.UserLifetime) {
		err = fmt.Errorf("user %s is expired", user.Uid)
	}
	// log.Printf("got user %v", user)
	return
}

func (user *User) toResponse(w http.ResponseWriter) error {
	b, err := user.Encode()
	if err != nil {
		log.Printf("marshal msgpack user ERR: %s", err)
		return err
	}
	value := base64.URLEncoding.EncodeToString(b)
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    value,
		MaxAge:   settings.UserLifetime,
		Path:     "/",
		HttpOnly: true,
	})
	return nil
}

func signinStaffGin(c *gin.Context, staff *models.Staff) {
	user := UserFromStaff(staff)
	user.Refresh()
	log.Printf("login ok %v", user)
	sess := ginSession(c)
	sess.Set(kAuthUser, user)
	user.toResponse(c.Writer)
	SessionSave(sess, c.Writer)
}

func signout(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})
}
