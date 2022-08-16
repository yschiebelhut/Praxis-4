package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
	"time"

	awsutil "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.wdf.sap.corp/ICN-ML/aicore/system-services/platform/pkg/log"
)

var (
	maxKeyAgeInDays = 7
	templateFile    = "terraform_policy.tpl"
)

type Aws struct {
	cfg                                                         awsutil.Config
	s3                                                          *s3.Client
	iam                                                         *iam.Client
	sts                                                         *sts.Client
	callerAccountID, policyArn                                  string
	clusterFullName, clusterName                                string
	terraformUsername, terraformPolicyName, terraformBucketName string
}

func New(clusterRegion, clusterName, clusterType string) (*Aws, error) {
	aws := new(Aws)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Sugar.Errorf("Error while loading aws config! Error: %v\n", err.Error())
		return nil, err
	}
	aws.cfg = cfg

	aws.cfg.Region = clusterRegion
	aws.clusterName = clusterName
	aws.clusterFullName = "aws." + clusterRegion + "." + clusterName

	aws.terraformUsername = fmt.Sprintf("tf-%v--bot", aws.clusterFullName)
	aws.terraformPolicyName = fmt.Sprintf("%v", aws.clusterFullName)
	aws.terraformBucketName = fmt.Sprintf("%v", aws.clusterFullName)

	aws.s3 = s3.NewFromConfig(aws.cfg)
	aws.iam = iam.NewFromConfig(aws.cfg)
	aws.sts = sts.NewFromConfig(aws.cfg)

	// set callerAccountID and policyArn
	callerIdentity, err := aws.sts.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Sugar.Errorf("Error while obtaining aws caller identity! Error: %v\n", err.Error())
		return nil, err
	}
	aws.callerAccountID = *callerIdentity.Account
	aws.policyArn = fmt.Sprintf("arn:aws:iam::%v:policy/%v", aws.callerAccountID, aws.terraformPolicyName)
	log.Sugar.Debug("created aws")

	return aws, nil
}

func (aws Aws) CreateTerraformStateBucket() error {
	log.Sugar.Debugf("Testing if Bucket %v already exists...\n", aws.terraformBucketName)
	_, err := aws.s3.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: &aws.terraformBucketName,
	})

	if err != nil {
		if strings.Contains(err.Error(), "StatusCode: 404") {
			log.Sugar.Debugf("Bucket %v does not exist.\n", aws.terraformBucketName)
			log.Sugar.Debugf("Creating bucket %v.\n", aws.terraformBucketName)

			// specify location constraint only if region is not us-east-1
			if aws.cfg.Region == "us-east-1" {
				_, err = aws.s3.CreateBucket(context.TODO(), &s3.CreateBucketInput{
					Bucket: &aws.terraformBucketName,
				})
			} else {
				_, err = aws.s3.CreateBucket(context.TODO(), &s3.CreateBucketInput{
					Bucket: &aws.terraformBucketName,
					CreateBucketConfiguration: &types.CreateBucketConfiguration{
						LocationConstraint: types.BucketLocationConstraint(aws.cfg.Region),
					},
				})
			}

			if err != nil {
				log.Sugar.Errorf("Error while creating bucket %v! Error: %v\n", aws.terraformBucketName, err.Error())
				return err
			}

			log.Sugar.Debugf("Bucket %v created successfully.\n", aws.terraformBucketName)
		} else {
			log.Sugar.Errorf("Error when checking for existance of bucket %v! Error: %v\n", aws.terraformBucketName, err.Error())
			return err
		}
	} else {
		log.Sugar.Debug("Bucket %v already exists. Proceeding with next step.\n", aws.terraformBucketName)
	}
	return nil
}

