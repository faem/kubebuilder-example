package v1

import (
	"github.com/robfig/cron"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var cronjoblog = logf.Log.WithName("cronjob-resource")

func (r *CronJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-batch-example-faem-io-v1-cronjob,mutating=true,failurePolicy=fail,groups=batch.example.faem.io,resources=cronjobs,verbs=create;update,versions=v1,name=mcronjob.kb.io

var _ webhook.Defaulter = &CronJob{}

func (r *CronJob) Default() {
	cronjoblog.Info("default", "name", r.Name)

	if r.Spec.ConcurrencyPolicy == "" {
		r.Spec.ConcurrencyPolicy = AllowConcurrent
	}
	if r.Spec.Suspend == nil {
		r.Spec.Suspend = new(bool)
	}
	if r.Spec.SuccessfulJobsHistoryLimit == nil {
		r.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*r.Spec.SuccessfulJobsHistoryLimit = 3
	}
	if r.Spec.FailedJobsHistoryLimit == nil {
		r.Spec.FailedJobsHistoryLimit = new(int32)
		*r.Spec.FailedJobsHistoryLimit = 1
	}
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-batch-example-faem-io-v1-cronjob,mutating=false,failurePolicy=fail,groups=batch.example.faem.io,resources=cronjobs,versions=v1,name=vcronjob.kb.io

var _ webhook.Validator = &CronJob{}

func (r *CronJob) ValidateCreate() error {
	cronjoblog.Info("validate create", "name", r.Name)

	return r.validateCronJob()
}

func (r *CronJob) ValidateUpdate(old runtime.Object) error {
	cronjoblog.Info("validate update", "name", r.Name)

	return r.validateCronJob()
}

func (r *CronJob) ValidateDelete() error {
	cronjoblog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *CronJob) validateCronJob() error {
	var allErrs field.ErrorList
	if err := r.validateCronJobName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateCronJobSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "batch.example.faem.io", Kind: "CronJob"},
		r.Name, allErrs)
}

func (r *CronJob) validateCronJobSpec() *field.Error {
	return validateScheduleFormat(
		r.Spec.Schedule,
		field.NewPath("spec").Child("schedule"))
}

func validateScheduleFormat(schedule string, fldPath *field.Path) *field.Error {
	if _, err := cron.ParseStandard(schedule); err != nil {
		return field.Invalid(fldPath, schedule, err.Error())
	}
	return nil
}

func (r *CronJob) validateCronJobName() *field.Error {
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 52 characters")
	}
	return nil
}
