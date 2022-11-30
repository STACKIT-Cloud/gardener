// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shoot

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	v1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/apis/operations"
	operationsv1alpha1 "github.com/gardener/gardener/pkg/apis/operations/v1alpha1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/client/kubernetes/clientmap"
	"github.com/gardener/gardener/pkg/client/kubernetes/clientmap/keys"
	"github.com/gardener/gardener/pkg/controllerutils"
	gardenerextensions "github.com/gardener/gardener/pkg/extensions"
	"github.com/gardener/gardener/pkg/features"
	"github.com/gardener/gardener/pkg/gardenlet/apis/config"
	gardenletfeatures "github.com/gardener/gardener/pkg/gardenlet/features"
	"github.com/gardener/gardener/pkg/operation"
	botanistpkg "github.com/gardener/gardener/pkg/operation/botanist"
	"github.com/gardener/gardener/pkg/operation/garden"
	seedpkg "github.com/gardener/gardener/pkg/operation/seed"
	shootpkg "github.com/gardener/gardener/pkg/operation/shoot"
	"github.com/gardener/gardener/pkg/utils"
	utilerrors "github.com/gardener/gardener/pkg/utils/errors"
	"github.com/gardener/gardener/pkg/utils/flow"
	gutil "github.com/gardener/gardener/pkg/utils/gardener"
	"github.com/gardener/gardener/pkg/utils/imagevector"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/gardener/gardener/pkg/utils/kubernetes/health"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/component-base/version"
	"k8s.io/utils/clock"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const reconcilerName = "shoot"

func (c *Controller) shootAdd(ctx context.Context, obj interface{}, resetRateLimiting bool) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		c.log.Error(err, "Could not get key", "obj", obj)
		return
	}

	if resetRateLimiting {
		c.getShootQueue(ctx, obj).Forget(key)
	}
	c.getShootQueue(ctx, obj).Add(key)
}

func (c *Controller) shootUpdate(ctx context.Context, oldObj, newObj interface{}) {
	var (
		oldShoot = oldObj.(*gardencorev1beta1.Shoot)
		newShoot = newObj.(*gardencorev1beta1.Shoot)
		log      = c.log.WithValues("shoot", client.ObjectKeyFromObject(newShoot))
	)

	// If the generation did not change for an update event (i.e., no changes to the .spec section have
	// been made), we do not want to add the Shoot to the queue. The period reconciliation is handled
	// elsewhere by adding the Shoot to the queue to dedicated times.
	if newShoot.Generation == newShoot.Status.ObservedGeneration {
		log.V(1).Info("Do not need to do anything as the Update event occurred due to .status field changes")
		return
	}

	// If the shoot's deletion timestamp is set then we want to forget about the potentially established exponential
	// backoff and enqueue it faster.
	resetRateLimiting := oldShoot.DeletionTimestamp == nil && newShoot.DeletionTimestamp != nil

	c.shootAdd(ctx, newObj, resetRateLimiting)
}

func (c *Controller) shootDelete(ctx context.Context, obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		c.log.Error(err, "Could not get key", "obj", obj)
		return
	}

	c.getShootQueue(ctx, obj).Add(key)
}

// NewShootReconciler returns a reconciler that implements the main shoot reconciliation logic, i.e creation,
// hibernation, migration and deletion.
func NewShootReconciler(
	clientMap clientmap.ClientMap,
	recorder record.EventRecorder,
	imageVector imagevector.ImageVector,
	identity *gardencorev1beta1.Gardener,
	gardenClusterIdentity string,
	config *config.GardenletConfiguration,
) reconcile.Reconciler {
	return &shootReconciler{
		clientMap:                     clientMap,
		recorder:                      recorder,
		imageVector:                   imageVector,
		identity:                      identity,
		gardenClusterIdentity:         gardenClusterIdentity,
		config:                        config,
		shootReconciliationDueTracker: newReconciliationDueTracker(),
	}
}

type shootReconciler struct {
	clientMap clientmap.ClientMap
	recorder  record.EventRecorder

	imageVector                   imagevector.ImageVector
	identity                      *gardencorev1beta1.Gardener
	gardenClusterIdentity         string
	config                        *config.GardenletConfiguration
	shootReconciliationDueTracker *reconciliationDueTracker
}

func (r *shootReconciler) reconcileInMaintenanceOnly() bool {
	return pointer.BoolDeref(r.config.Controllers.Shoot.ReconcileInMaintenanceOnly, false)
}

func (r *shootReconciler) respectSyncPeriodOverwrite() bool {
	return pointer.BoolDeref(r.config.Controllers.Shoot.RespectSyncPeriodOverwrite, false)
}

func confineSpecUpdateRollout(maintenance *gardencorev1beta1.Maintenance) bool {
	return maintenance != nil && maintenance.ConfineSpecUpdateRollout != nil && *maintenance.ConfineSpecUpdateRollout
}

func (r *shootReconciler) syncClusterResourceToSeed(ctx context.Context, shoot *gardencorev1beta1.Shoot, project *gardencorev1beta1.Project, cloudProfile *gardencorev1beta1.CloudProfile, seed *gardencorev1beta1.Seed) error {
	seedName := getActiveSeedName(shoot)
	if seedName == nil {
		return nil
	}
	seedClient, err := r.clientMap.GetClient(ctx, keys.ForSeedWithName(*seedName))
	if err != nil {
		return fmt.Errorf("could not initialize a new Kubernetes client for the seed cluster: %+v", err)
	}

	clusterName := shootpkg.ComputeTechnicalID(project.Name, shoot)
	return gardenerextensions.SyncClusterResourceToSeed(ctx, seedClient.Client(), clusterName, shoot, cloudProfile, seed, project)
}

func (r *shootReconciler) checkSeedAndSyncClusterResource(ctx context.Context, shoot *gardencorev1beta1.Shoot, project *gardencorev1beta1.Project, cloudProfile *gardencorev1beta1.CloudProfile, seed *gardencorev1beta1.Seed) error {
	seedName := getActiveSeedName(shoot)
	if seedName == nil || seed == nil {
		return nil
	}

	// Don't wait for the Seed to be ready if it is already marked for deletion. In this case
	// it will never get ready because the bootstrap loop is never executed again.
	// Don't block the Shoot deletion flow in this case to allow proper cleanup.
	if seed.DeletionTimestamp == nil {
		if err := health.CheckSeed(seed, r.identity); err != nil {
			return fmt.Errorf("seed is not yet ready: %w", err)
		}
	}

	if err := r.syncClusterResourceToSeed(ctx, shoot, project, cloudProfile, seed); err != nil {
		return fmt.Errorf("could not sync cluster resource to seed: %w", err)
	}

	return nil
}

