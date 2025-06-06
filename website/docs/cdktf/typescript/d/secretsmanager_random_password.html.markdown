---
subcategory: "Secrets Manager"
layout: "aws"
page_title: "AWS: aws_secretsmanager_random_password"
description: |-
  Generate a random password
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_secretsmanager_random_password

Generate a random password.

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsSecretsmanagerRandomPassword } from "./.gen/providers/aws/data-aws-secretsmanager-random-password";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataAwsSecretsmanagerRandomPassword(this, "test", {
      excludeNumbers: true,
      passwordLength: 50,
    });
  }
}

```

## Argument Reference

This data source supports the following arguments:

* `excludeCharacters` - (Optional) String of the characters that you don't want in the password.
* `excludeLowercase` - (Optional) Specifies whether to exclude lowercase letters from the password.
* `excludeNumbers` - (Optional) Specifies whether to exclude numbers from the password.
* `excludePunctuation` - (Optional) Specifies whether to exclude the following punctuation characters from the password: ``! " # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \ ] ^ _ ` { | } ~ .``
* `excludeUppercase` - (Optional) Specifies whether to exclude uppercase letters from the password.
* `includeSpace` - (Optional) Specifies whether to include the space character.
* `passwordLength` - (Optional) Length of the password.
* `requireEachIncludedType` - (Optional) Specifies whether to include at least one upper and lowercase letter, one number, and one punctuation.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `randomPassword` - Random password.

<!-- cache-key: cdktf-0.20.8 input-ae0a62972b9718316151ee45b599025ddb7bfd46a48b2935ac9b105e8dfa1b5b -->