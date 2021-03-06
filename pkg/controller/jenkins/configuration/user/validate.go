package user

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	virtuslabv1alpha1 "github.com/VirtusLab/jenkins-operator/pkg/apis/virtuslab/v1alpha1"
	"github.com/VirtusLab/jenkins-operator/pkg/controller/jenkins/configuration/base/resources"
	"github.com/VirtusLab/jenkins-operator/pkg/controller/jenkins/constants"
	"github.com/VirtusLab/jenkins-operator/pkg/log"

	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// Validate validates Jenkins CR Spec section
func (r *ReconcileUserConfiguration) Validate(jenkins *virtuslabv1alpha1.Jenkins) (bool, error) {
	valid, err := r.validateSeedJobs(jenkins)
	if !valid || err != nil {
		return valid, err
	}

	return r.verifyBackup()
}

func (r *ReconcileUserConfiguration) validateSeedJobs(jenkins *virtuslabv1alpha1.Jenkins) (bool, error) {
	valid := true
	if jenkins.Spec.SeedJobs != nil {
		for _, seedJob := range jenkins.Spec.SeedJobs {
			logger := r.logger.WithValues("seedJob", fmt.Sprintf("%+v", seedJob)).V(log.VWarn)

			// validate seed job id is not empty
			if len(seedJob.ID) == 0 {
				logger.Info("seed job id can't be empty")
				valid = false
			}

			// validate repository url match private key
			if strings.Contains(seedJob.RepositoryURL, "git@") {
				if seedJob.PrivateKey.SecretKeyRef == nil {
					logger.Info("private key can't be empty while using ssh repository url")
					valid = false
				}
			}

			// validate private key from secret
			if seedJob.PrivateKey.SecretKeyRef != nil {
				deployKeySecret := &v1.Secret{}
				namespaceName := types.NamespacedName{Namespace: jenkins.Namespace, Name: seedJob.PrivateKey.SecretKeyRef.Name}
				err := r.k8sClient.Get(context.TODO(), namespaceName, deployKeySecret)
				if err != nil && apierrors.IsNotFound(err) {
					logger.Info("secret not found")
					valid = false
				} else if err != nil {
					return false, err
				}

				privateKey := string(deployKeySecret.Data[seedJob.PrivateKey.SecretKeyRef.Key])
				if privateKey == "" {
					logger.Info("private key is empty")
					valid = false
				}

				if err := validatePrivateKey(privateKey); err != nil {
					logger.Info(fmt.Sprintf("private key is invalid: %s", err))
					valid = false
				}
			}
		}
	}
	return valid, nil
}

func validatePrivateKey(privateKey string) error {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return errors.New("failed to decode PEM block")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	err = priv.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileUserConfiguration) verifyBackup() (bool, error) {
	if r.jenkins.Spec.Backup == virtuslabv1alpha1.JenkinsBackupTypeAmazonS3 {
		return r.verifyBackupAmazonS3()
	}

	return true, nil
}

func (r *ReconcileUserConfiguration) verifyBackupAmazonS3() (bool, error) {
	backupSecretName := resources.GetBackupCredentialsSecretName(r.jenkins)
	backupSecret := &corev1.Secret{}
	err := r.k8sClient.Get(context.TODO(), types.NamespacedName{Namespace: r.jenkins.Namespace, Name: backupSecretName}, backupSecret)
	if err != nil {
		return false, err
	}

	if len(backupSecret.Data[constants.BackupAmazonS3SecretSecretKey]) == 0 {
		r.logger.V(log.VWarn).Info(fmt.Sprintf("Secret '%s' doesn't contains key: %s", backupSecretName, constants.BackupAmazonS3SecretSecretKey))
		return false, nil
	}

	if len(backupSecret.Data[constants.BackupAmazonS3SecretAccessKey]) == 0 {
		r.logger.V(log.VWarn).Info(fmt.Sprintf("Secret '%s' doesn't contains key: %s", backupSecretName, constants.BackupAmazonS3SecretAccessKey))
		return false, nil
	}

	return true, nil
}
