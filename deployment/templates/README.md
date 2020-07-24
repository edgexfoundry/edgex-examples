# Cloud Deployment Templates

## Quick Start

If you want to see how all of EdgeX works - you can leverage your own Azure or Amazon AWS accounts and deploy EdgeX to the cloud.

### Azure
This template leverages Azure Container Instances and will deploy a single group called "edgex-example" with 12 services deployed with 2.4 vCPUs allocated (0.2 vCPUs per service) and 6GB of RAM allocated (0.5 GB per service) with an estimated cost of $0.14904 / hour or $3.57696 / day. 
```
1 container groups * 3600 seconds * 2.4 vCPU * $0.0000135 per vCPU-s  = ~$0.11664

1 container groups * 3600 seconds * 6 GB * $0.0000015 per GB-s  = $0.0324

memory($0.0324) + cpu($0.11664) = $0.14904 / hour
= $3.57696 / day
```
[![Deploy to Azure](https://aka.ms/deploytoazurebutton)](https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fraw.githubusercontent.com%2Fedgexfoundry-holding%2Fapp-service-examples%2Fmaster%2Ftemplates%2Fazuredeploy.json)

### AWS
The sample stack can also be launched in AWS Fargate using the quick-create button below.

This Cloudformation template deploys twelve containers split in five Task Definitions, consuming a total of three vCPUs, six GBs of RAM and four Load Balancers.
The template will require you to pass a VPC, Private and Public subnets as parameters at launch time. If you need an example on how to build a VPC, AWS provides this [VPC Quick Start](https://aws.amazon.com/quickstart/architecture/vpc/).

```
Total cost per day, before traffic, would be about $5.71464

3 vCPUs * $0.04048/hr + 6 GB * $0.004445/hr + 4 NLBs + $0.0225/hr = $0.23811/hr
```
[![Deploy to AWS](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks/quickcreate?templateUrl=https%3A%2F%2Fraw.githubusercontent.com%2Fedgexfoundry-holding%2Fapp-service-examples%2Fmaster%2Ftemplates%2Faws-fargate.yaml&stackName=edgex-sample)

