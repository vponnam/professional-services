# Purpose
This repo contains helper code to support IAM recommender processing and analysis.  

## How to run the code
1. git clone https://github.com/vponnam/professional-services.git && cd professional-services/examples/iam-recommender-helpers/projectnumber-converter   

1. export the below variables
    ```
    export dataset_projectid="projectID"
    export dataset_name="recommenderDatasetName"
    export table_name="table_name" #Ex: "projectnumber_id"
    export project_sync_interval="10" #In minutes, default is set to 30.
    export asset_api_scope="Parent/ID" #Ex: organizations/number
    ```

1. gcloud auth login --update-adc 

1. Execute the appropriate binary from releases folder.
    ```
    ./releases/project_converter-{os-type}
    ```

## Run/Build/Compile from the source
1. If you have golang already installed on your machine or if you like to [install](https://go.dev/doc/install)  
   
    Run from source
    ```
    go mod install
    go run *.go
    ```

    Compile from source
    ```
    env GOOS=darwin GOARCH=amd64 go build -o releases/project_converter-darwin .
    env GOOS=linux GOARCH=amd64 go build -o releases/project_converter-linux .
    env GOOS=windows GOARCH=amd64 go build -o releases/project_converter-windows.exe .
    ```

## Deploy to Cloud Run/Functions
Coming soon..
