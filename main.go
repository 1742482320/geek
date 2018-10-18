package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	USER   string
	PASSWD string
)

func init() {
	flag.StringVar(&USER, "u", "", "-u user")
	flag.StringVar(&PASSWD, "p", "", "-p password")
}

func main() {
	flag.Parse()

	if len(USER) == 0 || len(PASSWD) == 0 {
		panic("请输入用户，密码")
	}

	header := http.Header{}

	header.Add("origin", "https://account.geekbang.org")
	header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")

	geekCli := NewGeekClient(header)

	_, cookies, err := geekCli.Login(USER, PASSWD)
	if err != nil {

		panic(err)
	}

	// log.Println(info)

	sks := []string{}
	for _, c := range cookies {
		sks = append(sks, c.Name+"="+c.Value)
	}

	header.Add("cookie", strings.Join(sks, "; "))

	geekCli = NewGeekClient(geekCli.Header)

	// cookies, err := LoginByChrome(USER, PASSWD)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(products.Data)
	do(geekCli)
}

func do(geekCli *GeekClient) {

	products, err := geekCli.MyProducts()
	if err != nil {
		panic(err)
	}

	articleIDs := []int{}

	for _, item := range products[0].List {
		articleList, err := geekCli.ColumnArticlesAll(item.Extra.ColumnID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("get articleList", len(articleList))
		// break

		for _, v := range articleList {
			info, err := geekCli.ArticleInfo(v.ID)
			if err != nil {
				log.Fatal(err)
			}

			// fmt.Println("get info", info)

			articleIDs = append(articleIDs, info.ID)

			item.Extra.ColumnTitle = strings.Replace(item.Extra.ColumnTitle, "/", "|", -1)
			info.ArticleTitle = strings.Replace(info.ArticleTitle, "/", "|", -1)

			dir := fmt.Sprintf("./data/columns/%s/%s", item.Extra.ColumnTitle, info.ArticleTitle)
			os.MkdirAll(dir, os.ModePerm)

			// err = SaveJSON(filepath.Join(dir, info.ArticleTitle+".json"), info)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			commentList, err := geekCli.ArticleCommentsAll(info.ID)
			if err != nil {
				log.Fatal(err)
			}

			// err = SaveJSON(filepath.Join(dir, "commentList.json"), commentList)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			html := TplArticleHTML(info, commentList)
			ioutil.WriteFile(filepath.Join(dir, info.ArticleTitle+".html"), []byte(html), os.ModePerm)

			// mp3, err := geekCli.GetResource(info.AudioDownloadURL)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// // mp3
			// err = ioutil.WriteFile(filepath.Join(dir, info.ArticleTitle+".mp3"), mp3, os.ModePerm)
			// if err != nil {
			// 	log.Fatal(err)
			// }

		}
	}

}

// SaveJSON SaveJSON
func SaveJSON(fpath string, data interface{}) error {

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fpath, json, os.ModePerm)
}
