// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// VSphere v sphere
//
// swagger:model VSphere
type VSphere struct {

	// BasePath configures a vCenter folder path that KKP will create an individual cluster folder in.
	// If it's an absolute path, the RootPath configured in the datacenter will be ignored. If it is a relative path,
	// the BasePath part will be appended to the RootPath to construct the full path.
	BasePath string `json:"basePath,omitempty"`

	// If datacenter is set, this preset is only applicable to the
	// configured datacenter.
	Datacenter string `json:"datacenter,omitempty"`

	// Datastore to be used for storing virtual machines and as a default for dynamic volume provisioning, it is mutually exclusive with DatastoreCluster.
	Datastore string `json:"datastore,omitempty"`

	// DatastoreCluster to be used for storing virtual machines, it is mutually exclusive with Datastore.
	DatastoreCluster string `json:"datastoreCluster,omitempty"`

	// Only enabled presets will be available in the KKP dashboard.
	Enabled bool `json:"enabled,omitempty"`

	// List of vSphere networks.
	Networks []string `json:"networks"`

	// The vSphere user password.
	Password string `json:"password,omitempty"`

	// ResourcePool is used to manage resources such as cpu and memory for vSphere virtual machines. The resource pool should be defined on vSphere cluster level.
	ResourcePool string `json:"resourcePool,omitempty"`

	// The vSphere user name.
	Username string `json:"username,omitempty"`

	// Deprecated: Use networks instead.
	VMNetName string `json:"vmNetName,omitempty"`
}

// Validate validates this v sphere
func (m *VSphere) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this v sphere based on context it is used
func (m *VSphere) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *VSphere) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VSphere) UnmarshalBinary(b []byte) error {
	var res VSphere
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
