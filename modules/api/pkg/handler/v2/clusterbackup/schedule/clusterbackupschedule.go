/*
Copyright 2023 The Kubermatic Kubernetes Platform contributors.

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

package clusterbackupschedule

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"k8c.io/dashboard/v2/pkg/provider"
)

func CreateEndpoint(userInfoGetter provider.UserInfoGetter, projectProvider provider.ProjectProvider,
	privilegedProjectProvider provider.PrivilegedProjectProvider) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return createEndpoint(ctx, request, userInfoGetter, projectProvider, privilegedProjectProvider)
	}
}

func DecodeCreateClusterBackupScheduleReq(c context.Context, r *http.Request) (interface{}, error) {
	return decodeCreateClusterBackupScheduleReq(c, r)
}

func ListEndpoint(userInfoGetter provider.UserInfoGetter, projectProvider provider.ProjectProvider,
	privilegedProjectProvider provider.PrivilegedProjectProvider) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return listEndpoint(ctx, request, userInfoGetter, projectProvider, privilegedProjectProvider)
	}
}

func DecodeListClusterBackupScheduleReq(c context.Context, r *http.Request) (interface{}, error) {
	return decodeListClusterBackupScheduleReq(c, r)
}

func GetEndpoint(userInfoGetter provider.UserInfoGetter, projectProvider provider.ProjectProvider,
	privilegedProjectProvider provider.PrivilegedProjectProvider) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return getEndpoint(ctx, request, userInfoGetter, projectProvider, privilegedProjectProvider)
	}
}

func DecodeGetRestoreBackupConfigReq(c context.Context, r *http.Request) (interface{}, error) {
	return decodeGetRestoreBackupConfigReq(c, r)

}

func DeleteEndpoint(userInfoGetter provider.UserInfoGetter, projectProvider provider.ProjectProvider,
	privilegedProjectProvider provider.PrivilegedProjectProvider) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return deleteEndpoint(ctx, request, userInfoGetter, projectProvider, privilegedProjectProvider)
	}
}

// check the name
func DecodeDeleteClusterBackupScheduleReq(c context.Context, r *http.Request) (interface{}, error) {
	return decodeDeleteClusterBackupScheduleReq(c, r)
}
