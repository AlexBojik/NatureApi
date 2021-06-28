package models

import (
	"encoding/base64"
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
	server := os.Getenv("MAIL_SERVER")
	mailTo := os.Getenv("MAIL_TO")
	imageSrc := os.Getenv("IMAGE_SRC")

	user := GetUser(m.Token)
	body := "<html>"
	body += "<body>"
	body += "<strong>Номер обращения:</strong> " + strconv.Itoa(int(m.Id)) + "<br>"
	body += "<strong>Время обращения:</strong> " + strconv.Itoa(int(m.Id)) + "<br>"
	body += "<strong>ФИО:</strong> " + user.Name + "<br>"
	body += "<strong>Текст обращения:</strong> " + m.Text + "<br>"
	body += "<strong>Координаты:</strong> " + fmt.Sprintf("%f", m.Lat) + ", " + fmt.Sprintf("%f", m.Lon) + "<br>"
	body += "<strong>Телефон:</strong> " + user.Phone + "<br>"
	body += "<strong>E-mail:</strong> " + user.Email + "<br>"
	body += "<strong>Адрес:</strong> " + user.ProAddr + "<br>"
	body += "<strong>Прикрепленные изображения:</strong><br>"
	for _, img := range m.Images {
		body += "<img src=\"" + imageSrc + img.Jpeg + "\" style=\"width: 200px; height: 200px;\">"
	}
	body += "</html>"
	body += "</body>"

	subject := "Природа26. Новое обращение!"

	var toMail []string
	toMail = strings.Split(mailTo, ",")

	err := sendMail(server+":25", username, subject, body, toMail)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(sql.MessageUpdateStatus, 1, m.Id)
}

func sendMail(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