// deleteClusterResourceFromSeed deletes the `Cluster` extension resource for the shoot in the seed cluster.
func (r *shootReconciler) deleteClusterResourceFromSeed(ctx context.Context, shoot *gardencorev1beta1.Shoot, projectName string) error {
	seedName := getActiveSeedName(shoot)
	if seedName == nil {
		return nil
	}

	seedClient, err := r.clientMap.GetClient(ctx, keys.ForSeedWithName(*seedName))
	if err != nil {
		return fmt.Errorf("could not initialize a new Kubernetes client for the seed cluster: %+v", err)
	}

	cluster := &extensionsv1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: shootpkg.ComputeTechnicalID(projectName, shoot),
		},
	}

	return client.IgnoreNotFound(seedClient.Client().Delete(ctx, cluster))
}

func getActiveSeedName(shoot *gardencorev1beta1.Shoot) *string {
	if shoot.Status.SeedName != nil {
		return shoot.Status.SeedName
	}
	return shoot.Spec.SeedName
}

func (r *shootReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := logf.FromContext(ctx)

	gardenClient, err := r.clientMap.GetClient(ctx, keys.ForGarden())
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get garden client: %w", err)
	}

	shoot := &gardencorev1beta1.Shoot{}
	if err := gardenClient.Client().Get(ctx, request.NamespacedName, shoot); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(1).Info("Object is gone, stop reconciling")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("error retrieving object from store: %w", err)
	}

	// fetch related objects required for shoot operation
	project, _, err := gutil.ProjectAndNamespaceFromReader(ctx, gardenClient.Client(), shoot.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	if project == nil {
		return reconcile.Result{}, fmt.Errorf("cannot find Project for namespace '%s'", shoot.Namespace)
	}

	cloudProfile := &gardencorev1beta1.CloudProfile{}
	if err := gardenClient.Client().Get(ctx, kutil.Key(shoot.Spec.CloudProfileName), cloudProfile); err != nil {
		return reconcile.Result{}, err
	}

	key := shootKey(shoot)
	if shoot.DeletionTimestamp != nil {
		log = log.WithValues("operation", "delete")
		r.shootReconciliationDueTracker.off(key)

		seedName := getActiveSeedName(shoot)
		seed := &gardencorev1beta1.Seed{}
		if err := gardenClient.Client().Get(ctx, client.ObjectKey{Name: *seedName}, seed); err != nil {
			return reconcile.Result{}, err
		}
		return r.deleteShoot(ctx, log, gardenClient, shoot, project, cloudProfile, seed)
	}

	seed := &gardencorev1beta1.Seed{}
	if err := gardenClient.Client().Get(ctx, client.ObjectKey{Name: *shoot.Spec.SeedName}, seed); err != nil {
		return reconcile.Result{}, err
	}

	if shouldPrepareShootForMigration(shoot) {
		log = log.WithValues("operation", "migrate")
		r.shootReconciliationDueTracker.off(key)

		if err := r.isSeedReadyForMigration(seed); err != nil {
			return reconcile.Result{}, fmt.Errorf("target Seed is not available to host the Control Plane of Shoot %s: %w", shoot.GetName(), err)
		}

		hasBastions, err := r.shootHasBastions(ctx, shoot, gardenClient)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to check for related Bastions: %w", err)
		}
		if hasBastions {
			hasBastionErr := errors.New("shoot has still Bastions")
			updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, hasBastionErr.Error(), gardencorev1beta1.LastOperationTypeMigrate, shoot.Status.LastErrors...)
			return reconcile.Result{}, utilerrors.WithSuppressed(hasBastionErr, updateErr)
		}

		sourceSeed := &gardencorev1beta1.Seed{}
		if err := gardenClient.Client().Get(ctx, client.ObjectKey{Name: *shoot.Status.SeedName}, sourceSeed); err != nil {
			return reconcile.Result{}, err
		}

		return r.prepareShootForMigration(ctx, log, gardenClient, shoot, project, cloudProfile, sourceSeed)
	}

	// if shoot is no longer managed by this gardenlet (e.g., due to migration to another seed) then don't requeue
	if !IsShootManagedByThisGardenlet(shoot, r.config) {
		log.V(1).Info("Skipping because Shoot is not managed by this gardenlet", "seedName", *shoot.Spec.SeedName)
		return reconcile.Result{}, nil
	}

	log = log.WithValues("operation", "reconcile")
	return r.reconcileShoot(ctx, log, gardenClient, shoot, project, cloudProfile, seed)
}

func shouldPrepareShootForMigration(shoot *gardencorev1beta1.Shoot) bool {
	return shoot.Status.SeedName != nil && shoot.Spec.SeedName != nil && *shoot.Spec.SeedName != *shoot.Status.SeedName
}

const taskID = "initializeOperation"

func (r *shootReconciler) initializeOperation(
	ctx context.Context,
	log logr.Logger,
	gardenClient kubernetes.Interface,
	shoot *gardencorev1beta1.Shoot,
	project *gardencorev1beta1.Project,
	cloudProfile *gardencorev1beta1.CloudProfile,
	seed *gardencorev1beta1.Seed,
) (
	*operation.Operation,
	error,
) {
	gardenSecrets, err := garden.ReadGardenSecrets(ctx, log, gardenClient.Client(), gutil.ComputeGardenNamespace(seed.Name), true)
	if err != nil {
		return nil, err
	}

	gardenObj, err := garden.
		NewBuilder().
		WithProject(project).
		WithInternalDomainFromSecrets(gardenSecrets).
		WithDefaultDomainsFromSecrets(gardenSecrets).
		Build(ctx)
	if err != nil {
		return nil, err
	}

	seedObj, err := seedpkg.
		NewBuilder().
		WithSeedObject(seed).
		Build(ctx)
	if err != nil {
		return nil, err
	}

	shootObj, err := shootpkg.
		NewBuilder().
		WithShootObject(shoot).
		WithCloudProfileObject(cloudProfile).
		WithShootSecretFrom(gardenClient.Client()).
		WithProjectName(project.Name).
		WithExposureClassFrom(gardenClient.Client()).
		WithDisableDNS(!seed.Spec.Settings.ShootDNS.Enabled).
		WithInternalDomain(gardenObj.InternalDomain).
		WithDefaultDomains(gardenObj.DefaultDomains).
		Build(ctx, gardenClient.Client())
	if err != nil {
		return nil, err
	}

	return operation.
		NewBuilder().
		WithLogger(log).
		WithConfig(r.config).
		WithGardenerInfo(r.identity).
		WithGardenClusterIdentity(r.gardenClusterIdentity).
		WithExposureClassHandlerFromConfig(r.config).
		WithSecrets(gardenSecrets).
		WithImageVector(r.imageVector).
		WithGarden(gardenObj).
		WithSeed(seedObj).
		WithShoot(shootObj).
		Build(ctx, r.clientMap)
}

