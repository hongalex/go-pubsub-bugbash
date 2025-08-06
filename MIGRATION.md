# **Migrating from Go PubSub v1 to v2**

This guide shows how to migrate from `cloud.google.com/go` (aka the “v1”) to `cloud.google.com/go/pubsub/v2` (aka the “v2"). In line with Google's [OSS Library Breaking Change Policy](https://opensource.google/documentation/policies/library-breaking-change), we plan to support the existing v1 for 12 months, until June 30th, 2026\. This includes a commitment to bug fixes and security patches for the v1, but it will not receive new features. We encourage all users to migrate to the new v2 by the above date.

Note that this is a major version bump that includes breaking changes for the Go library specifically, but the Pub/Sub API (as defined by the [proto file](https://github.com/googleapis/googleapis/blob/master/google/pubsub/v1/pubsub.proto)) is remaining the same.

## Overview

1. RPCs for managing topics, subscriptions, schemas, and IAM will be moved to new separate clients that handle these operations.  
2. The existing `Topic` and `Subscription` structs will be renamed to `Publisher` and `Subscriber`. `Publishing` and `Receiving` will be part of these respective structs.  
3. Removed settings: (`PublishSettings.BufferedByteLimit`, `ReceiveSettings.Synchronous`, `ReceiveSettings.UseLegacyFlowControl`)  
4. Renamed settings (e.g. `MaxExtensionPeriod` → `MaxDurationPerAckExtension`)  
5. Default value change: `ReceiveSettings.NumGoroutines` now defaults to 1\)  
6. Error types related to Publisher/Subscribers rename: (e.g. `ErrTopicStopped` \-\> `ErrPublisherStopped`)

## **New Imports**

There are two new packages:

1. [`cloud.google.com/go/v2`](http://cloud.google.com/go/v2)  
2. [`cloud.google.com/go/v2/apiv1/pubsubpb`](http://cloud.google.com/go/v2/apiv1/pubsubpb)

The first package is the new main v2 package. The second is the auto-generated protobuf Go types that will be used as arguments for admin operations. See Additional References below for other relevant packages.

## **A note about snippets**

The code snippets in this guide are meant to be a quick way of comparing the differences between the v1 and v2 packages and **will not compile as-is**. For a full list of samples, you can reference our [new updated samples](https://cloud.google.com/pubsub/docs/samples).

## **Admin operations**

The Pub/Sub admin plane is used to manage Pub/Sub resources like topics, subscriptions, and schemas. These admin operations include Create, Get, Update, List, and Delete.

One of the key differences between the v1 and v2 is the change to the admin API. We are adding a new `TopicAdminClient` and `SubscriptionAdminClient` which will handle these admin operations for topics and subscriptions respectively.

For topics and subscriptions, you can access these admin clients as fields of the main client:  `pubsub.Client.TopicAdminClient` and `pubsub.Client.SubscriptionAdminClient`. These clients are pre-initialized when calling `pubsub.NewClient`, and takes in the same ClientOptions when `NewClient` is called.

There is a mostly one-to-one mapping of existing admin methods to the new admin methods.

### General RPCs

The new gRPC-based admin client generally takes in go protobuf types and returns protobuf response types. If you have used other Google Cloud Go libraries like Compute Engine or Secret Manager, this should be familiar.

Here is an example comparing creating a Topic creation in v1 and v2 libraries. In this case, [`CreateTopic`](https://pkg.go.dev/cloud.google.com/go/pubsub/v2/apiv1#TopicAdminClient.CreateTopic) will now take in a generated protobuf type, [`pubsubpb.Topic`](https://pkg.go.dev/cloud.google.com/go/pubsub/v2/apiv1/pubsubpb#Topic) that is based on the Topic defined in [pubsub.proto](https://github.com/googleapis/googleapis/blob/3808680f22d715ef59493e67a6fe82e5ae3e00dd/google/pubsub/v1/pubsub.proto#L678). A key difference here is that the `Name` of the proto type is the fully qualified name for the topic, rather than just the project ID. In addition, specifying this name is part of the [`Topic`](https://pkg.go.dev/cloud.google.com/go/pubsub/v2/apiv1/pubsubpb#Topic) struct rather than an argument for `CreateTopic`.

```go
// v1 way to create a topic

import (
	pubsub "cloud.google.com/go/pubsub"
)
...
projectID := "my-project"
topicID := "my-topic"
client, err := pubsub.NewClient(ctx, projectID)

topic, err := client.CreateTopic(ctx, topicID)
```

```go
// v2 way to create a topic
import (
	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)
...
projectID := "my-project"
topicID := "my-topic"
client, err := pubsub.NewClient(ctx, projectID)

topicpb := &pubsubpb.Topic{
	Name: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
}
topic, err := client.TopicAdminClient.CreateTopic(ctx, topicpb)
```

Note that when migrating the v1 `CreateTopicWithConfig` which currently takes in a `TopicConfig` type, you will also use the [`pubsubpb.Topic`](https://pkg.go.dev/cloud.google.com/go/pubsub/v2/apiv1/pubsubpb#Topic) type.

```go
// v1 way to create a topic with settings

import (
	pubsub "cloud.google.com/go/pubsub"
)
...
projectID := "my-project"
topicID := "my-topic"
client, err := pubsub.NewClient(ctx, projectID)

// Create a new topic with the given name and config.
topicConfig := &pubsub.TopicConfig{
	RetentionDuration: 24 * time.Hour,
	MessageStoragePolicy: pubsub.MessageStoragePolicy{
		AllowedPersistenceRegions: []string{"us-east1"},
	},
}
topic, err := client.CreateTopicWithConfig(ctx, "topicName", topicConfig)
```

```go
// v2 way to create a topic with settings
import (
	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)
...
projectID := "my-project"
topicID := "my-topic"
client, err := pubsub.NewClient(ctx, projectID)

topicpb := &pubsubpb.Topic{
	Name: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
	MessageRetentionDuration: durationpb.New(24 * time.Hour),
	MessageStoragePolicy: &pubsubpb.MessageStoragePolicy{
		AllowedPersistenceRegions: []string{"us-central1"},
	},
}
topic, err := client.TopicAdminClient.CreateTopic(ctx, topicpb)
```

Code that creates a subscription should be migrated in a similar manner: using the `pubsubpb.Subscription` type and `SubscriptionAdminClient.CreateSubscription` method.

```
s := &pubsubpb.Subscription{
	Name: fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID),
}
topic, err := client.SubscriptionAdminClient.CreateSubscription(ctx, s)
```

The new proto types and their fields may differ slightly from the current v1 types. The new types are based on the Pub/Sub proto and can be found [here](https://pkg.go.dev/cloud.google.com/go/pubsub/v2/apiv1/pubsubpb). Let’s look at some examples:

In the above CreateTopic example, message retention duration was defined as `RetentionDuration` in the v1 as a Go duration, but in the v2 it is now `MessageRetentionDuration` of type [durationpb.Duration](https://pkg.go.dev/google.golang.org/protobuf/types/known/durationpb#hdr-Conversion_from_a_Go_Duration).

In another case, generated protobuf code doesn't follow Go styling guides for initialisms. For example, `KMSKeyName` is defined as `KmsKeyName` in the v2.

In addition, the v1 uses custom `optional` types for certain fields for durations and bools. Now some custom fields such as `time.Duration` now use a protobuf specific [durationpb.Duration](https://pkg.go.dev/google.golang.org/protobuf/types/known/durationpb) and bools directly.

```
s := &pubsubpb.Subscription{
	Name: fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID),
	TopicMessageRetentionDuration: durationpb.New(1 * time.Hour),
EnableExactlyOnceDelivery: true,
}
topic, err := client.SubscriptionAdminClient.CreateSubscription(ctx, s)

```

For more specific documentation, please refer to the method calls and arguments defined by the [new clients](https://pkg.go.dev/cloud.google.com/go/pubsub/v2/apiv1) and [Go protobuf types](https://pkg.go.dev/cloud.google.com/go/pubsub/v2/apiv1/pubsubpb).

### Delete RPCs

Let’s look at the differences for another operation: `DeleteTopic`.

```go
// v1 way to delete a topic
import (
	pubsub "cloud.google.com/go/pubsub"
)
...
projectID := "my-project"
topicID := "my-topic"
client, err := pubsub.NewClient(ctx, projectID)

topic := client.Topic(topicID)
topic.Delete(ctx)
```

```go
// v2 way to delete a topic
import (
	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)
...
projectID := "my-project"
topicID := "my-topic"
client, err := pubsub.NewClient(ctx, projectID)

req := &pubsubpb.DeleteTopicRequest{
	Topic: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
}
client.TopicAdminClient.DeleteTopic(ctx, req)
```

Note in this case, you have to instantiate a `DeleteTopicRequest` object and pass that into the `DeleteTopic` call. This includes specifying the **full path** of the topic, which includes the project ID, instead of just the topic ID.

### **Update RPCs**

Update RPCs usually will need a [FieldMask protobuf type](https://pkg.go.dev/google.golang.org/protobuf/types/known/fieldmaskpb) along with the resource you are modifying. The service uses the field mask to know which fields should be updated, as zero value fields could otherwise mean both “don’t edit” or “remove this value”. The strings to pass into the update field mask should be the name of the field of the resource you are editing, written in snake\_case (e.g. `message_storage_policy`). These should match the field names in the [proto definition](https://github.com/googleapis/googleapis/blob/master/google/pubsub/v1/pubsub.proto).

If a field mask is not present on update, the operation applies to all fields (as if a field mask of all fields has been specified) and overrides the entire resource.

```go
// v1 way to update subscriptions
projectID := "my-project"
subID := "my-subscription"
client, err := pubsub.NewClient(ctx, projectID)

cfg := pubsub.SubscriptionConfigToUpdate{EnableExactlyOnceDelivery: true}
subConfig, err := client.Subscription(subID).Update(ctx, cfg)
```

```go
// v2 way to update subscriptions
import (
	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	pb "cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

projectID := "my-project"
subID := "my-subscription"
client, err := pubsub.NewClient(ctx, projectID)
updateReq := &pb.UpdateSubscriptionRequest{
	Subscription: &pb.Subscription{EnableExactlyOnceDelivery: true},
	UpdateMask: &fieldmaskpb.FieldMask{
		Paths: []string{"enable_exactly_once_delivery"},
	},
}
sub, err := client.SubscriptionAdminClient.UpdateSubscription(ctx, updateReq)
```

### **Exists method removed**

The `Exists` method for topic, subscription and schema is removed in the v2. Checking if a topic or subscription exists can be done by performing a `GetTopic` or `GetSubscription` call. However, we recommend following the pattern of [optimistically expecting a resource exists](https://cloud.google.com/pubsub/docs/samples/pubsub-optimistic-subscribe#pubsub_optimistic_subscribe-go) and then handling the not found error, which saves a network call most of the time.

### **RPCs involving one-of fields**

RPCs that include one-of fields require instantiating specific Go generated protobuf structs that satisfy the interface type. This may involve generating structs that look duplicated. This is because in the generated code, the outer struct is the interface that satisfies the one-of condition while the inner struct is a wrapper around the actual one-of.

Let’s look at an example:

```go
// v1 way to create topic ingestion from kinesis

import (
	"cloud.google.com/go/pubsub"
)
...
cfg := &pubsub.TopicConfig{
	IngestionDataSourceSettings: &pubsub.IngestionDataSourceSettings{
		Source: &pubsub.IngestionDataSourceAWSKinesis{
			StreamARN:         streamARN,
			ConsumerARN:       consumerARN,
			AWSRoleARN:        awsRoleARN,
			GCPServiceAccount: gcpServiceAccount,
		},
	},
}

topic, err := client.CreateTopicWithConfig(ctx, topicID, cfg)
```

```go
// v2 way to create topic ingestion from kinesis

import (
	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	pb "cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)
...
topicpb := &pb.Topic{
	IngestionDataSourceSettings: &pb.IngestionDataSourceSettings{
		Source: &pb.IngestionDataSourceSettings_AwsKinesis_{
			AwsKinesis: &pb.IngestionDataSourceSettings_AwsKinesis{
	StreamArn:         streamARN,
	ConsumerArn:       consumerARN,
	AwsRoleArn:        awsRoleARN,
GcpServiceAccount: gcpServiceAccount,
}
,
		},
	},
}

topic, err := client.TopicAdminClient.CreateTopic(ctx, topicpb)
```

In the above example, `IngestionDataSourceSettings_AwsKinesis_` is a wrapper struct around `IngestionDataSourceSettings_AwsKinesis.` The former satisfies the interface type of being a ingestion data source, while the latter contains the actual fields of the settings.

Another example of a one of instantiation is with [Single Message Transforms](https://cloud.google.com/pubsub/docs/smts/smts-overview).

```go
import (
	"cloud.google.com/go/pubsub"
)
projectID := "my-project"
topicID := "my-topic"
client, err := pubsub.NewClient(ctx, projectID)
...

code := `function redactSSN(message, metadata) {...}`
transform := pubsub.MessageTransform{
	Transform: pubsub.JavaScriptUDF{
		FunctionName: "redactSSN",
		Code:         code,
	},
}
cfg := &pubsub.TopicConfig{
	MessageTransforms: []pubsub.MessageTransform{transform},
}
t, err := client.CreateTopicWithConfig(ctx, topicID, cfg)
```

```go
import (
	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	pb "cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)
...
projectID := "my-project"
topicID := "my-topic"
client, err := pubsub.NewClient(ctx, projectID)

code := `function redactSSN(message, metadata) {...}`
transform := pb.MessageTransform{
	Transform: &pb.MessageTransform_JavascriptUdf{
		JavascruptUdf: &pb.JavascriptUDF {
			FunctionName: "redactSSN",
			Code: 		  code,
		},
	},
}

topicpb := &pb.Topic{
	Name: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
	MessageTransforms: []*pb.MessageTransform{transform},
}
topic, err := client.TopicAdminClient.CreateTopic(ctx, topicpb)
```

In this case, `MessageTransform_JavascriptUdf` satisfies the interface, while `JavascriptUdf` holds the actual strings relevant for the message transform.

### **Schemas**

The existing Schema client will be replaced by a new SchemaClient, which behaves similarly to the topic and subscription admin clients above. Since schemas are less commonly used than publishing and subscribing, the Pub/Sub client will not preinitialize these for you. Instead, you need to call the `NewSchemaClient` method in [`cloud.google.com/go/pubsub/v2/apiv1`](http://cloud.google.com/go/pubsub/v2/apiv1). 

```go
import (
	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

projectID := "my-project-id"
schemaID := "my-schema"
ctx := context.Background()
client, err := pubsub.NewSchemaClient(ctx)
if err != nil {
	return fmt.Errorf("pubsub.NewSchemaClient: %w", err)
}
defer client.Close()

req := &pubsubpb.GetSchemaRequest{
	Name: fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
	View: pubsubpb.SchemaView_FULL,
}
s, err := client.GetSchema(ctx, req)

```

The main difference with the new autogenerated schema client is that you no longer will pass in a project ID at client instantiation. Instead, all references to schemas will be done by its fully qualified resource name (e.g. `projects/my-project/schemas/my-schema`).

## **Data Plane Operations**

In contrast with admin operations that deal with resource management, the data plane deals with **publishing** and **receiving** messages.

In the current v1, the data plane clients are intermixed with the admin plane structs: [`Topic`](https://pkg.go.dev/cloud.google.com/go/pubsub#Topic) and [`Subscription`](https://pkg.go.dev/cloud.google.com/go/pubsub#Subscription). For example, the `Topic` struct has the [`Publish`](https://pkg.go.dev/cloud.google.com/go/pubsub#Topic.Publish) method.

```go
// Simplified v1 code
client, err := pubsub.NewClient(ctx, projectID)
...
topic := client.Topic("my-topic")
topic.Publish(ctx, "message")
```

In the v2, replace `Topic` with `Publisher` to publish messages.

```go
// Simplified v2 code
client, err := pubsub.NewClient(ctx, projectID)
...
publisher := client.Publisher("my-topic")
publisher.Publish(ctx, "message")
```

Similarly, the v1 `Subscription` has [`Receive`](https://pkg.go.dev/cloud.google.com/go/pubsub#Subscription.Receive) for pulling messages. Replace `Subscription` with `Subscriber` to pull messages.

```go
client, err := pubsub.NewClient(ctx, projectID)
...
subscriber := client.Subscriber("my-subscription")
subscriber.Receive(ctx, ...)
```

### Instantiation from admin

In the v1, it is possible to call `CreateTopic` and then call `Publish` on the returned topic. Since the v2 `CreateTopic` returns a generated protobuf [topic](https://pkg.go.dev/cloud.google.com/go/pubsub/v2/apiv1/pubsubpb#Topic) that doesn’t have a `Publish` method, you will need to instantiate your own `Publisher` client to publish messages.

```go
// Simplified v2 code
client, err := pubsub.NewClient(ctx, projectID)
...

topicpb := &pb.Topic{
	Name: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
}
topic, err := client.TopicAdminClient.CreateTopic(ctx, topicpb)

// Instantiate the publisher from the topic name.
publisher := client.Publisher(topic.GetName())
publisher.Publish(ctx, "message")
```

### **TopicInProject and SubscriptionInProject removed**

To make this transition easier, the `Publisher` and `Subscriber` methods can take in either the resource ID (e.g. `my-topic`) or fully qualified name (e.g. `projects/p/topics/topic`) as arguments. This makes it easier to use the fully qualified topic name (accessible via `topic.GetName()`) rather than needing to parse out just the topicID. If you use the resource ID, the publisher and subscriber clients will assume you are referring to the project ID defined when instantiating the base pubsub client.

The previous `TopicInProject` and `SubscriptionInProject` methods have been removed from the v2. To create a publisher or subscriber in a different project, use the fully qualified name.

### **Renamed Settings**

Two subscriber flow control settings will be renamed:

* MinExtensionPeriod → MinDurationPerAckExtension  
* MaxExtensionPeriod → MaxDurationPerAckExtension

### **Default Settings Changes**

To align with other client libraries, we will be changing the default value for `ReceiveSettings.NumGoroutines` to 1\. This is a better default for most users as each stream can handle 10 MB/s and will reduce the number of idle streams for lower throughput applications.

### **Removed Settings**

`PublishSettings.BufferedByteLimit` is removed and replaced by the equivalent `PublishSettings.MaxOutstandingBytes`.

`ReceiveSettings.Synchronous` used to make the library use the synchronous `Pull` API for the mechanism to receive messages, but we are requiring only using the `StreamingPull` API in the v2.

Lastly, we will be removing `ReceiveSettings.UseLegacyFlowControl`, since server side flow control is now a mature feature and should be relied upon for managing flow control.

# **Additional References**

## Relevant packages

1. [cloud.google.com/go/pubsub/v2](http://cloud.google.com/go/pubsub/v2) is the base v2 package  
2. [cloud.google.com/go/pubsub/v2/apiv1](http://cloud.google.com/go/pubsub/v2/apiv1) is used for initializing SchemaClient  
3. [cloud.google.com/go/pubsub/v2/apiv1/pubsubpb](http://cloud.google.com/go/pubsub/v2/apiv1/pubsubpb) is used for creating admin protobuf requests  
4. [cloud.google.com/go/iam/apiv1/iampb](http://cloud.google.com/go/iam/apiv1/iampb) for IAM requests  
5. [google.golang.org/protobuf/types/known/durationpb](http://google.golang.org/protobuf/types/known/durationpb) for proto duration type in place of Go duration  
6. [google.golang.org/protobuf/types/known/fieldmaskpb](http://google.golang.org/protobuf/types/known/fieldmaskpb) for masking which fields are updated in update calls

# **FAQ**

**Q: Why does the new admin API package mention both v2 and apiv1?**

The new Pub/Sub v2 package is `cloud.google.com/go/v2`. All of the new v2 code will live in the v2 directory. `apiv1` denotes that the Pub/Sub server API is still under v1 and is **not** changing.

**Q: Why are you changing the admin API surface?**

One goal we had for this new Pub/Sub package is to reduce confusion between the data and admin plane surfaces. Particularly, the way that this package references topics and subscriptions was inconsistent with other Pub/Sub libraries in other languages. Creating a topic does not automatically create a publisher client in the Java or Python client libraries for example. Instead, we want it to be clear that creating a topic is a server side operation and creating a publisher client is a client operation.

In the past, we have seen users be confused about why setting `topic.PublishSettings` doesn't persist the settings across applications. This is because we are actually setting the ephemeral `PublishSettings` of the client, which isn't saved to the server.

Another goal is to improve development velocity by leveraging our autogeneration tools that already exist for other Go products. With this change, changes that only affect the admin plane (including recent features such as topic ingestion settings and export subscriptions) can be released sooner.

**Q: What do I have to do to migrate?**

Primarily, migration means the following

1. Import the new [`cloud.google.com/go/v2`](http://cloud.google.com/go/v2) package  
2. Replace all instances of `Topic()` and `Subscription()` calls with `Publisher()` and `Subscriber()` instead.  
3. Migrate admin operations (`CreateTopic`, `DeleteTopic`, etc) to the v2 admin API  
4. Change data plane client instantiation. If you currently call `CreateTopic` and use the returned `Topic` to call the `Publish` RPC, you will now need to instantiate a `Publisher` client, and then use that to call `Publish`.  
5. Change settings that have been renamed in the v2.  
6. Removing references to deprecating settings (`Synchronous, BufferedByteLimit, UseLegacyFlowControl`)  
7. Migrate references to migrated error types: `ErrTopicStopped` \-\> `ErrPublisherStopped`.
