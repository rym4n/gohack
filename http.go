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

//HTTPBuildQuery 将map转换为url查询参数形式
func HTTPBuildQuery(queryArr map[string]string) string {
	if queryArr == nil {
		return ""
	}
	q := make(url.Values)
	for k, v := range queryArr {
		q.Add(k, v)
	}
	return q.Encode()
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

// HTTPViaProxy 通过proxy发包，支持socks、http、https等类型的proxy
func HTTPViaProxy(method, url, data, proxy string, headers *http.Header, cookie *http.Cookie) (*http.Response, error) {
	if proxy == "" {
		proxy = "direct://0.0.0.0:0000"
	}
	dialer, err := proxyclient.NewProxyClient(proxy)
	if err != nil {
		log.Panicf("proxy 出错了: %s", err)
	}
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial:            dialer.Dial,
	}

	client := &http.Client{
		Transport: httpTransport,
		Timeout:   10 * time.Second,
	}

	var req *http.Request

	switch method {
	case "POST":
		req, err = http.NewRequest("POST", url, strings.NewReader(data))
	default:
		req, err = http.NewRequest("GET", url, nil)
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

	return client.Do(req)
}

/*HTTPRawViaProxy 用于发送原始报文
@raw  : burp 之类抓到的原始报文
*/
func HTTPRawViaProxy(protocol, host, port, raw, proxy string) (*http.Response, error) {
	sep := "\n"
	if strings.Contains(raw, "\r\n") {
		sep = "\r\n"
	}
	rawHTTP := strings.Split(raw, sep+sep)
	headers := strings.Split(rawHTTP[0], sep)
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

	url := fmt.Sprintf("%s://%s:%s%s", protocol, host, port, path)

	reqHeader := make(http.Header)
	for _, v := range headers[1:] {
		// headers := make(http.Header)
		h := strings.Split(v, ":")
		reqHeader.Add(h[0], h[1])
	}

	return HTTPViaProxy(firstHead[0], url, data, proxy, &reqHeader, nil)
}

// GetResponseText 从响应结构体中获取文本字符串，包括自动处理gzip
func GetResponseText(resp *http.Response) (string, error) {
	// 处理gzip
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
