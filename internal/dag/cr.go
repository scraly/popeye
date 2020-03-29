package dag

import (
	"context"
	"errors"

	"github.com/derailed/popeye/internal"
	"github.com/derailed/popeye/internal/client"
	"github.com/derailed/popeye/internal/dao"
	"github.com/derailed/popeye/pkg/config"
	"github.com/derailed/popeye/types"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ListClusterRoles list included ClusterRoles.
func ListClusterRoles(f types.Factory, cfg *config.Config) (map[string]*rbacv1.ClusterRole, error) {
	crs, err := listAllClusterRoles(f)
	if err != nil {
		return map[string]*rbacv1.ClusterRole{}, err
	}
	res := make(map[string]*rbacv1.ClusterRole, len(crs))
	for fqn, cr := range crs {
		res[fqn] = cr
	}

	return res, nil
}

// ListAllClusterRoles fetch all ClusterRoles on the cluster.
func listAllClusterRoles(f types.Factory) (map[string]*rbacv1.ClusterRole, error) {
	ll, err := fetchClusterRoles(f)
	if err != nil {
		return nil, err
	}

	crs := make(map[string]*rbacv1.ClusterRole, len(ll.Items))
	for i := range ll.Items {
		crs[metaFQN(ll.Items[i].ObjectMeta)] = &ll.Items[i]
	}

	return crs, nil
}

// FetchClusterRoles retrieves all ClusterRoles on the cluster.
func fetchClusterRoles(f types.Factory) (*rbacv1.ClusterRoleList, error) {
	var res dao.Resource
	res.Init(f, client.NewGVR("rbac.authorization.k8s.io/v1/clusterroles"))

	ctx := context.WithValue(context.Background(), internal.KeyFactory, f)
	oo, err := res.List(ctx, client.AllNamespaces)
	if err != nil {
		return nil, err
	}
	var ll rbacv1.ClusterRoleList
	for _, o := range oo {
		var cr rbacv1.ClusterRole
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &cr)
		if err != nil {
			return nil, errors.New("expecting clusterrole resource")
		}
		ll.Items = append(ll.Items, cr)
	}

	return &ll, nil
}
