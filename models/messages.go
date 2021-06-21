package models

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	geojson "github.com/paulmach/go.geojson"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
	"water-api/sql"
	"water-api/utils"
)

type Message struct {
	Id       int64   `json:"id"`
	UserName string  `json:"userName"`
	Status   int     `json:"status"`
	Images   []Image `json:"images"`
	Text     string  `json:"text"`
	Lat      float32 `json:"lat"`
	Lon      float32 `json:"lon"`
	Token    string  `json:"token"`
}

type Image struct {
	Jpeg string `json:"jpeg"`
}

func PostMessage(r *http.Request) []int64 {
	m := &Message{}
	_ = json.NewDecoder(r.Body).Decode(&m)
	point := "POINT(" + fmt.Sprintf("%f", m.Lat) + " " + fmt.Sprintf("%f", m.Lon) + ")"
	timestamp := time.Now().Unix()
	row, err := db.Exec(sql.MessageCreate, m.Token, m.Text, point, timestamp)
	var res = make([]int64, 0)
	if err != nil {
		return res
	}
	id, err := row.LastInsertId()
	if err != nil {
		return res
	}
	m.Id = id
	res = append(res, id)
	var ims []Image
	for _, img := range m.Images {
		fileName := utils.Base64toPng(img.Jpeg)
		im := Image{}
		im.Jpeg = fileName
		ims = append(ims, im)

		_, _ = db.Exec(sql.ImageCreate, id, fileName)
	}
	m.Images = ims
	go SendMail(m)
	return res
}

func GetMessagesCount() int {
	var count = 0

	row := db.QueryRow(sql.MessageCount)
	row.Scan(&count)
	return count
}

func SendMessages() int {
	rows, err := db.Query(sql.MessageNotSendList)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		var point *geojson.Geometry
		m := Message{}
		err = rows.Scan(&m.Id, &m.Text, &m.UserName, &m.Status, &point, &m.Token)

		m.Lat = float32(point.Point[0])
		m.Lon = float32(point.Point[1])
		m.Images = GetMessage(strconv.Itoa(int(m.Id)))

		SendMail(&m)
	}
	return 0
}

func GetMessageList() []*Message {
	res := make([]*Message, 0)

	rows, err := db.Query(sql.MessageList)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		m := Message{}
		err = rows.Scan(&m.Id, &m.Text, &m.UserName, &m.Status)
		res = append(res, &m)
	}

	return res
}

func GetMessage(id string) []Image {
	images := make([]Image, 0)
	idInt, _ := strconv.Atoi(id)
	rows, err := db.Query(sql.Message, idInt)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		image := Image{}
		err = rows.Scan(&image.Jpeg)
		images = append(images, image)
	}
	return images
}

func SendMail(m *Message) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("MAIL_AUTH")
	password := os.Getenv("MAIL_PASS")
	server := os.Getenv("MAIL_SERVER")
	//mailTo := os.Getenv("MAIL_TO")
	mailTo := getEnvAsSlice("MAIL_TO", ",")
	imageSrc := os.Getenv("IMAGE_SRC")

	auth := smtp.PlainAuth("", username, password, server)

	user := GetUser(m.Token)
	body := "<html>"
	body += "<body>"
	body += "<h3>Номер обращения: " + strconv.Itoa(int(m.Id)) + "</h3>"
	body += "<h3>ФИО: " + user.Name + "</h3>"
	body += "<h3>Текст обращения: " + m.Text + "</h3>"
	body += "<h3>Координаты: " + fmt.Sprintf("%f", m.Lat) + ", " + fmt.Sprintf("%f", m.Lon) + "</h3>"
	body += "<h3>Телефон: " + user.Phone + "</h3>"
	body += "<h3>E-mail: " + user.Email + "</h3>"
	body += "<h3>Адрес: " + user.ProAddr + "</h3>"
	body += "<h3>Прикрепленные изображения:</b></h3>"
	for _, img := range m.Images {
		body += "<img src=\"" + imageSrc + img.Jpeg + "\" style=\"width: 200px; height: 200px;\">"
	}
	body += "</html>"
	body += "</body>"

	contentType := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("From: Природа26 <" + username + ">\r\n" + "Subject: Природа26. Новое обращение! \r\n" + contentType + "\r\n\r\n" + body)

	err := smtp.SendMail(server+":25", auth, username, mailTo, msg)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(sql.MessageUpdateStatus, 1, m.Id)
}

func getEnvAsSlice(name string, sep string) []string {
	valStr := os.Getenv(name)
	val := strings.Split(valStr, sep)
	return val
}
