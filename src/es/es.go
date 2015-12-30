package es

import (
	. "github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/jeffail/gabs"
	"strconv"
	"strings"
)

func Search(c *gin.Context) {
	body, err := ReadBody(c)
	if err != nil {
		log.Error(err)
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("body parse json error: ", json)
		ReturnError(c, map[string]string{"error": "body parse json error"})
	}

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("can't find userid")
		ReturnError(c, map[string]string{"error": "can't find userid"})
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("can't find clustername")
		ReturnError(c, map[string]string{"error": "can't find clustername"})
	}

	var fields []string

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("can't find appname")
		ReturnError(c, map[string]string{"error": "can't find appname"})
	}

	fields = append(fields, `{"term":{"typename":"`+appname+`"}}`)

	hosts, err := json.S("hosts").Children()
	if err == nil && len(hosts) > 0 {
		fields = append(fields, `{"terms":{"ip": `+json.S("hosts").String()+`}}`)
	}

	start, ok := json.Path("start").Data().(string)
	if !ok {
		log.Error("can't find start time")
		ReturnError(c, map[string]string{"error": "can't find start time"})
	}

	end, ok := json.Path("end").Data().(string)
	if !ok {
		log.Error("can't find end time")
		ReturnError(c, map[string]string{"error": "can't find end time"})
	}

	fields = append(fields, `{"range":{"timestamp":{"gte":"`+start+`","lte":"`+end+`"}}}`)

	keyword, ok := json.Path("keyword").Data().(string)
	if ok {
		fields = append(fields, `{"match":{"msg":{"query":"`+keyword+`","analyzer":"ik"}}}`)
	}

	from, ok := json.Path("from").Data().(float64)
	if !ok {
		log.Error("can't find from")
		ReturnError(c, map[string]string{"error": "can't find from"})
	}

	size, ok := json.Path("size").Data().(float64)
	if !ok {
		log.Error("can't find size")
		ReturnError(c, map[string]string{"error": "can't find size"})
	}

	query := `{"fields":["timestamp","msg","ip","taskid"],
		"sort":{"timestamp":"asc"},"highlight":
		{"fields":{"msg":{"post_tags":["\u003c/em\u003e"],
		"pre_tags":["\u003cem style=\"color:red;\"\u003e"]}},
		"fragment_size":-1,"require_field_match":"true"},
		"query":{"bool":{"must":[` + strings.Join(fields, ",") + `]}},
		"from":` + strconv.Itoa(int(from)) + `,"size":` + strconv.Itoa(int(size)) + `}`
	index := "logstash-" + strconv.Itoa(int(userid)) + "-"
	ok, ymd := SameDay(start, end)
	if ok {
		index += ymd
	} else {
		index += "*"
	}
	estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname

	out, err := Conn.Search(index, estype, nil, query)
	content, _ := gabs.ParseJSON(out.RawJSON)

	ReturnOK(c, content.Data())
}

func SearchJump(c *gin.Context) {
	body, err := ReadBody(c)
	if err != nil {
		log.Error(err)
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("body parse json error: ", json)
		ReturnError(c, map[string]string{"error": "body parse json error"})
	}

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("can't find userid")
		ReturnError(c, map[string]string{"error": "can't find userid"})
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("can't find clustername")
		ReturnError(c, map[string]string{"error": "can't find clustername"})
	}

	var fields []string

	taskid, ok := json.Path("taskid").Data().(string)
	if !ok {
		log.Error("can't find taskid")
		ReturnError(c, map[string]string{"error": "can't find taskid"})
	}
	fields = append(fields, `{"term":{"taskid":"`+taskid+`"}}`)

	start, ok := json.Path("start").Data().(string)
	if !ok {
		log.Error("can't find start time")
		ReturnError(c, map[string]string{"error": "can't find start time"})
	}

	end, ok := json.Path("end").Data().(string)
	if !ok {
		log.Error("can't find end time")
		ReturnError(c, map[string]string{"error": "can't find end time"})
	}
	fields = append(fields, `{"range":{"timestamp":{"gte":"`+start+`","lte":"`+end+`"}}}`)

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("can't find appname")
		ReturnError(c, map[string]string{"error": "can't find appname"})
	}

	from, ok := json.Path("from").Data().(float64)
	if !ok {
		log.Error("can't find from")
		ReturnError(c, map[string]string{"error": "can't find from"})
	}

	size, ok := json.Path("size").Data().(float64)
	if !ok {
		log.Error("can't find size")
		ReturnError(c, map[string]string{"error": "can't find size"})
	}

	query := `{"fields":["timestamp","msg","ip","taskid"],
                "sort":{"timestamp":"asc"},"highlight":
                {"fields":{"msg":{"post_tags":["\u003c/em\u003e"],
                "pre_tags":["\u003cem style=\"color:red;\"\u003e"]}},
                "fragment_size":-1,"require_field_match":"true"},
                "query":{"bool":{"must":[` + strings.Join(fields, ",") + `]}},
                "from":` + strconv.Itoa(int(from)) + `,"size":` + strconv.Itoa(int(size)) + `}`

	index := "logstash-" + strconv.Itoa(int(userid)) + "-"
	ok, ymd := SameDay(start, end)
	if ok {
		index += ymd
	} else {
		index += "*"
	}
	estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname

	out, err := Conn.Search(index, estype, nil, query)
	content, _ := gabs.ParseJSON(out.RawJSON)

	ReturnOK(c, content.Data())
}
