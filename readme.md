# Connecting nodes

Copyright 2017 [Arnaud Poret](https://github.com/arnaudporet)

This work is licensed under the [BSD 2-Clause License](https://raw.githubusercontent.com/arnaudporet/konect/master/BSD_2_Clause_License.txt).

## konect

[konect](https://github.com/arnaudporet/konect) is a basic algorithm implemented in [Go](https://golang.org) for connecting nodes according to a reference network: the master network.

To do so, konect performs random walks in the master network from the nodes given as sources to reach the nodes given as targets.

Typical usage consists in extracting, from the master network, the paths connecting a couple of nodes belonging to it, thus returning a subnetwork of the master network.

konect handles networks encoded in the sif file format (see below): the provided master network must be encoded in the sif file format.

Together with the master network encoded in a sif file, konect requires the nodes to connect be listed in txt files (see below).

## The sif file format

In a sif file encoding a network, each line encodes an edge of the network as follow:
* `source \t interaction \t target`

Note that the separators are tabulations `\t`: the sif file format is the tab-separated values format (tsv) with exactly 3 columns.

For example, the edge representing the activation of RAF1 by HRAS is a line of a sif file encoded as follow:
* `HRAS \t activation \t RAF1`

## Usage

In a terminal emulator:
1. `go build konect.go`
2. `./konect networkFile sourceFile targetFile maxStep maxWalk selfConnect shortest`

or simply
* `go run konect.go networkFile sourceFile targetFile maxStep maxWalk selfConnect shortest`

Note that `go run` build konect each time before running it.

The Go package can have different names depending on your OS/Linux distribution. For example, with [Ubuntu](https://www.ubuntu.com/), the Go package is named golang: running a Go file with Ubuntu might be `golang-go run yourfile.go` instead of `go run yourfile.go` with [Arch Linux](https://www.archlinux.org).

Arguments:
* `networkFile`: the master network encoded in a sif file (see above)
* `sourceFile`: the source nodes listed in a txt file (one node per line)
* `targetFile`: the target nodes listed in a txt file (one node per line)
* `maxStep`: the maximum number of steps performed during a random walk starting from a source node in an attempt to reach a target node
* `maxWalk`: the maximum number of random walks performed in the master network to find paths from a source node to a target node
* `selfConnect` (`1` or `0`): allow `1` or not `0` konect to find paths connecting a node to itself if it belongs to both the source and target nodes
* `shortest` (`1` or `0`): among the found connecting paths, select only `1` or not only `0` the shortest

The returned file is a sif file encoding a subnetwork of the master network connecting the source nodes to the target nodes.

The lists of source and target nodes can overlap, or even be identical. If identical, `selfConnect` must be `1`.

## Cautions

* konect does not handle multigraphs (i.e. networks with nodes connected by more than one edge)
* the master network must be provided as a sif file (see above)
* in the txt files containing the node lists (see above): one node per line
* since konect uses random walks:
    * the results can be different between identical runs
    * returning all the possible connecting paths is not guaranteed
* setting `selfConnect` at `1` and/or `shortest` at `0` can greatly increase the size of the returned network
* increasing `maxWalk`:
    * increases the robustness of the results
    * but also increases the computational time

## The examples

All the master sif used in these examples are adapted from pathways coming from [KEGG Pathway](http://www.genome.jp/kegg/pathway.html).

* example 1: typical usage
    * `./konect MAPK_signaling_pathway.sif sources.txt targets.txt 100 1000000 0 1`
    * networkFile: the MAPK signaling pathway (1194 edges)
    * sourceFile: contains the nodes EGFR and IL1R1
    * targetFile: contains the nodes MAPK1 and MAPK14
    * maxStep: 100
    * maxWalk: 1000000
    * selfConnect: 0
    * shortest: 1
    * result: konected.sif (35 edges), also in svg for visualization

* example 2: not only the shortest paths
    * `./konect.go Toll_like_receptor_signaling_pathway.sif sources.txt targets.txt 100 1000000 0 0`
    * networkFile: the Toll-like receptor signaling pathway (392 edges)
    * sourceFile: contains the node MYD88
    * targetFile: contains the node TRAF6
    * maxStep: 100
    * maxWalk: 1000000
    * selfConnect: 0
    * shortest: 0
    * result: konected.sif (22 edges), also in svg for visualization

* example 3: allow self connections
    * `./konect.go cell_cycle.sif nodes.txt nodes.txt 100 1000000 1 0`
    * networkFile: the cell cycle (313 edges)
    * sourceFile: contains the node CCND1
    * targetFile: contains the node CCND1 (targetFile=sourceFile)
    * maxStep: 100
    * maxWalk: 1000000
    * selfConnect: 1
    * shortest: 0
    * result: konected.sif (9 edges), also in svg for visualization

## Forthcoming

* improving the code

## Go

How to get Go: https://golang.org/doc/install

Most [Linux distributions](https://distrowatch.com) provide Go in their official repositories. For example:
* go (Arch Linux)
* golang (Ubuntu)

## References

konect is inspired from [MCWalk](https://bitbucket.org/akittas/biosubg) [1].

1. Aristotelis Kittas, Aurelien Delobelle, Sabrina Schmitt, Kai Breuhahn, Carito Guziolowski, Niels Grabe (2016) Directed random walks and constraint programming reveal active pathways in hepatocyte growth factor signaling. FEBS journal 283(2):350-360.
