package gospider

import (
	"crypto/tls"
	"fmt"
	"go-crawler/client"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	tlsx "github.com/refraction-networking/utls"
)

// NewSpider 构造函数
func NewSpider(worker int) Spider {
	c, _ := tlsx.UTLSIdToSpec(tlsx.HelloRandomized)
	Transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         c.TLSVersMin,
			MaxVersion:         c.TLSVersMax,
			CipherSuites:       c.CipherSuites,
			ClientSessionCache: tls.NewLRUClientSessionCache(32),
		},
	}
	Jar, _ := cookiejar.New(nil)
	Client := &http.Client{
		Jar:       Jar,
		Transport: Transport,
		Timeout:   time.Second * 30}
	RChan := make(chan Request, worker)
	PChan := make(chan Response)
	return Spider{
		RequestQueue:  RChan,
		ResponseQueue: PChan,
		Client:        Client,
		Transport:     Transport,
		WorkerNum:     worker,
	}

}

// Spider 爬虫
type Spider struct {
	RequestQueue  chan Request
	ResponseQueue chan Response
	Client        *http.Client
	Transport     *http.Transport
	WorkerNum     int
}

// AddRequest 向请求队列添加新的请求
func (s *Spider) AddRequest(r Request) {
	s.RequestQueue <- r
}

// GetResponse 获取响应队列中的数据
func (s *Spider) GetResponse() (Response, error) {
	resp, ok := <-s.ResponseQueue
	if ok {
		return resp, nil
	}
	return Response{}, fmt.Errorf("Response queue closed.\n")
}

/*
Run 开始执行爬虫
*/
func (s *Spider) Run() {
	defer close(s.ResponseQueue)
	wg := sync.WaitGroup{}
	wg.Add(s.WorkerNum)
	for x := 0; x < s.WorkerNum; x++ {
		go s.spider(&wg)
	}
	wg.Wait()
}

func (s *Spider) spider(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		req, ok := <-s.RequestQueue
		if !ok {
			break
		}
		if req.Proxy != "" {
			proxy, err := url.Parse(req.Proxy)
			if err == nil {
				s.Transport.Proxy = http.ProxyURL(proxy)
			}
		}
		reqURL, _ := url.Parse(req.Url)
		if req.Cookie != nil {
			client.Jar.SetCookies(reqURL, req.Cookie)
		}
		clientReq, _ := http.NewRequest("GET", req.Url, nil)
		clientReq.Close = true
		for k, v := range req.Headers {
			clientReq.Header.Set(k, v)
		}
		clientResp, err := client.Client.Do(clientReq)
		resp := Response{
			Request: req,
			Error:   err,
		}
		if err == nil {
			defer clientResp.Body.Close()
			body, err := io.ReadAll(clientResp.Body)
			if err == nil {
				resp.Content = string(body)
			}
		}
		s.ResponseQueue <- resp
	}
}

/*
RandTransport 方法用于为Spider结构体生成随机的http.Transport

	该方法会设置Transport的DisableKeepAlives字段为true，禁用长连接
	同时会设置TLSClientConfig字段，包括跳过TLS证书验证、设置TLS协议版本范围、设置密码套件以及设置客户端会话缓存大小
*/
func (s *Spider) RandTransport() {
	c, _ := tlsx.UTLSIdToSpec(tlsx.HelloRandomized)
	transport := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         c.TLSVersMin,
			MaxVersion:         c.TLSVersMax,
			CipherSuites:       c.CipherSuites,
			ClientSessionCache: tls.NewLRUClientSessionCache(32),
		},
	}
	s.Transport = transport
	s.Client.Transport = transport
}

