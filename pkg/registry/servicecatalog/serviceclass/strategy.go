/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package serviceclass

// this was copied from where else and edited to fit our objects

import (
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/validation/field"

	"github.com/golang/glog"
	sc "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog"
	scv "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/validation"
)

// implements interfaces RESTCreateStrategy, RESTUpdateStrategy, RESTDeleteStrategy
type serviceclassRESTStrategy struct {
	runtime.ObjectTyper // inherit ObjectKinds method
	kapi.NameGenerator  // GenerateName method for CreateStrategy
}

var (
	serviceclassRESTStrategies = serviceclassRESTStrategy{
		// embeds to pull in existing code behavior from upstream

		ObjectTyper: kapi.Scheme,
		// use the generator from upstream k8s, or implement method
		// `GenerateName(base string) string`
		NameGenerator: kapi.SimpleNameGenerator,
	}
	_ rest.RESTCreateStrategy = serviceclassRESTStrategies
	_ rest.RESTUpdateStrategy = serviceclassRESTStrategies
	_ rest.RESTDeleteStrategy = serviceclassRESTStrategies
)

// Canonicalize does not transform a serviceclass.
func (serviceclassRESTStrategy) Canonicalize(obj runtime.Object) {
	_, ok := obj.(*sc.ServiceClass)
	if !ok {
		glog.Warning("received a non-serviceclass object to create")
	}
}

// NamespaceScoped returns false as serviceclasss are not scoped to a namespace.
func (serviceclassRESTStrategy) NamespaceScoped() bool {
	return false
}

// PrepareForCreate receives the incoming Serviceclass.
func (serviceclassRESTStrategy) PrepareForCreate(ctx kapi.Context, obj runtime.Object) {
	_, ok := obj.(*sc.ServiceClass)
	if !ok {
		glog.Warning("received a non-serviceclass object to create")
	}
	// service class is a data record and has no status to track
}

func (serviceclassRESTStrategy) Validate(ctx kapi.Context, obj runtime.Object) field.ErrorList {
	return scv.ValidateServiceclass(obj.(*sc.ServiceClass))
}

func (serviceclassRESTStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (serviceclassRESTStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (serviceclassRESTStrategy) PrepareForUpdate(ctx kapi.Context, new, old runtime.Object) {
	newServiceclass, ok := new.(*sc.ServiceClass)
	if !ok {
		glog.Warning("received a non-serviceclass object to update to")
	}
	oldServiceclass, ok := old.(*sc.ServiceClass)
	if !ok {
		glog.Warning("received a non-serviceclass object to update from")
	}
	// copy all fields individually?
	newServiceclass.BrokerName = oldServiceclass.BrokerName
}

func (serviceclassRESTStrategy) ValidateUpdate(ctx kapi.Context, new, old runtime.Object) field.ErrorList {
	newServiceclass, ok := new.(*sc.ServiceClass)
	if !ok {
		glog.Warning("received a non-serviceclass object to validate to")
	}
	oldServiceclass, ok := old.(*sc.ServiceClass)
	if !ok {
		glog.Warning("received a non-serviceclass object to validate from")
	}

	return scv.ValidateServiceclassUpdate(newServiceclass, oldServiceclass)
}
