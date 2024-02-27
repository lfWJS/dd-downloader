package datadog

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/girishg4t/dd-downloader/pkg/model"
)

func GetDataDogLogs(filter model.DataDogFilter, cur *string, limit int32) datadogV2.LogsListResponse {
	log.Printf("query=>%s, from %d, to %d\n", filter.Query, filter.From, filter.To)
	unixTo := time.Unix(int64(filter.To), 0)
	unixFrom := time.Unix(int64(filter.From), 0)
	body := datadogV2.LogsListRequest{
		Filter: &datadogV2.LogsQueryFilter{
			Query: datadog.PtrString(filter.Query),
			Indexes: []string{
				"main",
			},
			From: datadog.PtrString(unixFrom.GoString()),
			To:   datadog.PtrString(unixTo.GoString()),
		},
		Sort: datadogV2.LOGSSORT_TIMESTAMP_ASCENDING.Ptr(),
		Page: &datadogV2.LogsListRequestPage{
			Limit:  datadog.PtrInt32(limit),
			Cursor: cur,
		},
	}

	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)

	api := datadogV2.NewLogsApi(apiClient)
	resp, r, err := api.ListLogs(ctx, *datadogV2.NewListLogsOptionalParameters().WithBody(body))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LogsApi.ListLogs`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return datadogV2.LogsListResponse{
			Data:                 []datadogV2.Log{},
			Links:                &datadogV2.LogsListResponseLinks{},
			Meta:                 &datadogV2.LogsResponseMetadata{},
			UnparsedObject:       map[string]interface{}{},
			AdditionalProperties: map[string]interface{}{},
		}
	}

	return resp
}
