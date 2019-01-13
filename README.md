目前結構
```
gin-blog/
├── conf
│   └── app.ini            (應用程序的獨立值)
├── main.go
├── middleware
│   └── jwt.go             (建立一個handleFunc 在API進入分流之前 驗證token)
├── models                 (這裡的models 處理與Database之間的交互)
│   └── models.go          (關於Model通用性方法)
|     └── tag.go           (Tag相關database相關方法)
|     └── article.go       (Article相關database相關方法)
├── pkg
│   ├── e
│   │   ├── code.go
│   │   └── msg.go
│   ├── setting
│   │   └── setting.go     (程序設定檔)
│   └── util
│       └── jwt.go         (生成及驗證Token)
│── util
│    └── pagination.go
├── routers
│   ├── api
│   │   └── v1             (處理router分下來的分流處理)
│   │       └── tag.go  
│   │       └── article.go 
│   │   └── auth.go        (處理GetAuth分流方法)
│   └── router.go          (分流方法)
├── runtime
```

依賴包
```
go get -u github.com/dgrijalva/jwt-go
```