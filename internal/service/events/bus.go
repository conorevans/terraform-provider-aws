// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package events

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_cloudwatch_event_bus", name="Event Bus")
// @Tags(identifierAttribute="arn")
func resourceBus() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceBusCreate,
		ReadWithoutTimeout:   resourceBusRead,
		UpdateWithoutTimeout: resourceBusUpdate,
		DeleteWithoutTimeout: resourceBusDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dead_letter_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrARN: {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 1600),
								verify.ValidARN,
							),
						},
					},
				},
			},
			names.AttrDescription: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
			"event_source_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validSourceName,
			},
			"kms_key_identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 2048),
			},
			names.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validCustomEventBusName,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceBusCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EventsClient(ctx)

	eventBusName := d.Get(names.AttrName).(string)
	input := &eventbridge.CreateEventBusInput{
		Name: aws.String(eventBusName),
		Tags: getTagsIn(ctx),
	}

	if v, ok := d.GetOk("dead_letter_config"); ok && len(v.([]any)) > 0 {
		input.DeadLetterConfig = expandDeadLetterConfig(v.([]any)[0].(map[string]any))
	}

	if v, ok := d.GetOk(names.AttrDescription); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("event_source_name"); ok {
		input.EventSourceName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_identifier"); ok {
		input.KmsKeyIdentifier = aws.String(v.(string))
	}

	output, err := conn.CreateEventBus(ctx, input)

	// Some partitions (e.g. ISO) may not support tag-on-create.
	if input.Tags != nil && errs.IsUnsupportedOperationInPartitionError(meta.(*conns.AWSClient).Partition(ctx), err) {
		input.Tags = nil

		output, err = conn.CreateEventBus(ctx, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EventBridge Event Bus (%s): %s", eventBusName, err)
	}

	d.SetId(eventBusName)

	// For partitions not supporting tag-on-create, attempt tag after create.
	if tags := getTagsIn(ctx); input.Tags == nil && len(tags) > 0 {
		err := createTags(ctx, conn, aws.ToString(output.EventBusArn), tags)

		// If default tags only, continue. Otherwise, error.
		if v, ok := d.GetOk(names.AttrTags); (!ok || len(v.(map[string]any)) == 0) && errs.IsUnsupportedOperationInPartitionError(meta.(*conns.AWSClient).Partition(ctx), err) {
			return append(diags, resourceBusRead(ctx, d, meta)...)
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "setting EventBridge Event Bus (%s) tags: %s", d.Id(), err)
		}
	}

	return append(diags, resourceBusRead(ctx, d, meta)...)
}

func resourceBusRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EventsClient(ctx)

	output, err := findEventBusByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EventBridge Event Bus (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EventBridge Event Bus (%s): %s", d.Id(), err)
	}

	d.Set(names.AttrARN, output.Arn)
	d.Set("dead_letter_config", flattenDeadLetterConfig(output.DeadLetterConfig))
	d.Set(names.AttrDescription, output.Description)
	d.Set("kms_key_identifier", output.KmsKeyIdentifier)
	d.Set(names.AttrName, output.Name)

	return diags
}

func resourceBusUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EventsClient(ctx)

	if d.HasChanges("dead_letter_config", names.AttrDescription, "kms_key_identifier") {
		input := &eventbridge.UpdateEventBusInput{
			Name: aws.String(d.Get(names.AttrName).(string)),
		}

		if v, ok := d.GetOk("dead_letter_config"); ok && len(v.([]any)) > 0 {
			input.DeadLetterConfig = expandDeadLetterConfig(v.([]any)[0].(map[string]any))
		}

		// To unset the description, the only way is to explicitly set it to the empty string
		if v, ok := d.GetOk(names.AttrDescription); ok {
			input.Description = aws.String(v.(string))
		} else {
			input.Description = aws.String("")
		}

		if v, ok := d.GetOk("kms_key_identifier"); ok {
			input.KmsKeyIdentifier = aws.String(v.(string))
		}

		_, err := conn.UpdateEventBus(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating EventBridge Event Bus (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceBusRead(ctx, d, meta)...)
}

func resourceBusDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EventsClient(ctx)

	log.Printf("[INFO] Deleting EventBridge Event Bus: %s", d.Id())
	_, err := conn.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{
		Name: aws.String(d.Id()),
	})

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting EventBridge Event Bus (%s): %s", d.Id(), err)
	}

	return diags
}

func findEventBusByName(ctx context.Context, conn *eventbridge.Client, name string) (*eventbridge.DescribeEventBusOutput, error) {
	input := &eventbridge.DescribeEventBusInput{
		Name: aws.String(name),
	}

	output, err := conn.DescribeEventBus(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func expandDeadLetterConfig(tfMap map[string]any) *types.DeadLetterConfig {
	if tfMap == nil {
		return nil
	}
	apiObject := &types.DeadLetterConfig{}
	if v, ok := tfMap[names.AttrARN].(string); ok && v != "" {
		apiObject.Arn = aws.String(v)
	}
	return apiObject
}

func flattenDeadLetterConfig(apiObject *types.DeadLetterConfig) []map[string]any {
	if apiObject == nil {
		return nil
	}
	tfMap := map[string]any{}
	if v := apiObject.Arn; v != nil {
		tfMap[names.AttrARN] = aws.ToString(v)
	}
	return []map[string]any{tfMap}
}
