# mcp-weather-go

MCPサーバー開発における[チュートリアル](https://modelcontextprotocol.io/quickstart/server#python)をGo言語で書き直したもの

## ディレクトリ構成

```
mcp-weather-go/
├── main.go        // 登録呼び出し
├── http.go        // 定数・HTTPクライアント・getJSON
├── alerts.go      // get_alerts
└── forecast.go    // get_forecast
```

## 利用要件

- Go 1.24+
- モジュール: `github.com/smasui-0302/mcp-weather-go`
- 依存: `github.com/modelcontextprotocol/go-sdk v0.2.0`

## ビルド方法

```bash
git clone https://github.com/smasui-0302/mcp-weather-go
cd mcp-weather-go
go mod tidy
go build -o mcp-weather
```

## MCPクライアントでの利用方法

- Claude Desktopの利用を想定しています
- `claude_desktop_config.json` に絶対Pathで登録する

```json
{
  "mcpServers": {
    "weather-go": {
      "command": "/ABSOLUTE/PATH/TO/mcp-weather-go/mcp-weather",
      "args": []
    }
  }
}
```

## 利用例

- 「What are the active weather alerts in TX?」
- 「Forecast near 38.58,-121.49（サクラメント付近）」
