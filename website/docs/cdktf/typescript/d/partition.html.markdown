---
subcategory: "Meta Data Sources"
layout: "aws"
page_title: "AWS: aws_partition"
description: |-
  Get AWS partition identifier
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_partition

Use this data source to lookup information about the current AWS partition in
which Terraform is working.

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsIamPolicyDocument } from "./.gen/providers/aws/data-aws-iam-policy-document";
import { DataAwsPartition } from "./.gen/providers/aws/data-aws-partition";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    const current = new DataAwsPartition(this, "current", {});
    new DataAwsIamPolicyDocument(this, "s3_policy", {
      statement: [
        {
          actions: ["s3:ListBucket"],
          resources: ["arn:${" + current.partition + "}:s3:::my-bucket"],
          sid: "1",
        },
      ],
    });
  }
}

```

## Argument Reference

This data source does not support any arguments.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `dnsSuffix` - Base DNS domain name for the current partition (e.g., `amazonaws.com` in AWS Commercial, `amazonaws.com.cn` in AWS China).
* `id` - Identifier of the current partition (e.g., `aws` in AWS Commercial, `aws-cn` in AWS China).
* `partition` - Identifier of the current partition (e.g., `aws` in AWS Commercial, `aws-cn` in AWS China).
* `reverseDnsPrefix` - Prefix of service names (e.g., `com.amazonaws` in AWS Commercial, `cn.com.amazonaws` in AWS China).

<!-- cache-key: cdktf-0.20.8 input-c2808e72b3be3bcfa58fdfc5edfbbcb30356e9bd076385755b9662e03e3b477f -->