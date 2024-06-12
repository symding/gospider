## gospider

### Usage
```go
package main

import gospider 

func main() {
    spider := gospider.NewSpider(2)
    go spider.Run()
    go func(){
        for {
            request := gospider.Request{
                Url:"https://www.baidu.com",
            }
            spider.AddRequest(request)
        }
    }()
    for {
        response,err := spider.GetResponse()
        if err!=nil {
            break
        }
    }
}
```