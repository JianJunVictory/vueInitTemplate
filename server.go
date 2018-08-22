package main

import (
	"time"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
)
var db *sql.DB
type User struct{
	Email string `json:"email"`
	Password string `json:"password"`
}
type Goduser struct {
	CUST_ID int `json:"id"`
	NAME string `json:"name"`
	AGE int `json:"age"`
}
type RespData struct {
	Message string `json:"message"`
	Code string `json:"code"`
}
type JwtToken struct {
	Token string `json:"token"`
}
func CreateTokenEndpoint (user User) (string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": user.Email,
		"exp": time.Now().Add(2*time.Minute).Unix(),
	})
	tokenString, error := token.SignedString([]byte("secret"))
	return tokenString,error
}
func ValidateMiddleware(next http.HandlerFunc)http.HandlerFunc{
	return func(w http.ResponseWriter,req *http.Request){
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
                    json.NewEncoder(w).Encode(RespData{Message: error.Error(),Code:"-1"})
                    return
                }
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {		
					req.Header.Set("userEmail",claims["email"].(string))
					next(w, req)
                } else {
                    json.NewEncoder(w).Encode(RespData{Message: "Invalid authorization token",Code:"-1"})
                }
				
			}else{
				json.NewEncoder(w).Encode(RespData{Message: "Invalid authorization token",Code:"-1"})
			}
		}else {
			json.NewEncoder(w).Encode(RespData{Message: "token not exist",Code:"-1"})
		}
	}
}
func Login(w http.ResponseWriter,r *http.Request){
	if r.Method == "POST" {
		user := new(User)
		json.NewDecoder(r.Body).Decode(user)
		if user.Email == "" || user.Password == "" {
			respData := RespData{Message:"email and password is nil",Code:"-1"}
			errdata,_:=json.Marshal(respData)
			w.Write(errdata)
			return
		}
		w.Header().Set("Content-Type","Applocation/json")
		w.WriteHeader(http.StatusOK)

		tokenString,_ :=CreateTokenEndpoint(*user)
		json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}
func TestDb (w http.ResponseWriter,r *http.Request) {
	if r.Method == "POST" {
		email := r.Header.Get("userEmail")
		fmt.Println("这是从token中校验出来的")
		fmt.Println(email)
		var godusers []Goduser
		goduser := Goduser{}
		stmt,_ := db.Prepare(`SELECT * From customer where AGE < ?`)
		defer stmt.Close()
		rows,err := stmt.Query(33)
		if err != nil {
			fmt.Printf("insert data error: %v\n", err)
		}
		for rows.Next() {
			rows.Scan(&goduser.CUST_ID,&goduser.NAME,&goduser.AGE)
			godusers=append(godusers,goduser)
		}
		data,_:=json.Marshal(godusers)
		w.Write(data)
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
	
}
func init () {
	dbs,err := sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/jun")
	if err != nil {
		log.Panicln("db connect failed,error:",err)
	}
	db = dbs
}
func main(){
	http.HandleFunc("/login",Login)
	http.HandleFunc("/test",ValidateMiddleware(TestDb))
	http.ListenAndServe(":9001",nil)
}