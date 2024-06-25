# 获取最新的protobuf定义

官方QQ群咨询最新的proto定义文件

# 编译protobuf文件

按照protobuf官方安装go版本编译器
```
protoc --go_out=. --go_opt=module=github.com/hurisheng/go-futu-api/pb *.proto
```