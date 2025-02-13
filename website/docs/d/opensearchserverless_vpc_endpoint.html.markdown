---
subcategory: "OpenSearch Serverless"
layout: "aws"
page_title: "AWS: aws_opensearchserverless_vpc_endpoint"
description: |-
  Terraform data source for managing an AWS OpenSearch Serverless VPC Endpoint.
---

# Data Source: aws_opensearchserverless_vpc_endpoint

Terraform data source for managing an AWS OpenSearch Serverless VPC Endpoint.

## Example Usage

```terraform
data "aws_opensearchserverless_vpc_endpoint" "example" {
  id = "vpce-829a4487959e2a839"
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) The unique identifier of the endpoint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_date` - The date the endpoint was created.
* `name` - The name of the endpoint.
* `security_group_ids` - The IDs of the security groups that define the ports, protocols, and sources for inbound traffic that you are authorizing into your endpoint.
* `subnet_ids` - The IDs of the subnets from which you access OpenSearch Serverless.
* `vpc_id` - The ID of the VPC from which you access OpenSearch Serverless.
