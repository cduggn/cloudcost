# cloudcost

A `Go` command line tool which supports exploring AWS costs using the AWS Cost Explorer SDK. 

**Note**
The Cost Explorer API can access data for the last 12 months. This tool will only show data for the last 12 months.


## Prerequisites

The SDK uses the AWS credential chain to find AWS credentials. The SDK looks for credentials in the following order: environment variables, shared credentials file, and EC2 instance profile. For more information, see [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

## Installation

See the releases page for the latest version.

## Commands

**Note**
AWS billing information is updated up to three times daily. Querying the Cost Explorer is charged per paginated request. The default page size is 20. The default page size can be changed using the `--page-size` flag.

Cost Explorer supports the following commands:

### `get`

The `get` command returns the results of a Cost Explorer query. The `get` command supports the following sub flags:

`aws` - The AWS service to query. This is currently the only cloud service provider that can be queried.

The `aws` command provides the following sub flags:

Arguments:
   - --start-date string   The start date of the time period. The default is the current month.
   - --end-date string     The end date of the time period. The default is the current month.
   - --filter-by string    The filter to apply to the cost. The default is no filter. Used when the --group-by-tage flag is set.
   - --granularity string  The granularity of the cost. The default is DAILY.
   - --group-by-dimension string   The dimension to group the cost by. The default is [ SERVICE,USAGE_TYPE].
   - --group-by-tag string         The tag to group the cost by. The default is no grouping.

Basic usage:

The minimum required command is `get aws`. This will return the cost for the past 30 days. The default granularity is DAILY. The default group by dimension is [ SERVICE,USAGE_TYPE]. The default group by tag is no grouping.

    $ cloudcost get aws

Command returns the cost for the past 30 days grouped by the tag `ApplicationName` and the dimension `SERVICE`.
    
    $ cloudcost get aws --group-by-tag ApplicationName --group-by-dimension SERVICE

Command returns the cost for the past 30 days grouped by the tag `ApplicationName` and the dimension `SERVICE` and filter by the tag `ApplicationName` and the value `myapp`.
    
    $ cloudcost get aws --group-by-tag ApplicationName --filter-by myapp --group-by-dimension SERVICE

Command groups the cost by the dimension LINKED_ACCOUNT and filter by the tag `ApplicationName` and the value `myapp`.
    
    $ cloudcost get aws --group-by-dimension LINKED_ACCOUNT --group-by-tag ApplicationName--filter-by myapp

Command returns the cost for the past x days based on the provided start date. Refunds and credits are not filtered. UNBLENDED_COST cost is returned.

    $ cloudcost get aws  --group-by-tag ApplicationName --group-by-dimension SERVICE -r UNBLENDED_COST -g MONTHLY -s "2022-10-01"

Command returns the cost for the past x days based on the provided start date. Refunds and credits are filtered . UNBLENDED_COST costs are returned.

    $ cloudcost get aws  --group-by-tag ApplicationName --group-by-dimension SERVICE -r UNBLENDED_COST -g MONTHLY -s "2022-10-01" -c

Command returns the cost for the past x days based on the provided start date and groups by cost allocation tag filtered by specific value. Refunds and credits are filtered.

    $ cloudcost get aws  --group-by-tag ApplicationName --group-by-dimension SERVICE -r UNBLENDED_COST -g MONTHLY -s "2022-10-01" -c -f myapp

Dimension values include the following: AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, USAGE_TYPE_GROUP, RECORD_TYPE, and OPERATING_SYSTEM. For more information, see [Grouping and Filtering](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/billing-reports-costexplorer.html#ce-grouping-filtering).