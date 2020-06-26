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

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.name", "CPUUtilization"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.0.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.name", "DatabaseConnections"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.1.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.name", "DiskQueueDepth"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.2.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.name", "FreeStorageSpace"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.3.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.name", "FreeableMemory"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.4.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.name", "MaximumUsedTransactionIDs"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.5.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.name", "NetworkReceiveThroughput"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.6.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.name", "NetworkTransmitThroughput"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.7.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.name", "ReadIOPS"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.8.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.name", "ReadLatency"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.9.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.name", "ReadThroughput"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.10.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.name", "SwapUsage"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.11.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.name", "TransactionLogsDiskUsage"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.12.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.name", "TransactionLogsGeneration"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.13.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.name", "WriteIOPS"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.14.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.name", "WriteLatency"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.15.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.name", "WriteThroughput"),
					resource.TestCheckResourceAttr("circonus_check.rds_metrics", "metric.16.type", "numeric"),

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
    type = "numeric"
  }

  metric {
    name = "DatabaseConnections"
    type = "numeric"
  }

  metric {
    name = "DiskQueueDepth"
    type = "numeric"
  }

  metric {
    name = "FreeStorageSpace"
    type = "numeric"
  }

  metric {
    name = "FreeableMemory"
    type = "numeric"
  }

  metric {
    name = "MaximumUsedTransactionIDs"
    type = "numeric"
  }

  metric {
    name = "NetworkReceiveThroughput"
    type = "numeric"
  }

  metric {
    name = "NetworkTransmitThroughput"
    type = "numeric"
  }

  metric {
    name = "ReadIOPS"
    type = "numeric"
  }

  metric {
    name = "ReadLatency"
    type = "numeric"
  }

  metric {
    name = "ReadThroughput"
    type = "numeric"
  }

  metric {
    name = "SwapUsage"
    type = "numeric"
  }

  metric {
    name = "TransactionLogsDiskUsage"
    type = "numeric"
  }

  metric {
    name = "TransactionLogsGeneration"
    type = "numeric"
  }

  metric {
    name = "WriteIOPS"
    type = "numeric"
  }

  metric {
    name = "WriteLatency"
    type = "numeric"
  }

  metric {
    name = "WriteThroughput"
    type = "numeric"
  }

  tags = "${var.cloudwatch_rds_tags}"
}
`
