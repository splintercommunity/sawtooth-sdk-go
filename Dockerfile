# Copyright 2017 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------------------------

# Description:
#   Builds an image to be used when developing in Go. The default CMD is to run
#   build_go.
#
# Build:
#   $ cd sawtooth-sdk-go/docker
#   $ docker build . -f sawtooth-build-go-protos -t sawtooth-build-go-protos
#
# Run:
#   $ cd sawtooth-sdk-go
#   $ docker run -v $(pwd):/project/sawtooth-sdk-go sawtooth-build-go-protos

FROM ubuntu:bionic

RUN apt-get update \
 && apt-get install gnupg -y

LABEL "install-type"="mounted"

RUN echo "deb http://repo.sawtooth.me/ubuntu/ci bionic universe" >> /etc/apt/sources.list \
 && echo 'deb http://ppa.launchpad.net/gophers/archive/ubuntu bionic main' >> /etc/apt/sources.list \
 && (apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD \
 || apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 8AA7AF1F1091A5FD) \
 && (apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 308C15A29AD198E9 \
 || apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 308C15A29AD198E9) \
 && apt-get update \
 && apt-get install -y -q \
    git \
    libssl-dev \
    libzmq3-dev \
    openssl \
    protobuf-compiler \
    python3 \
    python3-pip \
    python3-pkg-resources \
    wget \
    pkg-config \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

RUN wget https://go.dev/dl/go1.16.15.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.16.15.linux-amd64.tar.gz

RUN pip3 install grpcio grpcio-tools

ENV GOPATH=/go:/project/:/go/src/github.com/hyperledger/sawtooth-sdk-go:/go/src/github.com/hyperledger/sawtooth-sdk-go/examples/smallbank/smallbank_go/:/go/src/github.com/hyperledger/sawtooth-sdk-go/protobuf

ENV PATH=$PATH:/go/bin:/usr/local/go/bin

RUN mkdir /go

COPY . /go/src/github.com/hyperledger/sawtooth-sdk-go

WORKDIR /go/src/github.com/hyperledger/sawtooth-sdk-go/

RUN go mod download

CMD go generate
