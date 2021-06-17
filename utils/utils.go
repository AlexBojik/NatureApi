package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"image/jpeg"
	"net/http"
	"os"
	"time"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Created(w http.ResponseWriter, id int64) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

func GetTimeStamp() string {
	return time.Now().Format("2006.01.02 15:04:05 -0700")
}

func IfNotEmptyConcat(s1 string, s2 string) string {
	res := ""
	if len(s1) > 0 {
		res = s2 + s1
	}
	return res
}

func Base64toPng(data string) string {
	uid, _ := uuid.NewRandom()
	fileName := uid.String() + ".jpg"

	unbased, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := jpeg.Decode(r)
	if err != nil {
		panic("Bad jpeg")
	}

	f, err := os.OpenFile("./images/"+fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}

	opt := jpeg.Options{Quality: 100}
	jpeg.Encode(f, im, &opt)
	return fileName
}
