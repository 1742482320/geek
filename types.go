package main

// MyProductsResp MyProductsResp
//easyjson:json
type MyProductsResp struct {
	Error []interface{}   `json:"error"`
	Data  []*ProductsData `json:"data"`
	Code  int             `json:"code"`
	Extra []interface{}   `json:"extra"`
}

// ProductsData ProductsData
//easyjson:json
type ProductsData struct {
	Page  Page          `json:"page"`
	Title string        `json:"title"`
	List  []*ColumnItem `json:"list"`
	ID    int           `json:"id"`
}

// Page Page
//easyjson:json
type Page struct {
	More  bool `json:"more"`
	Count int  `json:"count"`
}

// ColumnItem ColumnItem
//easyjson:json
type ColumnItem struct {
	Title string      `json:"title"`
	Extra *ColumnInfo `json:"extra"`
	Type  string      `json:"type"`
	Cover string      `json:"cover"`
	Score int64       `json:"score"`
}

// ColumnInfo ColumnInfo
//easyjson:json
type ColumnInfo struct {
	UpdateFrequency  string `json:"update_frequency"`
	ArticleCount     int    `json:"article_count"`
	ViewArticleCount int    `json:"view_article_count"`
	AuthorIntro      string `json:"author_intro"`
	Score            int64  `json:"score"`
	ColumnID         int    `json:"column_id"`
	HadSub           bool   `json:"had_sub"`
	ColumnType       int    `json:"column_type"`
	AuthorHeader     string `json:"author_header"`
	ColumnSubtitle   string `json:"column_subtitle"`
	ColumnSku        int    `json:"column_sku"`
	ColumnTitle      string `json:"column_title"`
	ColumnCover      string `json:"column_cover"`
	IsIncludeAudio   bool   `json:"is_include_audio"`
	SubTime          int    `json:"sub_time"`
	AuthorName       string `json:"author_name"`
}

// ArticlesParams ArticlesParams
//easyjson:json
type ArticlesParams struct {
	Cid    string `json:"cid"`
	Size   int    `json:"size"`
	Prev   int64  `json:"prev"`
	Order  string `json:"order"`
	Sample bool   `json:"sample"`
}

// ArticlesResp ArticlesResp
//easyjson:json
type ArticlesResp struct {
	Error []interface{} `json:"error"`
	Data  *ArticlesList `json:"data"`
	Code  int           `json:"code"`
	Extra []interface{} `json:"extra"`
}

// ArticlesList ArticlesList
//easyjson:json
type ArticlesList struct {
	List []*ArticleInfo `json:"list"`
	Page Page           `json:"page"`
}

// // ArticleItem ArticleItem
// //easyjson:json
// type ArticleItem struct {
// 	ArticleSubtitle     string `json:"article_subtitle"`
// 	ArticleCtime        int    `json:"article_ctime"`
// 	ID                  int    `json:"id"`
// 	ArticleCover        string `json:"article_cover"`
// 	ArticleTitle        string `json:"article_title"`
// 	ArticleSummary      string `json:"article_summary"`
// 	HadViewed           bool   `json:"had_viewed"`
// 	ArticleCouldPreview bool   `json:"article_could_preview"`
// 	ChapterID           string `json:"chapter_id"`
// 	Score               int64  `json:"score"`
// }

// ArticleInfoResp ArticleInfoResp
//easyjson:json
type ArticleInfoResp struct {
	Error []interface{} `json:"error"`
	Data  *ArticleInfo  `json:"data"`
	Code  int           `json:"code"`
	Extra []interface{} `json:"extra"`
}

// ArticleInfo ArticleInfo
//easyjson:json
type ArticleInfo struct {
	ArticleSubtitle     string         `json:"article_subtitle"`
	Sku                 string         `json:"sku"`
	ColumnHadSub        bool           `json:"column_had_sub"`
	AudioTitle          string         `json:"audio_title"`
	ViewCount           int            `json:"view_count"`
	VideoCover          string         `json:"video_cover"`
	AudioDownloadURL    string         `json:"audio_download_url"`
	AudioTime           string         `json:"audio_time"`
	VideoMedia          string         `json:"video_media"`
	ProductType         string         `json:"product_type"`
	ArticleContent      string         `json:"article_content"`
	LikeCount           int            `json:"like_count"`
	VideoHeight         int            `json:"video_height"`
	ArticleTitle        string         `json:"article_title"`
	AudioSize           int            `json:"audio_size"`
	ArticleSharetitle   string         `json:"article_sharetitle"`
	AuthorName          string         `json:"author_name"`
	ArticleCtime        int64          `json:"article_ctime"`
	ID                  int            `json:"id"`
	ArticleCover        string         `json:"article_cover"`
	AudioURL            string         `json:"audio_url"`
	VideoSize           int            `json:"video_size"`
	ChapterID           string         `json:"chapter_id"`
	HadLiked            bool           `json:"had_liked"`
	ColumnIsExperience  bool           `json:"column_is_experience"`
	HadViewed           bool           `json:"had_viewed"`
	Score               interface{}    `json:"score"`
	ColumnBgcolor       string         `json:"column_bgcolor"`
	ColumnCover         string         `json:"column_cover"`
	VideoTime           string         `json:"video_time"`
	AudioMd5            string         `json:"audio_md5"`
	AudioTimeArr        AudioTimeArr   `json:"audio_time_arr"`
	Cid                 int            `json:"cid"`
	ArticleCoverHidden  bool           `json:"article_cover_hidden"`
	ArticleSummary      string         `json:"article_summary"`
	ArticleCouldPreview bool           `json:"article_could_preview"`
	AudioDubber         string         `json:"audio_dubber"`
	VideoWidth          int            `json:"video_width"`
	ColumnID            int            `json:"column_id"`
	ArticlePosterWxlite string         `json:"article_poster_wxlite"`
	VideoMediaMap       *VideoMediaMap `json:"video_media_map"`
	VideoTotalSeconds   int            `json:"video_total_seconds"`
	VideoTimeArr        VideoTimeArr   `json:"video_time_arr"`
	VideoPlaySeconds    int            `json:"video_play_seconds"`
}

