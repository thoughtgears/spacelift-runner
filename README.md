# Spacelift Docker Runner

This repository contains a Docker image that can be used to run Spacelift with terraform after > 1.5.7 version since
Spacelift does not natively support any versions of terraform after 1.5.7. It will also contain the necessary tools to
run different context checks and tests.

### Tools

- tflint
- snyk
- huston

## Usage 

To use this image, you need to have a Spacelift account and a Spacelift workspace. You can need to reference the
image in your configuration file or in your stack configuration. 

`.spacelift/config.yaml`:
```yaml
version: "1"

stack_defaults:
    runner_image: ghcr.io/thoughtgears/spacelift-runner:latest
```

## Huston

Huston is an application that can run in the different contexts and create various tasks. Currently supported tasks are:

- backup state
- check latest module versions in terraform code

### Implemented tasks

- backup state: backs up your state to a bucket of your choice in GCP.
