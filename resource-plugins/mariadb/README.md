# HotROD MariaDB Plugin

This is a sample Signadot Resource Plugin to provision an empty MariaDB instance
for each Signadot Sandbox that requests one.

This directory contains the sources for the `signadot/hotrod-mariadb-plugin`
container image, which is used by the `resource-plugin.yaml` file.

It is not necessary to build this image yourself in order to try the demo.
However, this could be a starting point if you want to build your own plugin.

## Overview

- `Dockerfile` - This builds an image containing Helm, kubectl, and the two bash
  scripts in this directory, `provision` and `deprovision`.
- `resource-plugin.yaml` - This contains the Kubernetes manifests to install the
  plugin in a cluster that already has Signadot Operator installed. Once the
  plugin is installed, Signadot Sandboxes in that cluster can use the plugin by
  specifying `"plugin": "hotrod-mariadb"` in an entry in the `resources` list.
- `provision` - This is a bash script that runs Helm to install a MariaDB chart.
  It generates a unique Helm release name based on the Sandbox ID so that
  multiple Sandboxes can deploy DB instances without interference. Once the
  chart is installed, this script returns connection info to Signadot so that
  it can be injected as environment variables in forked workloads.
- `deprovision` - This is a bash script that runs Helm to uninstall a given
  Sandbox's Helm release.
