企业微信消息推送-Go

> Wechat-Work-Pusher Based On Golang

## 部署

#### 准备工作

1、注册企业微信： https://work.weixin.qq.com/ 注册完成后点击我的企业，获取企业ID(`cropId`)

2、进入应用管理，创建自建应用，获取`agentId`和`cropSecret`

3、进入我的企业-微信插件-邀请关注 扫码关注企业微信

#### config.json配置

1、将准备工作中获取到的cropId、cropSecret、agentId填入相应位置

2、receiver为默认接收者的微信号，对应通讯录中的帐号

3、填写token

> 说明：服务端路由基路径由 `rest.base` 控制，默认是 `wechat-work-pusher`，因此接口为：
> `POST /wechat-work-pusher/msg` 与 `POST /wechat-work-pusher/card`

#### Docker部署

    docker run -d -v /home/config.json:/root/config.json -p 9000:9000 --restart=always --name wechat-work-pusher myleo1/wechat-work-pusher

> 注意：替换/home/config.json 为config.json路径；也可通过环境变量 `CONFIG_PATH` 指定配置文件路径。

## 使用方法

以Go为例，其他语言类似

#### 文本消息（原生 net/http 示例）

```go
import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
)

func Push2Wechat(to, msg string) error {
    form := url.Values{}
    form.Set("to", to)
    form.Set("content", msg)

    req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9000/wechat-work-pusher/msg", strings.NewReader(form.Encode()))
    if err != nil { return err }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", 配置中的token))

    resp, err := http.DefaultClient.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK { return fmt.Errorf("http status: %d", resp.StatusCode) }
    return nil
}
```

#### 卡片消息（原生 net/http 示例）

卡片消息参数具体请参考：https://work.weixin.qq.com/api/doc/90000/90135/90236#%E6%96%87%E6%9C%AC%E5%8D%A1%E7%89%87%E6%B6%88%E6%81%AF

```go
import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
)

func Push2WechatCard(to, title, description, link string) error {
    form := url.Values{}
    form.Set("to", to)
    form.Set("title", title)
    form.Set("description", description)
    form.Set("url", link)

    req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9000/wechat-work-pusher/card", strings.NewReader(form.Encode()))
    if err != nil { return err }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", 配置中的token))

    resp, err := http.DefaultClient.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK { return fmt.Errorf("http status: %d", resp.StatusCode) }
    return nil
}
```

## 效果演示

![image](https://user-images.githubusercontent.com/66349676/111748431-7eadf380-88cb-11eb-8590-73e1414d98e6.png)

