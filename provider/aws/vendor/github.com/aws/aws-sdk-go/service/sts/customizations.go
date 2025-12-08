// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package sts

import "github.com/aws/aws-sdk-go/aws/request"

func init() {
	initRequest = customizeRequest
}

func customizeRequest(r *request.Request) {
	r.RetryErrorCodes = append(r.RetryErrorCodes, ErrCodeIDPCommunicationErrorException)
}
