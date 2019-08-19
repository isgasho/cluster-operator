package namespace

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller"
)

// ApplyUpdateChange is a no-op because the namespace in the tenant cluster is
// not updated.
func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	return nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	create, err := r.newCreateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := controller.NewPatch()
	patch.SetCreateChange(create)

	return patch, nil
}