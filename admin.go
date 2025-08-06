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
	"fmt"
	"log"
	"time"

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

	topicPath := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)
	topic, err := c.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: topicPath,
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

	subPath := fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID)
	sub, err := c.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:                      subPath,
		Topic:                     topicPath,
		EnableExactlyOnceDelivery: true,
		BigqueryConfig: &pubsubpb.BigQueryConfig{
			Table: "fake-project.fake-dataset.fake-table-id",
		},
	})
	if err != nil {
		return err
	}

	// We are removing ingestion from the topic and switching back to pull
	// based topic instead.
	_, err = c.TopicAdminClient.UpdateTopic(ctx, &pubsubpb.UpdateTopicRequest{
		Topic: &pubsubpb.Topic{
			Name:                        topicPath,
			IngestionDataSourceSettings: &pubsubpb.IngestionDataSourceSettings{},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"ingestion_data_source_settings"},
		},
	})
	if err != nil {
		return err
	}

	if err := c.TopicAdminClient.DeleteTopic(ctx, &pubsubpb.DeleteTopicRequest{
		Topic: topic.Name,
	}); err != nil {
		return err
	}

	if err := c.SubscriptionAdminClient.DeleteSubscription(ctx, &pubsubpb.DeleteSubscriptionRequest{
		Subscription: sub.Name,
	}); err != nil {
		return err
	}

	return nil
}
