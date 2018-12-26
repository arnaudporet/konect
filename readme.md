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

The returned file `konected.sif` is a SIF file encoding all the paths connecting the source nodes to the target nodes in the network.

The returned file `konected-shortest.sif` is a SIF file encoding only the shortest connecting paths contained in `konected.sif`.

## Cautions

* konect does not handle multi-edges (i.e. two or more edges having the same source node and the same target node)
* note that duplicated edges are multi-edges
* the network must be provided as a SIF file (see above)
* in the files containing the node lists (see above): one node per line

## Examples

All the networks used in these examples are adapted from pathways coming from [KEGG Pathway](https://www.genome.jp/kegg/pathway.html).

* ErbB signaling pathway
    * `konect ErbB_signaling_pathway.sif sources.txt targets.txt`
    * networkFile: the ErbB signaling pathway (239 edges)
    * sourceFile: contains the nodes EGFR (i.e. ERBB1), ERBB2, ERBB3 and ERBB4
    * targetFile: contains the node MTOR
    * results:
        * konected.sif (83 edges), also in svg for visualization
        * konected-shortest.sif (50 edges), also in svg for visualization

* Insulin signaling pathway
    * `konect Insulin_signaling_pathway.sif sources.txt targets.txt`
    * networkFile: the insulin signaling pathway (407 edges)
    * sourceFile: contains the node INSR
    * targetFile: contains the nodes GSK3B and MAPK1
    * results:
        * konected.sif (69 edges), also in svg for visualization
        * konected-shortest.sif (69 edges), also in svg for visualization

* Cell cycle
    * `konect Cell_cycle.sif nodes.txt nodes.txt`
    * networkFile: the cell cycle (650 edges)
    * sourceFile: contains the node RB1
    * targetFile=sourceFile: for getting the paths connecting RB1 to itself
    * results:
        * konected.sif (84 edges), also in svg for visualization
        * konected-shortest.sif (22 edges), also in svg for visualization

* Cell survival
    * to illustrate the advantage of also computing the shortest connecting paths, this example is voluntarily bigger
    * it is made of the following KEGG pathways: Apoptosis, Cell cycle, p53 signaling pathway, ErbB signaling pathway, TNF signaling pathway, TGF-beta signaling pathway, FoxO signaling pathway, Calcium signaling pathway, MAPK signaling pathway, PI3K-Akt signaling pathway and NF-kappa B signaling pathway
    * these pathways are involved in the cell growth/cell death balance
    * `konect Cell_survival.sif nodes.txt nodes.txt`
    * networkFile: some cell survival signaling pathways (11147 edges)
    * sourceFile: contains the nodes CASP3 (cell death effector), PIK3CA (involved in growth promoting signaling pathways) and TP53 (tumor suppressor)
    * targetFile=sourceFile: to see how these nodes interact with each other
    * results:
        * konected.sif (819 edges), also in svg for a quite challenging visualization
        * konected-shortest.sif (84 edges), also in svg for an easier visualization but only of the shortest connecting paths

## Forthcoming

## Go

Most [Linux distributions](https://distrowatch.com) provide Go in their official repositories. For example:
* go (Arch Linux)
* golang (Ubuntu)

Otherwise see https://golang.org/doc/install
