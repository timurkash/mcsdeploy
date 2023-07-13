package args

import "text/template"

var RepoTemp, _ = template.New("repo").Parse(repo)

const (
	repo = `
// --- services

func (s *{{ .Service }}Service) Act{{ .Single }}(ctx context.Context, req *pb.Act{{ .Single }}Request) (*pb.{{ .Single }}Reply, error) {
	return s.uc.Act{{ .Single }}(ctx, req.ActionId, req.{{ .Single }})
}

func (s *{{ .Service }}Service) List{{ .Plural }}(ctx context.Context, req *pb.List{{ .Plural }}Request) (*pb.List{{ .Plural }}Reply, error) {
	return s.uc.List{{ .Plural }}(ctx, req.Filter, req.Ool)
}

// --- biz interface

type {{ .Service }}Repo interface {
	Get{{ .Single }}(context.Context, uint32) (*pb.{{ .Single }}Reply, error)
	Create{{ .Single }}(context.Context, *pb.{{ .Single }}Info) (*pb.{{ .Single }}Reply, error)
	Update{{ .Single }}(context.Context, uint32, *pb.{{ .Single }}Info) (*pb.{{ .Single }}Reply, error)
	List{{ .Plural }}(context.Context, *common.Filter, *common.OrderOffsetLimit) ([]*pb.{{ .Single }}Reply, *common.Paging, error)
}

// --- biz implementation

func (uc *{{ .Service }}Usecase) checkPermission(ctx context.Context) error {
	return jwt.IsPermitted(ctx, jwt.{{ .Single }}Admin)
}

func (uc *{{ .Service }}Usecase) Act{{ .Single }}(ctx context.Context, actionId *common.ActionId, {{ .SingleLower }}Info *pb.{{ .Single }}Info) (*pb.{{ .Single }}Reply, error) {
	if err := cerrors.CheckActionId(actionId, {{ .SingleLower }}Info,
		common.Action_get,
		common.Action_insert,
		common.Action_update); err != nil {
		return nil, err
	}
	if err := uc.checkPermission(ctx); err != nil {
		return nil, err
	}
	switch actionId.Action {
	case common.Action_get:
		return uc.repo.Get{{ .Single }}(ctx, actionId.Id)
	case common.Action_insert:
		return uc.repo.Create{{ .Single }}(ctx, {{ .SingleLower }}Info)
	case common.Action_update:
		return uc.repo.Update{{ .Single }}(ctx, actionId.Id, {{ .SingleLower }}Info)
	}
	return nil, cerrors.GetWrongActionError(actionId.Action)
}

func (uc *{{ .Service }}Usecase) List{{ .Plural }}(
	ctx context.Context,
	filter *common.Filter,
	ool *common.OrderOffsetLimit,
) (*pb.List{{ .Plural }}Reply, error) {
	if err := uc.checkPermission(ctx); err != nil {
		return nil, err
	}
	{{ .PluralLower }}, paging, err := uc.repo.List{{ .Plural }}(ctx, filter, ool)
	if err != nil {
		return nil, err
	}
	return &pb.List{{ .Plural }}Reply{
		Items: items,
		Paging: paging,
	}, nil
}

// --- repo implementation

func (r *{{ .ServiceLower }}Repo) get{{ .Single }}Reply(record *ent.{{ .Single }}) *pb.{{ .Single }}Reply {
	return &pb.{{ .Single }}Reply{
	}
}

func (r *{{ .ServiceLower }}Repo) getMax{{ .Single }}(ctx context.Context) (uint32, error) {
	var v []struct {
		Max int
	}
	if err := r.relational.{{ .Single }}.Query().Aggregate(ent.Max(consts.Id)).Scan(ctx, &v); err != nil {
		return 0, err
	}
	return uint32(v[0].Max), nil
}

func (r *{{ .ServiceLower }}Repo) Get{{ .Single }}(ctx context.Context, id uint32) (*pb.{{ .Single }}Reply, error) {
	record, err := r.relational.{{ .Single }}.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.get{{ .Single }}Reply(record), nil
}

func (r *{{ .ServiceLower }}Repo) Create{{ .Single }}(ctx context.Context, info *pb.{{ .Single }}Info) (*pb.{{ .Single }}Reply, error) {
	max, err := r.getMax{{ .Single }}(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	create{{ .Single }}Record := r.relational.{{ .Single }}.Create().
		SetID(max + 1).
		SetCreatedAt(now).
		SetUpdatedAt(now)
	{{ .SingleLower }}Created, err := create{{ .Single }}Record.Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.get{{ .Single }}Reply({{ .SingleLower }}Created), nil
}

func (r *{{ .ServiceLower }}Repo) Update{{ .Single }}(ctx context.Context, id uint32, info *pb.{{ .Single }}Info) (*pb.{{ .Single }}Reply, error) {
	update{{ .Single }}Record := r.relational.{{ .Single }}.UpdateOneID(id).
		SetUpdatedAt(time.Now())
	{{ .SingleLower }}Updated, err := update{{ .Single }}Record.Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.get{{ .Single }}Reply({{ .SingleLower }}Updated), nil
}

func (r *{{ .ServiceLower }}Repo) List{{ .Plural }}(ctx context.Context, filter *common.Filter, ool *common.OrderOffsetLimit) ([]*pb.{{ .Single }}Reply, *common.Paging, error) {
	{{ .PluralLower }}Query := r.relational.{{ .Single }}.Query()
	if filter != nil {
		if filter.Name != nil {
			// name := filter.Name.Value
			{{ .PluralLower }}Query.Where(
				{{ .SingleLower_ }}.Or(
				),
			)
		}
	}
	total, err := {{ .PluralLower }}Query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	var (
		offset uint32
		limit  uint32
	)
	if ool != nil {
		offset = ool.Offset
		limit = ool.Limit
		{{ .PluralLower }}Query.Offset(int(offset)).Limit(int(limit))
	}
	{{ .PluralLower }}Query.Order(ent.Asc("id"))
	itemsAll, err := {{ .PluralLower }}Query.All(ctx)
	if err != nil {
		return nil, nil, err
	}
	{{ .PluralLower }} := make([]*pb.{{ .Single }}Reply, 0, limit)
	for _, item := range itemsAll {
		{{ .PluralLower }} = append({{ .PluralLower }}, r.get{{ .Single }}Reply(item))
	}
	return {{ .PluralLower }}, &common.Paging{
		Order:   "id",
		Offset:  offset,
		Limit:   limit,
		Total:   uint32(total),
		HasNext: len({{ .PluralLower }})+int(offset) < total,
	}, nil
}
`
)
