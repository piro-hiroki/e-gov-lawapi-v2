package egov

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
)

// LawDataParams are the query parameters for GET /law_data/{...}.
type LawDataParams struct {
	LawIDOrNumOrRevisionID      string `json:"law_id_or_num_or_revision_id" jsonschema:"法令ID/法令番号/法令履歴ID のいずれか（完全一致）。例: '411AC0000000127'"`
	Asof                        string `json:"asof,omitempty" jsonschema:"法令の時点 YYYY-MM-DD（law_revision_id 指定時は無視）"`
	Elm                         string `json:"elm,omitempty" jsonschema:"本文の一部のみ取得。例: 'MainProvision-Article[3]'"`
	JSONFormat                  string `json:"json_format,omitempty" jsonschema:"law_full_text の JSON 形式 (full|light)。light は簡易形式"`
	LawFullTextFormat           string `json:"law_full_text_format,omitempty" jsonschema:"law_full_text のレスポンス形式 (json|xml)"`
	OmitAmendmentSupplProvision *bool  `json:"omit_amendment_suppl_provision,omitempty" jsonschema:"true で改正法令の附則を含めない"`
	IncludeAttachedFileContent  *bool  `json:"include_attached_file_content,omitempty" jsonschema:"true で添付ファイルの image_data も返却"`
}

// ErrLawDataIDRequired is returned when LawDataParams.LawIDOrNumOrRevisionID is empty.
var ErrLawDataIDRequired = errors.New("law_id_or_num_or_revision_id is required")

// GetLawData calls GET /law_data/{law_id_or_num_or_revision_id}.
func (c *Client) GetLawData(ctx context.Context, p LawDataParams) (json.RawMessage, error) {
	if p.LawIDOrNumOrRevisionID == "" {
		return nil, ErrLawDataIDRequired
	}
	q := newQuery()
	q.set("asof", p.Asof)
	q.set("elm", p.Elm)
	q.set("json_format", p.JSONFormat)
	q.set("law_full_text_format", p.LawFullTextFormat)
	q.setBoolPtr("omit_amendment_suppl_provision", p.OmitAmendmentSupplProvision)
	q.setBoolPtr("include_attached_file_content", p.IncludeAttachedFileContent)
	return c.get(ctx, "/law_data/"+url.PathEscape(p.LawIDOrNumOrRevisionID), q)
}
