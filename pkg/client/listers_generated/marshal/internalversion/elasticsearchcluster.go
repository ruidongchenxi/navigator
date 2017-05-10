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

// This file was automatically generated by lister-gen

package internalversion

import (
	marshal "github.com/jetstack-experimental/navigator/pkg/apis/marshal"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ElasticsearchClusterLister helps list ElasticsearchClusters.
type ElasticsearchClusterLister interface {
	// List lists all ElasticsearchClusters in the indexer.
	List(selector labels.Selector) (ret []*marshal.ElasticsearchCluster, err error)
	// ElasticsearchClusters returns an object that can list and get ElasticsearchClusters.
	ElasticsearchClusters(namespace string) ElasticsearchClusterNamespaceLister
	ElasticsearchClusterListerExpansion
}

// elasticsearchClusterLister implements the ElasticsearchClusterLister interface.
type elasticsearchClusterLister struct {
	indexer cache.Indexer
}

// NewElasticsearchClusterLister returns a new ElasticsearchClusterLister.
func NewElasticsearchClusterLister(indexer cache.Indexer) ElasticsearchClusterLister {
	return &elasticsearchClusterLister{indexer: indexer}
}

// List lists all ElasticsearchClusters in the indexer.
func (s *elasticsearchClusterLister) List(selector labels.Selector) (ret []*marshal.ElasticsearchCluster, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*marshal.ElasticsearchCluster))
	})
	return ret, err
}

// ElasticsearchClusters returns an object that can list and get ElasticsearchClusters.
func (s *elasticsearchClusterLister) ElasticsearchClusters(namespace string) ElasticsearchClusterNamespaceLister {
	return elasticsearchClusterNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ElasticsearchClusterNamespaceLister helps list and get ElasticsearchClusters.
type ElasticsearchClusterNamespaceLister interface {
	// List lists all ElasticsearchClusters in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*marshal.ElasticsearchCluster, err error)
	// Get retrieves the ElasticsearchCluster from the indexer for a given namespace and name.
	Get(name string) (*marshal.ElasticsearchCluster, error)
	ElasticsearchClusterNamespaceListerExpansion
}

// elasticsearchClusterNamespaceLister implements the ElasticsearchClusterNamespaceLister
// interface.
type elasticsearchClusterNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ElasticsearchClusters in the indexer for a given namespace.
func (s elasticsearchClusterNamespaceLister) List(selector labels.Selector) (ret []*marshal.ElasticsearchCluster, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*marshal.ElasticsearchCluster))
	})
	return ret, err
}

// Get retrieves the ElasticsearchCluster from the indexer for a given namespace and name.
func (s elasticsearchClusterNamespaceLister) Get(name string) (*marshal.ElasticsearchCluster, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(marshal.Resource("elasticsearchcluster"), name)
	}
	return obj.(*marshal.ElasticsearchCluster), nil
}