func (aws Aws) DoesTerraformUserExists() (bool, error) {
	log.Sugar.Debugf("Checking if user %v exists...\n", aws.terraformUsername)
	input := &iam.GetUserInput{
		UserName: &aws.terraformUsername,
	}

	_, err := aws.iam.GetUser(context.TODO(), input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.(type) {
			case *iamtypes.NoSuchEntityException:
				log.Sugar.Debugf("User %v does not exist.\n", aws.terraformUsername)
				return false, nil
			}
		} else {
			log.Sugar.Errorf("Error when checking for existance of user %v. Error: %v\n", aws.terraformUsername, err.Error())
			return false, err
		}
	}

	log.Sugar.Debugf("User %v already exists.\n", aws.terraformUsername)
	return true, nil
}

func (aws Aws) CreateTerraformUser() (bool, error) {
	log.Sugar.Debugf("Attempting to create user %v...\n", aws.terraformUsername)
	input := &iam.CreateUserInput{
		UserName: &aws.terraformUsername,
	}

	_, err := aws.iam.CreateUser(context.TODO(), input)

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.(type) {
			case *iamtypes.EntityAlreadyExistsException:
				log.Sugar.Debugf("User %v already exists.\n", aws.terraformUsername)
				return false, nil
			}
		}
		log.Sugar.Errorf("Error while creating user %v! Error: %v\n", aws.terraformUsername, err.Error())
		return false, err
	}

	log.Sugar.Debugf("User %v created successfully!\n", aws.terraformUsername)
	return true, nil
}

// generate policy json as string from template
func (aws Aws) generatePolicy() (string, error) {
	policyTemplateString, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Sugar.Errorf("Error while loading policy template from file %v! Error: %v\n", templateFile, err.Error())
		return "", err
	}
	policyTemplate, err := template.New("bootstrap_aws_aicore_terraform_user_policy").Parse(string(policyTemplateString))
	if err != nil {
		log.Sugar.Errorf("Error when parsing policy template! Error: %v\n", err.Error())
		return "", err
	}

	var policy bytes.Buffer
	if err := policyTemplate.Execute(&policy, aws.clusterFullName); err != nil {
		log.Sugar.Errorf("Error when executing policy template! Error: %v\n", err.Error())
		return "", err
	}
	return policy.String(), nil
}

func (aws Aws) CreateUpdateTerraformPolicy() error {
	policyExists, err := aws.doesTerraformPolicyExists()
	if err != nil {
		return err
	}
	if !policyExists {
		log.Sugar.Debugf("Attempting to create policy %v.\n", aws.terraformPolicyName)
		policyString, err := aws.generatePolicy()
		if err != nil {
			return err
		}
		output, err := aws.iam.CreatePolicy(context.TODO(), &iam.CreatePolicyInput{
			PolicyName:     &aws.terraformPolicyName,
			PolicyDocument: &policyString,
			Description:    awsutil.String("access policy for technical terraform user"),
		})
		if err != nil {
			log.Sugar.Errorf("Error when creating policy %v! Error: %v\n", aws.terraformPolicyName, err.Error())
			return err
		}
		log.Sugar.Debugf("Created policy %v with arn: %v.\n", aws.terraformPolicyName, output.Policy.Arn)
	} else {
		log.Sugar.Debugf("Checking number of different versions of policy %v.\n", aws.terraformPolicyName)
		listOutput, err := aws.iam.ListPolicyVersions(context.TODO(), &iam.ListPolicyVersionsInput{
			PolicyArn: &aws.policyArn,
		})
		if err != nil {
			log.Sugar.Errorf("Error when listing verions of policy %v! Error: %v\n", aws.terraformPolicyName, err.Error())
			return err
		}

		if len(listOutput.Versions) >= 5 {
			log.Sugar.Debugf("Five or more versions detected for policy %v. Deleting latest!\n", aws.terraformPolicyName)
			for _, policyVersion := range listOutput.Versions {
				if !policyVersion.IsDefaultVersion {
					_, err = aws.iam.DeletePolicyVersion(context.TODO(), &iam.DeletePolicyVersionInput{
						PolicyArn: &aws.policyArn,
						VersionId: policyVersion.VersionId,
					})
					if err != nil {
						log.Sugar.Errorf("Error when deleting policy version in %v! Error: %v\n", aws.terraformPolicyName, err.Error())
						return err
					}
				}
			}
		}

		policyString, err := aws.generatePolicy()
		if err != nil {
			return err
		}
		_, err = aws.iam.CreatePolicyVersion(context.TODO(), &iam.CreatePolicyVersionInput{
			PolicyArn:      &aws.policyArn,
			PolicyDocument: &policyString,
			SetAsDefault:   true,
		})
		if err != nil {
			log.Sugar.Errorf("Error updating policy %v! Error: %v\n", aws.terraformPolicyName, err.Error())
			return err
		}
		log.Sugar.Debugf("Policy for terraform user %v has been updated.\n", aws.terraformUsername)
	}
	return nil
}

