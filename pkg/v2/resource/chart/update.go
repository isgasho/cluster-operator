package chart

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/context/resourcecanceledcontext"
	"k8s.io/helm/pkg/helm"

	"github.com/giantswarm/cluster-operator/pkg/v2/guestcluster"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	guestHelmClient, err := r.getGuestHelmClient(ctx, obj)
	if guestcluster.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "did not get a Helm client for the guest cluster")

		// A not found error here means that the cluster-operator certificate for
		// the current guest cluster was not found. We can't continue without a Helm
		// client. We will retry during the next execution, when the certificate
		// might be available.
		resourcecanceledcontext.SetCanceled(ctx)
		r.logger.LogCtx(ctx, "level", "debug", "message", "canceling resource reconciliation for custom object")

		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	updateState, err := toResourceState(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if !reflect.DeepEqual(updateState, ResourceState{}) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "updating chart-operator chart")

		tarballPath, err := r.apprClient.PullChartTarball(updateState.ReleaseName, chartOperatorChannel)
		if err != nil {
			return microerror.Mask(err)
		}
		defer func() {
			err := r.fs.Remove(tarballPath)
			if err != nil {
				r.logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("deletion of %q failed", tarballPath), "stack", fmt.Sprintf("%#v", err))
			}
		}()

		// We need to pass the UpdateValueOverrides option to make the install process
		// use the default values and prevent errors on nested values.
		//
		//     {
		//      rpc error: code = Unknown desc = render error in "cnr-server-chart/templates/deployment.yaml":
		//      template: cnr-server-chart/templates/deployment.yaml:20:26:
		//      executing "cnr-server-chart/templates/deployment.yaml" at <.Values.image.reposi...>: can't evaluate field repository in type interface {}
		//     }
		//
		err = guestHelmClient.UpdateReleaseFromTarball(updateState.ReleaseName, tarballPath,
			helm.UpdateValueOverrides([]byte("{}")),
			helm.UpgradeWait(true))
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "updated chart-operator chart")
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "not updating chart-operator chart")
	}

	return nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	create, err := r.newCreateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	update, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := controller.NewPatch()
	patch.SetCreateChange(create)
	patch.SetUpdateChange(update)

	return patch, nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentResourceState, err := toResourceState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	desiredResourceState, err := toResourceState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "finding out if chart-operator has to be updated")

	updateState := &ResourceState{}
	if !reflect.DeepEqual(currentResourceState, ResourceState{}) && !reflect.DeepEqual(currentResourceState, desiredResourceState) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "chart-operator has to be updated")

		updateState = &desiredResourceState
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "chart-operator does not have to be updated")
	}

	return updateState, nil
}