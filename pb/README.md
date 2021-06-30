# 编译protobuf文件

```
protoc --go_out=. --go_opt=module=github.com/hurisheng/go-futu-api/pb --go-futu_out=. --go-futu_opt=module=github.com/hurisheng/go-futu-api/pb *.proto
```