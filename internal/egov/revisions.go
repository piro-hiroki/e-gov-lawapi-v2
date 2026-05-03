package egov

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
)

// LawRevisionsParams are the query parameters for GET /law_revisions/{law_id_or_num}.
type LawRevisionsParams struct {
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

// ErrLawIDOrNumRequired is returned when a request requires law_id_or_num
// (or its revision-id variant) but it is empty.
var ErrLawIDOrNumRequired = errors.New("law_id_or_num is required")

// GetLawRevisions calls GET /law_revisions/{law_id_or_num}.
func (c *Client) GetLawRevisions(ctx context.Context, p LawRevisionsParams) (json.RawMessage, error) {
	if p.LawIDOrNum == "" {
		return nil, ErrLawIDOrNumRequired
	}
	q := newQuery()
	q.set("law_title", p.LawTitle)
	q.set("law_title_kana", p.LawTitleKana)
	q.set("amendment_law_id", p.AmendmentLawID)
	q.set("amendment_law_num", p.AmendmentLawNum)
	q.set("amendment_law_title", p.AmendmentLawTitle)
	q.set("amendment_date_from", p.AmendmentDateFrom)
	q.set("amendment_date_to", p.AmendmentDateTo)
	q.set("amendment_promulgate_date_from", p.AmendmentPromulgateDateFrom)
	q.set("amendment_promulgate_date_to", p.AmendmentPromulgateDateTo)
	q.setStringSlice("amendment_type", p.AmendmentType)
	q.setStringSlice("current_revision_status", p.CurrentRevisionStatus)
	q.setStringSlice("mission", p.Mission)
	q.setBoolPtr("remain_in_force", p.RemainInForce)
	q.set("repeal_date_from", p.RepealDateFrom)
	q.set("repeal_date_to", p.RepealDateTo)
	q.setStringSlice("repeal_status", p.RepealStatus)
	q.setStringSlice("category_cd", p.CategoryCD)
	q.set("updated_from", p.UpdatedFrom)
	q.set("updated_to", p.UpdatedTo)
	return c.get(ctx, "/law_revisions/"+url.PathEscape(p.LawIDOrNum), q)
}
