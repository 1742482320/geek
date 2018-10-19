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
	USER       string
	PASSWD     string
	SaveStatic bool
	SaveJSON   bool
	Force      bool
)

func init() {
	flag.StringVar(&USER, "u", "", "-u user")
	flag.StringVar(&PASSWD, "p", "", "-p password")
	flag.BoolVar(&SaveStatic, "save", false, "-save 存资源文件")
	flag.BoolVar(&SaveJSON, "json", false, "-json 存json")
	flag.BoolVar(&Force, "f", false, "-f 强制覆盖")
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

	// // log.Println(info)

	sks := []string{}
	for _, c := range cookies {
		sks = append(sks, c.Name+"="+c.Value)
	}

	// sks = append(sks, "GCESS=BAUEAAAAAAkBAQgBAwIEJ0DJWwoEAAAAAAEEu6APAAYEwnL9kAwBAQsCBAAEBIBRAQADBCdAyVsHBJ2WMgs-")
	// sks = append(sks, "GCID=7bff37f-22bec2f-2fbe82a-0203dde")
	// sks = append(sks, "SERVERID=fe79ab1762e8fabea8cbf989406ba8f4|1539915816|1539915647")
	// sks = append(sks, "_ga=GA1.2.1560661455.1535373390")
	// sks = append(sks, "_gid=GA1.2.1481632964.1539915651")

	header.Add("cookie", strings.Join(sks, "; "))

	geekCli = NewGeekClient(geekCli.Header)
	// geekCli := NewGeekClient(header)

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

	// articleIDs := []int{}

	for _, product := range products {

		fmt.Println("product.Title", product.ID)

		switch product.ID {
		case 1, 2, 3:
		default:
			fmt.Println("skip", product.Title)
			continue
		}

		for _, item := range product.List {

			fmt.Println("item.Title", item.Title)

			articleList, err := geekCli.ColumnArticlesAll(item.Extra.ColumnID)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("get articleList", len(articleList))
			// break

			for _, v := range articleList {

				item.Extra.ColumnTitle = strings.Replace(item.Extra.ColumnTitle, "/", "|", -1)
				v.ArticleTitle = strings.Replace(v.ArticleTitle, "/", "|", -1)

				dir := fmt.Sprintf("./data/%s/%s/%s", product.Title, item.Extra.ColumnTitle, v.ArticleTitle)

				// check exists
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					os.MkdirAll(dir, os.ModePerm)
				} else {
					if !Force {
						continue
					}

				}

				info, err := geekCli.ArticleInfo(v.ID)
				if err != nil {
					log.Fatal(err)
				}

				commentList, err := geekCli.ArticleCommentsAll(info.ID)
				if err != nil {
					log.Fatal(err)
				}

				if SaveJSON {
					err = SaveJSONInfo(filepath.Join(dir, "data.json"), v)
					if err != nil {
						log.Fatal(err)
					}
					infopath := filepath.Join(dir, "info.json")
					err = SaveJSONInfo(infopath, info)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println("write json ", infopath)

					comentPath := filepath.Join(dir, "commentList.json")
					err = SaveJSONInfo(comentPath, commentList)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println("write json ", comentPath)
				}

				if SaveStatic {

					if strings.HasPrefix(info.AudioDownloadURL, "http") {

						// mp3
						mp3Path := filepath.Join(dir, info.ArticleTitle+".mp3")
						if _, err := os.Stat(mp3Path); os.IsNotExist(err) {
							mp3, err := geekCli.GetResource(info.AudioDownloadURL)
							if err != nil {
								log.Fatal(err)
							}

							err = ioutil.WriteFile(mp3Path, mp3, os.ModePerm)
							if err != nil {
								log.Fatal(err)
							}

							fmt.Println("write mp3 ", mp3Path)
						}

						// info.AudioDownloadURL = "./" + filepath.Base(mp3Path)
					}

					if info.VideoMediaMap != nil {

						mp4Path := filepath.Join(dir, v.ArticleTitle+".mp4")

						if _, err := os.Stat(mp4Path); os.IsNotExist(err) {

							err = HLSdownload(info.VideoMediaMap.Hd.URL, mp4Path)
							if err != nil {
								log.Fatal(err)
							}

							fmt.Println("write mp4 ", mp4Path)
						}

						// info.VideoMediaMap.Hd.URL = "./" + filepath.Base(mp4Path)
					}
				}

				var html string
				if info.VideoMediaMap != nil {
					html = TplVideoHTML(info, commentList)
				} else {
					html = TplArticleHTML(info, commentList)
				}

				htmlPath := filepath.Join(dir, info.ArticleTitle+".html")
				ioutil.WriteFile(htmlPath, []byte(html), os.ModePerm)

				fmt.Println("write html ", htmlPath)

				break
			}
			break
		}

	}
}

// SaveJSONInfo SaveJSONInfo
func SaveJSONInfo(fpath string, data interface{}) error {

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fpath, json, os.ModePerm)
}
