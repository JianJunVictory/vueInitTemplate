package main

import (
	"crypto/sha1"
	"net/url"
	"sort"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"time"
	"math/rand"
	"encoding/base64"
)

var db *sql.DB

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   int    `json:"status"`
}

type Goduser struct {
	CUST_ID int    `json:"id"`
	NAME    string `json:"name"`
	AGE     int    `json:"age"`
}
type RespData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type JwtToken struct {
	Token string `json:"token"`
}
type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

type SendSmsReply struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}
func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func CreateTokenEndpoint(Id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  Id,
		"exp": time.Now().Add(30 * time.Minute).Unix(),
	})
	tokenString, error := token.SignedString([]byte("secret"))
	return tokenString, error
}
func ProtectedEndpoint(params string) (bool, int) {
	token, _ := jwt.Parse(params, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idStr := int(claims["id"].(float64))
		return true, idStr
	} else {
		return false, 0
	}
}
func ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
					ResponseWithJson(w, http.StatusOK, Response{
						Code:    -2001,
						Message: error.Error(),
					})
					return
				}
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					idStr := fmt.Sprintf("%v", claims["id"])
					req.Header.Set("uId", idStr)
					next(w, req)
				} else {
					ResponseWithJson(w, http.StatusOK, Response{
						Code:    -2001,
						Message: "Invalid authorization token",
					})
				}

			} else {
				ResponseWithJson(w, http.StatusOK, Response{
					Code:    -2001,
					Message: "Invalid authorization token",
				})
			}
		} else {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -2001,
				Message: "token not exist",
			})
		}
	}
}
func SendSms(phone string){
	// 第一步：请求参数
	accessKeyId := "LTAIzYQhVdws2SX5";
	accessKeySecret := "mju7qFeeZP4XhARdBvjbViYAX9Ka7o";
	paras := make(map[string]string)
	// 系统参数
	paras["SignatureMethod"]= "HMAC-SHA1"
	paras["SignatureNonce"]=fmt.Sprintf("%d", rand.Int63())
	paras["AccessKeyId"]= accessKeyId
	paras["SignatureVersion"]= "1.0"
	paras["Timestamp"]=time.Now().UTC().Format("2006-01-02T15:04:05Z")
	paras["Format"]= "JSON";
	// 业务API参数
	paras["Action"]= "SendSms"
	paras["Version"]="2017-05-25"
	paras["RegionId"]= "cn-hangzhou"
	paras["PhoneNumbers"]= phone
	paras["SignName"]="熵链科技"
	paras["TemplateParam"]="{\"code\":\"123456\"}"
	paras["TemplateCode"]="SMS_136450036"
	paras["OutId"]="yourOutId"
	
	// 第二步：根据参数Key排序（顺序）
	var keys []string
	for k := range paras {
        keys = append(keys, k)
    }
	sort.Strings(keys)
	//第三步：构造待签名的请求串
	var sortQueryString string
	for _, v := range keys {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, replace(v), replace(paras[v]))
	}
	stringToSign := fmt.Sprintf("GET&%s&%s", replace("/"), replace(sortQueryString[1:]))
	
	sign := Sign(accessKeySecret,stringToSign)

	apiUrl := fmt.Sprintf("http://dysmsapi.aliyuncs.com/?Signature=%s%s", sign, sortQueryString)

	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}
	ssr := new(SendSmsReply)
	json.NewDecoder(resp.Body).Decode(ssr)
	
	fmt.Printf("%#v\n",ssr)
}
func specialUrlEncode(str string)string{
	strTmp ,_ :=url.ParseQuery(str)
	return strTmp.Encode()
}
func Sign(accessKeySecret,stringToSign string)string{
	h :=hmac.New(sha1.New,[]byte(fmt.Sprintf("%s&", accessKeySecret)))
	h.Write([]byte(stringToSign))
	str := replace(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	return str
}
func replace(in string) string {
	rep := strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")
	return rep.Replace(url.QueryEscape(in))
}

func CheckEmail(email string) bool {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(email)

}
func SendEmail(email string, checkToken string) error {
	auth := smtp.PlainAuth("", "customer@chelecom.io", "Vhsej5UMWBcZQcwf", "smtp.exmail.qq.com")
	to := []string{email}
	nickname := "自由的空气"
	user := "customer@chelecom.io"
	subject := "邮箱验证码"
	content_type := "Content-Type: text/html; charset=UTF-8"
	body := "<html><body><a href=http://localhost:8888/#/activeAccount?token=" + checkToken + ">点击链接激活账号http://localhost:8888/#/activeAccount?token=" + checkToken + "</a><body></html>"
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	err := smtp.SendMail("smtp.exmail.qq.com:25", auth, user, to, msg)
	return err
}
func CryptoPassword(data string) string {
	h := hmac.New(sha256.New, []byte("userPassword"))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 注册
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := new(User)
		json.NewDecoder(r.Body).Decode(user)
		// 校验用户提交的用户名秘密是否为空
		if user.Email == "" || user.Password == "" {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -1,
				Message: "email and password is nil",
			})
			return
		}
		// 正则校验邮箱
		if !CheckEmail(user.Email) {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -1,
				Message: "Email format error",
			})
			return
		}
		// 从数据库中查询用户是否存在
		user1 := User{}
		stmt, _ := db.Prepare(`SELECT id,email,password,status From user where email = ?`)
		defer stmt.Close()

		row := stmt.QueryRow(user.Email)
		row.Scan(&user1.Id, &user1.Email, &user1.Password, &user1.Status)

		if user1.Id != 0 && user1.Email != "" && user1.Password != "" {
			if user1.Status != 0 { // 用户是激活用户
				ResponseWithJson(w, http.StatusOK, Response{
					Code:    -1,
					Message: "user existed",
				})
				return
			} else { // 用户是非激活用户
				checkToken, _ := CreateTokenEndpoint(int(user1.Id))
				err := SendEmail(user1.Email, checkToken)
				if err != nil { // 发送邮件失败
					ResponseWithJson(w, http.StatusOK, Response{
						Code:    -1,
						Message: err.Error(),
					})
					return
				}
				// 发送邮件成功
				ResponseWithJson(w, http.StatusOK, Response{
					Code:    0,
					Message: "已经重新发送激活邮件，请到邮箱激活账户",
				})
				return
			}
		}

		// 创建用户
		stmt, _ = db.Prepare(`INSERT INTO user (email, password) Values (?,?)`)
		defer stmt.Close()

		// 密码加密
		cryptoPassword := CryptoPassword(user.Password)
		result, err := stmt.Exec(user.Email, cryptoPassword)
		if err != nil {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -1,
				Message: "create user fail",
			})
			return
		}
		LastInsertId, err := result.LastInsertId()
		if nil != err {
			fmt.Println(err)
		}
		checkToken, _ := CreateTokenEndpoint(int(LastInsertId))
		err = SendEmail(user.Email, checkToken)
		if err != nil {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -1,
				Message: err.Error(),
			})
			return
		}

		ResponseWithJson(w, http.StatusOK, Response{
			Code:    0,
			Message: "用户注册成功，请到邮箱激活账户",
		})

	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
