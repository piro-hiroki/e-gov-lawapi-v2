# e-gov-lawapi-v2

[e-Gov 法令API v2](https://laws.e-gov.go.jp/api/2/swagger-ui) を [Model Context Protocol (MCP)](https://modelcontextprotocol.io) 経由で公開する MCP server です（Go 実装）。
Claude Desktop / Claude Code などの MCP 対応クライアントから、日本の法令を検索・取得できます。

## 提供ツール

| ツール名 | 概要 | 対応エンドポイント |
| --- | --- | --- |
| `search_laws` | 法令名や法令種別、公布日範囲などで法令一覧を検索 | `GET /laws` |
| `get_law_revisions` | 指定した法令の改正履歴一覧を取得 | `GET /law_revisions/{law_id_or_num}` |
| `get_law_data` | 法令本文（条文）を取得（条項単位の絞り込み可） | `GET /law_data/{law_id_or_num_or_revision_id}` |
| `keyword_search` | 法令本文を全文検索（AND/OR/NOT、ワイルドカード対応） | `GET /keyword` |

レスポンスはすべて JSON 文字列でテキストコンテンツとして返却されます。

## 必要要件

- Go 1.23 以上（開発時 1.26 で動作確認）

## セットアップ

```bash
git clone https://github.com/piro-hiroki/e-gov-lawapi-v2.git
cd e-gov-lawapi-v2
go build -o e-gov-lawapi-v2 ./...
```

`go install github.com/piro-hiroki/e-gov-lawapi-v2@latest` でも導入できます（`$GOBIN/e-gov-lawapi-v2`）。
バイナリは stdio 経由の MCP server として起動します。

## クライアントへの登録

### Claude Code

```bash
claude mcp add e-gov-lawapi-v2 -- /path/to/e-gov-lawapi-v2/e-gov-lawapi-v2
```

### Claude Desktop

`~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) に追加：

```json
{
  "mcpServers": {
    "e-gov-lawapi-v2": {
      "command": "/path/to/e-gov-lawapi-v2/e-gov-lawapi-v2"
    }
  }
}
```

## 使い方の例

```text
search_laws law_title="個人情報の保護に関する法律" limit=1
```

返却された `law_id` を使って本文を取得：

```text
get_law_data law_id_or_num_or_revision_id="415AC0000000057" json_format="light"
```

特定の条項のみ取得：

```text
get_law_data law_id_or_num_or_revision_id="415AC0000000057" elm="MainProvision-Article[3]"
```

法令本文を全文検索：

```text
keyword_search keyword="個人情報 AND 第三者提供" sentence_text_size=200
```

## 開発

```bash
go run ./...                         # ソースから直接起動
go build -o e-gov-lawapi-v2 ./...    # バイナリビルド
go vet ./...                         # 静的検査
```

stdio で対話する簡易動作確認：

```bash
(printf '%s\n' \
  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"0"}}}' \
  '{"jsonrpc":"2.0","method":"notifications/initialized"}' \
  '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'; sleep 1) \
| ./e-gov-lawapi-v2
```

## 備考

- 本サーバは e-Gov 法令API v2 のうち、JSON で返却される 4 エンドポイントをラップしています。バイナリを返す `/attachment/{law_revision_id}` および `/law_file/{file_type}/{law_id_or_num_or_revision_id}` は対象外です。
- 法令本文（`get_law_data`）はサイズが大きくなる場合があります。トークン消費を抑えるため `elm` で条項を絞り込むか、`json_format=light` の利用を推奨します。
- API の利用条件・制限・サーバメンテナンス情報は [e-Gov 法令検索](https://laws.e-gov.go.jp/) を参照してください。

## ライセンス

MIT
