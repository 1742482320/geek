package main

import (
	"bytes"
	"strconv"
	"time"

	"github.com/valyala/fasttemplate"
)

// TplVideoHTML TplVideoHTML
func TplVideoHTML(info *ArticleInfo, commentList []*Comment) string {
	t := fasttemplate.New(TplVideo, "{{", "}}")

	buf := new(bytes.Buffer)

	m := map[string]interface{}{}
	m["Title"] = info.ArticleTitle
	m["ArticleContent"] = info.ArticleContent
	m["VideoMediaURL"] = info.VideoMediaMap.Hd.URL
	// m["Column"] = info.AR
	m["comments"] = TplCommentList(commentList)

	buf.WriteString(t.ExecuteString(m))

	return buf.String()
}

// TplArticleHTML TplArticleHTML
func TplArticleHTML(info *ArticleInfo, commentList []*Comment) string {
	t := fasttemplate.New(articleHTML, "{{", "}}")

	buf := new(bytes.Buffer)

	m := map[string]interface{}{}
	m["AudioSize"] = strconv.Itoa(info.AudioSize / 1024)
	m["AudioTime"] = info.AudioTime
	m["ArticleTitle"] = info.ArticleTitle
	m["AuthorName"] = info.AuthorName
	m["ArticleCtime"] = time.Unix(info.ArticleCtime, 0).Format("2006-01-02")
	m["ArticleContent"] = info.ArticleContent
	m["AudioDownloadURL"] = info.AudioDownloadURL

	m["commentList"] = TplCommentList(commentList)

	buf.WriteString(t.ExecuteString(m))

	return buf.String()

}

// Replytemplate Replytemplate
var Replytemplate = `<div class="reply-hd"><i class="iconfont"></i> <span>{{UserName}}</span></div> 
	<p class="reply-content">{{Content}}</p> <p class="reply-time">{{Ctime}}</p>`

// TplReply TplReply
func TplReply(list []*Reply) string {
	t := fasttemplate.New(Replytemplate, "{{", "}}")

	buf := new(bytes.Buffer)
	for _, info := range list {
		m := map[string]interface{}{}
		m["UserName"] = info.UserName
		m["Content"] = info.Content
		m["Ctime"] = time.Unix(info.Ctime, 0).Format("2006-01-02")

		buf.WriteString(t.ExecuteString(m))
	}

	return buf.String()
}

// CommentInfotemplate CommentInfotemplate
var CommentInfotemplate = `<li data-v-87ffcada="" class="comment-item"><img src="{{UserHeader}}" class="avatar">
<div class="info">
	<div class="hd"><span class="username">{{UserName}}</span>
	<div class="control">
		 <a href="javascript:;" class="btn-praise"><i class="iconfont"></i> <span>{{LikeCount}}</span></a></div>
	</div>
	<div class="bd">{{comment_content}}</div> <span class="time">{{CommentCtime}}</span>
	<div class="reply">
	{{replies}}
	</div>
</div>
</li>`

// TplCommentList TplCommentList
func TplCommentList(list []*Comment) string {

	t := fasttemplate.New(CommentInfotemplate, "{{", "}}")

	buf := new(bytes.Buffer)
	for _, info := range list {
		m := map[string]interface{}{}
		m["UserHeader"] = info.UserHeader
		m["UserName"] = info.UserName
		m["LikeCount"] = strconv.Itoa(info.LikeCount)
		m["CommentContent"] = info.CommentContent
		m["CommentCtime"] = time.Unix(info.CommentCtime, 0).Format("2006-01-02")
		m["replies"] = TplReply(info.Replies)

		buf.WriteString(t.ExecuteString(m))
	}

	return buf.String()
}
