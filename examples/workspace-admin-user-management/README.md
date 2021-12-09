# Google Workspace user management

Programatic interaction with Google Workspace APIs to rotate SuperAdmin credential rotations. This repo includes a cloudbuild.yaml to automate the process.

## Setup instructions
### Step 1
Fork this repo

### Step 2
Update the "Variables" section in main.go

### Step 3
Setup [Google Workspace Domain-Wide Delegation](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#go) for the service account, but no keys needs to be downloaded, please delete any SA keys if already generated.

### Step 4: This is only required for demo/testing purposes.
Assign `roles/secretmanager.admin` IAM role to the service_account created in Step 3 above.

### Step 5: Create the go binary locally
./bin/build linux  

This will require a go runtime to be installed locally or you can build the binary based of a [temporary docker container](https://hub.docker.com/_/golang?tab=description) if you have the docker runtime locally.

##### Note: This is definetly an enhancement opportunity to build the binary as a cloudbuild step. Refer to [cloud-builders/go](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/go) for go images.  

### Step 6: Container build image for Cloudbuild pipeline
Please refer [here](https://github.com/terraform-google-modules/terraform-google-bootstrap/tree/master/modules/cloudbuild/cloudbuild_builder) for creating an image with gcloud and jq utilities, unless you already have an image ready for usage.  

Please update `name` property in your `cloudbuid.yaml` to pull the appropriate container image.

### Step 7
Please setup an appropriate CloudBuild trigger for kick off the pipeline, and please refer [here](https://cloud.google.com/build/docs/automating-builds/create-scheduled-triggers) for setting up scheduled trigger.  

Add relavent [substitution variables](https://cloud.google.com/build/docs/configuring-builds/substitute-variable-values?authuser=3) for the defaults defined in the cloudbuild.yaml file. Select `Cloud build configuration file` option to be able to add substitution variables. 


## References:  
[1] Setting up [Google Workspace Domain-Wide Delegation](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#go)  
[2] [Instantiating an Admin SDK Directory service object](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#instantiate_an_admin_sdk_directory_service_object) for interacting workspace APIs.  
[3] Please refer [here](https://gist.github.com/jay0lee/75cbcd8568633ea6efd013a938f3bf25) for bash use-cases to explicitly generating Bearer token for API calls through curl command. However using a [client library](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#instantiate_an_admin_sdk_directory_service_object) is recommended.  