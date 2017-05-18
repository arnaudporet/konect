// Copyright 2017 Arnaud Poret
// This work is licensed under the BSD 2-Clause License.
package main
import (
    "encoding/csv"
    "fmt"
    "math/rand"
    "os"
    "sort"
    "strconv"
    "strings"
    "time"
)
func main() {
    var (
        maxStep,maxWalk,shortest int64
        sources,targets []string
        allPaths [][]string
        succ map[string][]string
        edges map[string]map[string]string
    )
    if len(os.Args)==2 && os.Args[1]=="help" {
        fmt.Println(strings.Join([]string{
            "",
            "konect is a basic algorithm for connecting nodes according to a reference",
            "network.",
            "",
            "Typical usage consists in extracting, from the reference network, the paths",
            "connecting a couple of nodes belonging to it.",
            "",
            "konect handles networks encoded in the sif file format.",
            "",
            "konect does not handle multigraphs (i.e. networks with nodes connected by more",
            "than one edge).",
            "",
            "Usage: konect networkFile sourceFile targetFile maxStep maxWalk shortest",
            "",
            "    * networkFile: the reference network encoded in a sif file",
            "",
            "    * sourceFile: the source nodes listed in a csv file (one node per line)",
            "",
            "    * targetFile: the target nodes listed in a csv file (one node per line)",
            "",
            "    * maxStep: the maximum number of steps performed during a random walk",
            "      starting from a source node in an attempt to reach a target node",
            "",
            "    * maxWalk: the maximum number of random walks performed in the reference",
            "      network to find paths from a source node to a target node",
            "",
            "    * shortest (1 or 0): among the found connecting paths, select only (1) or",
            "      not only (0) the shortest",
            "",
            "The returned file is a sif file encoding a subnetwork of the reference network",
            "connecting the source nodes to the target nodes.",
            "",
            "For more information: https://github.com/arnaudporet/konect",
            "",
        },"\n"))
    } else if len(os.Args)==2 && os.Args[1]=="license" {
        fmt.Println(strings.Join([]string{
            "",
            "Copyright 2017 Arnaud Poret",
            "",
            "Redistribution and use in source and binary forms, with or without modification,",
            "are permitted provided that the following conditions are met:",
            "",
            "1. Redistributions of source code must retain the above copyright notice, this",
            "   list of conditions and the following disclaimer.",
            "",
            "2. Redistributions in binary form must reproduce the above copyright notice,",
            "   this list of conditions and the following disclaimer in the documentation",
            "   and/or other materials provided with the distribution.",
            "",
            "THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\" AND",
            "ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED",
            "WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE",
            "DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR",
            "ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES",
            "(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;",
            "LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON",
            "ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT",
            "(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS",
            "SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
            "",
        },"\n"))
    } else if len(os.Args)==2 && os.Args[1]=="usage" {
        fmt.Println(strings.Join([]string{
            "",
            "konect networkFile sourceFile targetFile maxStep maxWalk shortest",
            "",
        },"\n"))
    } else if len(os.Args)==7 {
        shortest,_=strconv.ParseInt(os.Args[6],10,0)
        maxWalk,_=strconv.ParseInt(os.Args[5],10,0)
        maxStep,_=strconv.ParseInt(os.Args[4],10,0)
        if int(shortest)!=0 && int(shortest)!=1 {
            panic("shortest must be 1 or 0")
        }
        if int(maxWalk)<1 {
            panic("maxWalk must 1 or more")
        }
        if int(maxStep)<1 {
            panic("maxStep must 1 or more")
        }
        rand.Seed(int64(time.Now().Nanosecond()))
        succ,edges=ReadNetwork(os.Args[1])
        sources=ReadNodes(os.Args[2],succ)
        targets=ReadNodes(os.Args[3],succ)
        allPaths=FindAllPaths(sources,targets,int(maxStep),int(maxWalk),int(shortest),succ)
        WriteNetwork("konected.sif",allPaths,edges)
    } else {
        fmt.Println(strings.Join([]string{
            "",
            "To print help:        konect help",
            "To print license:     konect license",
            "To print usage:       konect usage",
            "",
            "For more information: https://github.com/arnaudporet/konect",
            "",
        },"\n"))
    }
}
func CopyPath(path []string) []string {
    var y []string
    y=make([]string,len(path))
    copy(y,path)
    return y
}
func FindAllPaths(sources,targets []string,maxStep,maxWalk,shortest int,succ map[string][]string) [][]string {
    var (
        i1,i2,i3,imax int
        tail1,tail2 string
        paths,allPaths [][]string
    )
    tail1="/"+strconv.FormatInt(int64(len(sources)),10)+")"
    tail2="/"+strconv.FormatInt(int64(len(targets)),10)+")"
    for i1=range sources {
        fmt.Println("Sourcing "+sources[i1]+" ("+strconv.FormatInt(int64(i1+1),10)+tail1)
        for i2=range targets {
            fmt.Println("    Targeting "+targets[i2]+" ("+strconv.FormatInt(int64(i2+1),10)+tail2)
            paths=FindPaths(sources[i1],targets[i2],maxStep,maxWalk,succ)
            if shortest==1 {
                imax=sort.Search(len(paths),func(i int) bool {return len(paths[i])>len(paths[0])})
            } else {
                imax=len(paths)
            }
            for i3=0;i3<imax;i3++ {
                if !IsInPaths(allPaths,paths[i3]) {
                    allPaths=append(allPaths,CopyPath(paths[i3]))
                }
            }
        }
    }
    return allPaths
}
func FindPaths(source,target string,maxStep,maxWalk int,succ map[string][]string) [][]string {
    var (
        i int
        path []string
        paths [][]string
    )
    for i=0;i<maxWalk;i++ {
        path=RandomWalk(source,target,maxStep,succ)
        if len(path)!=0 && !IsInPaths(paths,path) {
            paths=append(paths,CopyPath(path))
        }
    }
    sort.Slice(paths,func(i,j int) bool {return len(paths[i])<len(paths[j])})
    return paths
}
func IsInPath(path []string,thatNode string) bool {
    var node string
    for _,node=range path {
        if node==thatNode {
            return true
        }
    }
    return false
}
func IsInPaths(paths [][]string,thatPath []string) bool {
    var path []string
    for _,path=range paths {
        if PathEq(path,thatPath) {
            return true
        }
    }
    return false
}
func IsInSucc(succ map[string][]string,thatNode string) bool {
    var node string
    for node=range succ {
        if node==thatNode {
            return true
        }
    }
    return false
}
func PathEq(path1,path2 []string) bool {
    var i int
    if len(path1)!=len(path2) {
        return false
    } else {
        for i=range path1 {
            if path1[i]!=path2[i] {
                return false
            }
        }
        return true
    }
}
func RandomWalk(source,target string,maxStep int,succ map[string][]string) []string {
    var (
        i int
        current string
        path []string
    )
    current=source
    path=append(path,source)
    for i=0;i<maxStep;i++ {
        if len(succ[current])==0 {
            break
        } else {
            current=succ[current][rand.Intn(len(succ[current]))]
            if IsInPath(path,current) && (current!=target) {
                break
            } else {
                path=append(path,current)
                if current==target {
                    return path
                }
            }
        }
    }
    return []string{}
}
func ReadNetwork(networkFile string) (map[string][]string,map[string]map[string]string) {
    var (
        err error
        node string
        line []string
        lines [][]string
        succ map[string][]string
        edges map[string]map[string]string
        reader *csv.Reader
        file *os.File
    )
    fmt.Println("Reading "+networkFile)
    file,_=os.Open(networkFile)
    defer file.Close()
    reader=csv.NewReader(file)
    reader.Comma='\t'
    reader.Comment=0
    reader.FieldsPerRecord=3
    reader.LazyQuotes=false
    reader.TrimLeadingSpace=true
    lines,err=reader.ReadAll()
    if err!=nil {
        fmt.Println("\nERROR: "+err.Error()+"\n")
        panic(networkFile+" is not properly formated")
    }
    succ=make(map[string][]string)
    edges=make(map[string]map[string]string)
    for _,line=range lines {
        for _,node=range []string{line[0],line[2]} {
            succ[node]=[]string{}
        }
        edges[line[0]]=make(map[string]string)
    }
    for _,line=range lines {
        succ[line[0]]=append(succ[line[0]],line[2])
        edges[line[0]][line[2]]=line[1]
    }
    return succ,edges
}
func ReadNodes(nodeFile string,succ map[string][]string) []string {
    var (
        err error
        nodes,line []string
        lines [][]string
        reader *csv.Reader
        file *os.File
    )
    fmt.Println("Reading "+nodeFile)
    file,_=os.Open(nodeFile)
    defer file.Close()
    reader=csv.NewReader(file)
    reader.Comma=','
    reader.Comment=0
    reader.FieldsPerRecord=1
    reader.LazyQuotes=false
    reader.TrimLeadingSpace=true
    lines,err=reader.ReadAll()
    if err!=nil {
        fmt.Println("\nERROR: "+err.Error()+"\n")
        panic(nodeFile+" is not properly formated")
    }
    for _,line=range lines {
        if !IsInPath(nodes,line[0]) && IsInSucc(succ,line[0]) {
            nodes=append(nodes,line[0])
        }
    }
    return nodes
}
func WriteNetwork(networkFile string,paths [][]string,edges map[string]map[string]string) {
    var (
        i int
        path,line []string
        lines [][]string
        writer *csv.Writer
        file *os.File
    )
    fmt.Println("Writing "+networkFile)
    for _,path=range paths {
        for i=0;i<len(path)-1;i++ {
            line=[]string{path[i],edges[path[i]][path[i+1]],path[i+1]}
            if !IsInPaths(lines,line) {
                lines=append(lines,CopyPath(line))
            }
        }
    }
    file,_=os.Create(networkFile)
    defer file.Close()
    writer=csv.NewWriter(file)
    writer.Comma='\t'
    writer.UseCRLF=false
    writer.WriteAll(lines)
}
