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
    "time"
)
func main() {
    var (
        maxStep,maxWalk,selfConnect,shortest int64
        sources,targets []string
        allPaths [][]string
        succ map[string][]string
        edges map[string]map[string]string
    )
    if len(os.Args)!=8 {
        fmt.Println("konect networkFile sourceFile targetFile maxStep maxWalk selfConnect shortest")
    } else {
        rand.Seed(int64(time.Now().Nanosecond()))
        succ,edges=ReadNetwork(os.Args[1])
        sources=ReadNodes(os.Args[2],succ)
        targets=ReadNodes(os.Args[3],succ)
        maxStep,_=strconv.ParseInt(os.Args[4],10,0)
        maxWalk,_=strconv.ParseInt(os.Args[5],10,0)
        selfConnect,_=strconv.ParseInt(os.Args[6],10,0)
        shortest,_=strconv.ParseInt(os.Args[7],10,0)
        allPaths=FindAllPaths(sources,targets,int(maxStep),int(maxWalk),int(selfConnect),int(shortest),succ)
        WriteNetwork("konected.sif",allPaths,edges)
    }
}
func CopyPath(path []string) []string {
    var y []string
    y=make([]string,len(path))
    copy(y,path)
    return y
}
func FindAllPaths(sources,targets []string,maxStep,maxWalk,selfConnect,shortest int,succ map[string][]string) [][]string {
    var (
        i1,i2,i3,imax int
        tail string
        paths,allPaths [][]string
    )
    tail="/"+strconv.FormatInt(int64(len(sources)),10)+")"
    for i1=range sources {
        fmt.Println("Sourcing "+sources[i1]+" ("+strconv.FormatInt(int64(i1+1),10)+tail)
        for i2=range targets {
            if (selfConnect==1) || (sources[i1]!=targets[i2]) {
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
func IsInEdges(edges map[string]map[string]string,thatSource string) bool {
    var source string
    for source=range edges {
        if source==thatSource {
            return true
        }
    }
    return false
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
    lines,_=reader.ReadAll()
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
    lines,_=reader.ReadAll()
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
