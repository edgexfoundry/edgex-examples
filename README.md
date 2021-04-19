# Official EdgeX Foundry Examples
[![Go Report Card](https://goreportcard.com/badge/github.com/edgexfoundry/edgex-examples)](https://goreportcard.com/report/github.com/edgexfoundry/edgex-examples) [![GitHub Pull Requests](https://img.shields.io/github/issues-pr-raw/edgexfoundry/edgex-examples)](https://github.com/edgexfoundry/edgex-examples/pulls) [![GitHub Contributors](https://img.shields.io/github/contributors/edgexfoundry/edgex-examples)](https://github.com/edgexfoundry/edgex-examples/contributors) [![GitHub Committers](https://img.shields.io/badge/team-committers-green)](https://github.com/orgs/edgexfoundry/teams/edgex-examples-committers/members) [![GitHub Commit Activity](https://img.shields.io/github/commit-activity/m/edgexfoundry/edgex-examples)](https://github.com/edgexfoundry/edgex-examples/commits)


The purpose of this repository is to provide centralized example and sample code for EdgeX adopters and users.  The folders help organize the examples to make it easier to find samples associated to various parts of EdgeX.

## Example "Rules"

- The TSC will review these examples with each release.
- The project does maintain the code in this repository (please file an issue for any broken code).
- An example is dropped when it is no longer used, maintained, or in sync with the current release as deemed by the TSC.
- Approval of the TSC is required for new folders in this repository.
- SDKs located in their respective repositories can also contain a small example or sample code folder to show developers how to use the SDK or SDK feature.
    - SDK examples include how to create a simple app service or device service.
    - SDK examples may also provide a how-to-guide on particular features (ex: how to implement an app function or how to handle automatic / dynamic provisioning).
    - SDK examples should remain small so as not to bloat the services created from the SDK.
    - Having the code with the SDK encourages the SDK maintainers to keep these up to date - critical for adoption.
    - The project has the goal to have automated testing in the future against examples in the SDKs (which is easier to do when the code is in the same repository as the SDK).
- Holding (http://github.com/edgexfoundry-holding) will revert to a “staging” area for code. Put examples in holding while under review.  
    - Like other project code, the goal of example code in holding would be to have it eventually approved and moved to this repository.
- EdgeX developers and writers should minimize examples in the docs and reference the code in this repo instead.
- TSC members are automatically committers to the repository. 
- TSC members are automatically committers to the repository. 
- TSC members can request committer rights for another contributor [see committer approval process](https://wiki.edgexfoundry.org/pages/viewpage.action?pageId=21823860).
    - Like other repositories, it requires the approval of 2 other TSC members.
- Examples will be tagged, as the project moves forward.
    - The project will update this READMe for how to handle the examples code with regard to releases
- The example code in this repository is not tested as part of the CI/CD process.
    - Examples are tested by developers but are provided without regular (nightly) testing as with the rest of EdgeX.
- Modules in these examples are exempt from 3rd party vetting per https://wiki.edgexfoundry.org/pages/viewpage.action?pageId=46760301.

## Example Folder Organization
- application-services – example app services to get data to cloud or enterprise systems
- analytics – examples around rules engine, or other analytics package integration or use
- device-services – portions of a device service or a device service that would not typically be used in production.
- deployment – Kubernetes, swarm, etc. deployment and orchestration examples
- hackathon-projects – projects collected during official EdgeX hackathons (these may be more temporal as they will likely not be maintained)
- miscellaneous – examples/samples that don’t belong in the other folders. Example scripts, makefiles, etc.
- platform-solutions – examples of getting EdgeX (or a portion of EdgeX) up and running on specific platforms like Raspberry Pi, Arduino
- product-integration – examples outlining how to integrate with or use EdgeX with other products whether proprietary or open (examples might include ObjectBox, Redis Streams, Foghorn, etc.) 