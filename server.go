package main

import (
	"time"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"regexp"
	"io"
	"crypto/hmac"
	"crypto/sha256"
)
var db *sql.DB
type User struct{
	Id int `json:"id"`
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
	Code int `json:"code"`
}
type JwtToken struct {
	Token string `json:"token"`
}
type Response struct{
	Message string `json:"message"`
	Code int `json:"code"`
	Data interface{} `json:"data"`
}
func ResponseWithJson(w http.ResponseWriter,code int, payload interface{}){
	response,_:=json.Marshal(payload)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func CreateTokenEndpoint (user User) (string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": user.Id,
		"exp": time.Now().Add(1*time.Minute).Unix(),
	})
	tokenString, error := token.SignedString([]byte("secret"))
	return tokenString,error
}
func ValidateTokenMiddleware(next http.HandlerFunc)http.HandlerFunc{
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
					ResponseWithJson(w,http.StatusOK,Response{
						Code:-2001,
						Message: error.Error(),
					})
                    return
                }
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					idStr := fmt.Sprintf("%v",claims["id"])
					req.Header.Set("uId",idStr)
					next(w, req)
                } else {
					ResponseWithJson(w,http.StatusOK,Response{
						Code:-2001,
						Message:  "Invalid authorization token",
					})
                }
				
			}else{
				ResponseWithJson(w,http.StatusOK,Response{
					Code:-2001,
					Message:  "Invalid authorization token",
				})
			}
		}else {
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-2001,
				Message:  "token not exist",
			})
		}
	}
}
func CheckEmail(email string)(bool){
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(email) 
		
}
func SendEmail(email string,code string)error{
	auth := smtp.PlainAuth("", "customer@chelecom.io", "Vhsej5UMWBcZQcwf", "smtp.exmail.qq.com")
    to := []string{email}
    nickname := "test"
    user := "customer@chelecom.io"
    subject := "邮箱验证码"
    content_type := "Content-Type: text/plain; charset=UTF-8"
    body := "验证码："+code
    msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	err := smtp.SendMail("smtp.exmail.qq.com:25", auth, user, to, msg)
	return err
}
func CryptoPassword(data string)string{
	h := hmac.New(sha256.New,[]byte("userPassword"))
	io.WriteString(h,data)
	return fmt.Sprintf("%x",h.Sum(nil))
}
func Register(w http.ResponseWriter,r *http.Request) {
	if r.Method == "POST" {
		user := new(User)
		json.NewDecoder(r.Body).Decode(user)
		// 校验用户提交的用户名秘密是否为空
		if user.Email == "" || user.Password == "" {
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-1,
				Message:  "email and password is nil",
			})
			return
		}
		// 正则校验邮箱
		if !CheckEmail(user.Email) {
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-1,
				Message:  "Email format error",
			})
			return
		}
		// 从数据库中查询用户是否存在
		user1 := User{}
		stmt,_ := db.Prepare(`SELECT id,email,password From user where email = ?`)
		defer stmt.Close()
		
		row :=stmt.QueryRow(user.Email)
		row.Scan(&user1.Id,&user1.Email,&user1.Password)
		
		if user1.Id != 0 && user1.Email != "" && user1.Password != ""{
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-1,
				Message:  "user existed",
			})
			return
		}

		// 创建用户
		stmt,_ = db.Prepare(`INSERT INTO user (email, password) Values (?,?)`)
		defer stmt.Close()

		// 密码加密
		cryptoPassword := CryptoPassword(user.Password)
		result,err := stmt.Exec(user.Email,cryptoPassword)
		if err != nil {
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-1,
				Message:  "create user fail",
			})
			return
		}
		if LastInsertId, err := result.LastInsertId(); nil == err {
			fmt.Println("LastInsertId:", LastInsertId)
		}
		if RowsAffected, err := result.RowsAffected(); nil == err {
			fmt.Println("RowsAffected:", RowsAffected)
		}
		ResponseWithJson(w,http.StatusOK,Response{
			Code:0,
			Message:  "create user success",
		})

	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}

func Login(w http.ResponseWriter,r *http.Request){
	if r.Method == "POST" {
		user := new(User)
		json.NewDecoder(r.Body).Decode(user)
		// 校验用户提交的用户名秘密是否为空
		if user.Email == "" || user.Password == "" {
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-1,
				Message:  "email and password is nil",
			})
			return
		}
		// 正则校验邮箱
		if !CheckEmail(user.Email) {
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-1,
				Message:  "Email format error",
			})
			return
		}

		//从数据库中查询用户是否存在
		user1 := User{}
		stmt,_ := db.Prepare(`SELECT id,email,password From user where email = ?`)
		defer stmt.Close()
		
		row :=stmt.QueryRow(user.Email)
		row.Scan(&user1.Id,&user1.Email,&user1.Password)
		
		if user1.Id == 0 && user1.Email == "" && user1.Password == ""{
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-1,
				Message:  "user not existed",
			})
			return
		}
		// TODO 从数据库中查询到用户 比较密码
		if CryptoPassword(user.Password) != user1.Password {
			ResponseWithJson(w,http.StatusOK,Response{
				Code:-1,
				Message:  "password error",
			})
			return
		}
		// TODO 生成Token返回到客户端
		tokenString,_ :=CreateTokenEndpoint(user1)
		ResponseWithJson(w,http.StatusOK,Response{
			Code:0,
			Message:  "OK",
			Data: JwtToken{Token: tokenString},
		})
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}
func Logout(w http.ResponseWriter,r *http.Request) {
	r.Header.Del("uId")
	ResponseWithJson(w,http.StatusOK,Response{
		Code:0,
		Message:  "OK",
	})
}
func TestDb (w http.ResponseWriter,r *http.Request) {
	if r.Method == "POST" {
		// uid := r.Header.Get("uId")
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
		ResponseWithJson(w,http.StatusOK,Response{
			Code:0,
			Message:  "OK",
			Data: godusers,
		})
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
	http.HandleFunc("/register",Register)
	http.HandleFunc("/login",Login)
	http.HandleFunc("/logout",Logout)
	http.HandleFunc("/test",ValidateTokenMiddleware(TestDb))
	http.ListenAndServe(":9001",nil)
}