func (r *shootReconciler) deleteShoot(
	ctx context.Context,
	log logr.Logger,
	gardenClient kubernetes.Interface,
	shoot *gardencorev1beta1.Shoot,
	project *gardencorev1beta1.Project,
	cloudProfile *gardencorev1beta1.CloudProfile,
	seed *gardencorev1beta1.Seed,
) (
	reconcile.Result,
	error,
) {
	var (
		err error

		operationType              = v1beta1helper.ComputeOperationType(shoot.ObjectMeta, shoot.Status.LastOperation)
		respectSyncPeriodOverwrite = r.respectSyncPeriodOverwrite()
		failed                     = gutil.IsShootFailed(shoot)
		ignored                    = gutil.ShouldIgnoreShoot(respectSyncPeriodOverwrite, shoot)
		failedOrIgnored            = failed || ignored
	)

	if !controllerutil.ContainsFinalizer(shoot, gardencorev1beta1.GardenerName) {
		return reconcile.Result{}, nil
	}

	// If the .status.uid field is empty, then we assume that there has never been any operation running for this Shoot
	// cluster. This implies that there can not be any resource which we have to delete.
	// We accept the deletion.
	if len(shoot.Status.UID) == 0 {
		log.Info("The `.status.uid` is empty, assuming Shoot cluster did never exist, deletion accepted")
		return r.finalizeShootDeletion(ctx, log, gardenClient, shoot, project.Name)
	}

	hasBastions, err := r.shootHasBastions(ctx, shoot, gardenClient)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to check for related Bastions: %w", err)
	}
	if hasBastions {
		hasBastionErr := errors.New("shoot has still Bastions")
		updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, hasBastionErr.Error(), operationType, shoot.Status.LastErrors...)
		return reconcile.Result{}, utilerrors.WithSuppressed(hasBastionErr, updateErr)
	}

	// If the .status.lastOperation already indicates that the deletion is successful then we finalize it immediately.
	if shoot.Status.LastOperation != nil && shoot.Status.LastOperation.Type == gardencorev1beta1.LastOperationTypeDelete && shoot.Status.LastOperation.State == gardencorev1beta1.LastOperationStateSucceeded {
		log.Info("The `.status.lastOperation` indicates a successful deletion, deletion accepted")
		return r.finalizeShootDeletion(ctx, log, gardenClient, shoot, project.Name)
	}

	// If shoot is failed or ignored then sync the Cluster resource so that extension controllers running in the seed
	// get to know of the shoot's status.
	if failedOrIgnored {
		if syncErr := r.syncClusterResourceToSeed(ctx, shoot, project, cloudProfile, seed); syncErr != nil {
			log.Error(syncErr, "Not allowed to update Shoot with error, trying to sync Cluster resource again")
			updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, syncErr.Error(), operationType, shoot.Status.LastErrors...)
			return reconcile.Result{}, utilerrors.WithSuppressed(syncErr, updateErr)
		}

		log.Info("Shoot is failed or ignored, do not start deletion flow")
		return reconcile.Result{}, nil
	}

	// Trigger regular shoot deletion flow.
	r.recorder.Event(shoot, corev1.EventTypeNormal, gardencorev1beta1.EventDeleting, "Deleting Shoot cluster")
	shootNamespace := shootpkg.ComputeTechnicalID(project.Name, shoot)
	if err = r.updateShootStatusOperationStart(ctx, gardenClient.Client(), shoot, shootNamespace, operationType); err != nil {
		return reconcile.Result{}, err
	}

	o, operationErr := r.initializeOperation(ctx, log, gardenClient, shoot, project, cloudProfile, seed)
	if operationErr != nil {
		updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, fmt.Sprintf("Could not initialize a new operation for Shoot cluster deletion: %s", operationErr.Error()), operationType, lastErrorsOperationInitializationFailure(shoot.Status.LastErrors, operationErr)...)
		return reconcile.Result{}, utilerrors.WithSuppressed(operationErr, updateErr)
	}

	// At this point the deletion is allowed, hence, check if the seed is up-to-date, then sync the Cluster resource
	// initialize a new operation and, eventually, start the deletion flow.
	if err := r.checkSeedAndSyncClusterResource(ctx, shoot, project, cloudProfile, seed); err != nil {
		return patchShootStatusAndRequeueOnSyncError(ctx, log, gardenClient.Client(), shoot, err)
	}

	if flowErr := r.runDeleteShootFlow(ctx, o); flowErr != nil {
		r.recorder.Event(shoot, corev1.EventTypeWarning, gardencorev1beta1.EventDeleteError, flowErr.Description)
		updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, flowErr.Description, operationType, flowErr.LastErrors...)
		return reconcile.Result{}, utilerrors.WithSuppressed(errors.New(flowErr.Description), updateErr)
	}

	r.recorder.Event(shoot, corev1.EventTypeNormal, gardencorev1beta1.EventDeleted, "Deleted Shoot cluster")
	return r.finalizeShootDeletion(ctx, log, gardenClient, shoot, project.Name)
}

func (r *shootReconciler) finalizeShootDeletion(ctx context.Context, log logr.Logger, gardenClient kubernetes.Interface, shoot *gardencorev1beta1.Shoot, projectName string) (reconcile.Result, error) {
	if cleanErr := r.deleteClusterResourceFromSeed(ctx, shoot, projectName); cleanErr != nil {
		lastErr := v1beta1helper.LastError(fmt.Sprintf("Could not delete Cluster resource in seed: %s", cleanErr))
		updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, lastErr.Description, gardencorev1beta1.LastOperationTypeDelete, *lastErr)
		r.recorder.Event(shoot, corev1.EventTypeWarning, gardencorev1beta1.EventDeleteError, lastErr.Description)
		return reconcile.Result{}, utilerrors.WithSuppressed(errors.New(lastErr.Description), updateErr)
	}

	return reconcile.Result{}, r.removeFinalizerFrom(ctx, log, gardenClient, shoot)
}

func (r *shootReconciler) isSeedReadyForMigration(seed *gardencorev1beta1.Seed) error {
	if seed.DeletionTimestamp != nil {
		return fmt.Errorf("seed is set for deletion")
	}

	return health.CheckSeedForMigration(seed, r.identity)
}

