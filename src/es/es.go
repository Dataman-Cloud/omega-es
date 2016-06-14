package es

import (
	"fmt"
	"github.com/Dataman-Cloud/omega-es/src/dao"
	. "github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/Jeffail/gabs"
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Health(c *echo.Context) error {
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
	return c.JSON(http.StatusOK, mall)
}

func SearchIndex(c *echo.Context) error {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("searchindex can't get request body")
		return ReturnError(c, map[string]string{"error": "searchindex can't get request body"})
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("searchindex param parse json error")
		return ReturnError(c, map[string]string{"error": "searchindex param parse json error"})
	}

	uid := c.Get("uid").(string)
	userid, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Error("serachindex invalid token")
		return ReturnError(c, map[string]string{"error": "searchindex inavlid token"})
	}
	/*userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found userid")
		return ReturnError(c, map[string]string{"error": "searchindex can't found userid"})
	}*/

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found clusterid")
		return ReturnError(c, map[string]string{"error": "searchindex can't found clusterid"})
	}

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("searchindex param can't found appname")
		return ReturnError(c, map[string]string{"error": "searchindex can't found appname"})
	}

	start, ok := json.Path("start").Data().(string)
	if !ok {
		log.Error("searchindex param can't found starttime")
		return ReturnError(c, map[string]string{"error": "searchindex can't found starttime"})
	}

	end, ok := json.Path("end").Data().(string)
	if !ok {
		log.Error("searchindex param can't found endtime")
		return ReturnError(c, map[string]string{"error": "searchindex can't found endtime"})
	}

	from, ok := json.Path("from").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found from")
		return ReturnError(c, map[string]string{"error": "searchindex can't found from"})
	}

	size, ok := json.Path("size").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found size")
		return ReturnError(c, map[string]string{"error": "searchindex can't found size"})
	}
	if size > 200 {
		size = 200
	}

	//ipport, iok := json.Path("ipport").Data().(string)
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
		return ReturnError(c, map[string]string{"error": "searchindex search es error"})
	}
	content, _ := gabs.ParseJSON(out.RawJSON)
	hits, err := content.Path("hits.hits").Children()
	if err != nil {
		log.Error("searchindex get hits error: ", err)
		return ReturnError(c, map[string]string{"error": "searchindex get hits error"})
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
	} else {
		log.Error("searchindex get hits error: ", err)
		return ReturnError(c, map[string]string{"error": "searchindex get hits error"})
	}
	return ReturnOK(c, content.Data())
}

func SearchContext(c *echo.Context) error {
	body, err := ReadBody(c)
	if err != nil {
		log.Error("searchcontext can't get request body")
		return ReturnError(c, map[string]string{"error": "searchcontext can't get request body"})
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("searchcontext param parse json error")
		return ReturnError(c, map[string]string{"error": "searchcontext param parse json error"})
	}

	uid := c.Get("uid").(string)
	userid, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Error("serachcontext invalid token")
		return ReturnError(c, map[string]string{"error": "searchcontext inavlid token"})
	}
	/*userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found userid")
		return ReturnError(c, map[string]string{"error": "searchcontext can't found userid"})
	}*/

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found clusterid")
		return ReturnError(c, map[string]string{"error": "searchcontext can't found clusterid"})
	}

	ipport, ok := json.Path("ipport").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found ipport")
		return ReturnError(c, map[string]string{"error": "searchcontext can't found ipport"})
	}

	source, sok := json.Path("source").Data().(string)
	if !sok {
		log.Error("searchcontext param can't found source")
		//return ReturnError(c, map[string]string{"error": "searchcontext can't found source"})
	}

	timestamp, ok := json.Path("timestamp").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found timestamp")
		return ReturnError(c, map[string]string{"error": "searchcontext can't found timestamp"})
	}

	counter, ok := json.Path("counter").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found counter")
		return ReturnError(c, map[string]string{"error": "searchcontext can't found counter"})
	}

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found appname")
		return ReturnError(c, map[string]string{"error": "searchcontext can't found appname"})
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
		log.Error("searchindex search es error: ", err)
		return ReturnError(c, map[string]string{"error": "searchindex search es error"})
	}
	content, _ := gabs.ParseJSON(out.RawJSON)
	hits, err := content.Path("hits.hits").Children()
	if err != nil {
		log.Error("searchindex get hits error: ", err)
		return ReturnError(c, map[string]string{"error": "searchindex get hits error"})
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
	} else {
		log.Error("searchindex get hits error: ", err)
		return ReturnError(c, map[string]string{"error": "searchindex get hits error"})
	}
	return ReturnOK(c, content.Data())
}
