package repos

import (
	"context"

	"github.com/uptrace/bun"
)

// SelectQueryTransformer is an interface for transforming bun.SelectQuery.
type SelectQueryTransformer interface {
	Transform(context.Context, *bun.SelectQuery) (*bun.SelectQuery, error)
}

// SelectQueryTransformerFunc is a function that transforms bun.SelectQuery.
type SelectQueryTransformerFunc func(context.Context, *bun.SelectQuery) (*bun.SelectQuery, error)

// SelectQueryTransformerFunc is a function that implements SelectQueryTransformer.
func (f SelectQueryTransformerFunc) Transform(ctx context.Context, query *bun.SelectQuery) (*bun.SelectQuery, error) {
	return f(ctx, query)
}

// TransformSelectQuery applies the given transformers to the query.
func TransformSelectQuery(ctx context.Context, query *bun.SelectQuery, sqts ...SelectQueryTransformer) (*bun.SelectQuery, error) {
	for _, sqt := range sqts {
		var err error
		query, err = sqt.Transform(ctx, query)
		if err != nil {
			return nil, err
		}
	}
	return query, nil
}

// WithPagination is a bun.SelectQuery transformer that adds paging to the query.
// The page and size can be specified by the given interface.
func WithPagination(p interface {
	GetPage() int32
	GetSize() int32
}) SelectQueryTransformer {
	return SelectQueryTransformerFunc(func(ctx context.Context, sq *bun.SelectQuery) (*bun.SelectQuery, error) {
		var size = 10
		if p.GetSize() > 0 {
			size = int(p.GetSize())
		}
		var page = 1
		if p.GetPage() > 1 {
			page = int(p.GetPage())
		}
		return sq.Offset((page - 1) * size).Limit(size), nil
	})
}

// WithColumns is a bun.SelectQuery transformer that adds columns to the query.
// The columns to include/exclude can be specified as variadic arguments.
// To exclude a column, prefix it with a minus sign.
//
// - Parameters:
//   - columns: The columns to include/exclude.
//
// Example 1: WithColumns("id", "name")
//
// Example 2: WithColumns("id", "-name")
//
// Example 3: WithColumns("-*", "name")
func WithColumns(columns ...string) SelectQueryTransformer {
	return SelectQueryTransformerFunc(func(ctx context.Context, sq *bun.SelectQuery) (*bun.SelectQuery, error) {
		if len(columns) > 0 {
			includes := make([]string, 0, len(columns))
			excludes := make([]string, 0, len(columns))
			for _, col := range columns {
				if col != "" {
					if col[0] == '-' {
						excludes = append(excludes, col[1:])
					} else {
						includes = append(includes, col)
					}
				}
			}
			return sq.Column(includes...).ExcludeColumn(excludes...), nil
		} else {
			return sq, nil
		}
	})
}

// WithRelation is a bun.SelectQuery transformer that adds a relation to the query.
// The columns to include/exclude can be specified as variadic arguments.
// To exclude a column, prefix it with a minus sign.
//
// - Parameters:
//   - relation: The relation name.
//   - columns: The columns to include/exclude.
//
// Example 1: WithRelation("Realm")
//
// Example 2: WithRelation("Realm", "name")
//
// Example 3: WithRelation("Realm", "-*", "name")
func WithRelation(relation string, columns ...string) SelectQueryTransformer {
	return SelectQueryTransformerFunc(func(ctx context.Context, sq *bun.SelectQuery) (*bun.SelectQuery, error) {
		if len(columns) > 0 {
			return sq.Relation(relation, func(ssq *bun.SelectQuery) *bun.SelectQuery {
				includes := make([]string, 0, len(columns))
				excludes := make([]string, 0, len(columns))
				for _, col := range columns {
					if col != "" {
						if col[0] == '-' {
							excludes = append(excludes, col[1:])
						} else {
							includes = append(includes, col)
						}
					}
				}
				return ssq.Column(includes...).ExcludeColumn(excludes...)
			}), nil
		} else {
			return sq.Relation(relation), nil
		}
	})
}
