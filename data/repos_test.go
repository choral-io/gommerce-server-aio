package data_test

import (
	"testing"
)

type baseRepo interface {
	mustBeBaseRepo()
	rename(string)
	Name() string
}

type BaseRepo struct {
	name string
}

func (r *BaseRepo) mustBeBaseRepo() {}

func (r *BaseRepo) rename(name string) {
	r.name = name
}

func (r *BaseRepo) Name() string {
	return r.name
}

type DemoRepo interface {
	baseRepo
}

type demoRepo struct {
	BaseRepo
}

func NewDemoRepo[R interface {
	DemoRepo
	*demoRepo
}](name string) R {
	return &demoRepo{
		BaseRepo: BaseRepo{
			name: name,
		},
	}
}

func shallowCopy[S any, P interface {
	*S
	baseRepo
}](r P) P {
	copy := *r
	return &copy
}

func TestShallowCopy(t *testing.T) {
	odr := NewDemoRepo("old")
	ndr := shallowCopy(odr)
	ndr.rename("new")
	if odr.Name() != "old" {
		t.Errorf("expected old name to be old, got %s", odr.Name())
	}
	if ndr.Name() != "new" {
		t.Errorf("expected new name to be new, got %s", ndr.Name())
	}
}
