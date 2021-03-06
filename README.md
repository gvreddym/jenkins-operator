# Jenkins Operator

[![Version](https://img.shields.io/badge/version-v0.0.3-brightgreen.svg)](https://github.com/VirtusLab/jenkins-operator/releases/tag/v0.0.3)
[![Build Status](https://travis-ci.org/VirtusLab/jenkins-operator.svg?branch=master)](https://travis-ci.org/VirtusLab/jenkins-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/jenkins-operator "Go Report Card")](https://goreportcard.com/report/github.com/VirtusLab/jenkins-operator)
[![Docker Pulls](https://img.shields.io/docker/pulls/virtuslab/jenkins-operator.svg)](https://hub.docker.com/r/virtuslab/jenkins-operator/tags)

## What's Jenkins Operator?

Jenkins operator is a Kubernetes native operator which fully manages Jenkins on Kubernetes.
It was built with immutability and declarative configuration as code in mind.

Out of the box it provides:
- integration with Kubernetes
- pipelines as code
- extensibility via groovy scripts or configuration as code plugin
- security and hardening

## Problem statement and goals

The main reason why we decided to write the **jenkins-operator** is the fact that we faced a lot of problems with standard Jenkins deployment.
We want to make Jenkins more robust, suitable for dynamic and multi-tenant environments. 

Some of the problems we want to solve:
- volumes handling (AWS EBS volume attach/detach issue when using PVC)
- installing plugins with incompatible versions or security vulnerabilities
- better configuration as code
- lack of end to end tests
- handle graceful shutdown properly
- security and hardening out of the box
- orphaned jobs with no jnlp connection
- make errors more visible for end users
- backup and restore for jobs history

## Documentation

1. [Installation][installation]
2. [Getting Started][getting_started]
3. [How it works][how_it_works]
4. [Security][security]
5. [Developer Guide][developer_guide]

## Contribution

Feel free to file [issues](https://github.com/VirtusLab/jenkins-operator/issues) or [pull requests](https://github.com/VirtusLab/jenkins-operator/pulls).    

## TODO

Common:
* simple API for generating Kubernetes events using one common format
* code clean up and more tests

Base configuration:
* TLS/SSL configuration

User configuration:
* backup and restore for Jenkins jobs running as standalone job (AWS, GCP, Azure)
* trigger backup job before pod deletion using preStop k8s hooks
* verify Jenkins configuration events

[installation]:docs/installation.md
[getting_started]:docs/getting-started.md
[how_it_works]:docs/how-it-works.md
[security]:docs/security.md
[developer_guide]:docs/developer-guide.md
