## go-1brc

`go-1brc` 是一個使用 Go 實作的 [1BRC](https://github.com/gunnarmorling/1brc)（One Billion Row Challenge）

專案中提供多個版本的解法（`solution/v1solution.go` 到 `solution/v7solution.go`），
可用來比較不同實作策略在效能上的差異。

## 前置作業

### 1) 建立資料集

[請參照](https://github.com/gunnarmorling/1brc?tab=readme-ov-file#prerequisites)

## 如何使用

### 1) 環境需求

- Go `1.26.0`（或相容版本）

### 2) 直接編譯

在專案根目錄執行：

```bash
go build -o go-1brc .
```

預設會使用：

- 輸入檔案：`measurements.txt`
- 解法版本：`-sv 7`

### 3) 指定檔案與版本

```bash
go-1brc -file measurements.txt -sv 7
```

常用參數：

- `-file`：輸入資料檔案（預設 `measurements.txt`）
- `-sv`：要執行的解法版本（`1`~`7`，預設 `7`）
- `-cpuprofile`：輸出 CPU profile 檔案
- `-memprofile`：輸出記憶體 profile 檔案

### 4) 產生與查看效能分析（pprof）

先執行：

```bash
go-1brc -cpuprofile cpuprofile.prof -memprofile memprofile.prof
```

再用 pprof UI 查看：

```bash
go tool pprof -http=:8080 cpuprofile.prof
```
