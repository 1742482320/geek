package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "c", "./config.toml", "-c config.toml")
}

func main() {
	flag.Parse()

	InitConfig(confPath)

	log.Println(Conf)

	if err := do(); err != nil {
		panic(err)
	}

	cron := cron.New()
	cron.AddFunc(Conf.CronEntry, func() {
		if err := do(); err != nil {
			log.Println(err)
		}
	})
	cron.Start()

	go func() {
		if err := StartHTTP(Conf); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cron.Stop()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func do() error {

	lines, err := downloadAll(Conf)
	if err != nil {
		return err
	}

	if err := updateIndex(Conf.DataDir); err != nil {
		return err
	}

	user := fmt.Sprintf("uc%d", rand.Intn(60))
	pass := RandStringBytesMaskImprSrc(6)

	body := new(bytes.Buffer)
	body.WriteString(user)
	body.WriteString(":")
	body.WriteString(pass)
	body.WriteString("\n")

	body.WriteString(strings.Join(lines, "\n"))

	// SendToMail
	if err := SendToMail("大婶，还学得动吗？", body.String()); err != nil {
		return errors.Wrap(err, "SendToMail")
	}

	log.Println("pass:", user, pass)

	Conf.HTTP.BasicAuth = []string{user + ":" + pass}

	if err := ioutil.WriteFile("./auth", []byte(user+":"+pass), os.ModePerm); err != nil {
		return errors.Wrap(err, "WriteFile")
	}
	return nil

}

// IndexNode IndexNode
type IndexNode struct {
	Text       string       `json:"text"`
	Href       string       `json:"href"`
	Nodes      []*IndexNode `json:"nodes"`
	Selectable bool         `json:"selectable"`
}

func updateIndex(dir string) error {

	root := &IndexNode{
		Text:  "/",
		Href:  "/",
		Nodes: []*IndexNode{},
	}

	dirs := map[string]*IndexNode{
		"/": root,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if path == dir {
			return nil
		}

		dir = strings.TrimLeft(dir, "./")
		dir = strings.TrimRight(dir, "/")
		// fmt.Println(path, filepath.Base(dir))
		cpath := strings.Replace(path, dir, "", 1)

		switch {
		case strings.HasPrefix(cpath, "/."):
			return nil
		case strings.HasPrefix(cpath, "/js"):
			return nil
		case strings.HasPrefix(cpath, "/css"):
			return nil
		case strings.HasPrefix(cpath, "/src"):
			return nil
		case strings.HasPrefix(cpath, "/index.html"):
			return nil
		case strings.HasPrefix(cpath, "/data.js"):
			return nil
		}

		// 父级名称
		dirName := filepath.Dir(cpath)

		log.Println("dir:", dir, "path:", path, "cpath:", cpath, "dirname:", dirName)

		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".html") {
			return nil
		}

		var parent *IndexNode
		var has bool

		if len(cpath) == 0 {
			return nil
		}

		parent, has = dirs[dirName]
		if !has {
			return errors.New("path not exist")
		}

		var (
			selectable bool
			href       string
		)

		if info.IsDir() {
			selectable = false
		} else {
			href = cpath
			selectable = true
		}

		// 前缀有没有，没有创建
		node := &IndexNode{
			Text:       filepath.Base(cpath),
			Href:       href,
			Nodes:      []*IndexNode{},
			Selectable: selectable,
		}
		parent.Nodes = append(parent.Nodes, node)

		if info.IsDir() {
			if _, has = dirs[cpath]; !has {
				dirs[cpath] = node
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	data, err := json.Marshal(root.Nodes)
	if err != nil {
		return err
	}

	l := len(data)

	buf := make([]byte, l+15)

	copy(buf, []byte("var IndexData="))
	copy(buf[14:], data)
	copy(buf[l+14:], []byte(";"))

	if err := ioutil.WriteFile(filepath.Join(dir, "data.js"), buf, os.ModePerm); err != nil {
		return errors.Wrap(err, "WriteFile")
	}

	return nil
}

func downloadAll(conf *Config) ([]string, error) {
	lines := []string{}
	for i := range conf.GeekUsers {
		res, err := doDownload(conf, conf.GeekUsers[i].User, conf.GeekUsers[i].Pass)
		if err != nil {
			return nil, err
		}

		lines = append(lines, res...)
	}
	return lines, nil
}

func doDownload(conf *Config, user, pass string) ([]string, error) {

	header := http.Header{}

	header.Add("origin", "https://account.geekbang.org")
	header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")

	geekCli := NewGeekClient(header)

	_, cookies, err := geekCli.Login(user, pass)
	if err != nil {
		return nil, err
	}

	sks := []string{}
	for _, c := range cookies {
		sks = append(sks, c.Name+"="+c.Value)
	}

	header.Add("cookie", strings.Join(sks, "; "))

	geekCli = NewGeekClient(geekCli.Header)

	products, err := geekCli.MyProducts()
	if err != nil {
		return nil, err
	}

	addHTML := []string{}

	for _, product := range products {

		log.Println("product.Title", product.ID)

		switch product.ID {
		case 1, 2, 3:
		default:
			log.Println("skip", product.Title)
			continue
		}

		for _, item := range product.List {

			log.Println("item.Title", item.Title)

			articleList, err := geekCli.ColumnArticlesAll(item.Extra.ColumnID)
			if err != nil {
				return nil, errors.Wrap(err, "ColumnArticlesAll")
			}

			log.Println("get articleList", len(articleList))
			// break

			for _, v := range articleList {

				item.Extra.ColumnTitle = strings.Replace(item.Extra.ColumnTitle, "/", "|", -1)
				v.ArticleTitle = strings.Replace(v.ArticleTitle, "/", "|", -1)

				dir := filepath.Join(conf.DataDir, product.Title, item.Extra.ColumnTitle, v.ArticleTitle)

				// check exists
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					os.MkdirAll(dir, os.ModePerm)
				} else {
					if !conf.Force {
						continue
					}

				}

				info, err := geekCli.ArticleInfo(v.ID)
				if err != nil {
					return nil, errors.Wrap(err, "ArticleInfo")
				}

				commentList, err := geekCli.ArticleCommentsAll(info.ID)
				if err != nil {
					return nil, errors.Wrap(err, "ArticleCommentsAll")
				}

				if conf.SaveJSON {
					err = SaveJSONInfo(filepath.Join(dir, "data.json"), v)
					if err != nil {
						return nil, errors.Wrap(err, "SaveJSONInfo")
					}
					infopath := filepath.Join(dir, "info.json")
					err = SaveJSONInfo(infopath, info)
					if err != nil {
						return nil, errors.Wrap(err, "SaveJSONInfo")
					}

					log.Println("write json ", infopath)

					comentPath := filepath.Join(dir, "commentList.json")
					err = SaveJSONInfo(comentPath, commentList)
					if err != nil {
						return nil, errors.Wrap(err, "SaveJSONInfo")
					}

					log.Println("write json ", comentPath)
				}

				if conf.SaveStatic {

					if strings.HasPrefix(info.AudioDownloadURL, "http") {

						// mp3
						mp3Path := filepath.Join(dir, v.ArticleTitle+".mp3")
						if _, err := os.Stat(mp3Path); os.IsNotExist(err) {
							mp3, err := geekCli.GetResource(info.AudioDownloadURL)
							if err != nil {
								return nil, errors.Wrap(err, "GetResource")
							}

							err = ioutil.WriteFile(mp3Path, mp3, os.ModePerm)
							if err != nil {
								return nil, errors.Wrap(err, "WriteFile")
							}

							log.Println("write mp3 ", mp3Path)
						}

						// info.AudioDownloadURL = "./" + filepath.Base(mp3Path)
					}

					if info.VideoMediaMap != nil {

						mp4Path := filepath.Join(dir, v.ArticleTitle+".mp4")

						if _, err := os.Stat(mp4Path); os.IsNotExist(err) {

							err = HLSdownload(info.VideoMediaMap.Hd.URL, mp4Path)
							if err != nil {
								return nil, errors.Wrap(err, "HLSdownload")
							}

							log.Println("write mp4 ", mp4Path)
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

				htmlPath := filepath.Join(dir, "index.html")
				if err := ioutil.WriteFile(htmlPath, []byte(html), os.ModePerm); err != nil {
					return nil, errors.Wrap(err, "WriteFile")
				}

				addHTML = append(addHTML, dir)
				log.Println("write html ", htmlPath)

			}

		}

	}
	return addHTML, nil
}

// SaveJSONInfo SaveJSONInfo
func SaveJSONInfo(fpath string, data interface{}) error {

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fpath, json, os.ModePerm)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc RandStringBytesMaskImprSrc
func RandStringBytesMaskImprSrc(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func genBaseAuth(user, pass string) *bytes.Buffer {
	b := new(bytes.Buffer)
	b.WriteString(user)
	b.WriteString(":")
	b.WriteString("{SHA}")
	b.WriteString(GetSha(pass))
	// b.WriteString("\n")
	return b
}

// GetSha GetSha
func GetSha(password string) string {
	s := sha1.New()
	s.Write([]byte(password))
	passwordSum := []byte(s.Sum(nil))
	return base64.StdEncoding.EncodeToString(passwordSum)
}

// SendToMail SendToMail
func SendToMail(subject, body string) error {
	if len(Conf.Emails) == 0 {
		return nil
	}
	hp := strings.Split(Conf.SMTP.Host, ":")
	auth := smtp.PlainAuth("", Conf.SMTP.User, Conf.SMTP.Pass, hp[0])

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", Conf.Emails[0], subject, body)

	return smtp.SendMail(Conf.SMTP.Host, auth, Conf.SMTP.User, Conf.Emails, []byte(msg))
}
