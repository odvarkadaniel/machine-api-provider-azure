/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resourceskus

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-03-01/compute"
	"github.com/openshift/machine-api-provider-azure/pkg/cloud/azure"
	"github.com/openshift/machine-api-provider-azure/pkg/cloud/azure/actuators"
	"github.com/pkg/errors"
)

// Client wraps go-sdk.
type Client interface {
	List(context.Context, string) ([]compute.ResourceSku, error)
}

// AzureClient contains the Azure go-sdk Client.
type AzureClient struct {
	skus compute.ResourceSkusClient
}

var _ Client = &AzureClient{}

// NewClient creates a new Resource SKUs client from subscription ID.
func NewClient(azureClients actuators.AzureClients) *AzureClient {
	return &AzureClient{
		skus: newResourceSkusClient(azureClients),
	}
}

// newResourceSkusClient creates a new Resource SKUs client from subscription ID.
func newResourceSkusClient(azureClients actuators.AzureClients) compute.ResourceSkusClient {
	resourceSkusClient := compute.NewResourceSkusClientWithBaseURI(azureClients.ResourceManagerEndpoint, azureClients.SubscriptionID)
	resourceSkusClient.Authorizer = azureClients.Authorizer
	resourceSkusClient.AddToUserAgent(azure.UserAgent)
	return resourceSkusClient
}

// List returns all Resource SKUs available to the subscription.
func (ac *AzureClient) List(ctx context.Context, filter string) ([]compute.ResourceSku, error) {
	iter, err := ac.skus.ListComplete(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "could not list resource skus")
	}

	var skus []compute.ResourceSku
	for iter.NotDone() {
		skus = append(skus, iter.Value())
		if err := iter.NextWithContext(ctx); err != nil {
			return skus, errors.Wrap(err, "could not iterate resource skus")
		}
	}

	return skus, nil
}
