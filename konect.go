// Copyright 2017-2018 Arnaud Poret
// This work is licensed under the BSD 2-Clause License.
package main
import (
    "encoding/csv"
    "fmt"
    "os"
    "strings"
)
func main() {
    var (
        sources,targets []string
        forward,backward,intersect,shortest [][]string
        nodeSucc,nodePred map[string][]string
        edgeNames map[string]map[string]string
        edgeSucc,edgePred map[string]map[string][][]string
    )
    if len(os.Args)==4 {
        nodeSucc,nodePred,edgeSucc,edgePred,edgeNames=ReadNetwork(os.Args[1])
        if len(edgeNames)==0 {
            fmt.Println("WARNING: "+os.Args[1]+" is empty after reading")
        } else {
            sources=ReadNodes(os.Args[2],nodeSucc)
            targets=ReadNodes(os.Args[3],nodeSucc)
            if len(sources)==0 {
                fmt.Println("WARNING: "+os.Args[2]+" is empty after reading")
            } else if len(targets)==0 {
                fmt.Println("WARNING: "+os.Args[3]+" is empty after reading")
            } else {
                forward=ForwardEdges(sources,nodeSucc,edgeSucc,true)
                backward=BackwardEdges(targets,nodePred,edgePred,true)
                if len(forward)==0 {
                    fmt.Println("WARNING: sources have no forward paths")
                } else if len(backward)==0 {
                    fmt.Println("WARNING: targets have no backward paths")
                } else {
                    intersect=IntersectEdges(forward,backward)
                    if len(intersect)==0 {
                        fmt.Println("WARNING: no connecting paths found")
                    } else {
                        WriteNetwork("konected.sif",intersect,edgeNames)
                        shortest=ShortestPaths(sources,targets,intersect)
                        WriteNetwork("konected-shortest.sif",shortest,edgeNames)
                    }
                }
            }
        }
    } else if (len(os.Args)==2) && (os.Args[1]=="help") {
        fmt.Println(strings.Join([]string{
            "",
            "konect is a tool for finding paths connecting a couple of nodes in a network.",
            "",
            "Typical use is to find in a network the paths connecting some source nodes to some target nodes.",
            "",
            "konect handles networks encoded in the SIF file format.",
            "",
            "konect does not handle multi-graphs (i.e. networks where nodes can be connected by more than one edge).",
            "",
            "Note that if a network contains duplicated edges then it is a multi-graph.",
            "",
            "Usage: konect networkFile sourceFile targetFile",
            "",
            "    * networkFile: the network encoded in a SIF file",
            "",
            "    * sourceFile:  the source nodes listed in a file (one node per line)",
            "",
            "    * targetFile:  the target nodes listed in a file (one node per line)",
            "",
            "The returned file \"konected.sif\" is a SIF file encoding all the paths connecting the source nodes to the target nodes in the network.",
            "",
            "The returned file \"konected-shortest.sif\" is a SIF file encoding only the shortest connecting paths contained in \"konected.sif\".",
            "",
            "For more information see https://github.com/arnaudporet/konect",
            "",
        },"\n"))
    } else if (len(os.Args)==2) && (os.Args[1]=="license") {
        fmt.Println(strings.Join([]string{
            "",
            "Copyright 2017-2018 Arnaud Poret",
            "",
            "Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:",
            "",
            "1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.",
            "",
            "2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.",
            "",
            "THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
            "",
        },"\n"))
    } else if (len(os.Args)==2) && (os.Args[1]=="usage") {
        fmt.Println(strings.Join([]string{
            "",
            "konect networkFile sourceFile targetFile",
            "",
        },"\n"))
    } else {
        fmt.Println(strings.Join([]string{
            "ERROR: wrong number of arguments",
            "",
            "To print help:    konect help",
            "To print license: konect license",
            "To print usage:   konect usage",
            "",
            "For more information see https://github.com/arnaudporet/konect",
            "",
        },"\n"))
    }
}
func BackwardEdges(targets []string,nodePred map[string][]string,edgePred map[string]map[string][][]string,verbose bool) [][]string {
    var (
        target,npred string
        eroot,check,epred []string
        backward,newCheck,toCheck [][]string
    )
    for _,target=range targets {
        if verbose {
            fmt.Println("    backwarding "+target)
        }
        for _,npred=range nodePred[target] {
            eroot=[]string{npred,target}
            if !IsInList2(backward,eroot) {
                backward=append(backward,CopyList(eroot))
                newCheck=[][]string{CopyList(eroot)}
                for {
                    toCheck=CopyList2(newCheck)
                    newCheck=[][]string{}
                    for _,check=range toCheck {
                        for _,epred=range edgePred[check[0]][check[1]] {
                            if !IsInList2(backward,epred) {
                                backward=append(backward,CopyList(epred))
                                newCheck=append(newCheck,CopyList(epred))
                            }
                        }
                    }
                    if len(newCheck)==0 {
                        break
                    }
                }
            }
        }
    }
    return backward
}
func CopyList(list []string) []string {
    var y []string
    y=make([]string,len(list))
    copy(y,list)
    return y
}
func CopyList2(list2 [][]string) [][]string {
    var (
        i int
        y [][]string
    )
    y=make([][]string,len(list2))
    for i=range list2 {
        y[i]=make([]string,len(list2[i]))
        copy(y[i],list2[i])
    }
    return y
}
func ForwardEdges(sources []string,nodeSucc map[string][]string,edgeSucc map[string]map[string][][]string,verbose bool) [][]string {
    var (
        source,nsucc string
        eroot,check,esucc []string
        forward,newCheck,toCheck [][]string
    )
    for _,source=range sources {
        if verbose {
            fmt.Println("    forwarding "+source)
        }
        for _,nsucc=range nodeSucc[source] {
            eroot=[]string{source,nsucc}
            if !IsInList2(forward,eroot) {
                forward=append(forward,CopyList(eroot))
                newCheck=[][]string{CopyList(eroot)}
                for {
                    toCheck=CopyList2(newCheck)
                    newCheck=[][]string{}
                    for _,check=range toCheck {
                        for _,esucc=range edgeSucc[check[0]][check[1]] {
                            if !IsInList2(forward,esucc) {
                                forward=append(forward,CopyList(esucc))
                                newCheck=append(newCheck,CopyList(esucc))
                            }
                        }
                    }
                    if len(newCheck)==0 {
                        break
                    }
                }
            }
        }
    }
    return forward
}
func IntersectEdges(edges1,edges2 [][]string) [][]string {
    var (
        edge []string
        intersect [][]string
    )
    fmt.Println("computing intersection")
    for _,edge=range edges1 {
        if IsInList2(edges2,edge) {
            intersect=append(intersect,CopyList(edge))
        }
    }
    return intersect
}
func IsInList(list []string,thatElement string) bool {
    var element string
    for _,element=range list {
        if element==thatElement {
            return true
        }
    }
    return false
}
func IsInList2(list2 [][]string,thatList []string) bool {
    var (
        found bool
        i int
        list []string
    )
    for _,list=range list2 {
        if len(list)==len(thatList) {
            found=true
            for i=range list {
                if list[i]!=thatList[i] {
                    found=false
                    break
                }
            }
            if found {
                return true
            }
        }
    }
    return false
}
func IsInNetwork(nodeSP map[string][]string,thatNode string) bool {
    var node string
    for node=range nodeSP {
        if node==thatNode {
            return true
        }
    }
    return false
}
func ReadNetwork(networkFile string) (map[string][]string,map[string][]string,map[string]map[string][][]string,map[string]map[string][][]string,map[string]map[string]string) {
    var (
        err error
        node,node2,node3 string
        line []string
        lines [][]string
        nodeSucc,nodePred map[string][]string
        edgeNames map[string]map[string]string
        edgeSucc,edgePred map[string]map[string][][]string
        file *os.File
        reader *csv.Reader
    )
    fmt.Println("reading "+networkFile)
    file,err=os.Open(networkFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+err.Error())
    } else {
        reader=csv.NewReader(file)
        reader.Comma='\t'
        reader.Comment=0
        reader.FieldsPerRecord=3
        reader.LazyQuotes=false
        reader.TrimLeadingSpace=true
        reader.ReuseRecord=true
        lines,err=reader.ReadAll()
        if err!=nil {
            fmt.Println("ERROR: "+err.Error())
        } else {
            nodeSucc=make(map[string][]string)
            nodePred=make(map[string][]string)
            edgeSucc=make(map[string]map[string][][]string)
            edgePred=make(map[string]map[string][][]string)
            edgeNames=make(map[string]map[string]string)
            for _,line=range lines {
                for _,node=range []string{line[0],line[2]} {
                    nodeSucc[node]=[]string{}
                    nodePred[node]=[]string{}
                }
                edgeSucc[line[0]]=make(map[string][][]string)
                edgePred[line[0]]=make(map[string][][]string)
                edgeNames[line[0]]=make(map[string]string)
            }
            for _,line=range lines {
                if IsInList(nodeSucc[line[0]],line[2]) {
                    fmt.Println("ERROR: multi-edges (or duplicated edges)")
                    nodeSucc=make(map[string][]string)
                    nodePred=make(map[string][]string)
                    edgeSucc=make(map[string]map[string][][]string)
                    edgePred=make(map[string]map[string][][]string)
                    edgeNames=make(map[string]map[string]string)
                    break
                } else {
                    nodeSucc[line[0]]=append(nodeSucc[line[0]],line[2])
                    nodePred[line[2]]=append(nodePred[line[2]],line[0])
                    edgeSucc[line[0]][line[2]]=[][]string{}
                    edgePred[line[0]][line[2]]=[][]string{}
                    edgeNames[line[0]][line[2]]=line[1]
                }
            }
            for node=range nodeSucc {
                for _,node2=range nodeSucc[node] {
                    for _,node3=range nodeSucc[node2] {
                        edgeSucc[node][node2]=append(edgeSucc[node][node2],[]string{node2,node3})
                    }
                }
            }
            for node=range nodePred {
                for _,node2=range nodePred[node] {
                    for _,node3=range nodePred[node2] {
                        edgePred[node2][node]=append(edgePred[node2][node],[]string{node3,node2})
                    }
                }
            }
        }
    }
    return nodeSucc,nodePred,edgeSucc,edgePred,edgeNames
}
func ReadNodes(nodeFile string,nodeSP map[string][]string) []string {
    var (
        err error
        line,nodes []string
        lines [][]string
        file *os.File
        reader *csv.Reader
    )
    fmt.Println("reading "+nodeFile)
    file,err=os.Open(nodeFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+err.Error())
    } else {
        reader=csv.NewReader(file)
        reader.Comma='\t'
        reader.Comment=0
        reader.FieldsPerRecord=1
        reader.LazyQuotes=false
        reader.TrimLeadingSpace=true
        reader.ReuseRecord=true
        lines,err=reader.ReadAll()
        if err!=nil {
            fmt.Println("ERROR: "+err.Error())
        } else {
            for _,line=range lines {
                if !IsInNetwork(nodeSP,line[0]) {
                    fmt.Println("WARNING: "+line[0]+" not in network")
                } else if !IsInList(nodes,line[0]) {
                    nodes=append(nodes,line[0])
                }
            }
        }
    }
    return nodes
}
func ShortestPaths(sources,targets []string,intersect [][]string) [][]string {
    var (
        found bool
        node,node2,node3,source,target string
        edge,edge2,visited []string
        edges,layer,paths,shortest [][]string
        layers [][][]string
        nsucc,npred map[string][]string
        esucc,epred map[string]map[string][][]string
    )
    fmt.Println("computing shortest paths")
    nsucc=make(map[string][]string)
    esucc=make(map[string]map[string][][]string)
    for _,edge=range intersect {
        for _,node=range edge {
            nsucc[node]=[]string{}
        }
        esucc[edge[0]]=make(map[string][][]string)
    }
    for _,edge=range intersect {
        nsucc[edge[0]]=append(nsucc[edge[0]],edge[1])
        esucc[edge[0]][edge[1]]=[][]string{}
    }
    for node=range nsucc {
        for _,node2=range nsucc[node] {
            for _,node3=range nsucc[node2] {
                esucc[node][node2]=append(esucc[node][node2],[]string{node2,node3})
            }
        }
    }
    for _,source=range sources {
        fmt.Println("    from "+source)
        if IsInNetwork(nsucc,source) {
            for _,target=range targets {
                fmt.Println("        to "+target)
                if IsInNetwork(nsucc,target) {
                    found=false
                    layers=[][][]string{}
                    layer=[][]string{}
                    visited=[]string{source}
                    for _,node=range nsucc[source] {
                        if node!=source {
                            layer=append(layer,[]string{source,node})
                        } else if source==target {
                            layer=append(layer,[]string{source,node})
                        }
                    }
                    for {
                        for _,edge=range layer {
                            visited=append(visited,edge[1])
                        }
                        layers=append(layers,CopyList2(layer))
                        layer=[][]string{}
                        for _,edge=range layers[len(layers)-1] {
                            if edge[1]==target {
                                found=true
                                break
                            } else {
                                for _,edge2=range esucc[edge[0]][edge[1]] {
                                    if !IsInList(visited,edge2[1]) {
                                        layer=append(layer,CopyList(edge2))
                                    } else if edge2[1]==target {
                                        layer=append(layer,CopyList(edge2))
                                    }
                                }
                            }
                        }
                        if found || len(layer)==0 {
                            break
                        }
                    }
                    if found {
                        edges=[][]string{}
                        npred=make(map[string][]string)
                        epred=make(map[string]map[string][][]string)
                        for _,layer=range layers {
                            edges=append(edges,CopyList2(layer)...)
                        }
                        for _,edge=range edges {
                            for _,node=range edge {
                                npred[node]=[]string{}
                            }
                            epred[edge[0]]=make(map[string][][]string)
                        }
                        for _,edge=range edges {
                            npred[edge[1]]=append(npred[edge[1]],edge[0])
                            epred[edge[0]][edge[1]]=[][]string{}
                        }
                        for node=range npred {
                            for _,node2=range npred[node] {
                                for _,node3=range npred[node2] {
                                    epred[node2][node]=append(epred[node2][node],[]string{node3,node2})
                                }
                            }
                        }
                        paths=BackwardEdges([]string{target},npred,epred,false)
                        for _,edge=range paths {
                            if !IsInList2(shortest,edge) {
                                shortest=append(shortest,CopyList(edge))
                            }
                        }
                    }
                }
            }
        }
    }
    return shortest
}
func WriteNetwork(networkFile string,edges [][]string,edgeNames map[string]map[string]string) {
    var (
        err error
        edge []string
        lines [][]string
        file *os.File
        writer *csv.Writer
    )
    fmt.Println("writing "+networkFile)
    file,err=os.Create(networkFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+err.Error())
    } else {
        for _,edge=range edges {
            lines=append(lines,[]string{edge[0],edgeNames[edge[0]][edge[1]],edge[1]})
        }
        writer=csv.NewWriter(file)
        writer.Comma='\t'
        writer.UseCRLF=false
        writer.WriteAll(lines)
    }
}
