package es

import (
	"fmt"
	"github.com/Dataman-Cloud/omega-es/src/dao"
	. "github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/Jeffail/gabs"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func Health(c *gin.Context) {
	astart := time.Now().UnixNano()
	ma := map[string]interface{}{
		"status": 0,
	}
	start := time.Now().UnixNano()
	_, err := Conn.AllNodesInfo()
	mes := make(map[string]interface{})
	if err == nil {
		mes["status"] = 0
	} else {
		mes["status"] = 1
		ma["status"] = 1
	}
	end := time.Now().UnixNano()
	mes["time"] = (end - start) / 1000000

	mall := map[string]interface{}{
		"elasticsearch": mes,
	}

	start = time.Now().UnixNano()
	conn := Open()
	defer conn.Close()
	mredis := make(map[string]interface{})
	_, err = redis.String(conn.Do("PING"))
	if err == nil {
		mredis["status"] = 0
	} else {
		mredis["status"] = 1
		ma["status"] = 1
	}
	end = time.Now().UnixNano()
	mredis["time"] = (end - start) / 1000000
	mall["redis"] = mredis
	aend := time.Now().UnixNano()
	ma["time"] = (aend - astart) / 1000000
	mall["omegaEs"] = ma

	mmsql := make(map[string]interface{})
	mysqlstart := time.Now().UnixNano()
	if err = dao.Ping(); err == nil {
		mmsql["status"] = 0
	} else {
		mmsql["status"] = 1
	}
	mmsql["time"] = (time.Now().UnixNano() - mysqlstart) / 1000000
	mall["mysql"] = mmsql
	ReturnOKGin(c, mall)
	return
}

func SearchIndex(c *gin.Context) {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("searchindex can't get request body")
		ReturnParamError(c, err.Error())
		return
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("searchindex param parse json error")
		ReturnParamError(c, err.Error())
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		log.Error("searchindex can't get uid")
		ReturnParamError(c, "searchindex can't get uid")
		return
	}
	userid, err := strconv.ParseInt(uid.(string), 10, 64)
	if err != nil {
		log.Error("serachindex invalid token")
		ReturnParamError(c, err.Error())
		return
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found clusterid")
		ReturnParamError(c, "searchindex param can't found clusterid")
		return
	}

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("searchindex param can't found appname")
		ReturnParamError(c, "searchindex param can't found appname")
		return
	}

	start, ok := json.Path("start").Data().(string)
	if !ok {
		log.Error("searchindex param can't found starttime")
		ReturnParamError(c, "searchindex param can't found starttime")
		return
	}

	end, ok := json.Path("end").Data().(string)
	if !ok {
		log.Error("searchindex param can't found endtime")
		ReturnParamError(c, "searchindex param can't found endtime")
		return
	}

	from, ok := json.Path("from").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found from")
		ReturnParamError(c, "searchindex param can't found from")
		return
	}

	size, ok := json.Path("size").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found size")
		ReturnParamError(c, "searchindex param can't found size")
		return
	}
	if size > 200 {
		size = 200
	}

	ipport, err := json.Path("ipport").Children()
	source, serr := json.Path("source").Children()
	keyword, kok := json.Path("keyword").Data().(string)

	query := `{
		    "query": {
		      "bool": {
		        "must": [
  	                  {
	                    "range": {
                              "timestamp": {
                                "gte": "` + start + `",
                                "lte": "` + end + `"
                              }
                            }
		          },
			  {
                            "term": {"typename": "` + appname + `"}
			  },
			  {
                            "term": {"clusterid": "` + strconv.Itoa(int(clusterid)) + `"}
			  }`
	if kok {
		query += `,
		          {
		            "match": {
		              "msg": {
			        "query": "` + keyword + `",
                                "analyzer": "ik"
			      }
                            }
			  }`
	}
	if err == nil && len(ipport) > 0 {
		var arr []string
		for _, ipp := range ipport {
			arr = append(arr, ipp.Data().(string))
		}
		query += `,
			  {
			    "terms": {
			      "ipport": ["` + strings.Join(arr, "\",\"") + `"]
			  }
			}`
	}
	if serr == nil && len(source) > 0 {
		var arr []string
		for _, sour := range source {
			arr = append(arr, sour.Data().(string))
		}
		query += `,
			  {
			    "terms": {
			      "source": ["` + strings.Join(arr, "\",\"") + `"]
			  }
			}`
	}
	query += `
		      ]
		    }
		  },
		"sort": {"timestamp.sort": "asc"},
		"from": ` + strconv.Itoa(int(from)) + `,
		"size": ` + strconv.Itoa(int(size)) + `,
		"fields": ["timestamp","msg","ipport","ip","taskid","counter", "typename", "source"],
		"highlight": {
	          "require_field_match": "true",
		  "fields": {
		    "msg": {
		      "pre_tags": ["<em style=\"color:red;\">"],
		      "post_tags": ["</em>"]
	            }
	          },
		  "fragment_size": -1
		}
	       }`
	esindex := "logstash-*" + fmt.Sprintf("%d", userid) + "-"
	estype := ""
	if start[:10] == end[:10] {
		esindex += start[:10]
		//estype = "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
	} else {
		esindex += "*"
	}
	esindex = "*"
	log.Debug(esindex, estype, query)
	out, err := Conn.Search(esindex, estype, nil, query)
	if err != nil {
		log.Error("searchindex search es error: ", err)
		ReturnEsError(c, err.Error())
		return
	}
	content, _ := gabs.ParseJSON(out.RawJSON)
	hits, err := content.Path("hits.hits").Children()
	if err != nil {
		log.Error("searchindex get hits error: ", err)
		ReturnEsError(c, err.Error())
		return
	}
	if err == nil {
		if len(hits) > 0 {
			for i, hit := range hits {
				msgs, err := hit.Path("fields.msg").Children()
				if err == nil {
					msg := msgs[0].Data().(string)
					msg = ReplaceHtml(msg)
					hits[i].Path("fields.msg").SetIndex(msg, 0)
				}
				msgh, err := hit.Path("highlight.msg").Children()
				if err == nil {
					msg := msgh[0].Data().(string)
					msg = strings.Replace(msg, "<em style=\"color:red;\">", "-emstart-", -1)
					msg = strings.Replace(msg, "</em>", "-emend-", -1)
					msg = ReplaceHtml(msg)
					msg = strings.Replace(msg, "-emstart-", "<em style=\"color:red;\">", -1)
					msg = strings.Replace(msg, "-emend-", "</em>", -1)
					hits[i].Path("highlight.msg").SetIndex(msg, 0)
					log.Debug("------:", msg)
				}
			}
		}
	}
	ReturnOKGin(c, content.Data())
	return
}

