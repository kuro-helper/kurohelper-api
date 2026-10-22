# KuroHelper API

## 專案用途
本專案為 **KuroHelper** 的後端 API 服務，基於 Go 語言與 Fiber (v3) 框架開發。主要提供前端及周邊工具所需的 RESTful API 服務，包含會員帳號與 Session 驗證管理、遊戲資訊維護、系統公告、統計數據。

## 啟動方式

1. **設定環境變數**  
   複製範例設定檔並設定相關配置（資料庫連線、連接埠等）：
   ```bash
   cp .env.example .env
   ```

2. **啟動服務**  
   直接透過 Go 執行主程式：
   ```bash
   go run main.go
   ```

## 編譯方式

透過 `go build` 編譯出 Linux 二進位執行檔：

```bash
go build -o kurohelper-api main.go
```

編譯完成後即可直接執行：
```bash
./kurohelper-api
```

> 執行環境目錄下需包含相應的 `.env` 設定檔。

