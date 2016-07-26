package es

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/Jeffail/gabs"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

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
	esindex := "dataman-app-" + fmt.Sprintf("%d", clusterid) + "-"
	estype := "dataman-" + appname
	//estype := ""
	if start[:10] == end[:10] {
		esindex += start[:10]
	} else {
		esindex += "*"
	}
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
	esindex := "dataman-app-" + fmt.Sprintf("%d", clusterid) + "-" + timestamp[:10]
	estype := "dataman-" + appname
	//esindex := "*"
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
