# 富途OpenAPI Golang SDK

## 简介

* Go语言封装的[富途牛牛OpenAPI](https://openapi.futunn.com/futu-api-doc/)。
* 尽量接近Python版本的使用方法。
* 利用Go语言特性，例如channel，goroutine等。
* 注意：该SDK仍在v0版本，不保证与旧版本兼容，也没有经过充分测试，使用时候请注意，可以直接向我反馈bug或者改进建议

## 代码结构

### TCP连接层

### 富途协议层

### API接口层

## 使用方法

1. import

    ```
    import "github.com/hurisheng/go-futu-api"
    ```

1. 创建API实例

    ```
    ft := futuapi.NewFutuAPI()
    ```

1. 设置连接参数

    ```
    ft.SetClientInfo("MyFutuAPI", 0)
    ```

1. 连接FutuOpenD

    ```
    ft.Connect(context.Background(), ":11111")
    ```

1. 调用获取方法

    ```
    sub, err := ft.QuerySubscription(context.Background(), true)
    ```

1. 接收推送

    ```
    ch, err := api.UpdateTicker(context.Background())
    // ch 为 channel类型，<- ch接收推送
    ```

## Change Log

* 更新了接口列表，删除了失效的接口，补充新接口，修改接口的参数等
* 去掉自定义的数据结构，直接使用protobuf编译出来的结构，update类型的数据，返回整个Response消息
* 对于个别get类型的接口，简化返回结构体嵌套，例如返回的数据中只有一项是有意义的，直接返回该字段
* 删除自动根据protobuf类型编译自定义类型的代码，改为reflection操作channel数据发送
