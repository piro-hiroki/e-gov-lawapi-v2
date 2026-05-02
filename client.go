package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL   = "https://laws.e-gov.go.jp/api/2"
	userAgent = "e-gov-mcp/0.1.0"
)

type apiClient struct {
	http *http.Client
}

func newAPIClient() *apiClient {
	return &apiClient{
		http: &http.Client{Timeout: 60 * time.Second},
	}
}

// apiError represents a non-2xx response from the e-Gov API.
type apiError struct {
	Status int
	Body   string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("e-Gov API HTTP %d: %s", e.Status, e.Body)
}

// queryBuilder accumulates query parameters, skipping zero values.
type queryBuilder struct {
	values url.Values
}

func newQuery() *queryBuilder {
	return &queryBuilder{values: url.Values{}}
}

func (q *queryBuilder) addString(key, v string) {
	if v != "" {
		q.values.Set(key, v)
	}
}

func (q *queryBuilder) addInt(key string, v int) {
	if v != 0 {
		q.values.Set(key, strconv.Itoa(v))
	}
}

// addBoolPtr writes the parameter only when the user explicitly set it.
func (q *queryBuilder) addBoolPtr(key string, v *bool) {
	if v != nil {
		q.values.Set(key, strconv.FormatBool(*v))
	}
}

func (q *queryBuilder) addStringSlice(key string, v []string) {
	if len(v) == 0 {
		return
	}
	q.values.Set(key, strings.Join(v, ","))
}

func (q *queryBuilder) encode() string {
	return q.values.Encode()
}

func (c *apiClient) get(ctx context.Context, path string, q *queryBuilder) (json.RawMessage, error) {
	q.addString("response_format", "json")
	full := baseURL + path
	if encoded := q.encode(); encoded != "" {
		full += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, full, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &apiError{Status: resp.StatusCode, Body: string(body)}
	}

	if !json.Valid(body) {
		return nil, fmt.Errorf("invalid JSON from e-Gov API: %s", string(body))
	}
	return body, nil
}

// --- Endpoint wrappers -----------------------------------------------------

type searchLawsArgs struct {
	LawTitle                string   `json:"law_title,omitempty" jsonschema:"法令名又は法令略称（部分一致）。例: '個人情報の保護に関する法律'"`
	LawTitleKana            string   `json:"law_title_kana,omitempty" jsonschema:"法令名読み（部分一致、ひらがな）"`
	LawID                   string   `json:"law_id,omitempty" jsonschema:"法令ID（部分一致）。例: '322CO0000000016'"`
	LawNum                  string   `json:"law_num,omitempty" jsonschema:"法令番号（部分一致）。例: '昭和二十二年政令第十六号'"`
	LawNumEra               string   `json:"law_num_era,omitempty" jsonschema:"法令番号の元号 (Meiji|Taisho|Showa|Heisei|Reiwa)"`
	LawNumYear              int      `json:"law_num_year,omitempty" jsonschema:"法令番号の年"`
	LawNumNum               string   `json:"law_num_num,omitempty" jsonschema:"法令番号の号数"`
	LawNumType              string   `json:"law_num_type,omitempty" jsonschema:"法令番号の法令種別"`
	LawType                 []string `json:"law_type,omitempty" jsonschema:"法令種別（複数指定可）。例: ['Act','Rule']"`
	CategoryCD              []string `json:"category_cd,omitempty" jsonschema:"事項別分類コード（複数指定可）"`
	AmendmentLawID          string   `json:"amendment_law_id,omitempty" jsonschema:"改正法令の法令ID（部分一致）"`
	Asof                    string   `json:"asof,omitempty" jsonschema:"法令の時点 YYYY-MM-DD"`
	PromulgationDateFrom    string   `json:"promulgation_date_from,omitempty" jsonschema:"公布日 下限 YYYY-MM-DD"`
	PromulgationDateTo      string   `json:"promulgation_date_to,omitempty" jsonschema:"公布日 上限 YYYY-MM-DD"`
	RepealStatus            []string `json:"repeal_status,omitempty" jsonschema:"廃止等の状態（複数指定可）。例: ['Repeal','Expire']"`
	Mission                 []string `json:"mission,omitempty" jsonschema:"新規制定/一部改正（複数指定可）。例: ['New','Partial']"`
	OmitCurrentRevisionInfo *bool    `json:"omit_current_revision_info,omitempty" jsonschema:"true で current_revision_info を省略"`
	Limit                   int      `json:"limit,omitempty" jsonschema:"取得件数上限（既定 100）"`
	Offset                  int      `json:"offset,omitempty" jsonschema:"取得開始位置（既定 0）"`
	Order                   string   `json:"order,omitempty" jsonschema:"並び順。例: '-revision_info.amendment_promulgate_date'"`
}