func (r *shootReconciler) shootHasBastions(ctx context.Context, shoot *gardencorev1beta1.Shoot, gardenClient kubernetes.Interface) (bool, error) {
	// list all bastions that reference this shoot
	bastionList := operationsv1alpha1.BastionList{}

	if err := gardenClient.Client().List(ctx, &bastionList, client.InNamespace(shoot.Namespace), client.MatchingFields{operations.BastionShootName: shoot.Name}); err != nil {
		return false, fmt.Errorf("failed to list related Bastions: %w", err)
	}

	return len(bastionList.Items) > 0, nil
}

func (r *shootReconciler) reconcileShoot(
	ctx context.Context,
	log logr.Logger,
	gardenClient kubernetes.Interface,
	shoot *gardencorev1beta1.Shoot,
	project *gardencorev1beta1.Project,
	cloudProfile *gardencorev1beta1.CloudProfile,
	seed *gardencorev1beta1.Seed,
) (
	reconcile.Result,
	error,
) {
	var (
		key                                        = shootKey(shoot)
		operationType                              = computeOperationType(shoot)
		isRestoring                                = operationType == gardencorev1beta1.LastOperationTypeRestore
		respectSyncPeriodOverwrite                 = r.respectSyncPeriodOverwrite()
		failed                                     = gutil.IsShootFailed(shoot)
		ignored                                    = gutil.ShouldIgnoreShoot(respectSyncPeriodOverwrite, shoot)
		failedOrIgnored                            = failed || ignored
		reconcileInMaintenanceOnly                 = r.reconcileInMaintenanceOnly()
		isUpToDate                                 = gutil.IsObservedAtLatestGenerationAndSucceeded(shoot)
		isNowInEffectiveShootMaintenanceTimeWindow = gutil.IsNowInEffectiveShootMaintenanceTimeWindow(shoot)
		alreadyReconciledDuringThisTimeWindow      = gutil.LastReconciliationDuringThisTimeWindow(shoot)
		regularReconciliationIsDue                 = isUpToDate && isNowInEffectiveShootMaintenanceTimeWindow && !alreadyReconciledDuringThisTimeWindow
		reconcileAllowed                           = !failedOrIgnored && ((!reconcileInMaintenanceOnly && !confineSpecUpdateRollout(shoot.Spec.Maintenance)) || !isUpToDate || (isNowInEffectiveShootMaintenanceTimeWindow && !alreadyReconciledDuringThisTimeWindow))
	)

	if !controllerutil.ContainsFinalizer(shoot, gardencorev1beta1.GardenerName) {
		log.Info("Adding finalizer")
		if err := controllerutils.AddFinalizers(ctx, gardenClient.Client(), shoot, gardencorev1beta1.GardenerName); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return reconcile.Result{}, nil
	}

	log.WithValues(
		"operationType", operationType,
		"respectSyncPeriodOverwrite", respectSyncPeriodOverwrite,
		"failed", failed,
		"ignored", ignored,
		"failedOrIgnored", failedOrIgnored,
		"reconcileInMaintenanceOnly", reconcileInMaintenanceOnly,
		"isUpToDate", isUpToDate,
		"isNowInEffectiveShootMaintenanceTimeWindow", isNowInEffectiveShootMaintenanceTimeWindow,
		"alreadyReconciledDuringThisTimeWindow", alreadyReconciledDuringThisTimeWindow,
		"regularReconciliationIsDue", regularReconciliationIsDue,
		"reconcileAllowed", reconcileAllowed,
	).Info("Checking if Shoot can be reconciled")

	// If shoot is failed or ignored then sync the Cluster resource so that extension controllers running in the seed
	// get to know of the shoot's status.
	if failedOrIgnored {
		if syncErr := r.syncClusterResourceToSeed(ctx, shoot, project, cloudProfile, seed); syncErr != nil {
			log.Error(syncErr, "Not allowed to update Shoot with error, trying to sync Cluster resource again")
			updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, syncErr.Error(), operationType, shoot.Status.LastErrors...)
			return reconcile.Result{}, utilerrors.WithSuppressed(syncErr, updateErr)
		}

		log.Info("Shoot is failed or ignored, do not start reconciliation flow")
		return reconcile.Result{}, nil
	}

	// If reconciliation is not allowed then compute the duration until the next sync and requeue.
	if !reconcileAllowed {
		return r.scheduleNextSync(log, shoot, regularReconciliationIsDue, "Reconciliation is not allowed"), nil
	}

	if regularReconciliationIsDue && !r.shootReconciliationDueTracker.tracked(key) {
		r.shootReconciliationDueTracker.on(key)
		return r.scheduleNextSync(log, shoot, regularReconciliationIsDue, "Reconciliation allowed and due but not yet tracked (jittering next sync)"), nil
	}

	r.recorder.Event(shoot, corev1.EventTypeNormal, gardencorev1beta1.EventReconciling, fmt.Sprintf("%s Shoot cluster", utils.IifString(isRestoring, "Restoring", "Reconciling")))
	shootNamespace := shootpkg.ComputeTechnicalID(project.Name, shoot)
	if err := r.updateShootStatusOperationStart(ctx, gardenClient.Client(), shoot, shootNamespace, operationType); err != nil {
		return reconcile.Result{}, err
	}

	o, operationErr := r.initializeOperation(ctx, log, gardenClient, shoot, project, cloudProfile, seed)
	if operationErr != nil {
		description := fmt.Sprintf("Could not initialize a new operation for Shoot cluster %s: %s", utils.IifString(isRestoring, "restoration", "reconciliation"), operationErr.Error())
		updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, description, operationType, lastErrorsOperationInitializationFailure(shoot.Status.LastErrors, operationErr)...)
		return reconcile.Result{}, utilerrors.WithSuppressed(operationErr, updateErr)
	}

	// write UID to status when operation was created successfully once
	if len(shoot.Status.UID) == 0 {
		patch := client.MergeFrom(shoot.DeepCopy())
		shoot.Status.UID = shoot.UID
		err := gardenClient.Client().Status().Patch(ctx, shoot, patch)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// At this point the reconciliation is allowed, hence, check if the seed is up-to-date, then sync the Cluster resource
	// initialize a new operation and, eventually, start the reconciliation flow.
	if err := r.checkSeedAndSyncClusterResource(ctx, shoot, project, cloudProfile, seed); err != nil {
		return patchShootStatusAndRequeueOnSyncError(ctx, log, gardenClient.Client(), shoot, err)
	}

	r.shootReconciliationDueTracker.off(key)

	if flowErr := r.runReconcileShootFlow(ctx, o, operationType); flowErr != nil {
		r.recorder.Event(shoot, corev1.EventTypeWarning, gardencorev1beta1.EventReconcileError, flowErr.Description)
		updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, flowErr.Description, operationType, flowErr.LastErrors...)
		return reconcile.Result{}, utilerrors.WithSuppressed(errors.New(flowErr.Description), updateErr)
	}

	r.recorder.Event(shoot, corev1.EventTypeNormal, gardencorev1beta1.EventReconciled, fmt.Sprintf("%s Shoot cluster", utils.IifString(isRestoring, "Restored", "Reconciled")))
	if err := r.patchShootStatusOperationSuccess(ctx, gardenClient.Client(), shoot, o.Shoot.SeedNamespace, &seed.Name, operationType); err != nil {
		return reconcile.Result{}, err
	}

	if syncErr := r.syncClusterResourceToSeed(ctx, shoot, project, cloudProfile, seed); syncErr != nil {
		log.Error(syncErr, "Cluster resource sync to seed failed")
		updateErr := r.patchShootStatusOperationError(ctx, gardenClient.Client(), shoot, syncErr.Error(), operationType, shoot.Status.LastErrors...)
		return reconcile.Result{}, utilerrors.WithSuppressed(syncErr, updateErr)
	}

	r.shootReconciliationDueTracker.on(key)
	return r.scheduleNextSync(log, shoot, false, fmt.Sprintf("%s finished successfully", utils.IifString(isRestoring, "Restoration", "Reconciliation"))), nil
}

