package middleware

import (
	"backend/controllers"
	"backend/loader"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		tokenString1:=r.Header.Get("Authorization")
		
		tokenString,err:=loader.Decrypt(tokenString1,os.Getenv("ENCODE_KEY"))
		fmt.Println("poda opoda",tokenString)	
		if !(tokenString!=""){
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
		// verifying token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
		})
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix())>claims["exp"].(float64){
				http.Error(w, "Forbidden", http.StatusForbidden)
			}else{
				exist:=controllers.CheckUserByEmail(claims["sub"].(string))
				if !exist{
					
					http.Error(w, "Forbidden", http.StatusForbidden)
				}else{
					next.ServeHTTP(w, r)
				}


				// fmt.Println(claims)
			}
		} else {
			fmt.Println(err,70)
			http.Error(w, "Forbidden", http.StatusForbidden)
		}

		// fmt.Println(auth.Value)
		
	})
	
	
}
