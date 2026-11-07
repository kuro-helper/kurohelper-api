# Fiber 請求資料與快取

## 為什麼需要複製字串

本專案使用 Fiber v3.6.1，目前以 `fiber.New()` 建立應用程式，未開啟 `Immutable`。

Fiber 預設透過 `unsafe`，讓部分回傳的字串直接共用請求緩衝區，以減少記憶體配置、文字複製與 GC 負擔。常見方法包括 `Query`、`Params`、`Get`（讀取標頭）及 `Cookies`；`Body`、`BodyRaw` 回傳的位元組切片也有生命週期限制。

這些資料只保證在當次 handler 內有效。handler 結束後，緩衝區可能被後續請求覆寫。如果直接把資料保留在快取、全域變數或背景工作中，之後讀到的內容可能改變，造成快取查找異常或資料競爭。

一般 Go 字串具有不可變的值語意；這裡的限制來自 Fiber 使用 `unsafe` 共用可修改的記憶體，不代表一般字串都需要複製才能保存。

## 本專案的處理方式

Music 搜尋會把關鍵字作為快取鍵，因此先複製文字，讓它可以安全保留到請求結束後：

```go
keyword := strings.Clone(strings.TrimSpace(c.Query("q")))
cache.ErogsMusicListStore.Set(keyword, list)
```

`TrimSpace` 不保證複製文字；`RuneCountInString` 只用來驗證字數，其回傳的整數並不是快取鍵。

Game 詳細資料則先將路徑參數解析成整數，再產生快取鍵，不再引用原本的請求緩衝區，因此不需要額外 Clone：

```go
id, err := strconv.Atoi(idText)
// 確認 err 為 nil 並驗證 id 後
cacheKey := strconv.Itoa(id)
```

## Immutable 與官方修正

Fiber 曾針對 issue #185 提供 `Immutable` 選項，但預設仍為 `false`，並未取消預設的緩衝區共用行為。

若設定 `fiber.Config{Immutable: true}`，相關 context 取值方法會處理資料複製，這裡便不需要再額外 Clone，代價是增加配置與複製成本。此設定不代表可以在 handler 結束後繼續使用 `fiber.Ctx` 本身。

本文件僅記錄行為，未調整應用程式設定或程式碼。

## 參考資料

- [Fiber：Zero Allocation](https://docs.gofiber.io/#zero-allocation)
- [Fiber：Immutable 設定](https://docs.gofiber.io/api/fiber/#config)
- [Issue #185](https://github.com/gofiber/fiber/issues/185)
- [v1.8.0 發布紀錄](https://github.com/gofiber/fiber/issues/191)
