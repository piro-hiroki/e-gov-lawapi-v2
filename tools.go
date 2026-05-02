package main

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/piro-hiroki/e-gov-lawapi-v2/internal/egov"
)

func registerTools(server *mcp.Server, api *egov.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:  "search_laws",
		Title: "法令一覧検索",
		Description: "e-Gov 法令API v2 の /laws を呼び出し、条件に該当する法令の一覧を取得します。" +
			"法令名（部分一致）、法令種別、公布日範囲、事項別分類等で絞り込めます。" +
			"返却された law_id / law_num / law_revision_id を後続のツールに渡して本文や改正履歴を取得できます。",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args egov.SearchLawsParams) (*mcp.CallToolResult, any, error) {
		return runJSONTool(ctx, args, api.SearchLaws)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:  "get_law_revisions",
		Title: "法令改正履歴の取得",
		Description: "/law_revisions/{law_id_or_num} を呼び出し、指定法令の改正履歴を取得します。" +
			"law_id_or_num は法令ID又は法令番号を完全一致で指定します。",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args egov.LawRevisionsParams) (*mcp.CallToolResult, any, error) {
		return runJSONTool(ctx, args, api.GetLawRevisions)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:  "get_law_data",
		Title: "法令本文の取得",
		Description: "/law_data/{...} を呼び出し、法令の本文（条文を含む）を取得します。" +
			"law_id, law_num, law_revision_id のいずれかを完全一致で指定します。" +
			"本文サイズが大きくなる場合があるため、必要に応じて elm パラメータで条項を絞り込み、" +
			"json_format='light' でパースしやすい簡易形式に切り替えるのを推奨します。",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args egov.LawDataParams) (*mcp.CallToolResult, any, error) {
		return runJSONTool(ctx, args, api.GetLawData)
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:  "keyword_search",
		Title: "法令本文のキーワード検索",
		Description: "/keyword を呼び出し、法令本文を全文検索します。AND/OR/NOT、ワイルドカード（* と ?）に対応。" +
			"ヒットした条文の位置（position）と一部テキスト（text）が返却されます。" +
			"sentence_text_size で text 長、highlight_tag でハイライトタグを制御できます。",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args egov.KeywordSearchParams) (*mcp.CallToolResult, any, error) {
		return runJSONTool(ctx, args, api.KeywordSearch)
	})
}

// runJSONTool invokes fn and packages the JSON response (or error) as a CallToolResult.
func runJSONTool[T any](
	ctx context.Context,
	args T,
	fn func(context.Context, T) (json.RawMessage, error),
) (*mcp.CallToolResult, any, error) {
	body, err := fn(ctx, args)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		}, nil, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(body)}},
	}, nil, nil
}
