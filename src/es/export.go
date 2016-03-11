package es

import (
	"encoding/json"
	. "github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
	"github.com/jeffail/gabs"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
)

func ExportIndex(c *echo.Context) error {
	uid := c.Query("userid")
	cid := c.Query("clusterid")
	appname := c.Query("appname")
	start := c.Query("start")
	end := c.Query("end")
	ipport := strings.Split(c.Query("ipport"), ",")
	source := strings.Split(c.Query("source"), ",")
	keyword := c.Query("keyword")
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
	                            "term": {"clusterid": "` + cid + `"}
				  }`
	if keyword != "" {
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
	if c.Query("ipport") != "" {
		var arr []string
		for _, ipp := range ipport {
			arr = append(arr, ipp)
		}
		query += `,
				  {
				    "terms": {
				      "ipport": ["` + strings.Join(arr, "\",\"") + `"],
				      "minimum_match": 1
				  }
				}`
	}
	if c.Query("source") != "" {
		var arr []string
		for _, sour := range source {
			arr = append(arr, sour)
		}
		query += `,
				  {
				    "terms": {
				      "source": ["` + strings.Join(arr, "\",\"") + `"],
				      "minimum_match": 1
				  }
				}`
	}
	query += `
			      ]
			    }
			  },
			"sort": {"timestamp.sort": "asc"},
			"from": 0,
			"size": 5000,
			"fields": ["timestamp","ip","ipport","msg"],
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

	esindex := "logstash-*" + uid + "-"
	estype := ""
	if start[:10] == end[:10] {
		esindex += start[:10]
		//estype = "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
	} else {
		esindex += "*"
	}
	log.Debug(esindex, estype, query)
	out, err := Conn.Search(esindex, estype, nil, query)
	if err != nil {
		log.Error("searchindex search es error: ", err)
		ReturnError(c, map[string]interface{}{"error": "searchindex search es error"})
		return nil
	}
	content, err := gabs.ParseJSON(out.RawJSON)
	if err == nil {
		hits, _ := content.Path("hits.hits").Children()
		var rdata []map[string]interface{}
		for _, v := range hits {
			rdata = append(rdata, map[string]interface{}{
				"msg":       v.Path("fields.msg").Index(0).Data(),
				"timestamp": v.Path("fields.timestamp").Index(0).Data(),
				"ip":        v.Path("fields.ip").Index(0).Data(),
				"ipport":    v.Path("fields.ipport").Index(0).Data(),
			})
		}

		if len(rdata) == 0 {
			rdata = append(rdata, map[string]interface{}{
				"data": 0,
			})
		}
		b, _ := json.Marshal(rdata)
		c.Response().Header()["Content-Type"] = []string{"text/csv"}
		c.Response().Header()["Content-Disposition"] = []string{"attachment;filename=log.json"}
		c.Response().Write(b)
		return nil
		//ReturnOK(w, rdata)
		//return nil
	} else {
		log.Error("searchindex get hits error: ", err)
		ReturnError(c, map[string]interface{}{"error": "searchindex get hits error"})
		return nil
	}
	return nil
}

func IndexExport(w http.ResponseWriter, h *http.Request) {
	body, err := ReadBodyRequest(h)
	if err != nil {
		log.Error("searchindex can't get request body")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex can't get request body"})
		return
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("searchindex param parse json error")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex param parse json error"})
		return
	}

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found userid")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex can't found userid"})
		return
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found clusterid")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex can't found clusterid"})
		return
	}

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("searchindex param can't found appname")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex can't found appname"})
		return
	}

	start, ok := json.Path("start").Data().(string)
	if !ok {
		log.Error("searchindex param can't found starttime")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex can't found starttime"})
		return
	}

	end, ok := json.Path("end").Data().(string)
	if !ok {
		log.Error("searchindex param can't found endtime")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex can't found endtime"})
		return
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
			      "ipport": ["` + strings.Join(arr, "\",\"") + `"],
			      "minimum_match": 1
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
			      "source": ["` + strings.Join(arr, "\",\"") + `"],
			      "minimum_match": 1
			  }
			}`
	}
	query += `
		      ]
		    }
		  },
		"sort": {"timestamp.sort": "asc"},
		"from": 0,
		"size": 5000,
		"fields": ["timestamp","ip","ipport","msg"],
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
	esindex := "logstash-*" + strconv.Itoa(int(userid)) + "-"
	estype := ""
	if start[:10] == end[:10] {
		esindex += start[:10]
		//estype = "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
	} else {
		esindex += "*"
	}
	log.Debug(esindex, estype, query)
	out, err := Conn.Search(esindex, estype, nil, query)
	if err != nil {
		log.Error("searchindex search es error: ", err)
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex search es error"})
		return
	}
	content, err := gabs.ParseJSON(out.RawJSON)
	if err == nil {
		hits, _ := content.Path("hits.hits").Children()
		var rdata []map[string]interface{}
		for _, v := range hits {
			rdata = append(rdata, map[string]interface{}{
				"msg":       v.Path("fields.msg").Index(0).Data(),
				"timestamp": v.Path("fields.timestamp").Index(0).Data(),
				"ip":        v.Path("fields.ip").Index(0).Data(),
				"ipport":    v.Path("fields.ipport").Index(0).Data(),
			})
		}

		if len(rdata) == 0 {
			rdata = append(rdata, map[string]interface{}{
				"data": 0,
			})
		}
		ReturnOKResponse(w, rdata)
		return
	} else {
		log.Error("searchindex get hits error: ", err)
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex get hits error"})
		return
	}
	return
}

func ExportContext(c *echo.Context) error {
	uid := c.Query("userid")
	cid := c.Query("clusterid")
	ipport := c.Query("ipport")
	source := c.Query("source")
	timestamp := c.Query("timestamp")
	counter, err := strconv.Atoi(c.Query("counter"))
	if err != nil {
		log.Error("searchcontext counter strconv int error : ", err)
		ReturnError(c, map[string]interface{}{"error": "searchcontext counter strconv int error"})
		return nil
	}
	appname := c.Query("appname")
	log.Debug(uid, cid, ipport, source, counter, appname)

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
			    "term": {"ipport": "` + ipport + `"}
			  },
			  {
			    "term": {"source": "` + source + `"}
			  },
			  {
			    "term": {"clusterid": "` + cid + `"}
			  }
			]
		      }
		    },
		    "sort": {"timestamp.sort": "asc"},
		    "from": 0,
		    "size": 5000,
		    "fields": ["timestamp","ip","ipport","msg"],
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
	esindex := "logstash-*" + uid + "-" + timestamp[:10]
	estype := ""
	//estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
	out, err := Conn.Search(esindex, estype, nil, query)
	if err != nil {
		log.Error("searchindex search es error: ", err)
		ReturnError(c, map[string]interface{}{"error": "searchindex search es error"})
		return nil
	}
	content, _ := gabs.ParseJSON(out.RawJSON)
	if err == nil {
		hits, _ := content.Path("hits.hits").Children()
		var rdata []map[string]interface{}
		for _, v := range hits {
			rdata = append(rdata, map[string]interface{}{
				"msg":       v.Path("fields.msg").Index(0).Data(),
				"timestamp": v.Path("fields.timestamp").Index(0).Data(),
				"ip":        v.Path("fields.ip").Index(0).Data(),
				"ipport":    v.Path("fields.ipport").Index(0).Data(),
			})
		}
		if len(rdata) == 0 {
			rdata = append(rdata, map[string]interface{}{
				"data": 0,
			})
		}
		b, _ := json.Marshal(rdata)
		c.Response().Header()["Content-Type"] = []string{"text/csv"}
		c.Response().Header()["Content-Disposition"] = []string{"attachment;filename=log.json"}
		c.Response().Write(b)
		return nil
	} else {
		log.Error("searchindex get hits error: ", err)
		ReturnError(c, map[string]interface{}{"error": "searchindex get hits error"})
		return nil
	}
	return nil
}
func ContextExport(w http.ResponseWriter, h *http.Request) {
	body, err := ReadBodyRequest(h)
	if err != nil {
		log.Error("searchcontext can't get request body")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext can't get request body"})
		return
	}
	json, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("searchcontext param parse json error")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext param parse json error"})
		return
	}

	userid, ok := json.Path("userid").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found userid")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext can't found userid"})
		return
	}

	clusterid, ok := json.Path("clusterid").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found clusterid")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext can't found clusterid"})
		return
	}

	ipport, ok := json.Path("ipport").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found ipport")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext can't found ipport"})
		return
	}

	source, ok := json.Path("source").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found source")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext can't found source"})
		return
	}

	timestamp, ok := json.Path("timestamp").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found timestamp")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext can't found timestamp"})
		return
	}

	counter, ok := json.Path("counter").Data().(float64)
	if !ok {
		log.Error("searchcontext param can't found counter")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext can't found counter"})
		return
	}

	appname, ok := json.Path("appname").Data().(string)
	if !ok {
		log.Error("searchcontext param can't found appname")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchcontext can't found appname"})
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
			    "term": {"ipport": "` + ipport + `"}
			  },
			  {
			    "term": {"source": "` + source + `"}
			  },
			  {
			    "term": {"clusterid": "` + strconv.Itoa(int(clusterid)) + `"}
			  }
			]
		      }
		    },
		    "sort": {"timestamp.sort": "asc"},
		    "from": 0,
		    "size": 5000,
		    "fields": ["timestamp","ip","ipport","msg"],
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
	esindex := "logstash-*" + strconv.Itoa(int(userid)) + "-" + timestamp[:10]
	estype := ""
	//estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
	out, err := Conn.Search(esindex, estype, nil, query)
	if err != nil {
		log.Error("searchindex search es error: ", err)
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex search es error"})
		return
	}
	content, _ := gabs.ParseJSON(out.RawJSON)
	if err == nil {
		hits, _ := content.Path("hits.hits").Children()
		var rdata []map[string]interface{}
		for _, v := range hits {
			rdata = append(rdata, map[string]interface{}{
				"msg":       v.Path("fields.msg").Index(0).Data(),
				"timestamp": v.Path("fields.timestamp").Index(0).Data(),
				"ip":        v.Path("fields.ip").Index(0).Data(),
				"ipport":    v.Path("fields.ipport").Index(0).Data(),
			})
		}
		if len(rdata) == 0 {
			rdata = append(rdata, map[string]interface{}{
				"data": 0,
			})
		}
		ReturnOKResponse(w, rdata)
		return
	} else {
		log.Error("searchindex get hits error: ", err)
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex get hits error"})
		return
	}
	return
}
