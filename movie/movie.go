package movie

import (
	"movie_store/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Movie struct {
	Id           int    `db:"id"`
	Title        string `db:"title"`
	Area         string `db:"area"`
	Which_type   string `db:"which_type"`
	Date         string `db:"date"`
	Rate         string `db:"rate"`
	Language     string `db:"language"`
	Director     string `db:"director"`
	Actor        string `db:"actor"`
	Image_path   string `db:"image_url"`
	Introduction string `db:"introduction"`
}

var db *sqlx.DB
var Movies []*Movie
var IdToMovie map[int]*Movie

func init() {
	IdToMovie = make(map[int]*Movie)

	database, err := sqlx.Open("mysql", "root:your_password@tcp(localhost:3306)/movies")
	utils.CheckError(err, true)
	err = database.Ping()
	utils.CheckError(err, true)
	db = database

	sql := "select id, title, area, which_type, date, rate, language, director, actor, image_url, introduction from movies"
	err = db.Select(&Movies, sql)
	utils.CheckError(err, true)

	for _, movie := range Movies {
		IdToMovie[movie.Id] = movie
	}
}