func (c *apiClient) searchLaws(ctx context.Context, a searchLawsArgs) (json.RawMessage, error) {
	q := newQuery()
	q.addString("law_title", a.LawTitle)
	q.addString("law_title_kana", a.LawTitleKana)
	q.addString("law_id", a.LawID)
	q.addString("law_num", a.LawNum)
	q.addString("law_num_era", a.LawNumEra)
	q.addInt("law_num_year", a.LawNumYear)
	q.addString("law_num_num", a.LawNumNum)
	q.addString("law_num_type", a.LawNumType)
	q.addStringSlice("law_type", a.LawType)
	q.addStringSlice("category_cd", a.CategoryCD)
	q.addString("amendment_law_id", a.AmendmentLawID)
	q.addString("asof", a.Asof)
	q.addString("promulgation_date_from", a.PromulgationDateFrom)
	q.addString("promulgation_date_to", a.PromulgationDateTo)
	q.addStringSlice("repeal_status", a.RepealStatus)
	q.addStringSlice("mission", a.Mission)
	q.addBoolPtr("omit_current_revision_info", a.OmitCurrentRevisionInfo)
	q.addInt("limit", a.Limit)
	q.addInt("offset", a.Offset)
	q.addString("order", a.Order)
	return c.get(ctx, "/laws", q)
}

type getLawRevisionsArgs struct {
	LawIDOrNum                  string   `json:"law_id_or_num" jsonschema:"法令ID又は法令番号（完全一致）。例: '503AC0000000036'"`
	LawTitle                    string   `json:"law_title,omitempty" jsonschema:"法令名（部分一致または /.../ で正規表現完全一致）"`
	LawTitleKana                string   `json:"law_title_kana,omitempty" jsonschema:"法令名読み（部分一致）"`
	AmendmentLawID              string   `json:"amendment_law_id,omitempty" jsonschema:"改正法令の法令ID（部分一致）"`
	AmendmentLawNum             string   `json:"amendment_law_num,omitempty" jsonschema:"改正法令の法令番号（部分一致）"`
	AmendmentLawTitle           string   `json:"amendment_law_title,omitempty" jsonschema:"改正法令の法令名（部分一致または正規表現）"`
	AmendmentDateFrom           string   `json:"amendment_date_from,omitempty" jsonschema:"改正法令施行日 下限 YYYY-MM-DD"`
	AmendmentDateTo             string   `json:"amendment_date_to,omitempty" jsonschema:"改正法令施行日 上限 YYYY-MM-DD"`
	AmendmentPromulgateDateFrom string   `json:"amendment_promulgate_date_from,omitempty" jsonschema:"改正法令公布日 下限 YYYY-MM-DD"`
	AmendmentPromulgateDateTo   string   `json:"amendment_promulgate_date_to,omitempty" jsonschema:"改正法令公布日 上限 YYYY-MM-DD"`
	AmendmentType               []string `json:"amendment_type,omitempty" jsonschema:"改正種別（複数指定可）"`
	CurrentRevisionStatus       []string `json:"current_revision_status,omitempty" jsonschema:"履歴の状態（複数指定可）。例: ['CurrentEnforced']"`
	Mission                     []string `json:"mission,omitempty" jsonschema:"新規制定/一部改正（複数指定可）"`
	RemainInForce               *bool    `json:"remain_in_force,omitempty" jsonschema:"廃止後も効力を有するもの"`
	RepealDateFrom              string   `json:"repeal_date_from,omitempty" jsonschema:"廃止日 下限 YYYY-MM-DD"`
	RepealDateTo                string   `json:"repeal_date_to,omitempty" jsonschema:"廃止日 上限 YYYY-MM-DD"`
	RepealStatus                []string `json:"repeal_status,omitempty" jsonschema:"廃止等の状態（複数指定可）"`
	CategoryCD                  []string `json:"category_cd,omitempty" jsonschema:"事項別分類コード（複数指定可）"`
	UpdatedFrom                 string   `json:"updated_from,omitempty" jsonschema:"更新日 下限 YYYY-MM-DD"`
	UpdatedTo                   string   `json:"updated_to,omitempty" jsonschema:"更新日 上限 YYYY-MM-DD"`
}