func (r *shootReconciler) scheduleNextSync(
	log logr.Logger,
	shoot *gardencorev1beta1.Shoot,
	regularReconciliationIsDue bool,
	reason string,
) reconcile.Result {
	durationUntilNextSync := r.durationUntilNextShootSync(shoot, regularReconciliationIsDue)
	nextReconciliation := time.Now().UTC().Add(durationUntilNextSync)

	log.Info("Scheduled next reconciliation for Shoot", "reason", reason, "duration", durationUntilNextSync, "nextReconciliation", nextReconciliation)
	return reconcile.Result{RequeueAfter: durationUntilNextSync}
}

func (r *shootReconciler) durationUntilNextShootSync(shoot *gardencorev1beta1.Shoot, regularReconciliationIsDue bool) time.Duration {
	syncPeriod := gutil.SyncPeriodOfShoot(r.respectSyncPeriodOverwrite(), r.config.Controllers.Shoot.SyncPeriod.Duration, shoot)
	if !r.reconcileInMaintenanceOnly() && !confineSpecUpdateRollout(shoot.Spec.Maintenance) {
		return syncPeriod
	}

	now := time.Now()
	window := gutil.EffectiveShootMaintenanceTimeWindow(shoot)

	return window.RandomDurationUntilNext(now, regularReconciliationIsDue)
}

func patchShootStatusAndRequeueOnSyncError(
	ctx context.Context,
	log logr.Logger,
	c client.StatusClient,
	shoot *gardencorev1beta1.Shoot,
	err error,
) (
	reconcile.Result,
	error,
) {
	log.Error(err, "Shoot cannot be synced with seed")

	patch := client.MergeFrom(shoot.DeepCopy())
	shoot.Status.LastOperation.Description = fmt.Sprintf("Shoot cannot be synced with Seed: %v", err)
	shoot.Status.LastOperation.LastUpdateTime = metav1.NewTime(time.Now().UTC())
	if patchErr := c.Status().Patch(ctx, shoot, patch); patchErr != nil {
		return reconcile.Result{}, utilerrors.WithSuppressed(err, patchErr)
	}

	return reconcile.Result{
		RequeueAfter: 15 * time.Second, // prevent ddos-ing the seed
	}, nil
}

func (r *shootReconciler) newProgressReporter(reporterFn flow.ProgressReporterFn) flow.ProgressReporter {
	if r.config.Controllers.Shoot != nil && r.config.Controllers.Shoot.ProgressReportPeriod != nil {
		return flow.NewDelayingProgressReporter(clock.RealClock{}, reporterFn, r.config.Controllers.Shoot.ProgressReportPeriod.Duration)
	}
	return flow.NewImmediateProgressReporter(reporterFn)
}

