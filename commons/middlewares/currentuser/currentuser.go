package currentuser

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/somtooo/Chit-Slip-Commons/commons/errors"
	"net/http"
	"os"
)

type user struct {
	CurrentUser jwt.MapClaims `json:"currentUser"`
}

type currentUser string

// Key is the key to access the currentUser ctx on req
const Key currentUser = "currentUser"

//CurrentUser is a Middleware that checks if user is logged in
func CurrentUser(handler http.Handler) http.Handler {
	ctx := context.WithValue(context.Background(), Key, user{CurrentUser: nil})
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie := req.Cookies()
		if len(cookie) == 0 {
			req = req.Clone(ctx)
			handler.ServeHTTP(res, req)
			return
		}

		token, err := jwt.Parse(cookie[0].Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if token == nil {
			fmt.Println("Token is null")
			req = req.Clone(ctx)
			handler.ServeHTTP(res, req)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(context.Background(), Key, user{CurrentUser: claims})
			req = req.Clone(ctx)
			handler.ServeHTTP(res, req)
		} else {
			fmt.Println("Token verify Error: ", err)
			req = req.Clone(ctx)
			handler.ServeHTTP(res, req)
		}
	})

}

//RequireAuth throws an error is user is not authorized
func RequireAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if user := req.Context().Value(Key).(user); user.CurrentUser == nil {
			var notAuthorized errors.BadRequestError = "Not authorized"
			errors.HTTPError(res, notAuthorized, http.StatusUnauthorized)
		}

	})
}
