FROM public.ecr.aws/spacelift/runner-terraform:latest as spacelift

ARG TERRAFORM_VERSION=1.7.5

WORKDIR /tmp

RUN curl -O -L https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && chmod +x terraform

# Build the final image
FROM alpine:3.19

RUN apk -U upgrade && apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    git \
    jq \
    openssh \
    openssh-keygen \
    tzdata

COPY --from=spacelift /bin/infracost /bin/infracost
COPY --from=spacelift /tmp/terraform /bin/terraform

RUN echo "hosts: files dns" > /etc/nsswitch.conf && \
    adduser --disabled-password --no-create-home --uid=1983 spacelift

USER spacelift