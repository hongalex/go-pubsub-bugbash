# gcloud MCP Extension for Gemini CLI

This document provides instructions for an AI agent on how to use the available tools to manage Google Cloud resources using the gcloud CLI.

## Guiding Principles

*   **Prefer Specific, Native Tools**: Always prefer to use the most specific tool available. This means a GKE-specific or Cloud Run-specific tool should be used over a gcloud tool for the same functionality. This ensures better-structured data and more reliable execution.
*   **Prefer Native Tools:** Prefer to use the dedicated tools provided by this extension instead of a generic tool for shelling out to `gcloud` or `kubectl` for the same functionality. 
*   **Clarify Ambiguity:** Do not guess or assume values for required parameters like cluster names or locations. If the user's request is ambiguous, ask clarifying questions to confirm the exact resource they intend to interact with.
*   **Use Defaults:** If a `project_id` is not specified by the user, you can use the default value configured in the environment.

## gcloud Reference Documentation

General information on gcloud CLI can be found at https://cloud.google.com/sdk/gcloud

If additional context or information is needed on a gcloud command or command group, the reference documentation can be found at https://cloud.google.com/sdk/gcloud/reference.

For example, documentation on `gcloud compute instances list` can be found  at https://cloud.google.com/sdk/gcloud/reference/compute/instances/list

### gcloud --help

Reference for a specific command can also be found by appending the `--help` flag to the command. The docs returned there will be specific to the gcloud binary the user has installed, and potentially more relevant for their use-case.


## gcloud Environment Variables and Properties

Local gcloud configuration properties can be set via environment variables that gcloud commands may reference at runtime. The local configuration of a user can be viewed by running `gcloud config list`. For more information on managing gcloud CLI properties see https://cloud.google.com/sdk/gcloud/reference/config and https://cloud.google.com/sdk/docs/properties. 