func (c *apiClient) getLawRevisions(ctx context.Context, a getLawRevisionsArgs) (json.RawMessage, error) {
	if a.LawIDOrNum == "" {
		return nil, fmt.Errorf("law_id_or_num is required")
	}
	q := newQuery()
	q.addString("law_title", a.LawTitle)
	q.addString("law_title_kana", a.LawTitleKana)
	q.addString("amendment_law_id", a.AmendmentLawID)
	q.addString("amendment_law_num", a.AmendmentLawNum)
	q.addString("amendment_law_title", a.AmendmentLawTitle)
	q.addString("amendment_date_from", a.AmendmentDateFrom)
	q.addString("amendment_date_to", a.AmendmentDateTo)
	q.addString("amendment_promulgate_date_from", a.AmendmentPromulgateDateFrom)
	q.addString("amendment_promulgate_date_to", a.AmendmentPromulgateDateTo)
	q.addStringSlice("amendment_type", a.AmendmentType)
	q.addStringSlice("current_revision_status", a.CurrentRevisionStatus)
	q.addStringSlice("mission", a.Mission)
	q.addBoolPtr("remain_in_force", a.RemainInForce)
	q.addString("repeal_date_from", a.RepealDateFrom)
	q.addString("repeal_date_to", a.RepealDateTo)
	q.addStringSlice("repeal_status", a.RepealStatus)
	q.addStringSlice("category_cd", a.CategoryCD)
	q.addString("updated_from", a.UpdatedFrom)
	q.addString("updated_to", a.UpdatedTo)
	return c.get(ctx, "/law_revisions/"+url.PathEscape(a.LawIDOrNum), q)
}

type getLawDataArgs struct {
	LawIDOrNumOrRevisionID      string `json:"law_id_or_num_or_revision_id" jsonschema:"法令ID/法令番号/法令履歴ID のいずれか（完全一致）。例: '411AC0000000127'"`
	Asof                        string `json:"asof,omitempty" jsonschema:"法令の時点 YYYY-MM-DD（law_revision_id 指定時は無視）"`
	Elm                         string `json:"elm,omitempty" jsonschema:"本文の一部のみ取得。例: 'MainProvision-Article[3]'"`
	JSONFormat                  string `json:"json_format,omitempty" jsonschema:"law_full_text の JSON 形式 (full|light)。light は簡易形式"`
	LawFullTextFormat           string `json:"law_full_text_format,omitempty" jsonschema:"law_full_text のレスポンス形式 (json|xml)"`
	OmitAmendmentSupplProvision *bool  `json:"omit_amendment_suppl_provision,omitempty" jsonschema:"true で改正法令の附則を含めない"`
	IncludeAttachedFileContent  *bool  `json:"include_attached_file_content,omitempty" jsonschema:"true で添付ファイルの image_data も返却"`
}

func (c *apiClient) getLawData(ctx context.Context, a getLawDataArgs) (json.RawMessage, error) {
	if a.LawIDOrNumOrRevisionID == "" {
		return nil, fmt.Errorf("law_id_or_num_or_revision_id is required")
	}
	q := newQuery()
	q.addString("asof", a.Asof)
	q.addString("elm", a.Elm)
	q.addString("json_format", a.JSONFormat)
	q.addString("law_full_text_format", a.LawFullTextFormat)
	q.addBoolPtr("omit_amendment_suppl_provision", a.OmitAmendmentSupplProvision)
	q.addBoolPtr("include_attached_file_content", a.IncludeAttachedFileContent)
	return c.get(ctx, "/law_data/"+url.PathEscape(a.LawIDOrNumOrRevisionID), q)
}

type keywordSearchArgs struct {
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

func (c *apiClient) keywordSearch(ctx context.Context, a keywordSearchArgs) (json.RawMessage, error) {
	if a.Keyword == "" {
		return nil, fmt.Errorf("keyword is required")
	}
	q := newQuery()
	q.addString("keyword", a.Keyword)
	q.addStringSlice("law_type", a.LawType)
	q.addString("law_num", a.LawNum)
	q.addString("law_num_era", a.LawNumEra)
	q.addInt("law_num_year", a.LawNumYear)
	q.addString("law_num_num", a.LawNumNum)
	q.addString("law_num_type", a.LawNumType)
	q.addStringSlice("category_cd", a.CategoryCD)
	q.addString("asof", a.Asof)
	q.addString("promulgation_date_from", a.PromulgationDateFrom)
	q.addString("promulgation_date_to", a.PromulgationDateTo)
	q.addInt("limit", a.Limit)
	q.addInt("offset", a.Offset)
	q.addString("order", a.Order)
	q.addInt("sentences_limit", a.SentencesLimit)
	q.addInt("sentence_text_size", a.SentenceTextSize)
	q.addString("highlight_tag", a.HighlightTag)
	return c.get(ctx, "/keyword", q)
}
