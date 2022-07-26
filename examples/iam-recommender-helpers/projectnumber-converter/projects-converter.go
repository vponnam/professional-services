/*
# Copyright 2022 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

/* Converts a list of project numbers to IDs.
See: https://cloud.google.com/go/docs/reference/cloud.google.com/go/asset/latest/apiv1
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

var (
	// Varibales for BQ detection
	dataset_projectid     = os.Getenv("dataset_projectid")
	dataset_name          = os.Getenv("dataset_name")
	table_name            = os.Getenv("table_name")
	project_sync          = os.Getenv("project_sync")          //Set to enable for scheduled execution. Default is enable
	project_sync_interval = os.Getenv("project_sync_interval") // How frequent the BQ data task is executed. Default is 30(mins)
	asset_api_scope       = os.Getenv("asset_api_scope")       //Ex: "organizations/614859914915"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}

/*
See:
- https://cloud.google.com/go/docs/reference/cloud.google.com/go/asset/latest/apiv1#cloud_google_com_go_asset_apiv1_Client_SearchAllResources
- https://cloud.google.com/go/docs/reference/cloud.google.com/go/asset/latest/apiv1
*/
func projectIDConverter() {
	prjIDmap := make(map[string]string)

	ctx := context.Background()

	client, err := asset.NewClient(ctx)
	if err != nil {
		checkError(err)
	}

	defer client.Close()

	req := &assetpb.SearchAllResourcesRequest{
		Scope:      asset_api_scope, // Change this to a variable
		AssetTypes: []string{"compute.googleapis.com/Project"},
	}

	resp := client.SearchAllResources(ctx, req)

	for {
		result, err := resp.Next()
		if err == iterator.Done {
			break
		}

		// fmt.Printf("%v, \t%v\n", result.DisplayName, strings.Split(result.Project, "/")[1])
		// Build the hashmap with Number as key and Name as value. projects/1234
		prjIDmap[strings.Split(result.Project, "/")[1]] = result.DisplayName
	}

	// fmt.Println(prjIDmap)
	bqWrite(prjIDmap)
}

/*
Write to BQ table
See:
https://cloud.google.com/bigquery/docs/quickstarts/quickstart-client-libraries#client-libraries-install-go
https://cloud.google.com/bigquery/docs/writing-results#writing_query_results
DDL: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language
*/
func bqWrite(prjmappingData map[string]string) {

	// fmt.Println(dataset_projectid, dataset_name, table_name)
	// validation
	if dataset_projectid == "" || dataset_name == "" || table_name == "" {
		fmt.Println("One of more of these variable are not found: dataset_projectid|dataset_name|table_name")
		return
	}
	ctx := context.Background()
	table := fmt.Sprintf("%v.%v", dataset_name, table_name)
	bqData := make(map[string]string)

	client, err := bigquery.NewClient(ctx, dataset_projectid)
	checkError(err)
	defer client.Close()

	// Check table
	createTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (projectNumber STRING, projectID STRING)", table)
	createTable := client.Query(createTableQuery)

	job, err := createTable.Run(ctx)
	checkError(err)

	status, err := job.Wait(ctx)
	checkError(err)
	if status.Err() != nil {
		checkError(status.Err())
	}

	//  Read from BQ first to determine missing projects
	queryRead := fmt.Sprintf("SELECT projectNumber, projectID from %v", table)
	getProjects := client.Query(queryRead)

	getJob, err := getProjects.Run(ctx)
	checkError(err)
	getStatus, err := getJob.Wait(ctx)
	checkError(err)
	if getStatus.Err() != nil {
		checkError(getStatus.Err())
	}

	rows, err := getJob.Read(ctx)
	checkError(err)

	for {
		var row []bigquery.Value
		err := rows.Next(&row)
		if err == iterator.Done {
			break
		}
		checkError(err)
		// For type conversion. bigquery.Value => string
		number := fmt.Sprintf("%v", row[0])
		prjID := fmt.Sprintf("%v", row[1])
		bqData[number] = prjID
	}

	// Define unique map.
	missingProjects := removeDups(prjmappingData, bqData)

	// Write missing projects to BQ
	insertValues := []string{}

	for number, id := range missingProjects {
		// Build multiple values to insert as a single BQ write operation
		insertValues = append(insertValues, fmt.Sprintf("(%q, %q)", number, id))
	}

	if len(insertValues) >= 1 {

		queryInsert := fmt.Sprintf("INSERT INTO %v (projectNumber, projectID) VALUES %v", table, strings.Join(insertValues[:], ","))
		// fmt.Println("Insert Query: " + queryInsert)
		addProject := client.Query(queryInsert)

		insertJob, err := addProject.Run(ctx)
		checkError(err)
		insertStatus, err := insertJob.Wait(ctx)
		checkError(err)
		if insertStatus.Err() != nil {
			checkError(insertStatus.Err())
		}

		fmt.Printf("Wrote %v projects to BQ", len(missingProjects))
	} else {
		fmt.Println("All projects are in sync, no new projects found.")
	}
}

// Returns a unique map of key-values that are not in BQ, given the 2 maps.
func removeDups(APImap, BQmap map[string]string) (uniqueMap map[string]string) {
	uniqueMap = make(map[string]string)

	for number, id := range APImap {
		if id != BQmap[number] {
			uniqueMap[number] = id
		}
	}
	return uniqueMap
}