func (r *shootReconciler) updateShootStatusOperationStart(ctx context.Context, gardenClient client.Client, shoot *gardencorev1beta1.Shoot, shootNamespace string, operationType gardencorev1beta1.LastOperationType) error {
	var (
		now                   = metav1.NewTime(time.Now().UTC())
		operationTypeSwitched bool
		description           string
	)

	switch operationType {
	case gardencorev1beta1.LastOperationTypeCreate, gardencorev1beta1.LastOperationTypeReconcile:
		description = "Reconciliation of Shoot cluster initialized."
		operationTypeSwitched = false

	case gardencorev1beta1.LastOperationTypeRestore:
		description = "Restoration of Shoot cluster initialized."
		operationTypeSwitched = false

	case gardencorev1beta1.LastOperationTypeMigrate:
		description = "Preparation of Shoot cluster for migration initialized."
		operationTypeSwitched = false

	case gardencorev1beta1.LastOperationTypeDelete:
		description = "Deletion of Shoot cluster in progress."
		operationTypeSwitched = shoot.Status.LastOperation != nil && shoot.Status.LastOperation.Type != gardencorev1beta1.LastOperationTypeDelete
	}

	if shoot.Status.RetryCycleStartTime == nil ||
		shoot.Generation != shoot.Status.ObservedGeneration ||
		shoot.Status.Gardener.Version != version.Get().GitVersion ||
		operationTypeSwitched {

		shoot.Status.RetryCycleStartTime = &now
	}

	if len(shoot.Status.TechnicalID) == 0 {
		shoot.Status.TechnicalID = shootNamespace
	}

	if !equality.Semantic.DeepEqual(shoot.Status.SeedName, shoot.Spec.SeedName) &&
		operationType != gardencorev1beta1.LastOperationTypeMigrate &&
		operationType != gardencorev1beta1.LastOperationTypeDelete {
		shoot.Status.SeedName = shoot.Spec.SeedName
	}

	shoot.Status.LastErrors = v1beta1helper.DeleteLastErrorByTaskID(shoot.Status.LastErrors, taskID)
	shoot.Status.Gardener = *r.identity
	shoot.Status.ObservedGeneration = shoot.Generation
	shoot.Status.LastOperation = &gardencorev1beta1.LastOperation{
		Type:           operationType,
		State:          gardencorev1beta1.LastOperationStateProcessing,
		Progress:       0,
		Description:    description,
		LastUpdateTime: now,
	}

	var mustRemoveOperationAnnotation bool

	switch shoot.Annotations[v1beta1constants.GardenerOperation] {
	case v1beta1constants.ShootOperationRotateCredentialsStart:
		mustRemoveOperationAnnotation = true
		if gardenletfeatures.FeatureGate.Enabled(features.ShootCARotation) {
			startRotationCA(shoot, &now)
		}
		if gardenletfeatures.FeatureGate.Enabled(features.ShootSARotation) {
			startRotationServiceAccountKey(shoot, &now)
		}
		startRotationKubeconfig(shoot, &now)
		startRotationSSHKeypair(shoot, &now)
		startRotationObservability(shoot, &now)
		startRotationETCDEncryptionKey(shoot, &now)
	case v1beta1constants.ShootOperationRotateCredentialsComplete:
		mustRemoveOperationAnnotation = true
		if gardenletfeatures.FeatureGate.Enabled(features.ShootCARotation) {
			completeRotationCA(shoot)
		}
		if gardenletfeatures.FeatureGate.Enabled(features.ShootSARotation) {
			completeRotationServiceAccountKey(shoot)
		}
		completeRotationETCDEncryptionKey(shoot)

	case v1beta1constants.ShootOperationRotateCAStart:
		mustRemoveOperationAnnotation = true
		startRotationCA(shoot, &now)
	case v1beta1constants.ShootOperationRotateCAComplete:
		mustRemoveOperationAnnotation = true
		completeRotationCA(shoot)

	case v1beta1constants.ShootOperationRotateKubeconfigCredentials:
		mustRemoveOperationAnnotation = true
		startRotationKubeconfig(shoot, &now)

	case v1beta1constants.ShootOperationRotateSSHKeypair:
		mustRemoveOperationAnnotation = true
		startRotationSSHKeypair(shoot, &now)

	case v1beta1constants.ShootOperationRotateObservabilityCredentials:
		mustRemoveOperationAnnotation = true
		startRotationObservability(shoot, &now)

	case v1beta1constants.ShootOperationRotateServiceAccountKeyStart:
		mustRemoveOperationAnnotation = true
		startRotationServiceAccountKey(shoot, &now)
	case v1beta1constants.ShootOperationRotateServiceAccountKeyComplete:
		mustRemoveOperationAnnotation = true
		completeRotationServiceAccountKey(shoot)

	case v1beta1constants.ShootOperationRotateETCDEncryptionKeyStart:
		mustRemoveOperationAnnotation = true
		startRotationETCDEncryptionKey(shoot, &now)
	case v1beta1constants.ShootOperationRotateETCDEncryptionKeyComplete:
		mustRemoveOperationAnnotation = true
		completeRotationETCDEncryptionKey(shoot)
	}

	if err := gardenClient.Status().Update(ctx, shoot); err != nil {
		return err
	}

	if mustRemoveOperationAnnotation {
		patch := client.MergeFrom(shoot.DeepCopy())
		delete(shoot.Annotations, v1beta1constants.GardenerOperation)
		return gardenClient.Patch(ctx, shoot, patch)
	}

	return nil
}

func (r *shootReconciler) patchShootStatusOperationSuccess(
	ctx context.Context,
	gardenClient client.Client,
	shoot *gardencorev1beta1.Shoot,
	shootSeedNamespace string,
	seedName *string,
	operationType gardencorev1beta1.LastOperationType,
) error {
	var (
		now                        = metav1.NewTime(time.Now().UTC())
		description                string
		setConditionsToProgressing bool
	)

	switch operationType {
	case gardencorev1beta1.LastOperationTypeCreate, gardencorev1beta1.LastOperationTypeReconcile:
		description = "Shoot cluster has been successfully reconciled."
		setConditionsToProgressing = true

	case gardencorev1beta1.LastOperationTypeMigrate:
		description = "Shoot cluster has been successfully prepared for migration."
		setConditionsToProgressing = false

	case gardencorev1beta1.LastOperationTypeRestore:
		description = "Shoot cluster has been successfully restored."
		setConditionsToProgressing = true

	case gardencorev1beta1.LastOperationTypeDelete:
		description = "Shoot cluster has been successfully deleted."
		setConditionsToProgressing = false
	}

	patch := client.StrategicMergeFrom(shoot.DeepCopy())

	if len(shootSeedNamespace) > 0 && seedName != nil {
		isHibernated, err := r.isHibernationActive(ctx, shootSeedNamespace, seedName)
		if err != nil {
			return fmt.Errorf("error updating Shoot (%s/%s) after successful reconciliation when checking for active hibernation: %w", shoot.Namespace, shoot.Name, err)
		}
		shoot.Status.IsHibernated = isHibernated
	}

	if setConditionsToProgressing {
		for i, cond := range shoot.Status.Conditions {
			switch cond.Type {
			case gardencorev1beta1.ShootAPIServerAvailable,
				gardencorev1beta1.ShootControlPlaneHealthy,
				gardencorev1beta1.ShootEveryNodeReady,
				gardencorev1beta1.ShootSystemComponentsHealthy:
				if cond.Status != gardencorev1beta1.ConditionFalse {
					shoot.Status.Conditions[i].Status = gardencorev1beta1.ConditionProgressing
					shoot.Status.Conditions[i].LastUpdateTime = metav1.Now()
				}
			}
		}
		for i, constr := range shoot.Status.Constraints {
			switch constr.Type {
			case gardencorev1beta1.ShootHibernationPossible,
				gardencorev1beta1.ShootMaintenancePreconditionsSatisfied:
				if constr.Status != gardencorev1beta1.ConditionFalse {
					shoot.Status.Constraints[i].Status = gardencorev1beta1.ConditionProgressing
					shoot.Status.Conditions[i].LastUpdateTime = metav1.Now()
				}
			}
		}
	}

	if operationType != gardencorev1beta1.LastOperationTypeDelete {
		shoot.Status.SeedName = shoot.Spec.SeedName
	}

	shoot.Status.RetryCycleStartTime = nil
	shoot.Status.LastErrors = nil
	shoot.Status.LastOperation = &gardencorev1beta1.LastOperation{
		Type:           operationType,
		State:          gardencorev1beta1.LastOperationStateSucceeded,
		Progress:       100,
		Description:    description,
		LastUpdateTime: now,
	}

	switch v1beta1helper.GetShootCARotationPhase(shoot.Status.Credentials) {
	case gardencorev1beta1.RotationPreparing:
		v1beta1helper.MutateShootCARotation(shoot, func(rotation *gardencorev1beta1.ShootCARotation) {
			rotation.Phase = gardencorev1beta1.RotationPrepared
		})

	case gardencorev1beta1.RotationCompleting:
		v1beta1helper.MutateShootCARotation(shoot, func(rotation *gardencorev1beta1.ShootCARotation) {
			rotation.Phase = gardencorev1beta1.RotationCompleted
			rotation.LastCompletionTime = &now
		})
	}

	switch v1beta1helper.GetShootServiceAccountKeyRotationPhase(shoot.Status.Credentials) {
	case gardencorev1beta1.RotationPreparing:
		v1beta1helper.MutateShootServiceAccountKeyRotation(shoot, func(rotation *gardencorev1beta1.ShootServiceAccountKeyRotation) {
			rotation.Phase = gardencorev1beta1.RotationPrepared
		})

	case gardencorev1beta1.RotationCompleting:
		v1beta1helper.MutateShootServiceAccountKeyRotation(shoot, func(rotation *gardencorev1beta1.ShootServiceAccountKeyRotation) {
			rotation.Phase = gardencorev1beta1.RotationCompleted
			rotation.LastCompletionTime = &now
		})
	}

	switch v1beta1helper.GetShootETCDEncryptionKeyRotationPhase(shoot.Status.Credentials) {
	case gardencorev1beta1.RotationPreparing:
		v1beta1helper.MutateShootETCDEncryptionKeyRotation(shoot, func(rotation *gardencorev1beta1.ShootETCDEncryptionKeyRotation) {
			rotation.Phase = gardencorev1beta1.RotationPrepared
		})

	case gardencorev1beta1.RotationCompleting:
		v1beta1helper.MutateShootETCDEncryptionKeyRotation(shoot, func(rotation *gardencorev1beta1.ShootETCDEncryptionKeyRotation) {
			rotation.Phase = gardencorev1beta1.RotationCompleted
			rotation.LastCompletionTime = &now
		})
	}

	if v1beta1helper.IsShootKubeconfigRotationInitiationTimeAfterLastCompletionTime(shoot.Status.Credentials) {
		v1beta1helper.MutateShootKubeconfigRotation(shoot, func(rotation *gardencorev1beta1.ShootKubeconfigRotation) {
			rotation.LastCompletionTime = &now
		})
	}

	if v1beta1helper.IsShootSSHKeypairRotationInitiationTimeAfterLastCompletionTime(shoot.Status.Credentials) {
		v1beta1helper.MutateShootSSHKeypairRotation(shoot, func(rotation *gardencorev1beta1.ShootSSHKeypairRotation) {
			rotation.LastCompletionTime = &now
		})
	}

	if v1beta1helper.IsShootObservabilityRotationInitiationTimeAfterLastCompletionTime(shoot.Status.Credentials) {
		v1beta1helper.MutateObservabilityRotation(shoot, func(rotation *gardencorev1beta1.ShootObservabilityRotation) {
			rotation.LastCompletionTime = &now
		})
	}

	if pointer.BoolEqual(shoot.Spec.Kubernetes.EnableStaticTokenKubeconfig, pointer.Bool(false)) && shoot.Status.Credentials != nil && shoot.Status.Credentials.Rotation != nil {
		shoot.Status.Credentials.Rotation.Kubeconfig = nil
	}

	return gardenClient.Status().Patch(ctx, shoot, patch)
}

