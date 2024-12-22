/*
Copyright 2024.

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

package controller

import (
	"context"

	sbatchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "github.com/blackhorseya/webscraper-operator/api/v1"
)

// WebScraperReconciler reconciles a WebScraper object
type WebScraperReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=batch.seancheng.space,resources=webscrapers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.seancheng.space,resources=webscrapers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch.seancheng.space,resources=webscrapers/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=cronjobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebScraper object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *WebScraperReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get the WebScraper object
	webScraper := &batchv1.WebScraper{}
	if err := r.Get(ctx, req.NamespacedName, webScraper); err != nil {
		logger.Error(err, "unable to fetch WebScraper")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get or create the CronJob
	cronJob := &sbatchv1.CronJob{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: webScraper.Namespace,
		Name:      webScraper.Name,
	}, cronJob)
	if err != nil && errors.IsNotFound(err) {
		var newCronJob *sbatchv1.CronJob
		newCronJob, err = generateCronJob(webScraper, r.Scheme)
		if err != nil {
			logger.Error(err, "unable to generate CronJob")
			return ctrl.Result{}, err
		}
		if err = r.Create(ctx, newCronJob); err != nil {
			logger.Error(err, "unable to create CronJob")
			return ctrl.Result{}, err
		}
		logger.Info("CronJob created successfully", "name", newCronJob.Name)
	} else if err != nil {
		logger.Error(err, "unable to fetch CronJob")
		return ctrl.Result{}, err
	}

	// Check Jobs associated with the CronJob
	jobList := &sbatchv1.JobList{}
	listOpts := []client.ListOption{
		client.InNamespace(webScraper.Namespace),
		client.MatchingLabels{"webscraper-name": webScraper.Name},
	}
	if err := r.List(ctx, jobList, listOpts...); err != nil {
		logger.Error(err, "unable to list Jobs")
		return ctrl.Result{}, err
	}

	// Update WebScraper status based on Job status
	for _, job := range jobList.Items {
		if job.Status.Failed > 0 {
			if err = r.Get(ctx, req.NamespacedName, webScraper); err != nil {
				logger.Error(err, "unable to fetch WebScraper")
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			webScraper.Status.LastRunTime = metav1.Now()
			webScraper.Status.Success = false
			webScraper.Status.Message = "Task failed. Retrying..."
			if err := r.Status().Update(ctx, webScraper); err != nil {
				logger.Error(err, "unable to update WebScraper status for failed task")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		if job.Status.Succeeded > 0 {
			if err = r.Get(ctx, req.NamespacedName, webScraper); err != nil {
				logger.Error(err, "unable to fetch WebScraper")
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}

			webScraper.Status.LastRunTime = metav1.Now()
			webScraper.Status.Success = true
			webScraper.Status.Message = "Task succeeded"
			if err := r.Status().Update(ctx, webScraper); err != nil {
				logger.Error(err, "unable to update WebScraper status for succeeded task")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	// Default status update if no Jobs are found
	if len(jobList.Items) == 0 {
		if err = r.Get(ctx, req.NamespacedName, webScraper); err != nil {
			logger.Error(err, "unable to fetch WebScraper")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		webScraper.Status.Message = "No jobs found for this CronJob"
		if err := r.Status().Update(ctx, webScraper); err != nil {
			logger.Error(err, "unable to update WebScraper status for no jobs")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebScraperReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.WebScraper{}).
		Owns(&sbatchv1.CronJob{}).
		Complete(r)
}

func generateCronJob(webScraper *batchv1.WebScraper, scheme *runtime.Scheme) (*sbatchv1.CronJob, error) {
	cronJob := &sbatchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webScraper.Name,
			Namespace: webScraper.Namespace,
			Labels:    map[string]string{"webscraper-name": webScraper.Name},
		},
		Spec: sbatchv1.CronJobSpec{
			Schedule: webScraper.Spec.Schedule,
			JobTemplate: sbatchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"webscraper-name": webScraper.Name},
				},
				Spec: sbatchv1.JobSpec{
					BackoffLimit: ptr.To(webScraper.Spec.Retries),
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:      "scraper",
									Image:     webScraper.Spec.Image,
									Command:   webScraper.Spec.Command,
									Resources: webScraper.Spec.Resources,
								},
							},
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}

	// 设置 OwnerReference
	if err := ctrl.SetControllerReference(webScraper, cronJob, scheme); err != nil {
		return nil, err
	}

	return cronJob, nil
}
