package args

import "text/template"

var ProtoTemp, _ = template.New("proto").Parse(proto)

const (
	proto = `
// --- services
  // {{ .Single }}
  rpc Act{{ .Single }} (Act{{ .Single }}Request) returns ({{ .Single }}Reply);
  rpc List{{ .Plural }} (List{{ .Plural }}Request) returns (List{{ .Plural }}Reply);

// --- messages

// {{ .Single }}
message Act{{ .Single }}Request {
  common.ActionId action_id = 1;
  {{ .Single }}Info {{ .SingleLower }} = 2;
}

message {{ .Single }}Reply {
  {{ .Single }}Info {{ .SingleLower }} = 1;
  common.IdTimestamps id_timestamps = 2;
}

message {{ .Single }}Info {
  string name = 1;
}

message List{{ .Plural }}Request {
  common.Filter filter = 1;
  common.OrderOffsetLimit ool = 2;
}

message List{{ .Plural }}Reply {
  repeated {{ .Single }}Reply {{ .PluralLower }} = 1;
  common.Paging paging = 2;
}
`
)
