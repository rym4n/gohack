package gohack

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gamexg/proxyclient"
)

// HTTPBuildQuery 将map转换为url查询参数形式
func HTTPBuildQuery(queryArr map[string]string) (queryString string) {
	if queryArr == nil {
		return ""
	}
	q := make(url.Values)
	for k, v := range queryArr {
		q.Add(k, v)
	}
	queryString = q.Encode()
	return
}

// QueryStringToMap 将url查询参数转换为map
func QueryStringToMap(query string) (ret map[string]string, err error) {
	ret = make(map[string]string)
	m, err := url.ParseQuery(query)
	if err != nil {
		return ret, err
	}

	for k, v := range m {
		if len(v) > 0 {
			ret[k] = v[0]
		}
	}
	return ret, nil
}

/*
HTTPViaProxy 通过proxy发包，支持socks、http、https等类型的proxy
@method: 值为 OPTIONS 或 GET 或 HEAD 或 POST 或 PUT 或 DELETE 或 TRACE 或 CONNECT
@url: 请求URL
@data: HTTP Body 中的数据
@proxy: 指定代理，格式为 type://host:port, 例如: socks5://127.0.0.1:1080, 无代理则传入空字符串
@timeout: 超时时间，请求并发量较大时，timeout最好大一点
@headers: http.Header类型指针，表示请求头
@cookie: http.Cookie类型指针，表示请求cookie
*/
func HTTPViaProxy(method, url, data, proxy string, timeout int, headers *http.Header, cookie *http.Cookie) (response *http.Response, err error) {
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if proxy != "" {
		dialer, err := proxyclient.NewProxyClient(proxy)
		if err != nil {
			log.Panicf("proxy 出错了: %s", err)
		}
		httpTransport.Dial = dialer.Dial
	}

	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: httpTransport,
	}

	var req *http.Request

	method = strings.ToUpper(method)
	if data == "" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(data))
	}

	if err != nil {
		log.Panicf("构造HTTP请求出错: %s", err)
	}

	if headers != nil {
		req.Header = *headers
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	response, err = client.Do(req)
	return
}

/*
HTTPRawViaProxy 用于发送原始报文
@protocol: http 或 https
@host: 主机IP或域名
@port: 端口
@raw: burpsuite之类的工具抓到的原始报文
@proxy: 指定代理，格式为 type://host:port, 例如: socks5://127.0.0.1:1080, 无代理则传入空字符串
@timeout: 超时时间，请求并发量较大时，timeout最好大一点
*/
func HTTPRawViaProxy(protocol, host, port, raw, proxy string, timeout int) (response *http.Response, err error) {
	sep := "\n"
	if strings.Contains(raw, "\r\n") {
		sep = "\r\n"
	}
	rawHTTP := strings.Split(raw, sep+sep)
	headers := strings.Split(strings.TrimSpace(rawHTTP[0]), sep)
	data := ""
	if len(rawHTTP) > 1 {
		data = strings.Join(rawHTTP[1:], sep+sep)
	}

	firstHead := strings.Split(headers[0], " ")
	path := firstHead[1]

	if strings.HasPrefix(path, "http") {
		u, err := url.Parse(path)
		if err != nil {
			log.Panicf("path解析出错: %s", err)
		}
		path = strings.Split(path, u.Host)[1]
	}

	url := fmt.Sprintf("%s://%s:%s%s", strings.ToLower(protocol), host, port, path)

	reqHeader := make(http.Header)
	for _, v := range headers[1:] {
		h := strings.Split(v, ":")
		reqHeader.Add(h[0], h[1])
	}

	response, err = HTTPViaProxy(firstHead[0], url, data, proxy, timeout, &reqHeader, nil)
	return
}

// GetResponseText 从响应结构体中获取文本字符串，包括自动处理gzip
func GetResponseText(resp *http.Response) (string, error) {
	//处理gzip
	var (
		body []byte
		err  error
	)
	defer resp.Body.Close()
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		defer reader.Close()
		if err != nil {
			log.Fatalf("读取HTTP响应错误: %s", err)
			return "", err
		}
		body, err = ioutil.ReadAll(reader)
		if err != nil {
			log.Fatalf("读取HTTP响应错误: %s", err)
			return "", err
		}
	} else {
		body, err = ioutil.ReadAll(resp.Body)
	}

	return string(body), err
}
