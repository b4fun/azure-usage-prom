package usagelister

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	azureusageprom "github.com/b4fun/azure-usage-prom"
)

// AzureResourceUsageLister implements azure resource usage lister.
type AzureResourceUsageLister struct {
	cloud      azure.Environment
	authorizer autorest.Authorizer
}

// New creates an Azure usage lister.
func New(
	cloud azure.Environment,
	authorizer autorest.Authorizer,
) azureusageprom.AzureUsageLister {
	return &AzureResourceUsageLister{
		cloud:      cloud,
		authorizer: authorizer,
	}
}

var _ azureusageprom.AzureUsageLister = (*AzureResourceUsageLister)(nil)

// ListByResourceProvider - implements azureusageprom.AzureUsageLister.
func (l *AzureResourceUsageLister) ListByResourceProvider(
	ctx context.Context,
	target azureusageprom.QueryTarget,
) (azureusageprom.AzureUsageList, error) {
	// TODO(hbc): cache usage

	req, err := autorest.Prepare(
		&http.Request{},
		l.authorizer.WithAuthorization(),
		autorest.WithBaseURL(l.cloud.ResourceManagerEndpoint),
		autorest.AsGet(),
		autorest.WithPathParameters(
			"/subscriptions/{subscriptionID}/providers/{resourceProvider}/locations/{location}/usages",
			map[string]interface{}{
				"subscriptionID":   target.SubscriptionID,
				"resourceProvider": target.ResourceProvider,
				"location":         target.Location,
			},
		),
		autorest.WithQueryParameters(map[string]interface{}{
			"api-version": target.ResourceProvider.APIVersion(),
		}),
	)
	if err != nil {
		return azureusageprom.AzureUsageList{}, fmt.Errorf("prepare request: %w", err)
	}

	resp, err := autorest.Send(req)
	if err != nil {
		return azureusageprom.AzureUsageList{}, fmt.Errorf("send request: %w", err)
	}

	var rv azureusageprom.AzureUsageList
	err = autorest.Respond(
		resp,
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(&rv),
		autorest.ByClosing(),
	)
	if err != nil {
		return azureusageprom.AzureUsageList{}, fmt.Errorf("respond: %w", err)
	}

	return rv, nil
}
