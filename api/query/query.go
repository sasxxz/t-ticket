package query

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Query struct {
	Email        string `json:"email"`
	TimeStamp    int64  `json:"timestamp"`
	Sign         string `json:"sign"`
	Nonce        string `json:"nonce"`
	Sign_version string `json:"sign_version"`
}

func GetQuery() string {
	nonce := uuid.New().String()
	url := "https://servicecenter-alauda.udesk.cn/open_api_v1/log_in"
	contentType := "application/json"
	date := `{"email": "hbqi@alauda.io","password": "Ye_qiu@123"}`
	resp, err := http.Post(url, contentType, strings.NewReader(date))
	if err != nil {
		log.Fatalf("post failed, err:%v\n", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("get resp failed, err:%v\n", err)
	}
	var btmp struct {
		UdeskToken string `json:"open_api_auth_token"`
	}
	err = json.Unmarshal([]byte(string(b)), &btmp)
	if err != nil {
		log.Fatalf("json Unmarshal failed!")
	}
	tmpsign := fmt.Sprintf("hbqi@alauda.io&%v&%v&%v&v2", btmp.UdeskToken, time.Now().Unix(), nonce)
	sign := fmt.Sprintf("%x", sha256.Sum256([]byte(tmpsign)))
	var apiquery = &Query{
		Email:        "hbqi@alauda.io",
		TimeStamp:    time.Now().Unix(),
		Sign:         sign,
		Nonce:        nonce,
		Sign_version: "v2",
	}
	var builder strings.Builder
	builder.WriteString("email=")
	builder.WriteString(apiquery.Email)
	builder.WriteString("&timestamp=")
	builder.WriteString(strconv.Itoa(int(apiquery.TimeStamp)))
	builder.WriteString("&sign=")
	builder.WriteString(apiquery.Sign)
	builder.WriteString("&nonce=")
	builder.WriteString(apiquery.Nonce)
	builder.WriteString("&sign_version=")
	builder.WriteString(apiquery.Sign_version)
	builder.WriteString("&id=1310155774")
	return builder.String()
}
