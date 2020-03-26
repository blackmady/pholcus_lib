package pholcus_lib

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析

	// "github.com/henrylee2cn/pholcus/logs"               //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包

	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	"regexp"
	// "strconv"
	"strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	WangyiWorld.Register()
}

var WangyiWorld = &Spider{
	Name:        "网易国际新闻",
	Description: "网易国际新闻-分三大块",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "https://news.163.com/world/", Rule: "国际新闻首页"})
		},

		Trunk: map[string]*Rule{
			"国际新闻首页": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find(".today_news a,.mt23.mod_jsxw a").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							ctx.AddQueue(&request.Request{Url: url, Rule: "今日推荐-即时新闻-缩略图"})
						}
					})
					// query.Find(".mt23.mod_jsxw a").Each(func(i int, s *goquery.Selection) {
					// 	if url, ok := s.Attr("href"); ok {
					// 		ctx.AddQueue(&request.Request{Url: url, Rule: "今日推荐-即时新闻"})
					// 	}
					// })
					query.Find(".ndi_main .news_article").Each(func(i int, s *goquery.Selection) {
						thumb := query.Find("img")
						a := query.Find("h3 a")
						if url, ok := a.Attr("href"); ok {
							ctx.AddQueue(&request.Request{Url: url, Rule: "今日推荐-即时新闻-缩略图", Temp: map[string]interface{}{
								"thumb": thumb,
							}})
						}
					})
				},
			},
			"今日推荐-即时新闻-缩略图": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"title",
					"content",
					"ReleaseTime",
					"url",
					"thumb",
					"extra",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					// 获取标题
					title := query.Find("h1").Text()
					// 获取内容
					content := query.Find("#endText").Text()
					re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
					// content = re.ReplaceAllStringFunc(content, strings.ToLower)
					content = re.ReplaceAllString(content, "")

					// 获取发布日期
					extra := query.Find(".post_time_source").Text()
					aextra := strings.Split(extra, "来源:")
					release := aextra[0]
					extra = aextra[1]
					release = strings.Trim(release, " \t\n")

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: title,
						1: content,
						2: release,
						3: ctx.Request.Url,
						4: ctx.GetTemp("thumb", ""),
						5: extra,
					})
				},
			},
		},
	},
}
