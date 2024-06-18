## gospider

### Usage
`go get -u github.com/symding/gospider`
```go
package main

import "github.com/symding/gospider"

func main() {
    // create a spider with 2 download goroutine
    spider := gospider.NewSpider(2)
    // initial spider
    go spider.Run()
    // put request task
    go func(){
        for {
            request := gospider.Request{
                Url:"https://www.baidu.com",
            }
            spider.AddRequest(request)
        }
    }()
    // response loop
    for {
        response,err := spider.GetResponse()
        if err!=nil {
            break
        }
        response.Xpath.ExtractFirst("//div[@class='mnav']/a/@href")
        // parse response code here
    }
}
```