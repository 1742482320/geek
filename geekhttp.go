package main

import (
	"bytes"
	json "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// GeekClient GeekClient
type GeekClient struct {
	Header http.Header
	url    string
}

// NewGeekClient NewGeekClient
func NewGeekClient(h http.Header) *GeekClient {
	return &GeekClient{
		url:    "https://time.geekbang.org",
		Header: h,
	}
}

func (p *GeekClient) doHTTP(method, url string, params interface{}) ([]byte, error) {

	fmt.Println("cli -> ", url, params)

	var body io.Reader

	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		fmt.Println(string(data))

		body = bytes.NewBuffer(data)
	}

	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("%s%s", p.url, url)
	}
	//提交请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = p.Header
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// MyProducts MyProducts
func (p *GeekClient) MyProducts() ([]*ProductsData, error) {
	var res *MyProductsResp
	data, err := p.doHTTP("POST", "/serv/v1/my/products/all", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &res)
	// err = res.UnmarshalJSON(data)
	if err != nil {
		fmt.Println(string(data))
		return nil, err
	}

	if res.Code != 0 {
		return nil, fmt.Errorf("%v", res.Error)
	}

	if len(res.Data) != 4 {
		return nil, fmt.Errorf("我的订阅页面有变化")
	}

	return res.Data, nil
}

// ColumnArticlesAll ColumnArticlesAll
func (p *GeekClient) ColumnArticlesAll(id int) ([]*ArticleItem, error) {

	list := []*ArticleItem{}

	var (
		preID int64
	)

	for {
		res, err := p.ColumnArticles(id, preID)
		if err != nil {
			return nil, err
		}

		list = append(list, res.Data.List...)

		if !res.Data.Page.More {
			break
		}

		l := len(res.Data.List)
		preID = res.Data.List[l-1].Score
	}

	return list, nil
}

// ColumnArticles ColumnArticles
// id column id
// prev pre scourc
func (p *GeekClient) ColumnArticles(id int, prev int64) (*ArticlesResp, error) {
	var res *ArticlesResp

	args := &ArticlesParams{}
	args.Cid = strconv.Itoa(id)
	args.Prev = prev

	if args.Cid == "" {
		return nil, fmt.Errorf("empty Cid")
	}

	if args.Order == "" {
		args.Order = "newest"
	}

	args.Size = 100

	args.Sample = true

	data, err := p.doHTTP("POST", "/serv/v1/column/articles", args)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &res)
	// err = res.UnmarshalJSON(data)
	if err != nil {
		fmt.Println(string(data))
		return nil, err
	}

	if res.Code != 0 {
		return nil, fmt.Errorf("%v", res.Error)
	}

	return res, nil
}

// ArticleInfo ArticleInfo
func (p *GeekClient) ArticleInfo(id int) (*ArticleInfo, error) {
	var res *ArticleInfoResp

	args := &ID{}
	args.ID = id

	data, err := p.doHTTP("POST", "/serv/v1/article", args)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &res)
	// err = res.UnmarshalJSON(data)
	if err != nil {
		fmt.Println(string(data))
		return nil, err
	}

	if res.Code != 0 {
		return nil, fmt.Errorf("%v", res.Error)
	}

	return res.Data, nil
}

// GetResource GetResource
func (p *GeekClient) GetResource(url string) ([]byte, error) {
	return p.doHTTP("GET", url, nil)
}

// ArticleCommentsAll ArticleCommentsAll
func (p *GeekClient) ArticleCommentsAll(id int) ([]*Comment, error) {

	list := []*Comment{}

	var (
		preID = "0"
	)

	for {
		res, err := p.ArticleComments(id, preID)
		if err != nil {
			return nil, err
		}

		list = append(list, res.Data.List...)

		if !res.Data.Page.More {
			break
		}

		l := len(res.Data.List)
		preID = res.Data.List[l-1].Score
	}

	return list, nil
}

// ArticleComments ArticleComments
// id column id
// prev pre scourc
func (p *GeekClient) ArticleComments(id int, prev string) (*CommentsResp, error) {
	var res *CommentsResp

	args := &CommentsParams{}
	args.Aid = strconv.Itoa(id)
	args.Prev = prev
	args.Size = 100

	data, err := p.doHTTP("POST", "/serv/v1/comments", args)
	if err != nil {
		return nil, err
	}

	fmt.Println("comment data", string(data))

	err = json.Unmarshal(data, &res)
	// err = res.UnmarshalJSON(data)
	if err != nil {
		fmt.Println(string(data))
		return nil, err
	}

	if res.Code != 0 {
		return nil, fmt.Errorf("%v", res.Error)
	}

	return res, nil
}
