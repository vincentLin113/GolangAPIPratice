目前結構
```
gin-blog/
├── conf
│   └── app.ini (應用程序的獨立值)
├── main.go
├── middleware
├── models    (這裡的models 處理與Database之間的交互)
│   └── models.go    (關於Model通用性方法)
 |     └── tag.go    (Tag相關database相關方法)
├── pkg
│   ├── e
│   │   ├── code.go
│   │   └── msg.go
│   ├── setting
│   │   └── setting.go   (程序設定檔)
│   └── util
│       └── pagination.go
├── routers
│   ├── api
│   │   └── v1      (這裡的tag.go 是處理router分下來的分流處理)
│   │       └── tag.go  
│   └── router.go   (分流方法)
├── runtime
```
