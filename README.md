目前結構
```
gin-blog/
├── conf
│   └── app.ini               (應用程序的獨立值)
├── main.go
├── middleware
│   └── jwt
│      └── jwt.go             (建立一個handleFunc 在API進入分流之前 驗證token)
├── models                    (這裡的models 處理與Database之間的交互)
│   └── models.go             (關於Model通用性方法)
|   └── tag.go                (Tag與database相關方法)
|   └── article.go            (Article與database相關方法)
|   └── auth.go               (Auth與database相關方法)
├── pkg
│   ├── e
│   │   ├── code.go
│   │   └── msg.go
│   ├── setting
│   │   └── setting.go        (程序設定檔)
│   └── util
│   │   └── jwt.go            (生成及驗證Token)
│   └── upload
│   │   └── image.go          (照片檔案相關方法)
│   └── logging
│       └── file.go           (檔案相關方法)
│       └── log.go            (記錄在logfile相關方法)
├── pkg
│   ├── article_service
│            ├──  article.go  (Article的獲Get/Add/Updata方法)
│   ├── tag_service
│            ├──  tag.go
│   ├── cache_service
│            ├──  article.go  (產生獨特Key的方法)
│── util
│    └── pagination.go
├── routers生成及驗證Token
│   ├── api
│   ├── upload                (處理router分下來的分流處理)
│   │       └── tag.go  
│   │       └── article.go 
│   │   └── auth.go           (處理GetAuth分流方法)
│   │   └── auth.go           (UploadImage方法)
│   └── router.go             (分流方法)
├── runtime
│   ├── logs
│         ├── log.log         (工作日誌儲存位置)
│   ├── upload
│         ├── images          (上傳照片後的儲存位置)
```

依賴包
```
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/fvbock/endless
```