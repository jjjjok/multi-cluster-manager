/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "harmonycloud.cn/multi-cluster-manager/pkg/apis/multicluster/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ClusterResourceLister helps list ClusterResources.
// All objects returned here must be treated as read-only.
type ClusterResourceLister interface {
	// List lists all ClusterResources in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ClusterResource, err error)
	// ClusterResources returns an object that can list and get ClusterResources.
	ClusterResources(namespace string) ClusterResourceNamespaceLister
	ClusterResourceListerExpansion
}

// clusterResourceLister implements the ClusterResourceLister interface.
type clusterResourceLister struct {
	indexer cache.Indexer
}

// NewClusterResourceLister returns a new ClusterResourceLister.
func NewClusterResourceLister(indexer cache.Indexer) ClusterResourceLister {
	return &clusterResourceLister{indexer: indexer}
}

// List lists all ClusterResources in the indexer.
func (s *clusterResourceLister) List(selector labels.Selector) (ret []*v1alpha1.ClusterResource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ClusterResource))
	})
	return ret, err
}

// ClusterResources returns an object that can list and get ClusterResources.
func (s *clusterResourceLister) ClusterResources(namespace string) ClusterResourceNamespaceLister {
	return clusterResourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ClusterResourceNamespaceLister helps list and get ClusterResources.
// All objects returned here must be treated as read-only.
type ClusterResourceNamespaceLister interface {
	// List lists all ClusterResources in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ClusterResource, err error)
	// Get retrieves the ClusterResource from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ClusterResource, error)
	ClusterResourceNamespaceListerExpansion
}

// clusterResourceNamespaceLister implements the ClusterResourceNamespaceLister
// interface.
type clusterResourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ClusterResources in the indexer for a given namespace.
func (s clusterResourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ClusterResource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ClusterResource))
	})
	return ret, err
}

// Get retrieves the ClusterResource from the indexer for a given namespace and name.
func (s clusterResourceNamespaceLister) Get(name string) (*v1alpha1.ClusterResource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("clusterresource"), name)
	}
	return obj.(*v1alpha1.ClusterResource), nil
}
