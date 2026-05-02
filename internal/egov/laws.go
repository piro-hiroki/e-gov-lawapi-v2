package egov

import (
	"context"
	"encoding/json"
)

// SearchLawsParams are the query parameters for GET /laws.
//
// All fields are optional. The e-Gov API supports partial matches on string
// fields and "OR" semantics across multiple values for slice fields.
type SearchLawsParams struct {
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

// SearchLaws calls GET /laws.
func (c *Client) SearchLaws(ctx context.Context, p SearchLawsParams) (json.RawMessage, error) {
	q := newQuery()
	q.set("law_title", p.LawTitle)
	q.set("law_title_kana", p.LawTitleKana)
	q.set("law_id", p.LawID)
	q.set("law_num", p.LawNum)
	q.set("law_num_era", p.LawNumEra)
	q.setInt("law_num_year", p.LawNumYear)
	q.set("law_num_num", p.LawNumNum)
	q.set("law_num_type", p.LawNumType)
	q.setStringSlice("law_type", p.LawType)
	q.setStringSlice("category_cd", p.CategoryCD)
	q.set("amendment_law_id", p.AmendmentLawID)
	q.set("asof", p.Asof)
	q.set("promulgation_date_from", p.PromulgationDateFrom)
	q.set("promulgation_date_to", p.PromulgationDateTo)
	q.setStringSlice("repeal_status", p.RepealStatus)
	q.setStringSlice("mission", p.Mission)
	q.setBoolPtr("omit_current_revision_info", p.OmitCurrentRevisionInfo)
	q.setInt("limit", p.Limit)
	q.setInt("offset", p.Offset)
	q.set("order", p.Order)
	return c.get(ctx, "/laws", q)
}
