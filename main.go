package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
)

var rootJson = `
{
  "kind":"APIResourceList",
  "apiVersion":"v1",
  "groupVersion":"apis.jtthink.com/v1beta1",
  "resources":[
     {"name":"mypods","singularName":"mypod","shortNames":["mp"],"namespaced":true,"kind":"MyPod","verbs":["get","list"]}
  ]}
`
var podsList = `
{
  "kind": "MyPodList",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {},
  "items":[
    {
	  "metadata": {
        "name": "testpod1",
        "namespace": "default"
       }
    },
    {
	  "metadata": {
        "name": "testpod2",
        "namespace": "default"
       }
    }
   ]
}
`
var podDetail = `
{
  "kind": "MyPod",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {"name":"{name}","namespace":"{namespace}"},
  "spec":{"属性":"你懂的"}
}
`

func main() {

	r := gin.New()
	r.Use(func(c *gin.Context) {
		fmt.Println(c.Request.URL.Path)
		c.Next()
	})

	r.GET("/apis/apis.jtthink.com/v1beta1", func(c *gin.Context) {
		c.Header("content-type", "application/json")
		c.String(200, rootJson)
	})

	//列表  （根据ns)
	r.GET("/apis/apis.jtthink.com/v1beta1/namespaces/:ns/mypods", func(c *gin.Context) {
		c.Header("content-type", "application/json")
		json := strings.Replace(podsList, "default", c.Param("ns"), -1)
		c.String(200, json)
	})

	//列表  （所有 ) kb get mp -A
	r.GET("/apis/apis.jtthink.com/v1beta1/mypods", func(c *gin.Context) {
		c.Header("content-type", "application/json")
		json := strings.Replace(podsList, "default", "all", -1)
		c.String(200, json)
	})

	//详细 （根据ns)  kb get mp testpod1
	r.GET("/apis/apis.jtthink.com/v1beta1/namespaces/:ns/mypods/:name", func(c *gin.Context) {
		//c.Header("content-type", "application/json")
		//json := strings.Replace(podDetail, "{namespace}", c.Param("ns"), -1)
		//json = strings.Replace(json, "{name}", c.Param("name"), -1)
		//c.String(200, json)

		// 自定义字段
		t := metav1.Table{}
		t.Kind = "Table"
		t.APIVersion = "meta.k8s.io/v1"
		// 列
		t.ColumnDefinitions = []metav1.TableColumnDefinition{
			{Name: "name", Type: "string"},
			{Name: "命令空间", Type: "string"},
			{Name: "状态", Type: "string"},
		}

		// 内容
		t.Rows = []metav1.TableRow{
			{Cells: []interface{}{c.Param("name"), c.Param("ns"), "ready"}},
		}
		c.JSON(200, t)
	})

	//  8443  没有为啥
	if err := r.RunTLS(":8443",
		"certs/aaserver.crt", "certs/aaserver.key"); err != nil {
		log.Fatalln(err)
	}
}