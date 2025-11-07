package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"udesk/api/query"
	"udesk/server"
)

type NolinePer struct {
	Items []struct {
		Name             string `json:"name"`
		Agent_work_state string `json:"agent_work_state"`
	} `json:"items"`
}

var (
	per = map[int]string{
		1775184: "王杭渝(Adrian)",
		1775174: "孙凯(Max)",
		1775164: "李晓升(Mazuki)",
		425424:  "王奥(Gallagher)",
		401041:  "崔星林(Layne)",
		401031:  "王硕(Patt)",
		314001:  "李斌(Leven)",
	}
	OnPerId int
	Url     string
)

type Index struct {
	Agents []struct {
		Id        int64  `json:"id"`
		Nick_name string `json:"nick_name"`
	} `json:"agents"`
}

func UdeskApi() {

	var builder strings.Builder
	builder.WriteString("https://servicecenter-alauda.udesk.cn/open_api_v1/callcenter_analysis/agents_of_group?&page=1&perpage=30&")
	builder.WriteString(query.GetQuery())
	// fmt.Println(builder.String())
	resp, err := http.Get(builder.String())
	if err != nil {
		log.Fatalf("central get failed,err:%v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read from central resp.body failed, err:%v\n", err)
	}
	var nper = &NolinePer{}
	err = json.Unmarshal([]byte(string(body)), nper)
	if err != nil {
		log.Fatalf("json Unmarshal failed..%v\n", err)
	}
	for _, v := range nper.Items {
		for i, j := range per {
			if v.Agent_work_state == "idle" && v.Name == j {
				OnPerId = i
				break
			}
		}

	}
}

func ReAgent() {
	var builder strings.Builder
	builder.WriteString("https://servicecenter-alauda.udesk.cn/open_api_v1/tickets/")
	builder.WriteString(<-server.AgentId)
	builder.WriteString("?")
	builder.WriteString(query.GetQuery())

	data := struct {
		Ticket struct {
			Agent_id int `json:"agent_id"`
		} `json:"ticket"`
	}{
		Ticket: struct {
			Agent_id int "json:\"agent_id\""
		}{OnPerId},
	}
	Body, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("json marshal:%v\n", err)
	}
	fmt.Println(string(Body))
	req, err := http.NewRequest("PUT", builder.String(), bytes.NewBuffer(Body))
	if err != nil {
		log.Fatalf("put failed err,%v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
