// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package profile

type HttpProfile struct {
	ReqMethod  string
	ReqTimeout int
	Endpoint   string
	Protocol   string
}

func NewHttpProfile() *HttpProfile {
	return &HttpProfile{
		ReqMethod:  "POST",
		ReqTimeout: 60,
		Endpoint:   "",
		Protocol:   "HTTPS",
	}
}
