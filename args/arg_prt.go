package args

import (
	"github.com/stoewer/go-strcase"
	"os"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
  message Filter {
    common.String name = 1;
  }
  Filter filter = 1;
  common.OrderOffsetLimit ool = 2;
}

message List{{ .Plural }}Reply {
  repeated {{ .Single }}Reply {{ .PluralLower }} = 1;
  common.Paging paging = 2;
}
`
)

type SinglePlural struct {
	Single      string
	SingleLower string
	Plural      string
	PluralLower string
}

var cas = cases.Title(language.English)

func ArgProto(single, plural string) error {
	if plural == "" {
		plural = single + "s"
	}
	sp := &SinglePlural{
		Single:      strcase.UpperCamelCase(cas.String(single)),
		SingleLower: single,
		Plural:      strcase.UpperCamelCase(cas.String(plural)),
		PluralLower: plural,
	}
	tempProto, err := template.New("proto").Parse(proto)
	if err != nil {
		return err
	}
	if err := tempProto.Execute(os.Stdout, sp); err != nil {
		return err
	}
	return nil
}
