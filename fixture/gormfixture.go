package fixture

import (
	"errors"

	"gorm.io/gorm"
)

type GormFixture interface {
	DB() *gorm.DB
	Reference() (string, error)
	SetReference(reference string)
	Create() (Model, error)
}

type NestedBuilder interface {
	AddDependencies(dependencies ...NestedBuilder)
	Build() error
}

type BaseNestedFixture struct {
	GormFixture
	db           *gorm.DB
	store        Store
	reference    string
	dependencies []NestedBuilder
}

func NewBaseNestedFixture(db *gorm.DB, store Store) *BaseNestedFixture {
	return &BaseNestedFixture{
		db:    db,
		store: store,
	}
}

func (b *BaseNestedFixture) SetFoxFixture(foxFixture GormFixture) {
	b.GormFixture = foxFixture
}

func (b *BaseNestedFixture) DB() *gorm.DB {
	return b.db
}

func (b *BaseNestedFixture) Reference() (string, error) {
	if b.reference == "" {
		return "", errors.New("reference of this fixture is not defined yet")
	}

	return b.reference, nil
}

func (b *BaseNestedFixture) SetReference(reference string) {
	b.reference = reference
}

func (b *BaseNestedFixture) AddDependencies(dependencies ...NestedBuilder) {
	b.dependencies = append(b.dependencies, dependencies...)
}

func (b *BaseNestedFixture) BuildDependencies() error {
	for _, dependency := range b.dependencies {
		err := dependency.Build()
		if nil != err {
			return err
		}
	}

	return nil
}

func (b *BaseNestedFixture) Build() error {
	_, err := b.GetFixture(b.reference)
	if err == nil {
		return nil
	}
	err = b.BuildDependencies()
	if nil != err {
		return err
	}
	fixture, err := b.Create()
	if nil != err {
		return err
	}

	return b.store.Set(b.reference, fixture)
}

func (b *BaseNestedFixture) GetFixture(reference string) (Model, error) {
	fixture, err := b.store.Get(reference)
	if err != nil {
		return nil, err
	}

	return fixture, nil
}