// VideoTimeArr VideoTimeArr
//easyjson:json
type VideoTimeArr struct {
	M string `json:"m"`
	S string `json:"s"`
	H string `json:"h"`
}

// VideoMediaMap VideoMediaMap
//easyjson:json
type VideoMediaMap struct {
	Hd Media `json:"hd"`
	Sd Media `json:"sd"`
}

// Media Media
//easyjson:json
type Media struct {
	Size int    `json:"size"`
	URL  string `json:"url"`
}

// AudioTimeArr AudioTimeArr
//easyjson:json
type AudioTimeArr struct {
	M string `json:"m"`
	S string `json:"s"`
	H string `json:"h"`
}

// ID ID
//easyjson:json
type ID struct {
	ID int `json:"id"`
}

// CommentsResp CommentsResp
//easyjson:json
type CommentsResp struct {
	Error []interface{} `json:"error"`
	Data  CommentList   `json:"data"`
	Code  int           `json:"code"`
	Extra []interface{} `json:"extra"`
}

// CommentList CommentList
//easyjson:json
type CommentList struct {
	List []*Comment `json:"list"`
	Page Page       `json:"page"`
}

// Comment Comment
//easyjson:json
type Comment struct {
	UserHeader     string   `json:"user_header"`
	UserName       string   `json:"user_name"`
	ID             int      `json:"id"`
	LikeCount      int      `json:"like_count"`
	CommentIsTop   bool     `json:"comment_is_top"`
	HadLiked       bool     `json:"had_liked"`
	CommentCtime   int64    `json:"comment_ctime"`
	CommentContent string   `json:"comment_content"`
	Score          string   `json:"score"`
	Replies        []*Reply `json:"replies,omitempty"`
}

// Reply Reply
//easyjson:json
type Reply struct {
	CommentID    int    `json:"comment_id"`
	Content      string `json:"content"`
	Utype        int    `json:"utype"`
	ID           string `json:"id"`
	UserName     string `json:"user_name"`
	UserNameReal string `json:"user_name_real"`
	Ctime        int64  `json:"ctime"`
	UID          string `json:"uid"`
}

// CommentsParams CommentsParams
//easyjson:json
type CommentsParams struct {
	Aid  string `json:"aid"`
	Prev string `json:"prev"`
	Size int    `json:"size"`
}

// LoginResp LoginResp
//easyjson:json
type LoginResp struct {
	Code  int           `json:"code"`
	Data  *UserInfo     `json:"data"`
	Error []interface{} `json:"error"`
	Extra struct {
		Cost      float64 `json:"cost"`
		RequestID string  `json:"request-id"`
	} `json:"extra"`
}

// UserInfo UserInfo
//easyjson:json
type UserInfo struct {
	UID       int    `json:"uid"`
	Type      int    `json:"type"`
	Cellphone string `json:"cellphone"`
	Country   string `json:"country"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	// Gender     string `json:"gender"`
	// Birthday   string `json:"birthday"`
	// Graduation string `json:"graduation"`
	// Profession string `json:"profession"`
	// Industry   string `json:"industry"`
	// Email      string `json:"email"`
	// Name       string `json:"name"`
	// Address    string `json:"address"`
	// Mobile     string `json:"mobile"`
	// Contact    string `json:"contact"`
	// Position   string `json:"position"`
	// Passworded bool   `json:"passworded"`
	// CreateTime int64  `json:"create_time"`
	OssToken string `json:"oss_token"`
}

// LoginParams LoginParams
//easyjson:json
type LoginParams struct {
	Country   int    `json:"country"`
	Cellphone string `json:"cellphone"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	Remember  int    `json:"remember"`
	Platform  int    `json:"platform"`
	Appid     int    `json:"appid"`
}

// MyProductListResp MyProductListResp
//easyjson:json
type MyProductListResp struct {
	Error []interface{}  `json:"error"`
	Data  *MyProductList `json:"data"`
	Code  int            `json:"code"`
	Extra []interface{}  `json:"extra"`
}

// MyProductList MyProductList
//easyjson:json
type MyProductList struct {
	List []*ColumnItem `json:"list"`
	Page Page          `json:"page"`
}

// MyProductListParams MyProductListParams
//easyjson:json
type MyProductListParams struct {
	NavID int   `json:"nav_id"`
	Prev  int64 `json:"prev"`
	Size  int   `json:"size"`
}
