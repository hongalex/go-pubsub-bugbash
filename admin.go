// Copyright 2025 Google LLC
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

package bugbash

import (
	"context"
	"log"
	"time"

	// TODO: change this import
	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func setupAdmin(opts ...option.ClientOption) error {
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	defer c.Close()

	// TODO: change all of the admin operations
	_, err = c.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: fullTopicName,
		MessageStoragePolicy: &pubsubpb.MessageStoragePolicy{
			AllowedPersistenceRegions: []string{"us-central1"},
		},
		MessageRetentionDuration: durationpb.New(24 * time.Hour),
		IngestionDataSourceSettings: &pubsubpb.IngestionDataSourceSettings{
			Source: &pubsubpb.IngestionDataSourceSettings_AwsKinesis_{
				AwsKinesis: &pubsubpb.IngestionDataSourceSettings_AwsKinesis{
					StreamArn:         "fake-stream-arn",
					ConsumerArn:       "fake-consumer-arn",
					AwsRoleArn:        "fake-aws-role-arn",
					GcpServiceAccount: "fake-service-account",
				},
			},
		},
	})
	if err != nil {
		return err
	}

	// TODO: change this call
	_, err = c.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:                      fullSubName,
		Topic:                     fullTopicName,
		EnableExactlyOnceDelivery: true,
		BigqueryConfig: &pubsubpb.BigQueryConfig{
			Table: "fake-project.fake-dataset.fake-table-id",
		},
	})
	if err != nil {
		return err
	}

	// TODO: change this call
	// We are removing ingestion from the topic and switching back to pull
	// based topic instead.
	_, err = c.TopicAdminClient.UpdateTopic(ctx, &pubsubpb.UpdateTopicRequest{
		Topic: &pubsubpb.Topic{
			Name:                        fullTopicName,
			IngestionDataSourceSettings: &pubsubpb.IngestionDataSourceSettings{},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"ingestion_data_source_settings"},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
