package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
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

	cookies, err := LoginByChrome(USER, PASSWD)
	if err != nil {
		panic(err)
	}

	// fmt.Println(products.Data)
	do(cookies)
}

// LoginByChrome LoginByChrome
func LoginByChrome(user, pass string) ([]*network.Cookie, error) {
	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf), chromedp.WithRunnerOptions(
		runner.Headless,
		runner.DisableGPU))
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// run task list
	var cookies []*network.Cookie

	err = c.Run(ctxt, geekLogin(user, pass, &cookies))
	if err != nil {
		return nil, err
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		return nil, err
	}

	// wait for chrome to finish
	c.Wait()
	// if err != nil {
	// 	return nil, err
	// }

	return cookies, nil
}

func geekLogin(user, pass string, cookies *[]*network.Cookie) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(`https://account.geekbang.org/signin`),
		chromedp.WaitVisible(`.nw-phone-container`, chromedp.ByQuery),
		chromedp.SendKeys(".nw-phone-wrap input", user, chromedp.ByQuery),
		chromedp.SendKeys(".input-wrap input", pass, chromedp.ByQuery),
		chromedp.Click(".mybtn", chromedp.ByQuery),
		// chromedp.Sleep(2 * time.Second), // wait for animation to finish
		chromedp.WaitVisible(".account-sidebar", chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {

			cs, err := network.GetCookies().Do(ctx, h)
			if err != nil {
				return err
			}

			*cookies = cs

			return nil
		}),
	}

}

func do(cookies []*network.Cookie) {
	header := http.Header{}
	sks := []string{}
	for _, c := range cookies {
		sks = append(sks, c.Name+"="+c.Value)
	}

	header.Add("cookie", strings.Join(sks, "; "))
	header.Add("origin", "https://time.geekbang.org")
	header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")

	fmt.Println("header", header)

	geekCli := NewGeekClient(header)
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
