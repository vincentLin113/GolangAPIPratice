目前結構
```
gin-blog/
├── conf
│   └── app.ini
├── main.go
├── middleware
├── models    (這裡的models 處理與Database之間的交互)
│   └── models.go
 |     └── tag.go        
├── pkg
│   ├── e
│   │   ├── code.go
│   │   └── msg.go
│   ├── setting
│   │   └── setting.go
│   └── util
│       └── pagination.go
├── routers
│   ├── api
│   │   └── v1      (這裡的tag.go 是處理router分下來的分流處理)
│   │       └── tag.go  
│   └── router.go
├── runtime
```
