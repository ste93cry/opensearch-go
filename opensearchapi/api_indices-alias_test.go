// SPDX-License-Identifier: Apache-2.0
//
// The OpenSearch Contributors require contributions made to
// this file be licensed under the Apache-2.0 license or a
// compatible open source license.
//
//go:build integration && (core || opensearchapi)

package opensearchapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ostest "github.com/opensearch-project/opensearch-go/v4/internal/test"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	osapitest "github.com/opensearch-project/opensearch-go/v4/opensearchapi/internal/test"
)

func TestAliasClient(t *testing.T) {
	client, err := ostest.NewClient()
	require.Nil(t, err)
	failingClient, err := osapitest.CreateFailingClient()
	require.Nil(t, err)

	index := "test-index-test"
	alias := "test-alias-test"

	t.Cleanup(func() {
		client.Indices.Delete(nil, opensearchapi.IndicesDeleteReq{Indices: []string{index}})
		client.Indices.Alias.Delete(nil, opensearchapi.AliasDeleteReq{Alias: []string{alias}})
	})

	_, err = client.Indices.Create(nil, opensearchapi.IndicesCreateReq{Index: index})
	require.Nil(t, err)

	type aliasTests struct {
		Name    string
		Results func() (osapitest.Response, error)
	}

	testCases := []struct {
		Name  string
		Tests []aliasTests
	}{
		{
			Name: "Put",
			Tests: []aliasTests{
				{
					Name: "with request",
					Results: func() (osapitest.Response, error) {
						return client.Indices.Alias.Put(nil, opensearchapi.AliasPutReq{
							Indices: []string{index},
							Alias:   alias,
						})
					},
				},
				{
					Name: "inspect",
					Results: func() (osapitest.Response, error) {
						return failingClient.Indices.Alias.Put(nil, opensearchapi.AliasPutReq{
							Indices: []string{index},
							Alias:   alias,
						})
					},
				},
			},
		},
		{
			Name: "Get",
			Tests: []aliasTests{
				{
					Name: "with request",
					Results: func() (osapitest.Response, error) {
						return client.Indices.Alias.Get(nil, opensearchapi.AliasGetReq{
							Indices: []string{index},
							Alias:   []string{alias},
						})
					},
				},
				{
					Name: "with request without indices",
					Results: func() (osapitest.Response, error) {
						return client.Indices.Alias.Get(nil, opensearchapi.AliasGetReq{
							Alias: []string{alias},
						})
					},
				},
				{
					Name: "inspect",
					Results: func() (osapitest.Response, error) {
						return failingClient.Indices.Alias.Get(nil, opensearchapi.AliasGetReq{
							Indices: []string{index},
							Alias:   []string{alias},
						})
					},
				},
			},
		},
		{
			Name: "Exists",
			Tests: []aliasTests{
				{
					Name: "with request",
					Results: func() (osapitest.Response, error) {
						var (
							resp osapitest.DummyInspect
							err  error
						)
						resp.Response, err = client.Indices.Alias.Exists(nil, opensearchapi.AliasExistsReq{
							Indices: []string{index},
							Alias:   []string{alias},
						})
						return resp, err
					},
				},
				{
					Name: "inspect",
					Results: func() (osapitest.Response, error) {
						var (
							resp osapitest.DummyInspect
							err  error
						)
						resp.Response, err = failingClient.Indices.Alias.Exists(nil, opensearchapi.AliasExistsReq{
							Indices: []string{index},
							Alias:   []string{alias},
						})
						return resp, err
					},
				},
			},
		},
		{
			Name: "Delete",
			Tests: []aliasTests{
				{
					Name: "with request",
					Results: func() (osapitest.Response, error) {
						return client.Indices.Alias.Delete(nil, opensearchapi.AliasDeleteReq{
							Indices: []string{index},
							Alias:   []string{alias},
						})
					},
				},
				{
					Name: "inspect",
					Results: func() (osapitest.Response, error) {
						return failingClient.Indices.Alias.Delete(nil, opensearchapi.AliasDeleteReq{
							Indices: []string{index},
							Alias:   []string{alias},
						})
					},
				},
			},
		},
	}
	for _, value := range testCases {
		t.Run(value.Name, func(t *testing.T) {
			for _, testCase := range value.Tests {
				t.Run(testCase.Name, func(t *testing.T) {
					res, err := testCase.Results()
					if testCase.Name == "inspect" {
						assert.NotNil(t, err)
						assert.NotNil(t, res)
						osapitest.VerifyInspect(t, res.Inspect())
					} else {
						require.Nil(t, err)
						require.NotNil(t, res)
						assert.NotNil(t, res.Inspect().Response)
						if value.Name != "Get" && value.Name != "Exists" {
							ostest.CompareRawJSONwithParsedJSON(t, res, res.Inspect().Response)
						}
					}
				})
			}
		})
	}
}
