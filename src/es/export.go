package es

import (
	. "github.com/Dataman-Cloud/omega-es/src/util"
	log "github.com/cihub/seelog"
	"github.com/jeffail/gabs"
	"net/http"
	"strconv"
	"strings"
)

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

	from, ok := json.Path("from").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found from")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex can't found from"})
		return
	}

	size, ok := json.Path("size").Data().(float64)
	if !ok {
		log.Error("searchindex param can't found size")
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex can't found size"})
		return
	}
	if size > 200 {
		size = 200
	}

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
		"from": 0,
		"size": 10000,
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
	esindex := "logstash-" + strconv.Itoa(int(userid)) + "-"
	if start[:10] == end[:10] {
		esindex += start[:10]
	} else {
		esindex += "*"
	}
	estype := ""
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
		ReturnOKResponse(w, rdata)
		return
	} else {
		log.Error("searchindex get hits error: ", err)
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex get hits error"})
		return
	}
	return
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
			    "term": {"ipport": "` + ipport + `"}
			  },
			  {
			    "term": {"typename": "` + appname + `"}
			  },
			  {
			    "term": {"clusterid": "` + strconv.Itoa(int(clusterid)) + `"}
			  }
			]
		      }
		    },
		    "sort": {"timestamp.sort": "asc"},
		    "from": 0,
		    "size": 10000,
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
	esindex := "logstash-" + strconv.Itoa(int(userid)) + "-" + timestamp[:10]
	estype := ""
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
		ReturnOKResponse(w, rdata)
		return
	} else {
		log.Error("searchindex get hits error: ", err)
		ReturnErrorResponse(w, map[string]interface{}{"error": "searchindex get hits error"})
		return
	}
	return
}
