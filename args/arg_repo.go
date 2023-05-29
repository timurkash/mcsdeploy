package args

import (
	"text/template"
)

const (
	repo = `type {{ .Single }}Repo interface {
	Get{{ .Single }}(context.Context, uint32) (*pb.{{ .Single }}Reply, error)
	Create{{ .Single }}(context.Context, *pb.{{ .Single }}Info) (*pb.{{ .Single }}Reply, error)
	Update{{ .Single }}(context.Context, uint32, *pb.{{ .Single }}Info) (*pb.{{ .Single }}Reply, error)
	List{{ .Plural }}(context.Context, *common.Filter, *common.OrderOffsetLimit) ([]*pb.{{ .Single }}Reply, *common.Paging, error)
}
`
)

var RepoTemp, _ = template.New("repo").Parse(repo)
