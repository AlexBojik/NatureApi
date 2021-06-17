package models

import (
	"log"
	"time"
	"water-api/sql"
)

type News struct {
	Id          int       `json:"id"`
	Created     time.Time `json:"created"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
}

func GetNews(filtered bool) []*News {
	res := make([]*News, 0)

	now := time.Now()
	sqlRequest := sql.NewsList

	if filtered {
		sqlRequest = sql.NewsFilteredList
	}

	rows, err := db.Query(sqlRequest, now, now)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		var created string
		var start string
		var end string
		var descr string
		nw := News{}
		err = rows.Scan(&nw.Id, &descr, &created, &start, &end)
		nw.Created, _ = time.Parse("2006-01-02 15:04:05", created)
		nw.Start, _ = time.Parse("2006-01-02 15:04:05", start)
		nw.End, _ = time.Parse("2006-01-02 15:04:05", end)
		nw.Description = descr
		res = append(res, &nw)
	}

	return res
}

func CreateNews(news *News) int64 {
	res, err := db.Exec(sql.NewsCreate, news.Created.Format("2006-01-02 15:04:05"), news.Start.Format("2006-01-02 15:04:05"), news.End.Format("2006-01-02 15:04:05"), news.Description)
	if err != nil {
		print(err.Error())
		return 0
	}
	id, _ := res.LastInsertId()
	return id
}

func UpdateNews(news *News) {
	db.Exec(sql.NewsUpdate, news.Created.Format("2006-01-02 15:04:05"), news.Start.Format("2006-01-02 15:04:05"), news.End.Format("2006-01-02 15:04:05"), news.Description, news.Id)
}

func DeleteNews(id int) {
	db.Exec(sql.NewsDelete, id)
}
