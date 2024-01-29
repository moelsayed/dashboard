//go:build ee

/*
                  Kubermatic Enterprise Read-Only License
                         Version 1.0 ("KERO-1.0”)
                     Copyright © 2024 Kubermatic GmbH

   1.	You may only view, read and display for studying purposes the source
      code of the software licensed under this license, and, to the extent
      explicitly provided under this license, the binary code.
   2.	Any use of the software which exceeds the foregoing right, including,
      without limitation, its execution, compilation, copying, modification
      and distribution, is expressly prohibited.
   3.	THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
      EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
      MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
      IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
      CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
      TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
      SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

   END OF TERMS AND CONDITIONS
*/

package backuplocation

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	apiv2 "k8c.io/dashboard/v2/pkg/api/v2"
	"k8c.io/dashboard/v2/pkg/handler/middleware"
	"k8c.io/dashboard/v2/pkg/handler/v1/common"
	"k8c.io/dashboard/v2/pkg/provider"
)

func ListEndpoint(userInfoGetter provider.UserInfoGetter, projectProvider provider.ProjectProvider,
	privilegedProjectProvider provider.PrivilegedProjectProvider) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listCbslLReq)

		// check if user has access to the project
		_, err := common.GetProject(ctx, userInfoGetter, projectProvider, privilegedProjectProvider, req.ProjectID, nil)
		if err != nil {
			return nil, common.KubernetesErrorToHTTPError(err)
		}
		userInfo, err := userInfoGetter(ctx, req.ProjectID)
		if err != nil {
			return nil, common.KubernetesErrorToHTTPError(err)
		}
		cbslProvider := ctx.Value(middleware.CBSLProviderContextKey).(provider.ClusterBackupStorageLocationProvider)
		cbslList, err := cbslProvider.List(ctx, userInfo, provider.CSBSLListOptions{ProjectID: req.ProjectID})
		if err != nil {
			return nil, common.KubernetesErrorToHTTPError(err)
		}
		var respList []*apiv2.ClusterBackupStorageLocation
		for _, cbsl := range cbslList {
			respList = append(respList, &apiv2.ClusterBackupStorageLocation{
				Name:   cbsl.Name,
				Spec:   *cbsl.Spec.DeepCopy(),
				Status: *cbsl.Status.DeepCopy(),
			})
		}
		return respList, nil
	}
}

// listCbslLReq represents a request for listing ClusterBackupStorageLocations
// swagger:parameters listCbsl
type listCbslLReq struct {
	common.ProjectReq
}

func DecodeListCBSLReq(c context.Context, r *http.Request) (interface{}, error) {
	var req listCbslLReq

	pr, err := common.DecodeProjectRequest(c, r)
	if err != nil {
		return nil, err
	}

	req.ProjectReq = pr.(common.ProjectReq)
	return req, nil
}
