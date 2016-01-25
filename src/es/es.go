package es

import (
	. "github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
	"github.com/jeffail/gabs"
	"github.com/labstack/echo"
	"strconv"
	"strings"
)

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

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found userid")
		return ReturnError(c, map[string]string{"error": "searchindex can't found userid"})
	}

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
			      "ipport": ["` + strings.Join(arr, ",") + `"],
			      "minimum_match": 1
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
		"fields": ["timestamp","msg","ipport","ip","taskid","counter", "typename"],
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
	esindex := "logstash-" + strconv.Itoa(int(userid)) + "-"
	if start[:10] == end[:10] {
		esindex += start[:10]
	} else {
		esindex += "*"
	}
	//estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
	estype := ""
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
					msg = strings.Replace(msg, "&", "&amp;", -1)
					msg = strings.Replace(msg, "<", "&lt;", -1)
					msg = strings.Replace(msg, ">", "&gt;", -1)
					msg = strings.Replace(msg, "\"", "&quot;", -1)
					msg = strings.Replace(msg, " ", "&nbsp;", -1)
					hits[i].Path("fields.msg").SetIndex(msg, 0)
				} else {
					continue
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

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found userid")
		return ReturnError(c, map[string]string{"error": "searchcontext can't found userid"})
	}

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
	esindex := "logstash-" + strconv.Itoa(int(userid)) + "-" + timestamp[:10]
	estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
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
					log.Debug(msg)
					msg = strings.Replace(msg, "&", "&amp;", -1)
					msg = strings.Replace(msg, "<", "&lt;", -1)
					msg = strings.Replace(msg, ">", "&gt;", -1)
					msg = strings.Replace(msg, "\"", "&quot;", -1)
					msg = strings.Replace(msg, " ", "&nbsp;", -1)
					hits[i].Path("fields.msg").SetIndex(msg, 0)
				} else {
					continue
				}
			}
		}
	} else {
		log.Error("searchindex get hits error: ", err)
		return ReturnError(c, map[string]string{"error": "searchindex get hits error"})
	}
	return ReturnOK(c, content.Data())
}
