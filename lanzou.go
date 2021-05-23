package lanZouApi

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var header = req.Header{
	"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
	"Referer":         "https://lanzous.com",
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)" +
		" Chrome/89.0.4389.114 Safari/537.36",
}

type LanZouCloud struct {
	Url       string
	Pwd       string
	DirectUrl string
	client    *req.Req
}

func New(u, p string) *LanZouCloud {
	lanZouCloud := new(LanZouCloud)
	lanZouCloud.Url = u
	lanZouCloud.Pwd = p
	return lanZouCloud
}

func (c *LanZouCloud) init() {
	c.client = req.New()
	c.client.SetClient(&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	})
}

func (c *LanZouCloud) Do() {
	c.init()
	if len(c.Pwd) == 0 {
		c.getDirectByNoPwd()
	} else {
		c.getDirectByPwd()
	}

}

func (c *LanZouCloud) postData(par interface{}) {
	tmpHeader := header.Clone()
	tmpHeader["Content-Type"] = "application/x-www-form-urlencoded"
	resp, _ := c.client.Post("https://lanzoui.com/ajaxm.php", tmpHeader, par)
	var res = gjson.GetMany(resp.String(), "dom", "url")
	var builder strings.Builder
	builder.WriteString(res[0].String())
	builder.WriteString(`/file/`)
	builder.WriteString(res[1].String())
	fakeUrl := builder.String()
	rr, err := c.client.Get(fakeUrl, header)
	if err != nil {
		fmt.Println(err)
		return
	}
	location, _ := rr.Response().Location()
	c.DirectUrl = location.String()

}
func (c *LanZouCloud) getDirectByPwd() {
	resp, err := c.client.Get(c.Url, header)
	if err != nil {
		log.Println(err)
		return
	}
	reg, _ := regexp.Compile(`data : '(.+)'\+pwd`)
	par := reg.FindStringSubmatch(resp.String())[1] + c.Pwd
	c.postData(par)
}
func (c *LanZouCloud) getDirectByNoPwd() {
	resp, err := c.client.Get(c.Url, header)
	if err != nil {
		log.Println(err)
		return
	}

	reg, _ := regexp.Compile(`src="(.{20,})" frameborder`)
	par := reg.FindStringSubmatch(resp.String())[1]
	r, _ := c.client.Get("https://lanzous.com" + par)
	tmp := r.String()

	regData, _ := regexp.Compile(`var (.*?) = '(.*?)';`)
	datas := regData.FindAllStringSubmatch(tmp, -1)

	regReplaceAjaxData, _ := regexp.Compile(`'signs':(.*?),`)
	tmp = regReplaceAjaxData.ReplaceAllString(tmp, "'signs':'"+datas[0][2]+"',")
	regReplacePostDown, _ := regexp.Compile(`'sign':(.*?),`)
	tmp = regReplacePostDown.ReplaceAllString(tmp, "'sign':'"+datas[1][2]+"',")
	reg, _ = regexp.Compile(`[^/]{2,}data : ({.+})`)
	//data : { 'action':'downprocess','signs':ajaxdata,'sign':postdown,'ves':1,'websign':'','websignkey':'s3tz' }
	pars := reg.FindStringSubmatch(tmp)[1]
	fmt.Println(pars)
	params := new(req.Param)
	pars = strings.Replace(pars, `'`, `"`, -1)
	err = json.Unmarshal([]byte(pars), params)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.postData(*params)
	return
}