func (aws Aws) AttachTerraformPolicyToTerraformUser() error {
	policyAttached, err := aws.isTerraformPolicyAttached()
	if err != nil {
		return err
	}
	if policyAttached {
		log.Sugar.Debugf("Policy %v is already attached to %v.\n", aws.terraformPolicyName, aws.terraformUsername)
	} else {
		_, err := aws.iam.AttachUserPolicy(context.TODO(), &iam.AttachUserPolicyInput{
			PolicyArn: &aws.policyArn,
			UserName:  &aws.terraformUsername,
		})
		if err != nil {
			log.Sugar.Errorf("Error when attaching policy %v to user %v! Error: %v\n", aws.terraformPolicyName, aws.terraformUsername, err.Error())
			return err
		}
		log.Sugar.Debugf("Attached policy %v to user %v.\n", aws.terraformPolicyName, aws.terraformUsername)
	}
	return nil
}

func (aws Aws) doesTerraformPolicyExists() (bool, error) {
	log.Sugar.Debugf("Checking for existance of policy %v.\n", aws.policyArn)
	_, err := aws.iam.GetPolicy(context.TODO(), &iam.GetPolicyInput{
		PolicyArn: &aws.policyArn,
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.(type) {
			case *iamtypes.NoSuchEntityException:
				log.Sugar.Debugf("Policy %v does not exist.\n", aws.policyArn)
				return false, nil
			}
		}
		log.Sugar.Errorf("Error occured when checking for existance of policy %v! Error: %v\n", aws.policyArn, err)
		return false, err
	}
	log.Sugar.Debugf("Policy %v exists.\n", aws.policyArn)
	return true, nil
}

func (aws Aws) isTerraformPolicyAttached() (bool, error) {
	log.Sugar.Debugf("Checking if policy %v is attached to user %v.\n", aws.terraformPolicyName, aws.terraformUsername)
	output, err := aws.iam.ListAttachedUserPolicies(context.TODO(), &iam.ListAttachedUserPoliciesInput{
		UserName: &aws.terraformUsername,
	})
	if err != nil {
		log.Sugar.Errorf("Error when listing policies of user %v! Error: %v\n", aws.terraformUsername, err.Error())
		return false, err
	}

	for _, policy := range output.AttachedPolicies {
		if policy.PolicyName == &aws.terraformPolicyName {
			log.Sugar.Debugf("Policy %v is already attached to user %v.\n", aws.terraformPolicyName, aws.terraformUsername)
			return true, nil
		}
	}
	log.Sugar.Debugf("Policy %v is not yet attached to user %v.\n", aws.terraformPolicyName, aws.terraformUsername)
	return false, nil
}

func (aws Aws) getAccessKeys() ([]iamtypes.AccessKeyMetadata, error) {
	log.Sugar.Debugf("Retrieving access keys of user %v.\n", aws.terraformUsername)
	result, err := aws.iam.ListAccessKeys(context.TODO(), &iam.ListAccessKeysInput{
		UserName: &aws.terraformUsername,
	})

	if err != nil {
		log.Sugar.Errorf("Error when listing access keys of user %v! Error: %v\n", aws.terraformUsername, err.Error())
		return nil, err
	}
	return result.AccessKeyMetadata, nil
}

