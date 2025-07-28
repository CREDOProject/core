FROM ubuntu:24.04 AS base

ENV \
	DEBIAN_FRONTEND=noninteractive \
	LANG="C.UTF-8"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

FROM base AS final

RUN set -x \
	&& apt-get update -yq --no-install-recommends \
	&& apt-get install -yq --no-install-recommends \
	build-essential \
	gfortran \
	ca-certificates

COPY ./credo /usr/bin/credo

RUN set -x \
	&& chmod +x /usr/bin/credo \
	&& mkdir -p /workdir

WORKDIR /workdir
