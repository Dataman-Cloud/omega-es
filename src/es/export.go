package es

import (
	"encoding/json"
	. "github.com/Dataman-Cloud/omega-es/src/util"
	"github.com/Jeffail/gabs"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func ExportIndex(c *gin.Context) {
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
				      "ipport": ["` + strings.Join(arr, "\",\"") + `"]
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
				      "source": ["` + strings.Join(arr, "\",\"") + `"]
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
	esindex = "*"
	log.Debug(esindex, estype, query)
	out, err := Conn.Search(esindex, estype, nil, query)
	if err != nil {
		log.Error("searchindex search es error: ", err)
		ReturnEsError(c, err.Error())
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
		b, _ := json.Marshal(rdata)
		c.Writer.Header()["Content-Type"] = []string{"text/csv"}
		c.Writer.Header()["Content-Disposition"] = []string{"attachment;filename=log.json"}
		b, _ = FormatJson(b)
		c.Writer.Write(b)
		return
	}
	return
}

func ExportContext(c *gin.Context) {
	uid := c.Query("userid")
	cid := c.Query("clusterid")
	ipport := c.Query("ipport")
	source := c.Query("source")
	timestamp := c.Query("timestamp")
	counter, err := strconv.Atoi(c.Query("counter"))
	if err != nil {
		log.Error("searchcontext counter strconv int error : ", err)
		ReturnParamError(c, err.Error())
		return
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
			  },`
	if source != "" {
		query += `
			  {
			    "term": {"source": "` + source + `"}
			  },`
	}
	query += `
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
	esindex = "*"
	estype := ""
	//estype := "logstash-" + strconv.Itoa(int(clusterid)) + "-" + appname
	out, err := Conn.Search(esindex, estype, nil, query)
	if err != nil {
		log.Error("searchindex search es error: ", err)
		ReturnEsError(c, err.Error())
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
		b, _ := json.Marshal(rdata)
		b, _ = FormatJson(b)
		c.Writer.Header()["Content-Type"] = []string{"text/csv"}
		c.Writer.Header()["Content-Disposition"] = []string{"attachment;filename=log.json"}
		c.Writer.Write(b)
		return
	}
	return
}