func (aws Aws) checkDeleteAccessKeys(accessKeys []iamtypes.AccessKeyMetadata) (createNewKey bool, err error) {
	log.Sugar.Debugf("Checking existing keys of user %v.\n", aws.terraformUsername)
	// if there are no access keys, a new one has to be created
	if len(accessKeys) == 0 {
		log.Sugar.Debugf("User %v does not yet have access keys. Creating new key.\n", aws.terraformUsername)
		createNewKey = true
		return
	}

	if len(accessKeys) == 1 {
		// if the key is older than maxKeyAgeInDays create new key
		if time.Since(*accessKeys[0].CreateDate).Hours() > float64(maxKeyAgeInDays)*24*time.Hour.Hours() {
			log.Sugar.Debugf("Existing key of user %v exceeds maximum age of %v days. Crating new key.\n", aws.terraformUsername, maxKeyAgeInDays)
			createNewKey = true
			return
		}
		log.Sugar.Debugf("1 access key exists for user %v. Key is younger than maximum age.\n", aws.terraformUsername)
		createNewKey = false
		return
	}

	log.Sugar.Debugf("2 access keys exist for user %v. Determining older key.\n", aws.terraformUsername)
	// two keys exist
	// indentify the olderKey key
	olderKey := accessKeys[0]
	if accessKeys[1].CreateDate.Before(*olderKey.CreateDate) {
		olderKey = accessKeys[1]
	}
	// if the older key is to old, delete old and create new
	if time.Since(*olderKey.CreateDate).Hours() > float64(maxKeyAgeInDays)*24*time.Hour.Hours() {
		log.Sugar.Debugf("Older access key exceeded maximum age of %v days. Deleting key %v.\n", maxKeyAgeInDays, olderKey.AccessKeyId)
		_, err = aws.iam.DeleteAccessKey(context.TODO(), &iam.DeleteAccessKeyInput{
			AccessKeyId: olderKey.AccessKeyId,
			UserName:    olderKey.UserName,
		})

		if err != nil {
			log.Sugar.Errorf("Error when deleting access key %v of user %v! Error: %v\n", olderKey.AccessKeyId, aws.terraformUsername, err.Error())
			createNewKey = false
			return
		}

		log.Sugar.Debugf("Deleted access key %v. Now creating new key.\n", olderKey.AccessKeyId)
		createNewKey = true
		return
	}

	log.Sugar.Debugf("Older key is younger than minimum age.\n")
	createNewKey = false
	return

	// if no key, make one and save to vault
	// if one key, make a new one if key older than 7 days and save to vault
	// if two keys, delete older one (if key older than 7 days) and save new to vault

}

func (aws Aws) createAccessKey() (iamtypes.AccessKey, error) {
	log.Sugar.Debugf("Attempting to create access key for user %v.\n", aws.terraformUsername)
	result, err := aws.iam.CreateAccessKey(context.TODO(), &iam.CreateAccessKeyInput{
		UserName: &aws.terraformUsername,
	})

	if err != nil {
		log.Sugar.Errorf("Error when creating access key for user %v! Error: %v\n", aws.terraformUsername, err.Error())
		return iamtypes.AccessKey{}, err
	}

	return *result.AccessKey, nil
}

func (aws Aws) CheckCreateAccessSecretKeyTerraformUser() (newKey iamtypes.AccessKey, created bool, err error) {
	accessKeys, err := aws.getAccessKeys()
	if err != nil {
		return
	}
	createdNewKey, err := aws.checkDeleteAccessKeys(accessKeys)
	if err != nil {
		return
	}
	if createdNewKey {
		newKey, err = aws.createAccessKey()
		if err != nil {
			return
		}
		created = true
		return
	}

	created = false
	return
}

func (aws Aws) getVaultResourcePath() string {
	return fmt.Sprintf("product/foundation/environments/%v/terraform-aws-credential", aws.clusterFullName)
}
