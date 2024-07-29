// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package csm

type metricException interface {
	Exception() string
	Message() string
}

type requestException struct {
	exception string
	message   string
}

func (e requestException) Exception() string {
	return e.exception
}
func (e requestException) Message() string {
	return e.message
}

type awsException struct {
	requestException
}

type sdkException struct {
	requestException
}
