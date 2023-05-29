package args

import (
	"text/template"
)

const (
	proto = `
  // {{ .Single }}
  rpc Act{{ .Single }} ({{ .Single }}Request) returns ({{ .Single }}Reply);
  rpc List{{ .Plural }} (List{{ .Plural }}Request) returns (List{{ .Plural }}Reply);

// {{ .Single }}
message {{ .Single }}Request {
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

var ProtoTemp, _ = template.New("proto").Parse(proto)
