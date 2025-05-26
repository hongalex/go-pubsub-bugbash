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
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func setupAdmin(opts ...option.ClientOption) error {
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	defer c.Close()

	// TODO: change all of the admin operations
	topic, err := c.CreateTopicWithConfig(ctx, topicID, &pubsub.TopicConfig{
		MessageStoragePolicy: pubsub.MessageStoragePolicy{
			AllowedPersistenceRegions: []string{"us-central1"},
		},
		RetentionDuration: 24 * time.Hour,
		IngestionDataSourceSettings: &pubsub.IngestionDataSourceSettings{
			Source: &pubsub.IngestionDataSourceAWSKinesis{
				StreamARN:         "fake-stream-arn",
				ConsumerARN:       "fake-consumer-arn",
				AWSRoleARN:        "fake-aws-role-arn",
				GCPServiceAccount: "fake-service-account",
			},
		},
	})
	if err != nil {
		return err
	}

	// TODO: change this call
	sub, err := c.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic:                     topic,
		EnableExactlyOnceDelivery: true,
		BigQueryConfig: pubsub.BigQueryConfig{
			Table: "fake-project.fake-dataset.fake-table-id",
		},
	})
	if err != nil {
		return err
	}

	// TODO: change this call
	// We are removing ingestion from the topic and switching back to pull
	// based topic instead.
	_, err = topic.Update(ctx, pubsub.TopicConfigToUpdate{
		IngestionDataSourceSettings: &pubsub.IngestionDataSourceSettings{},
	})
	if err != nil {
		return err
	}

	// TODO: change this call
	if err := topic.Delete(ctx); err != nil {
		return err
	}

	// TODO: change this call
	if err := sub.Delete(ctx); err != nil {
		return err
	}

	return nil
}
