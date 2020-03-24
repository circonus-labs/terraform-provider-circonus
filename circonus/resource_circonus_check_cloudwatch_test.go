package circonus

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCirconusCheckCloudWatch_basic(t *testing.T) {
	{
		// fudge the aws credentials required for the purposes of testing the creation...
		envVar := "AWS_ACCESS_KEY_ID"
		if os.Getenv(envVar) == "" {
			os.Setenv(envVar, "test_key")
		}
		envVar = "AWS_SECRET_ACCESS_KEY"
		if os.Getenv(envVar) == "" {
			os.Setenv(envVar, "test_secret")
		}
	}
	checkName := fmt.Sprintf("Terraform test: RDS Metrics via CloudWatch - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckCloudWatchConfigFmt, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "collector.2388330941.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "cloudwatch.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "cloudwatch.46539847.dimmensions.%", "1"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "cloudwatch.46539847.dimmensions.DBInstanceIdentifier", "atlas-production"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "cloudwatch.46539847.metric.#", "17"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "cloudwatch.46539847.namespace", "AWS/RDS"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "cloudwatch.46539847.version", "2010-08-01"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "cloudwatch.46539847.url", "https://monitoring.us-east-1.amazonaws.com"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "notes", "Collect all the things exposed"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.#", "17"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.name", "ReadLatency"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.unit", "seconds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.name", "TransactionLogsGeneration"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.unit", ""),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.name", "WriteIOPS"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.unit", "iops"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.name", "FreeStorageSpace"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.unit", ""),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.name", "WriteLatency"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.unit", "seconds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.name", "DatabaseConnections"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.unit", "connections"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.name", "FreeableMemory"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.unit", "bytes"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.name", "MaximumUsedTransactionIDs"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.unit", ""),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.name", "ReadThroughput"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.unit", "bytes"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.name", "ReadIOPS"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.unit", "iops"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.name", "NetworkReceiveThroughput"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.unit", "bytes"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.name", "TransactionLogsDiskUsage"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.unit", "bytes"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.name", "CPUUtilization"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.unit", "%"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.name", "SwapUsage"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.unit", "bytes"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.name", "NetworkTransmitThroughput"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.unit", "bytes"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.name", "DiskQueueDepth"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.unit", ""),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.name", "WriteThroughput"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.unit", "bytes"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "tags.2964981562", "app:postgresql"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "tags.1313458811", "app:rds"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "tags.4259413593", "source:cloudwatch"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "target", "atlas-production.us-east-1.rds._aws"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "type", "cloudwatch"),
				),
			},
		},
	})
}

const testAccCirconusCheckCloudWatchConfigFmt = `
variable "cloudwatch_rds_tags" {
  type = "list"
  default = [
    "app:postgresql",
    "app:rds",
    "lifecycle:unittest",
    "source:cloudwatch",
  ]
}

resource "circonus_check" "rds_metrics" {
  active = true
  name = "%s"
  notes = "Collect all the things exposed"
  period = "60s"

  collector {
    id = "/broker/1"
  }

  target = "atlas-production.us-east-1.rds._aws"

  cloudwatch {
    dimmensions = {
      DBInstanceIdentifier = "atlas-production",
    }

    metric = [
      "CPUUtilization",
      "DatabaseConnections",
      "DiskQueueDepth",
      "FreeStorageSpace",
      "FreeableMemory",
      "MaximumUsedTransactionIDs",
      "NetworkReceiveThroughput",
      "NetworkTransmitThroughput",
      "ReadIOPS",
      "ReadLatency",
      "ReadThroughput",
      "SwapUsage",
      "TransactionLogsDiskUsage",
      "TransactionLogsGeneration",
      "WriteIOPS",
      "WriteLatency",
      "WriteThroughput",
    ]

    namespace = "AWS/RDS"
    url = "https://monitoring.us-east-1.amazonaws.com"
  }

  metric {
    name = "CPUUtilization"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "%%"
  }

  metric {
    name = "DatabaseConnections"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "connections"
  }

  metric {
    name = "DiskQueueDepth"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
  }

  metric {
    name = "FreeStorageSpace"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
  }

  metric {
    name = "FreeableMemory"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "bytes"
  }

  metric {
    name = "MaximumUsedTransactionIDs"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
  }

  metric {
    name = "NetworkReceiveThroughput"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "bytes"
  }

  metric {
    name = "NetworkTransmitThroughput"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "bytes"
  }

  metric {
    name = "ReadIOPS"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "iops"
  }

  metric {
    name = "ReadLatency"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "seconds"
  }

  metric {
    name = "ReadThroughput"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "bytes"
  }

  metric {
    name = "SwapUsage"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "bytes"
  }

  metric {
    name = "TransactionLogsDiskUsage"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "bytes"
  }

  metric {
    name = "TransactionLogsGeneration"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
  }

  metric {
    name = "WriteIOPS"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "iops"
  }

  metric {
    name = "WriteLatency"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "seconds"
  }

  metric {
    name = "WriteThroughput"
    tags = "${var.cloudwatch_rds_tags}"
    type = "numeric"
    unit = "bytes"
  }

  tags = "${var.cloudwatch_rds_tags}"
}
`
