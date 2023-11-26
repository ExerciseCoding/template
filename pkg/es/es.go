package es

import (
	"context"
	"errors"
	"log"
	"os"
	"github.com/gogf/gf/frame/g"
	"github.com/olivere/elastic/v7"
	"time"
)

func connection() (search *elastic.Client, err error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL("http://" + g.Cfg().Get("Es.Url").(string)),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(g.Cfg().Get("Es.UserName").(string), g.Cfg().Get("Es.Password").(string)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Printf("Es 链接失败 error:【%v】", err)
	}
	return client, err
}


type GetDocReq struct {
	DocId     string `json:"doc_id"`
	IndexName string `json:"index_name"`
}

func GetDoc(req GetDocReq) (docContent *elastic.GetResult, err error) {
	client, err := connection()
	if err != nil || client == nil {
		err = errors.New("elastic连接失败:" + err.Error())
		return
	}
	defer client.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	docContent, err = client.Get().Index(req.IndexName).Id(req.DocId).Do(ctx)
	if err != nil {
		err = errors.New("获取文档失败,错误原因：" + err.Error())
		log.Printf("获取文档失败 错误原因：【%s】", err.Error())
		return
	}
	return
}
