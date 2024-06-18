package gospider

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	RChan := make(chan Request, worker*2)
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
	Stat          spiderStat
	close         bool
}

/*
AddRequest 向请求队列添加新的请求

	返回任务队列状态，
*/
func (s *Spider) AddRequest(r Request) bool {
	s.RequestQueue <- r
	s.Stat.RequestIncr()
	return !s.close
}

// GetResponse 获取响应队列中的数据
func (s *Spider) GetResponse() (Response, bool) {
	resp, ok := <-s.ResponseQueue
	if ok {
		s.Stat.ResponseIncr(resp)
		return resp, ok
	}
	s.Stat.Stop()
	return Response{}, ok
}

/*
Run 开始执行爬虫
*/
func (s *Spider) Run() {
	s.Stat = newSpiderStat()
	go s.Signal()
	defer log.Println("Response queue closed.")
	defer close(s.ResponseQueue)
	wg := sync.WaitGroup{}
	wg.Add(s.WorkerNum)
	for x := 0; x < s.WorkerNum; x++ {
		go s.spider(&wg)
	}
	wg.Wait()
}

// Signal 捕获信号量并处理
func (s *Spider) Signal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Printf("Received signal: %s\n", sig)
	s.close = true
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
			s.Client.Jar.SetCookies(reqURL, req.Cookie)
		}
		clientReq, _ := http.NewRequest("GET", req.Url, nil)
		clientReq.Close = true
		for k, v := range req.Headers {
			clientReq.Header.Set(k, v)
		}
		clientResp, err := s.Client.Do(clientReq)
		resp := Response{
			Request:    req,
			Error:      err,
			Meta:       req.Meta,
			StatusCode: clientResp.StatusCode,
		}
		if err == nil {
			defer clientResp.Body.Close()
			body, err := io.ReadAll(clientResp.Body)
			if err == nil {
				resp.Content = string(body)
				resp.Xpath = NewXpathParser(body)
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

func (s *Spider) Close() {
	close(s.RequestQueue)
	log.Println("Request queue closed.")
}
