package main

import (
	"github.com/joho/godotenv"
	"os"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"time"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	_ "github.com/jinzhu/gorm/dialects/mysql"

)

type qiitaArticleJson struct {
	CreatedAt time.Time `json:"created_at"`
	ID string `json:"id"`
	Private bool `json:"private"`
	Tags []struct {
		Name string `json:"name"`
		Versions []string `json:"versions"`
	} `json:"tags"`
	Title string `json:"title"`
	URL string `json:"url"`
}

type qiitaArticle struct {
	gorm.Model
	Id    int `gorm:"primary_key"`
	Title string
	Link string
	Datetime string
}

type qiitaTag struct {
	gorm.Model
	Tag_Id int `gorm:"primary_key"`
	Article_Id int
	Tag string
}


// qiita apiを叩いて自分の投稿一覧とストック/いいね一覧を取得してDBに格納する
func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Can't Read env")
		os.Exit(1)
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_CHARSET := os.Getenv("DB_CHARSET")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	DB_CONNECT := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME +"?charset=" + DB_CHARSET + "&parseTime=true&loc=Asia%2FTokyo"
	db, err := gorm.Open("mysql", DB_CONNECT)
	defer db.Close()


	qiitaEndpoint := "https://qiita.com/api/v2/"
	qiitaAccessToken := os.Getenv("ACCESS_TOKEN")


	request := gorequest.New()
	// 汎用的に繰り返しページ数を取得してくるべきだが、件数少ないので一旦これで
	_, body ,errs := request.Get(qiitaEndpoint + "authenticated_user/items?per_page=100").Set("Authorization", "Bearer " + qiitaAccessToken).End()

	if errs != nil {
		fmt.Println("Request Err")
		os.Exit(1)
	}


	jsonBytes := []byte(body)

	var data []qiitaArticleJson
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}


	qiitaArticleIn := qiitaArticle{}
	for _, d := range data {
		qiitaArticleIn.ID = 0
		qiitaArticleIn.Title = d.Title
		qiitaArticleIn.Link = d.URL
		qiitaArticleIn.Datetime = d.CreatedAt.Format("2006-01-02 15:04:05")

		log.Println("ID : " + strconv.Itoa(int(qiitaArticleIn.ID)))
		log.Println("TITLE : " + d.Title)
		log.Println("Link : " + d.URL)
		log.Println("date : " + d.CreatedAt.Format("2006-01-02 15:04:05"))
		db.Create(&qiitaArticleIn)

		for _, t := range d.Tags {
			tagIn := qiitaTag{}
			tagIn.Article_Id = int(qiitaArticleIn.ID)
			tagIn.Tag = t.Name
			tagIn.Tag_Id = 0

			log.Println("ARTICLEID : " + strconv.Itoa(int(qiitaArticleIn.ID)))
			log.Println("TAG : " + t.Name)
			db.Create(&tagIn)

		}
	}






}