func (r *shootReconciler) patchShootStatusOperationError(
	ctx context.Context,
	gardenClient client.Client,
	shoot *gardencorev1beta1.Shoot,
	description string,
	operationType gardencorev1beta1.LastOperationType,
	lastErrors ...gardencorev1beta1.LastError,
) error {
	var (
		now          = metav1.NewTime(time.Now().UTC())
		state        = gardencorev1beta1.LastOperationStateError
		willNotRetry = v1beta1helper.HasNonRetryableErrorCode(lastErrors...) || utils.TimeElapsed(shoot.Status.RetryCycleStartTime, r.config.Controllers.Shoot.RetryDuration.Duration)
	)

	statusPatch := client.StrategicMergeFrom(shoot.DeepCopy())

	if willNotRetry {
		state = gardencorev1beta1.LastOperationStateFailed
		shoot.Status.RetryCycleStartTime = nil
	} else {
		description += " Operation will be retried."
	}

	shoot.Status.Gardener = *r.identity
	shoot.Status.LastErrors = lastErrors

	if shoot.Status.LastOperation == nil {
		shoot.Status.LastOperation = &gardencorev1beta1.LastOperation{}
	}
	shoot.Status.LastOperation.Type = operationType
	shoot.Status.LastOperation.State = state
	shoot.Status.LastOperation.Description = description
	shoot.Status.LastOperation.LastUpdateTime = now

	return gardenClient.Status().Patch(ctx, shoot, statusPatch)
}

func lastErrorsOperationInitializationFailure(lastErrors []gardencorev1beta1.LastError, err error) []gardencorev1beta1.LastError {
	var incompleteDNSConfigError *shootpkg.IncompleteDNSConfigError

	if errors.As(err, &incompleteDNSConfigError) {
		return v1beta1helper.UpsertLastError(lastErrors, gardencorev1beta1.LastError{
			TaskID:      pointer.String(taskID),
			Description: err.Error(),
			Codes:       []gardencorev1beta1.ErrorCode{gardencorev1beta1.ErrorConfigurationProblem},
		})
	}

	return lastErrors
}

// isHibernationActive uses the Cluster resource in the Seed to determine whether the Shoot is hibernated
// The Cluster contains the actual or "active" spec of the Shoot resource for this reconciliation
// as the Shoot resources field `spec.hibernation.enabled` might have changed during the reconciliation
func (r *shootReconciler) isHibernationActive(ctx context.Context, shootSeedNamespace string, seedName *string) (bool, error) {
	if seedName == nil {
		return false, nil
	}

	seedClient, err := r.clientMap.GetClient(ctx, keys.ForSeedWithName(*seedName))
	if err != nil {
		return false, err
	}

	cluster, err := gardenerextensions.GetCluster(ctx, seedClient.Client(), shootSeedNamespace)
	if err != nil {
		return false, err
	}

	shoot := cluster.Shoot
	if shoot == nil {
		return false, fmt.Errorf("shoot is missing in cluster resource: %w", err)
	}

	return v1beta1helper.HibernationIsEnabled(shoot), nil
}

