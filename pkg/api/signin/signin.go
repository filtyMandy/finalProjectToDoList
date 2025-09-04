package signin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"finalProjectToDoList/pkg/util.go"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

var jmtKey = []byte("secret")

func makeJWT(password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"check": Hash(password),
		"exp":   time.Now().Add(8 * time.Hour).Unix(),
	})
	return token.SignedString(jmtKey)
}

func Hash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func validateJWT(tokenString, password string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jmtKey, nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		hashInTocken, _ := claims["check"].(string)
		return hashInTocken == Hash(password), nil
	}
	return false, err
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Password string `json:"password"`
	}
	type Resp struct {
		Token string `json:"token,omitempty"`
		Error string `json:"error,omitempty"`
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util_go.WriteJSONError(w, http.StatusBadRequest, "Error JSON")
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if len(password) == 0 || req.Password != password {
		util_go.WriteJSONError(w, http.StatusUnauthorized, "Error Password")
		return
	}
	token, err := makeJWT(password)
	if err != nil {
		util_go.WriteJSONError(w, http.StatusInternalServerError, "Error JWT generation")
		return
	}
	util_go.WriteJSON(w, Resp{Token: token})
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			util_go.WriteJSONError(w, http.StatusUnauthorized, "Authentification required")
			return
		}
		valid, err := validateJWT(cookie.Value, pass)
		if !valid || err != nil {
			util_go.WriteJSONError(w, http.StatusUnauthorized, "Authentification required")
			return
		}
		next(w, r)
	}
}
