package esia

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"water-api/models"
	"water-api/utils"

	"github.com/google/uuid"
)

type Params struct {
	Id          string
	State       string
	Scope       string
	TimeStamp   string
	Secret      string
	AccessToken string  `json:"access_token"`
	OpenId      float64 `json:"urn:esia:sbj_id"`
	RrnsUrl     string
}

func (p *Params) New() {
	p.Id = os.Getenv("ESIA_ID")
	p.State = uuid.New().String()
	p.Scope = os.Getenv("ESIA_SCOPE")
	p.TimeStamp = utils.GetTimeStamp()
	p.Secret = crypt(p)
}

func (p *Params) getToken(code string, mobile bool) {
	params := getBaseUrlParams(mobile)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	params.Set("token_type", "Bearer")

	client := &http.Client{}

	encodedParams := params.Encode()
	r, _ := http.NewRequest("POST", os.Getenv("ESIA_TOKEN_URL"), strings.NewReader(encodedParams))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(encodedParams)))

	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &p)
	if err != nil {
		return
	}

	chunk := strings.Split(p.AccessToken, ".")
	if len(chunk) < 2 {
		return
	}
	decoded, _ := base64.RawURLEncoding.DecodeString(chunk[1])
	_ = json.Unmarshal(decoded, &p)
	p.RrnsUrl = os.Getenv("ESIA_PRNS_URL") + fmt.Sprintf("%.0f", p.OpenId)
}

func GetAuthUrl(mobile bool) string {
	authUrl := os.Getenv("ESIA_AUTH_URL") + "?" + getBaseUrlParams(mobile).Encode()
	return authUrl
}

func GetRedirectUrl(r *http.Request, mobile bool) string {
	redirectUrl := os.Getenv("REDIRECT_URL")
	if !mobile {
		redirectUrl = os.Getenv("WEB_REDIRECT_URL")
	}

	p := Params{}
	p.New()
	p.getToken(getCodeFromUrl(r), mobile)

	fmt.Println("user with openid:", p.OpenId)

	if p.OpenId > 0 {
		user := models.User{}
		user.OpenId = p.OpenId
		getPersonInfo(&user, &p)
		getContactInfo(&user, &p)
		getDocumentInfo(&user, &p)
		getAddressInfo(&user, &p)

		fmt.Println("Token: ", user.Token)

		if len(user.Token) > 0 {
			fmt.Println("Create/renew user")
			models.CreateUser(&user)
			redirectUrl += "?t=" + user.Token
		}
	}
	return redirectUrl
}

func getAddressInfo(user *models.User, p *Params) {
	e := models.Elements{}
	doRequest(p.RrnsUrl+"/addrs", p.AccessToken, &e)

	for _, addr := range e.Elements {
		ad := models.Address{}
		doRequest(addr, p.AccessToken, &ad)

		value := ad.ZipCode + ", " + ad.AddressStr
		value += utils.IfNotEmptyConcat(ad.Building, ", стр ")
		value += utils.IfNotEmptyConcat(ad.Frame, ", к ")
		value += utils.IfNotEmptyConcat(ad.House, ", д ")
		value += utils.IfNotEmptyConcat(ad.Flat, ", кв ")
		if ad.Type == "PLV" {
			user.ProAddr = value
		} else if ad.Type == "PRG" {
			user.RegAddr = value
		}
	}
}

func getDocumentInfo(user *models.User, p *Params) {
	e := models.Elements{}
	doRequest(p.RrnsUrl+"/docs", p.AccessToken, &e)

	for _, addr := range e.Elements {
		doc := models.Documents{}
		doRequest(addr, p.AccessToken, &doc)
		if doc.Type == "RF_PASSPORT" {
			user.Doc = doc.Series + " " + doc.Number + ", выдан: " + doc.IssuedBy
			user.Doc += "(" + doc.IssueId + "), дата выдачи:" + doc.IssueDate
		}
	}
}

func getContactInfo(user *models.User, p *Params) {
	e := models.Elements{}
	doRequest(p.RrnsUrl+"/ctts", p.AccessToken, &e)

	for _, ctt := range e.Elements {
		ct := models.Contact{}
		doRequest(ctt, p.AccessToken, &ct)

		if ct.Type == "MBT" {
			user.Phone = ct.Value
		} else if ct.Type == "EML" {
			user.Email = ct.Value
		}
	}
}

func getPersonInfo(user *models.User, p *Params) {
	if p.OpenId == 0 {
		return
	}

	info := models.Info{}
	doRequest(p.RrnsUrl, p.AccessToken, &info)
	user.Name = strings.Join([]string{info.LastName, info.FirstName, info.MiddleName}, " ")
	user.Token = info.ETag
	user.Snils = info.Snils
}

func getCodeFromUrl(r *http.Request) string {
	code := ""
	codes, ok := r.URL.Query()["code"]
	if ok && len(codes) > 0 {
		code = codes[0]
	}

	return code
}

func getBaseUrlParams(mobile bool) url.Values {
	p := Params{}
	p.New()

	redirectUrl := os.Getenv("ESIA_REDIRECT")
	if !mobile {
		redirectUrl = os.Getenv("ESIA_WEB_REDIRECT")
	}

	urlParams := url.Values{}
	urlParams.Set("redirect_uri", redirectUrl)
	urlParams.Set("response_type", "code")
	urlParams.Set("access_type", "offline")
	urlParams.Set("client_id", p.Id)
	urlParams.Set("scope", p.Scope)
	urlParams.Set("state", p.State)
	urlParams.Set("client_secret", p.Secret)
	urlParams.Set("timestamp", p.TimeStamp)

	return urlParams
}

func doRequest(url string, token string, js interface{}) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &js)
}

func crypt(e *Params) string {
	body := "body=" + url.QueryEscape(e.Scope+e.TimeStamp+e.Id+e.State)
	resp, _ := http.Post(os.Getenv("CRYPTO_URL"), "text", strings.NewReader(body))
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return strings.ReplaceAll(string(data), "\n", "")
}
