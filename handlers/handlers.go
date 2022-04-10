package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"movie_store/movie"
	"movie_store/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var ptTemplate *template.Template

func init() {
	var err error
	ptTemplate, err = template.New("").ParseGlob("./tpls/*.tpl")
	template.ParseFiles()
	utils.CheckError(err, true)
}

// 获取图片
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	// 需要先登录
	image_filename := "./images/" + mux.Vars(r)["image_filename"]
	if utils.FileExist(image_filename) {
		data, _ := ioutil.ReadFile(image_filename)
		w.Write(data)
	} else {
		http.Error(w, "not found", http.StatusNotFound)
	}
}

// 主页
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	ptTemplate.ExecuteTemplate(w, "homePage.tpl", movie.Movies)
}

// 详情页
func MoviePageHandler(w http.ResponseWriter, r *http.Request) {
	// 需要先登录
	movie_id, err := strconv.Atoi(mux.Vars(r)["movie_id"])
	if err != nil {
		w.Write([]byte("404"))
	} else {
		movie, ok := movie.IdToMovie[movie_id]
		if ok {
			ptTemplate.ExecuteTemplate(w, "movie.tpl", movie)
		} else {
			w.Write([]byte("404"))
		}
	}
}

// 登录页面
func Login(w http.ResponseWriter, r *http.Request) {
	ptTemplate.ExecuteTemplate(w, "login.tpl", nil)
}

// 处理登录表单
func DoLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username != "xxx" || password != "xxx" {
		fmt.Println("username or password is not right.")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	token := fmt.Sprintf("username=%s&password=%s", username, password)
	data := base64.StdEncoding.EncodeToString([]byte(token))
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    data,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
