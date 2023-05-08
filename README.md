# Akita Demo

This repository hosts the code for demoing the Akita Agent. Before running the demo, please ensure your system meets the following requirements

## Prerequisites
1. **Docker**: Please ensure that you have Docker installed on your machine. You can download it from the official [Docker Website](https://docs.docker.com/engine/install/)
2. **Akita User Account**: You must have an active user account with Akita. If you don't have one yet, please [sign up here](https://app.akita.software/login?sign_up)

## Getting started
### Step 1: Obtain Akita API Credentials and Project Name
Go through initial onboarding in [Akita's web dashboard](https://app.akita.software) and store your API credentials and created project name for later use

### Step 2: Clone the repository
Clone this repository to your local machine:
```
git clone https://github.com/akitasoftware/akita-demo.git
```
Then `cd` into the demo's directory:
```
cd akita-demo
```
### Step 3: Run the demo
Run the following commands to start the demo:
```
./run.sh
```

To stop the demo, run `./stop.sh`. To restart the demo, run `./restart.sh`

## Getting involved

* Please file bugs as issues to this repository.
* We welcome contributions! If you want to make changes please see our [contributing guide](CONTRIBUTING.md).
* We're always happy to answer any questions about the Docker extension, or about how you
  can contribute. Email us at `opensource [at] akitasoftware [dot] com` or
  [request to join our Slack](https://docs.google.com/forms/d/e/1FAIpQLSfF-Mf4Li_DqysCHy042IBfvtpUDHGYrV6DOHZlJcQV8OIlAA/viewform?usp=sf_link)!

## What is Akita?

Drop-in API monitoring, no code changes necessary.

Built for busy developer teams who don't have time to become experts in monitoring and observability, Akita is the
fastest, easiest way to see and monitor your API endpoints which makes it possible to quickly track API endpoints and
their usage in real time.

* **See API endpoints.** Automatically get a searchable map of your API endpoints in use. Explore by latency, errors,
  and usage. Export as OpenAPI specs.
* **Get drop-in API monitoring.** Get a drop-in view of volume, latency, and errors, updated in near real-time. Set
  per-endpoint alerts.
* **Quickly understand the impact of changes.** Keep track of the endpoints you care about and identify how new
  deployments impact your endpoints.

Simply drop Akita into your system to understand your system behavior, without having to instrument code or build your
own dashboards.

[Create an account in the Akita App](https://app.akita.software/login?sign_up) to get started!

## Related links
* [Akita blog](https://www.akitasoftware.com/blog)
* [Akita docs](https://docs.akita.software/)



