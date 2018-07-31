# Connecting nodes

Copyright 2017-2018 [Arnaud Poret](https://github.com/arnaudporet)

This work is licensed under the [BSD 2-Clause License](https://raw.githubusercontent.com/arnaudporet/konect/master/license.txt).

## konect

[konect](https://github.com/arnaudporet/konect) is a tool implemented in [Go](https://golang.org) for finding paths connecting a couple of nodes in a network.

Typical use is to find in a network the paths connecting some source nodes to some target nodes.

konect handles networks encoded in the SIF file format (see below): the provided network must be encoded in the SIF file format.

Together with the network encoded in a SIF file, konect requires the nodes to connect to be listed in files (see below).

## The SIF file format

In a SIF file encoding a network, each line encodes an edge of the network as follows:
* `source \t interaction \t target`

Note that the field separator is the tabulation `\t`: the SIF file format is the tab-separated values format (TSV) with exactly 3 columns.

For example, the edge representing the activation of RAF1 by HRAS is a line of a SIF file encoded as follows:
* `HRAS \t activation \t RAF1`

## Usage

In a terminal emulator:
1. `go build konect.go`
2. `./konect networkFile sourceFile targetFile`

or simply
* `go run konect.go networkFile sourceFile targetFile`

Note that `go run` builds konect each time before running it.

The Go package can have different names depending on your operating system. For example, with [Ubuntu](https://www.ubuntu.com), the Go package is named golang. Consequently, running a Go file with Ubuntu might be `golang-go run yourfile.go` instead of `go run yourfile.go` with [Arch Linux](https://www.archlinux.org).

Arguments:
* `networkFile`: the network encoded in a SIF file (see above)
* `sourceFile`: the source nodes listed in a file (one node per line)
* `targetFile`: the target nodes listed in a file (one node per line)

The returned file is a SIF file encoding the paths connecting the source nodes to the target nodes in the network.

## Cautions

* konect does not handle multi-graphs (i.e. networks where nodes can be connected by more than one edge)
* note that if a network contains duplicated edges then it is a multi-graph
* the network must be provided as a SIF file (see above)
* in the files containing the node lists (see above): one node per line

## Examples

All the networks used in these examples are adapted from pathways coming from [KEGG Pathway](https://www.genome.jp/kegg/pathway.html).

* Cell cycle
    * `konect Cell_cycle.sif nodes.txt nodes.txt`
    * networkFile: the cell cycle (650 edges)
    * sourceFile: contains the node RB1
    * targetFile=sourceFile: for getting paths connecting RB1 to itself
    * result: konected.sif (84 edges), also in svg for visualization

* ErbB signaling pathway
    * `konect ErbB_signaling_pathway.sif sources.txt targets.txt`
    * networkFile: the ErbB signaling pathway (239 edges)
    * sourceFile: contains the nodes EGFR (i.e. ERBB1), ERBB2, ERBB3 and ERBB4
    * targetFile: contains the node MTOR
    * result: konected.sif (83 edges), also in svg for visualization

* Insulin signaling pathway
    * `konect Insulin_signaling_pathway.sif sources.txt targets.txt`
    * networkFile: the insulin signaling pathway (407 edges)
    * sourceFile: contains the node INSR
    * targetFile: contains the nodes GSK3B and MAPK1
    * result: konected.sif (69 edges), also in svg for visualization

## Forthcoming

## Go

Most [Linux distributions](https://distrowatch.com) provide Go in their official repositories. For example:
* go (Arch Linux)
* golang (Ubuntu)

Otherwise see https://golang.org/doc/install
