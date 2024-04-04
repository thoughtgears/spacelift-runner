FROM public.ecr.aws/spacelift/runner-terraform:latest as spacelift

ARG TERRAFORM_VERSION=1.7.5
ARG TFLINT_VERSION=0.50.3
ARG SNYK_VERSION=1.666.0

WORKDIR /tmp

RUN curl -O -L https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && chmod +x terraform

RUN curl -O -L https://github.com/terraform-linters/tflint/releases/download/v${TFLINT_VERSION}/tflint_linux_amd64.zip \
  && unzip tflint_linux_amd64.zip \
  && chmod +x tflint

RUN curl -O -L https://static.snyk.io/cli/v${SNYK_VERSION}/snyk-alpine \
  && chmod +x snyk-alpine

# Build the final image
FROM alpine:3.19

RUN apk -U upgrade && apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    git \
    jq \
    tzdata \
    libstdc++

COPY --from=spacelift /tmp/terraform /bin/terraform
COPY --from=spacelift /tmp/tflint /bin/tflint
COPY --from=spacelift /tmp/snyk-alpine /bin/snyk
COPY builds/huston-linux-amd64 /bin/huston

RUN echo "hosts: files dns" > /etc/nsswitch.conf && \
    adduser --disabled-password --uid=1983 spacelift && \
    ln -s /bin/terraform /usr/local/bin/terraform

USER spacelift
