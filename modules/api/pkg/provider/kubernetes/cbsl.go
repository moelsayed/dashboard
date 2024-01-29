/*
Copyright 2024 The Kubermatic Kubernetes Platform contributors.

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
package kubernetes

import (
	"context"
	"fmt"

	"k8c.io/dashboard/v2/pkg/provider"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/log"
	"k8s.io/apimachinery/pkg/api/meta"
	kubenetutil "k8s.io/apimachinery/pkg/util/net"

	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type CBSLProvider struct {
	// createMasterImpersonatedClient is used as a ground for impersonation
	// whenever a connection to Seed API server is required
	seedImpersonatedClients map[string]ImpersonationClient
	clientsPrivileged       map[string]ctrlruntimeclient.Client
}

type PrivilegedCBSLProvider struct {
	// treat clientPrivileged as a privileged user and use wisely
	clientsPrivileged map[string]ctrlruntimeclient.Client
}

var _ provider.PrivilegedClusterBackupStorageLocationProvider = &PrivilegedCBSLProvider{}
var _ provider.ClusterBackupStorageLocationProvider = &CBSLProvider{}

func NewPrivilegedClusterBackupStorageLocationProvider(clients map[string]ctrlruntimeclient.Client) *PrivilegedCBSLProvider {
	return &PrivilegedCBSLProvider{
		clientsPrivileged: clients,
	}
}

func NewClusterBackupStorageLocationProvider(seedImpersonatedClients map[string]ImpersonationClient, clients map[string]ctrlruntimeclient.Client) *CBSLProvider {
	return &CBSLProvider{
		seedImpersonatedClients: seedImpersonatedClients,
		clientsPrivileged:       clients,
	}
}

func (p *PrivilegedCBSLProvider) ListUnsecured(ctx context.Context, options provider.CSBSLListOptions) ([]kubermaticv1.ClusterBackupStorageLocation, error) {
	respList := []kubermaticv1.ClusterBackupStorageLocation{}

	for _, clientPrivileged := range p.clientsPrivileged {

		cbslList := &kubermaticv1.ClusterBackupStorageLocationList{}
		labelSelector := ctrlruntimeclient.MatchingLabels{
			kubermaticv1.ProjectIDLabelKey: options.ProjectID,
		}
		if err := clientPrivileged.List(ctx, cbslList, labelSelector); err != nil {
			if kubenetutil.IsConnectionRefused(err) {
				continue
			}
			return nil, fmt.Errorf("failed to list ClusterBackupStorageLocations: %w", err)
		}
		respList = append(respList, cbslList.Items...)
	}

	return respList, nil
}

func (p *CBSLProvider) List(ctx context.Context, userInfo *provider.UserInfo, options provider.CSBSLListOptions) ([]kubermaticv1.ClusterBackupStorageLocation, error) {
	respList := []kubermaticv1.ClusterBackupStorageLocation{}

	for _, cseedImpersonationClient := range p.seedImpersonatedClients {
		impersonationClient, err := createImpersonationClientWrapperFromUserInfo(userInfo, cseedImpersonationClient)
		if err != nil {
			return nil, err
		}

		cbslList := &kubermaticv1.ClusterBackupStorageLocationList{}
		labelSelector := ctrlruntimeclient.MatchingLabels{
			kubermaticv1.ProjectIDLabelKey: options.ProjectID,
		}
		if err := impersonationClient.List(ctx, cbslList, labelSelector, ctrlruntimeclient.InNamespace("kubermatic")); err != nil {
			if kubenetutil.IsConnectionRefused(err) {
				continue
			}
			return nil, fmt.Errorf("failed to list ClusterBackupStorageLocations: %w", err)
		}
		respList = append(respList, cbslList.Items...)
	}

	return respList, nil
}

func CBSLProjectProviderFactory(mapper meta.RESTMapper, seedKubeconfigGetter provider.SeedKubeconfigGetter) provider.CBSLProviderGetter {
	return func(seeds map[string]*kubermaticv1.Seed) (provider.ClusterBackupStorageLocationProvider, error) {
		clientsPrivileged := make(map[string]ctrlruntimeclient.Client)
		seedImpersonationClients := make(map[string]ImpersonationClient)

		for seedName, seed := range seeds {
			if seed.Status.Phase == kubermaticv1.SeedInvalidPhase {
				log.Logger.Warnf("skipping seed %s as it is in an invalid phase", seedName)
				continue
			}

			cfg, err := seedKubeconfigGetter(seed)
			if err != nil {
				// if one or more Seeds are bad, continue with the request, log that a Seed is in error
				log.Logger.Warnw("error getting seed kubeconfig", "seed", seedName, "error", err)
				continue
			}
			seedImpersonationClients[seedName] = NewImpersonationClient(cfg, mapper).CreateImpersonatedClient
			clientPrivileged, err := ctrlruntimeclient.New(cfg, ctrlruntimeclient.Options{Mapper: mapper})
			if err != nil {
				return nil, err
			}
			clientsPrivileged[seedName] = clientPrivileged
		}

		return NewClusterBackupStorageLocationProvider(
			seedImpersonationClients,
			clientsPrivileged,
		), nil
	}
}
