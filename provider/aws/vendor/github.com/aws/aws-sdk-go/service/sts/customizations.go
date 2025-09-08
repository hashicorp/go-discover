// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sts

import "github.com/aws/aws-sdk-go/aws/request"

func init() {
	initRequest = customizeRequest
}

func customizeRequest(r *request.Request) {
	r.RetryErrorCodes = append(r.RetryErrorCodes, ErrCodeIDPCommunicationErrorException)
}
