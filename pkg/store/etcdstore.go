package store

import (
	"context"
	"fmt"
	"github.com/boyfoo/k8s-aa-basis/pkg/apis/myingress/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

// REST implements a RESTStorage for API services against etcd
type REST struct {
	*genericregistry.Store //默认抽插的 ETCD
}

func (*REST) ShortNames() []string {
	return []string{"mi"}
}

// 没啥特别之处， 就是加了个判断 而已
func RESTInPeace(storage rest.StandardStorage, err error) rest.StandardStorage {
	if err != nil {
		err = fmt.Errorf("unable to create REST storage for a resource due to %v, will die", err)
		panic(err)
	}
	return storage
}

// 构建 myIngress增删改查策略 就是怎么新增、怎么删除、怎么修改
func NewStrategy(typer runtime.ObjectTyper) MyIngressStrategy {
	return MyIngressStrategy{typer, names.SimpleNameGenerator}
}
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {

	apiserver, ok := obj.(*v1beta1.MyIngress)
	if !ok {
		return nil, nil, fmt.Errorf(" object is not a MyIngress")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// 标签 和字段 匹配器
func MatchMyIngress(label labels.Selector, field fields.Selector) storage.SelectionPredicate {

	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *v1beta1.MyIngress) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type MyIngressStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

//更新时发出的警告
func (s MyIngressStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	//TODO implement me
	return []string{}
}

//创建时 是否要发出警告--- 发出个屁
func (s MyIngressStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	//TODO implement me
	return []string{}
}

func (MyIngressStrategy) NamespaceScoped() bool {
	return true
}

//Validate 之前调用
func (MyIngressStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {

}

func (MyIngressStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

//这是字段验证相关
func (MyIngressStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {

	return field.ErrorList{} //不高兴验证
}

func (MyIngressStrategy) AllowCreateOnUpdate() bool {
	return true
}

func (MyIngressStrategy) AllowUnconditionalUpdate() bool {
	return true
}

func (MyIngressStrategy) Canonicalize(obj runtime.Object) {
}

//依然不高兴验证
func (MyIngressStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*REST, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &v1beta1.MyIngress{}
		},
		NewListFunc: func() runtime.Object {
			return &v1beta1.MyIngressList{}
		},
		PredicateFunc:            MatchMyIngress,
		DefaultQualifiedResource: v1beta1.SchemeGroupResource,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: rest.NewDefaultTableConvertor(v1beta1.SchemeGroupResource),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
