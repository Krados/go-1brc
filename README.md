## Go 實作 1BRC（One Billion Row Challenge）：執行時間從 v1 的 75 秒大幅提升到 v8 的 2 秒

`go-1brc` 是一個使用 Go 實作的 [1BRC](https://github.com/gunnarmorling/1brc)（One Billion Row Challenge）

專案中提供多個版本的解法（`solution/v1solution.go` 到 `solution/v8solution.go`），
可用來比較不同實作策略在效能上的差異。

## 配備介紹
OS: Win11

CPU: AMD Ryzen 9 9950X3D 16-Core Processor

RAM: Kingston FURY Beast DDR5–6000 64GB

SSD: KINGSTON SFYRD2000G

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
- 解法版本：`-sv 8`

### 3) 指定檔案與版本

```bash
go-1brc -file measurements.txt -sv 8
```

常用參數：

- `-file`：輸入資料檔案（預設 `measurements.txt`）
- `-sv`：要執行的解法版本（`1`~`8`，預設 `8`）
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

## 結果比較
| 版本 | 方法/特色 | 執行時間 | 相對 v1 提升倍率 |
|------|-----------|----------|------------------|
| v1   | 簡單暴力 | 75s   | 1× |
| v2   | 使用 `scanner.Buffer` | 67s   | 1.12× |
| v3   | 自訂 `CustomByteSplit` 函數 | 46s   | 1.63× |
| v4   | 自訂 `ByteFloat64ToInt64V1` 函數 | 36s   | 2.08× |
| v5   | v4 的多工版本 | 3.6s  | 20.8× |
| v6   | 移除 leftover buf，復用原先的 buf | 2.5s  | 30× |
| v7   | 自訂 `ByteFloat64ToInt64V2` 函數 | 2.38s | 31.5× |
| v8   | `map` 的 key = `fnv1a(name)` | 2s   | 37.5× |
