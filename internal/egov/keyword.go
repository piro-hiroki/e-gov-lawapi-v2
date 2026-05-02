package egov

import (
	"context"
	"encoding/json"
	"errors"
)

// KeywordSearchParams are the query parameters for GET /keyword.
type KeywordSearchParams struct {
	Keyword              string   `json:"keyword" jsonschema:"検索キーワード（必須）。AND/OR/NOT、ワイルドカード * ?"`
	LawType              []string `json:"law_type,omitempty" jsonschema:"法令種別（複数指定可）"`
	LawNum               string   `json:"law_num,omitempty" jsonschema:"法令番号（部分一致）"`
	LawNumEra            string   `json:"law_num_era,omitempty" jsonschema:"法令番号の元号"`
	LawNumYear           int      `json:"law_num_year,omitempty" jsonschema:"法令番号の年"`
	LawNumNum            string   `json:"law_num_num,omitempty" jsonschema:"法令番号の号数"`
	LawNumType           string   `json:"law_num_type,omitempty" jsonschema:"法令番号の法令種別"`
	CategoryCD           []string `json:"category_cd,omitempty" jsonschema:"事項別分類コード（複数指定可）"`
	Asof                 string   `json:"asof,omitempty" jsonschema:"法令の時点 YYYY-MM-DD"`
	PromulgationDateFrom string   `json:"promulgation_date_from,omitempty" jsonschema:"公布日 下限 YYYY-MM-DD"`
	PromulgationDateTo   string   `json:"promulgation_date_to,omitempty" jsonschema:"公布日 上限 YYYY-MM-DD"`
	Limit                int      `json:"limit,omitempty" jsonschema:"position 数の総和の上限（既定 100、最大 1000）"`
	Offset               int      `json:"offset,omitempty" jsonschema:"取得開始位置（既定 0）"`
	Order                string   `json:"order,omitempty" jsonschema:"並び順"`
	SentencesLimit       int      `json:"sentences_limit,omitempty" jsonschema:"sentences の表示件数制限"`
	SentenceTextSize     int      `json:"sentence_text_size,omitempty" jsonschema:"text の表示文字数（既定 100）"`
	HighlightTag         string   `json:"highlight_tag,omitempty" jsonschema:"ヒット箇所を囲む HTML タグ名（既定 'span'）"`
}

// ErrKeywordRequired is returned when KeywordSearchParams.Keyword is empty.
var ErrKeywordRequired = errors.New("keyword is required")

// KeywordSearch calls GET /keyword.
func (c *Client) KeywordSearch(ctx context.Context, p KeywordSearchParams) (json.RawMessage, error) {
	if p.Keyword == "" {
		return nil, ErrKeywordRequired
	}
	q := newQuery()
	q.set("keyword", p.Keyword)
	q.setStringSlice("law_type", p.LawType)
	q.set("law_num", p.LawNum)
	q.set("law_num_era", p.LawNumEra)
	q.setInt("law_num_year", p.LawNumYear)
	q.set("law_num_num", p.LawNumNum)
	q.set("law_num_type", p.LawNumType)
	q.setStringSlice("category_cd", p.CategoryCD)
	q.set("asof", p.Asof)
	q.set("promulgation_date_from", p.PromulgationDateFrom)
	q.set("promulgation_date_to", p.PromulgationDateTo)
	q.setInt("limit", p.Limit)
	q.setInt("offset", p.Offset)
	q.set("order", p.Order)
	q.setInt("sentences_limit", p.SentencesLimit)
	q.setInt("sentence_text_size", p.SentenceTextSize)
	q.set("highlight_tag", p.HighlightTag)
	return c.get(ctx, "/keyword", q)
}