func computeOperationType(shoot *gardencorev1beta1.Shoot) gardencorev1beta1.LastOperationType {
	lastOperation := shoot.Status.LastOperation
	if lastOperation != nil && lastOperation.Type == gardencorev1beta1.LastOperationTypeMigrate &&
		(lastOperation.State == gardencorev1beta1.LastOperationStateSucceeded || lastOperation.State == gardencorev1beta1.LastOperationStateAborted) {
		return gardencorev1beta1.LastOperationTypeRestore
	}
	return v1beta1helper.ComputeOperationType(shoot.ObjectMeta, shoot.Status.LastOperation)
}

func needsControlPlaneDeployment(ctx context.Context, o *operation.Operation, kubeAPIServerDeploymentFound bool, infrastructure *extensionsv1alpha1.Infrastructure) (bool, error) {
	var (
		namespace = o.Shoot.SeedNamespace
		name      = o.Shoot.GetInfo().Name
	)

	// If the `ControlPlane` resource and the kube-apiserver deployment do no longer exist then we don't want to re-deploy it.
	// The reason for the second condition is that some providers inject a cloud-provider-config into the kube-apiserver deployment
	// which is needed for it to run.
	exists, markedForDeletion, err := extensionResourceStillExists(ctx, o.K8sSeedClient.APIReader(), &extensionsv1alpha1.ControlPlane{}, namespace, name)
	if err != nil {
		return false, err
	}

	switch {
	// treat `ControlPlane` in deletion as if it is already gone. If it is marked for deletion, we also shouldn't wait
	// for it to be reconciled, as it can potentially block the whole deletion flow (deletion depends on other control
	// plane components like kcm and grm) which are scaled up later in the flow
	case !exists && !kubeAPIServerDeploymentFound || markedForDeletion:
		return false, nil
	// The infrastructure resource has not been found, no need to redeploy the control plane
	case infrastructure == nil:
		return false, nil
	// The infrastructure resource has been found with a non-nil provider status, so redeploy the control plane
	case infrastructure.Status.ProviderStatus != nil:
		return true, nil
	default:
		return false, nil
	}
}

func extensionResourceStillExists(ctx context.Context, reader client.Reader, obj client.Object, namespace, name string) (bool, bool, error) {
	if err := reader.Get(ctx, kutil.Key(namespace, name), obj); err != nil {
		if apierrors.IsNotFound(err) {
			return false, false, nil
		}
		return false, false, err
	}
	return true, obj.GetDeletionTimestamp() != nil, nil
}

func checkIfSeedNamespaceExists(ctx context.Context, o *operation.Operation, botanist *botanistpkg.Botanist) error {
	botanist.SeedNamespaceObject = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: o.Shoot.SeedNamespace}}
	if err := botanist.K8sSeedClient.APIReader().Get(ctx, client.ObjectKeyFromObject(botanist.SeedNamespaceObject), botanist.SeedNamespaceObject); err != nil {
		if apierrors.IsNotFound(err) {
			o.Logger.Info("Did not find namespace in the Seed cluster - nothing to be done", "namespace", client.ObjectKeyFromObject(o.SeedNamespaceObject))
			return utilerrors.Cancel()
		}
		return err
	}
	return nil
}

func shootKey(shoot *gardencorev1beta1.Shoot) string {
	return shoot.Namespace + "/" + shoot.Name
}

type reconciliationDueTracker struct {
	lock    sync.Mutex
	tracker map[string]bool
}

func newReconciliationDueTracker() *reconciliationDueTracker {
	return &reconciliationDueTracker{tracker: make(map[string]bool)}
}

func (t *reconciliationDueTracker) on(key string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.tracker[key] = true
}

func (t *reconciliationDueTracker) off(key string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.tracker, key)
}

func (t *reconciliationDueTracker) tracked(key string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	v, ok := t.tracker[key]
	return ok && v
}

func startRotationCA(shoot *gardencorev1beta1.Shoot, now *metav1.Time) {
	v1beta1helper.MutateShootCARotation(shoot, func(rotation *gardencorev1beta1.ShootCARotation) {
		rotation.Phase = gardencorev1beta1.RotationPreparing
		rotation.LastInitiationTime = now
	})
}

func completeRotationCA(shoot *gardencorev1beta1.Shoot) {
	v1beta1helper.MutateShootCARotation(shoot, func(rotation *gardencorev1beta1.ShootCARotation) {
		rotation.Phase = gardencorev1beta1.RotationCompleting
	})
}

func startRotationServiceAccountKey(shoot *gardencorev1beta1.Shoot, now *metav1.Time) {
	v1beta1helper.MutateShootServiceAccountKeyRotation(shoot, func(rotation *gardencorev1beta1.ShootServiceAccountKeyRotation) {
		rotation.Phase = gardencorev1beta1.RotationPreparing
		rotation.LastInitiationTime = now
	})
}

func completeRotationServiceAccountKey(shoot *gardencorev1beta1.Shoot) {
	v1beta1helper.MutateShootServiceAccountKeyRotation(shoot, func(rotation *gardencorev1beta1.ShootServiceAccountKeyRotation) {
		rotation.Phase = gardencorev1beta1.RotationCompleting
	})
}

func startRotationETCDEncryptionKey(shoot *gardencorev1beta1.Shoot, now *metav1.Time) {
	v1beta1helper.MutateShootETCDEncryptionKeyRotation(shoot, func(rotation *gardencorev1beta1.ShootETCDEncryptionKeyRotation) {
		rotation.Phase = gardencorev1beta1.RotationPreparing
		rotation.LastInitiationTime = now
	})
}

func completeRotationETCDEncryptionKey(shoot *gardencorev1beta1.Shoot) {
	v1beta1helper.MutateShootETCDEncryptionKeyRotation(shoot, func(rotation *gardencorev1beta1.ShootETCDEncryptionKeyRotation) {
		rotation.Phase = gardencorev1beta1.RotationCompleting
	})
}

func startRotationKubeconfig(shoot *gardencorev1beta1.Shoot, now *metav1.Time) {
	v1beta1helper.MutateShootKubeconfigRotation(shoot, func(rotation *gardencorev1beta1.ShootKubeconfigRotation) {
		rotation.LastInitiationTime = now
	})
}

func startRotationSSHKeypair(shoot *gardencorev1beta1.Shoot, now *metav1.Time) {
	v1beta1helper.MutateShootSSHKeypairRotation(shoot, func(rotation *gardencorev1beta1.ShootSSHKeypairRotation) {
		rotation.LastInitiationTime = now
	})
}

func startRotationObservability(shoot *gardencorev1beta1.Shoot, now *metav1.Time) {
	v1beta1helper.MutateObservabilityRotation(shoot, func(rotation *gardencorev1beta1.ShootObservabilityRotation) {
		rotation.LastInitiationTime = now
	})
}
