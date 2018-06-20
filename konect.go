// Copyright 2017-2018 Arnaud Poret
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
        maxStep,maxWalk,shortest,selfConnect int64
        sources,targets []string
        allPaths [][]string
        succ map[string][]string
        edges map[string]map[string]string
    )
    if (len(os.Args)==2) && (os.Args[1]=="help") {
        fmt.Println(strings.Join([]string{
            "",
            "konect is a tool for connecting nodes according to a reference network.",
            "",
            "Typical usage consists in extracting, from the reference network, the paths connecting a couple of nodes.",
            "",
            "konect handles networks encoded in the sif file format.",
            "",
            "konect does not handle multi-graphs (i.e. networks where nodes can be connected by more than one edge).",
            "",
            "Usage: konect networkFile sourceFile targetFile maxStep maxWalk shortest selfConnect",
            "",
            "    * networkFile: the reference network encoded in a sif file",
            "",
            "    * sourceFile: the source nodes listed in a file (one node per line)",
            "",
            "    * targetFile: the target nodes listed in a file (one node per line)",
            "",
            "    * maxStep (>0): the maximum number of steps performed during a random walk when searching for a path connecting a source node to a target node",
            "",
            "    * maxWalk (>0): the maximum number of random walks performed in the reference network when searching for paths connecting a source node to a target node",
            "",
            "    * shortest (0 or 1): among the found connecting paths, selects only the shortest ones (1) or not (0)",
            "",
            "    * selfConnect (0 or 1): if a node belongs to both the source and target nodes, allows to find paths connecting it to itself (1) or not (0)",
            "",
            "The returned file is a sif file encoding a subnetwork (of the reference network) connecting the source nodes to the target nodes.",
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
            "konect networkFile sourceFile targetFile maxStep maxWalk shortest selfConnect",
            "",
        },"\n"))
    } else if len(os.Args)==8 {
        selfConnect,_=strconv.ParseInt(os.Args[7],10,0)
        shortest,_=strconv.ParseInt(os.Args[6],10,0)
        maxWalk,_=strconv.ParseInt(os.Args[5],10,0)
        maxStep,_=strconv.ParseInt(os.Args[4],10,0)
        if (int(selfConnect)!=0) && (int(selfConnect)!=1) {
            fmt.Println("ERROR: selfConnect must be 0 or 1")
        } else if (int(shortest)!=0) && (int(shortest)!=1) {
            fmt.Println("ERROR: shortest must be 0 or 1")
        } else if int(maxWalk)<1 {
            fmt.Println("ERROR: maxWalk must 1 or more")
        } else if int(maxStep)<1 {
            fmt.Println("ERROR: maxStep must 1 or more")
        } else {
            succ,edges=ReadNetwork(os.Args[1])
            sources=ReadNodes(os.Args[2],succ)
            targets=ReadNodes(os.Args[3],succ)
            if len(edges)==0 {
                fmt.Println("ERROR: "+os.Args[1]+" is empty after reading")
            } else if len(sources)==0 {
                fmt.Println("ERROR: "+os.Args[2]+" is empty after reading")
            } else if len(targets)==0 {
                fmt.Println("ERROR: "+os.Args[3]+" is empty after reading")
            } else {
                rand.Seed(int64(time.Now().Nanosecond()))
                allPaths=FindAllPaths(sources,targets,int(maxStep),int(maxWalk),int(shortest),int(selfConnect),succ)
                if len(allPaths)==0 {
                    fmt.Println("WARNING: no connecting paths found")
                } else {
                    WriteNetwork("konected.sif",allPaths,edges)
                }
            }
        }
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
func CopyPath(path []string) []string {
    var y []string
    y=make([]string,len(path))
    copy(y,path)
    return y
}
func FindAllPaths(sources,targets []string,maxStep,maxWalk,shortest,selfConnect int,succ map[string][]string) [][]string {
    var (
        i1,i2,i3,imax int
        tail1,tail2 string
        paths,allPaths [][]string
    )
    tail1="/"+strconv.FormatInt(int64(len(sources)),10)+")"
    tail2="/"+strconv.FormatInt(int64(len(targets)),10)+")"
    for i1=range sources {
        fmt.Println("sourcing "+sources[i1]+" ("+strconv.FormatInt(int64(i1+1),10)+tail1)
        for i2=range targets {
            if (sources[i1]!=targets[i2]) || (selfConnect==1) {
                fmt.Println("    targeting "+targets[i2]+" ("+strconv.FormatInt(int64(i2+1),10)+tail2)
                paths=FindPaths(sources[i1],targets[i2],maxStep,maxWalk,succ)
                if len(paths)!=0 {
                    if shortest==1 {
                        sort.Slice(paths,func(i,j int) bool {return len(paths[i])<len(paths[j])})
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
        if (len(path)!=0) && !IsInPaths(paths,path) {
            paths=append(paths,CopyPath(path))
        }
    }
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
    fmt.Println("reading "+networkFile)
    file,_=os.Open(networkFile)
    reader=csv.NewReader(file)
    reader.Comma='\t'
    reader.Comment=0
    reader.FieldsPerRecord=3
    reader.LazyQuotes=false
    reader.TrimLeadingSpace=true
    reader.ReuseRecord=true
    lines,err=reader.ReadAll()
    file.Close()
    succ=make(map[string][]string)
    edges=make(map[string]map[string]string)
    if err!=nil {
        fmt.Println("ERROR: "+networkFile+" "+err.Error())
    } else {
        for _,line=range lines {
            for _,node=range []string{line[0],line[2]} {
                succ[node]=[]string{}
            }
            edges[line[0]]=make(map[string]string)
        }
        for _,line=range lines {
            if IsInPath(succ[line[0]],line[2]) {
                fmt.Println("ERROR: "+networkFile+" contains multi-edges")
                succ=make(map[string][]string)
                edges=make(map[string]map[string]string)
                break
            } else {
                succ[line[0]]=append(succ[line[0]],line[2])
                edges[line[0]][line[2]]=line[1]
            }
        }
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
    fmt.Println("reading "+nodeFile)
    file,_=os.Open(nodeFile)
    reader=csv.NewReader(file)
    reader.Comma=','
    reader.Comment=0
    reader.FieldsPerRecord=1
    reader.LazyQuotes=false
    reader.TrimLeadingSpace=true
    reader.ReuseRecord=true
    lines,err=reader.ReadAll()
    file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+nodeFile+" "+err.Error())
    } else {
        for _,line=range lines {
            if !IsInSucc(succ,line[0]) {
                fmt.Println("WARNING: "+line[0]+" in "+nodeFile+" but not in network")
            } else if !IsInPath(nodes,line[0]) {
                nodes=append(nodes,line[0])
            }
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
    fmt.Println("writing "+networkFile)
    for _,path=range paths {
        for i=0;i<len(path)-1;i++ {
            line=[]string{path[i],edges[path[i]][path[i+1]],path[i+1]}
            if !IsInPaths(lines,line) {
                lines=append(lines,CopyPath(line))
            }
        }
    }
    file,_=os.Create(networkFile)
    writer=csv.NewWriter(file)
    writer.Comma='\t'
    writer.UseCRLF=false
    writer.WriteAll(lines)
    file.Close()
}