func SearchContext(c *gin.Context) {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("searchcontext can't get request body")
		ReturnParamError(c, err.Error())
		return
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("searchcontext param parse json error")
		ReturnParamError(c, err.Error())
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		log.Error("searchcontext can't get uid")
		ReturnParamError(c, "searchcontext can't get uid")
		return
	}
	userid, err := strconv.ParseInt(uid.(string), 10, 64)
	if err != nil {
		log.Error("serachcontext invalid token")
		ReturnParamError(c, err.Error())
		return
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found clusterid")
		ReturnParamError(c, "searchcontext can't found clusterid")
		return
	}

	ipport, ok := json.Path("ipport").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found ipport")
		ReturnParamError(c, "searchcontext can't found ipport")
		return
	}

	source, sok := json.Path("source").Data().(string)
	if !sok {
		log.Error("searchcontext param can't found source")
	}

	timestamp, ok := json.Path("timestamp").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found timestamp")
		ReturnParamError(c, "searchcontext can't found timestamp")
		return
	}

	counter, ok := json.Path("counter").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found counter")
		ReturnParamError(c, "searchcontext can't found counter")
		return
	}

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found appname")
		ReturnParamError(c, "searchcontext can't found appname")
		return
	}

	countergte := int(counter) - 100
	if countergte < 0 {
		countergte = 0
	}
	counterlte := int(counter) + 100

	query := `{
		    "query": {
		      "bool": {
	                "must": [
			  {
			    "range": {
			      "counter": {
			        "gte": ` + strconv.Itoa(countergte) + `,
				"lte": ` + strconv.Itoa(counterlte) + `
			      }
			    }
			  },
			  {
			    "term": {"typename": "` + appname + `"}
			  },
			  {
			    "term": {"clusterid": "` + strconv.Itoa(int(clusterid)) + `"}
			  },
			  `
	if sok {
		query += `
			  {
			    "term": {"source": "` + source + `"}
			  },
			  `
	}
	query += `
			  {
			    "term": {"ipport": "` + ipport + `"}
			  }
			]
		      }
		    },
		    "sort": {"timestamp.sort": "asc"},
		    "from": 0,
		    "size": 200,
		    "fields": ["timestamp","msg","ipport","ip","taskid","counter"],
		    "highlight": {
	              "require_field_match": "true",
		      "fields": {
		        "msg": {
	                  "pre_tags": ["<em style=\"color:red;\">"],
			  "post_tags": ["</em>"]
		        }
		      },
		      "fragment_size": -1
	            }
		  }`
	esindex := "logstash-*" + fmt.Sprintf("%d", userid) + "-" + timestamp[:10]
	//estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
	estype := ""
	esindex = "*"
	log.Debug(esindex, estype, query)
	out, err := Conn.Search(esindex, estype, nil, query)
	if err != nil {
		log.Error("searchcontext search es error: ", err)
		ReturnEsError(c, err.Error())
		return
	}
	content, _ := gabs.ParseJSON(out.RawJSON)
	hits, err := content.Path("hits.hits").Children()
	if err != nil {
		log.Error("searchindex get hits error: ", err)
		ReturnEsError(c, err.Error())
		return
	}
	if err == nil {
		if len(hits) > 0 {
			for i, hit := range hits {
				msgs, err := hit.Path("fields.msg").Children()
				if err == nil {
					msg := msgs[0].Data().(string)
					msg = ReplaceHtml(msg)
					hits[i].Path("fields.msg").SetIndex(msg, 0)
				}
				msgh, err := hit.Path("highlight.msg").Children()
				if err == nil {
					msg := msgh[0].Data().(string)
					msg = strings.Replace(msg, "<em style=\"color:red;\">", "-emstart-", -1)
					msg = strings.Replace(msg, "</em>", "-emend-", -1)
					msg = ReplaceHtml(msg)
					msg = strings.Replace(msg, "-emstart-", "<em style=\"color:red;\">", -1)
					msg = strings.Replace(msg, "-emend-", "</em>", -1)
					hits[i].Path("highlight.msg").SetIndex(msg, 0)
					log.Debug("------:", msg)
				}
			}
		}
	}
	ReturnOKGin(c, content.Data())
	return
}
