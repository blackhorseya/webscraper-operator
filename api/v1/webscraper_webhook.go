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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var webscraperlog = logf.Log.WithName("webscraper-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *WebScraper) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-batch-seancheng-space-v1-webscraper,mutating=true,failurePolicy=fail,sideEffects=None,groups=batch.seancheng.space,resources=webscrapers,verbs=create;update,versions=v1,name=mwebscraper.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &WebScraper{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *WebScraper) Default() {
	webscraperlog.Info("default", "name", r.Name)

	if r.Spec.Retries == 0 {
		r.Spec.Retries = 3
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-batch-seancheng-space-v1-webscraper,mutating=false,failurePolicy=fail,sideEffects=None,groups=batch.seancheng.space,resources=webscrapers,verbs=create;update,versions=v1,name=vwebscraper.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &WebScraper{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *WebScraper) ValidateCreate() (admission.Warnings, error) {
	webscraperlog.Info("validate create", "name", r.Name)

	if r.Spec.Retries < 0 {
		return admission.Warnings{"Retries cannot be negative"}, nil
	}

	if r.Spec.Retries > 10 {
		return admission.Warnings{"Retries cannot be greater than 10"}, nil
	}

	if r.Spec.Schedule == "" {
		return admission.Warnings{"Schedule cannot be empty"}, nil
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *WebScraper) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	webscraperlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *WebScraper) ValidateDelete() (admission.Warnings, error) {
	webscraperlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