//激活用户存在问题，应该点击一次，只能生效一次后就失效，此处仅作参考，激活示例
func ActiveAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		token := new(JwtToken)
		json.NewDecoder(r.Body).Decode(token)
		result, uid := ProtectedEndpoint(token.Token)
		if result {
			newToken, _ := CreateTokenEndpoint(uid)
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    0,
				Message: "OK",
				Data:    JwtToken{Token: newToken},
			})
		} else {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -2001,
				Message: "Token is expired",
			})
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := new(User)
		json.NewDecoder(r.Body).Decode(user)
		// 校验用户提交的用户名秘密是否为空
		if user.Email == "" || user.Password == "" {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -1,
				Message: "email and password is nil",
			})
			return
		}
		// 正则校验邮箱
		if !CheckEmail(user.Email) {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -1,
				Message: "Email format error",
			})
			return
		}

		//从数据库中查询用户是否存在
		user1 := User{}
		stmt, _ := db.Prepare(`SELECT id,email,password,status From user where email = ?`)
		defer stmt.Close()

		row := stmt.QueryRow(user.Email)
		row.Scan(&user1.Id, &user1.Email, &user1.Password, &user1.Status)

		if user1.Id == 0 && user1.Email == "" && user1.Password == "" {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -1,
				Message: "user not existed",
			})
			return
		}
		// TODO 从数据库中查询到用户 比较密码
		if CryptoPassword(user.Password) != user1.Password {
			ResponseWithJson(w, http.StatusOK, Response{
				Code:    -1,
				Message: "password error",
			})
			return
		}
		// TODO 生成Token返回到客户端
		tokenString, _ := CreateTokenEndpoint(user1.Id)
		ResponseWithJson(w, http.StatusOK, Response{
			Code:    0,
			Message: "OK",
			Data:    JwtToken{Token: tokenString},
		})
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
func Logout(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("uId")
	ResponseWithJson(w, http.StatusOK, Response{
		Code:    0,
		Message: "OK",
	})
}
func TestDb(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// uid := r.Header.Get("uId")
		var godusers []Goduser
		goduser := Goduser{}
		stmt, _ := db.Prepare(`SELECT * From customer where AGE < ?`)
		defer stmt.Close()
		rows, err := stmt.Query(33)
		if err != nil {
			fmt.Printf("insert data error: %v\n", err)
		}
		for rows.Next() {
			rows.Scan(&goduser.CUST_ID, &goduser.NAME, &goduser.AGE)
			godusers = append(godusers, goduser)
		}
		ResponseWithJson(w, http.StatusOK, Response{
			Code:    0,
			Message: "OK",
			Data:    godusers,
		})
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
func testFunc(w http.ResponseWriter, r *http.Request){
	SendSms("13402651404")
	w.Write([]byte("OK"))
}
func init() {
	dbs, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/jun")
	if err != nil {
		log.Panicln("db connect failed,error:", err)
	}
	db = dbs
}
func main() {
	http.HandleFunc("/register", Register)
	http.HandleFunc("/active", ActiveAccount)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/testDb", ValidateTokenMiddleware(TestDb))
	http.HandleFunc("/test", testFunc)
	http.ListenAndServe(":9001", nil)
}
