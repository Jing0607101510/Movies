package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"movie_store/utils"
	"net/http"
	"net/http/cookiejar"
	"path"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/publicsuffix"
)

func myCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 5 {
		return errors.New("stopped after 5 redirects")
	}
	return nil
}

func CreateClient() *http.Client {
	cookie_options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&cookie_options)
	utils.CheckError(err, true)

	client := &http.Client{
		Timeout:       15 * time.Second,
		Jar:           jar,
		CheckRedirect: myCheckRedirect,
	}
	return client
}

func CrawlPage(client *http.Client, startPage int, wg *sync.WaitGroup) {
	defer wg.Done()
	targetUrl := "https://www.lfmsj.com/category/71/page/" + strconv.Itoa(startPage+1) + ".html"
	fmt.Println(targetUrl)
	// req, err := http.NewRequest("Get", targetUrl, nil)
	// if utils.CheckError(err, false) {
	// 	return
	// }
	// // "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36 Edg/99.0.1150.55"
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	// resp, err := client.Do(req)
	resp, err := client.Get(targetUrl)
	if utils.CheckError(err, false) {
		return
	}
	defer resp.Body.Close()

	HandlePage(client, resp.Body)
}

func HandlePage(client *http.Client, body io.Reader) {
	doc, err := goquery.NewDocumentFromReader(body)
	if utils.CheckError(err, false) {
		return
	}

	doc.Find("ul.vodlist li.vodlist_item").Each(
		func(i int, s *goquery.Selection) {
			href, exist := s.Find("div.vodlist_titbox p.vodlist_title a").Attr("href")
			wholeHref := "https://www.lfmsj.com" + href
			if exist {
				HandleMovieInfo(client, wholeHref)
			}
		},
	)
}

func HandleMovieInfo(client *http.Client, href string) {
	resp, err := client.Get(href)
	if utils.CheckError(err, false) {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if utils.CheckError(err, false) {
		return
	}

	title := doc.Find("h1.title").Text()

	info1_content := doc.Find("li.data").First().Text()
	info1_reg := regexp.MustCompile(`^地区：(.*?)类型：(.*?)上映时间：(.*?)豆瓣评分：(.*?)$`)
	info1 := info1_reg.FindStringSubmatch(info1_content)
	area := info1[1]
	whichType := info1[2]
	date := info1[3]
	rate := info1[4]

	info2_content := doc.Find("li.data").Eq(1).Text()
	info2_reg := regexp.MustCompile(`^语言：(.*?)状态：`)
	info2 := info2_reg.FindStringSubmatch(info2_content)
	language := info2[1]

	info3_content := doc.Find("li.data").Eq(2).Text()
	info3_reg := regexp.MustCompile(`^导演：(.*?)$`)
	info3 := info3_reg.FindStringSubmatch(info3_content)
	director := info3[1]

	info4_content := doc.Find("li.data").Eq(3).Text()
	info4_reg := regexp.MustCompile(`^主演：(.*?)$`)
	info4 := info4_reg.FindStringSubmatch(info4_content)
	actor := info4[1]

	image_url, _ := doc.Find("div.content_thumb a.vodlist_thumb").Attr("data-original")

	introduction := doc.Find("body > div.container > div.left_row.fl > div:nth-child(1) > div > div > section:nth-child(1) > div.content_desc.full_text.clearfix > span").Text()

	// DownloadImage(client, image_url)

	image_path := path.Base(image_url)

	_, err = db.Exec(sql, title, area, whichType, date, rate, language, director, actor, image_path, introduction)
	if utils.CheckError(err, false) {
		return
	}

	fmt.Println(title, area, whichType, date, rate, language, director, actor, image_url, introduction)
}

func DownloadImage(client *http.Client, image_url string) {
	resp, err := client.Get(image_url)
	if utils.CheckError(err, false) {
		return
	}
	defer resp.Body.Close()

	image_name := "../images/" + path.Base(image_url)
	data, _ := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile(image_name, data, 0666)
}

var db *sqlx.DB
var sql = "insert into movies(title, area, which_type, date, rate, language, director, actor, image_url, introduction) values (?,?,?,?,?,?,?,?,?,?)"

func main() {
	database, err := sqlx.Open("mysql", "root:your_password@tcp(localhost)/movies")
	if err != nil {
		fmt.Println("创建数据库入库对象sql.DB失败！")
		log.Fatal(err)
	} else {
		fmt.Println("创建数据库入库对象sql.DB成功！")
	}
	if err := database.Ping(); err != nil {
		fmt.Println("连接数据库失败！")
		log.Fatal(err)
	} else {
		fmt.Println("连接数据库成功！")
	}
	utils.CheckError(err, true)
	db = database

	client := CreateClient()
	nPages := 5
	wg := sync.WaitGroup{}
	for i := 0; i < nPages; i++ {
		wg.Add(1)
		go CrawlPage(client, i, &wg)
	}
	wg.Wait()
}
