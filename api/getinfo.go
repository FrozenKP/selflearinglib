package api

import (
  "io"
  "os"
  "strings"
  "fmt"
  "net/http"
  "encoding/json"
  "encoding/csv"
  "strconv"
  "time"
)

type Getinfo struct {
}

type IndexInfo struct{
  Title   string
  Serial  string
  Cover   string
  Begin   string
  End     string
}

type IndexInfoPkg struct{
  Status  int
  Data    []IndexInfo
}

//index info
func (gihandle Getinfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  //GET
  if r.Method != "GET" {
    fmt.Println(w,"404 page not found")
    return
  }

  //pre process
  r.ParseForm()
  var out IndexInfoPkg
  out.Status=200
  pathS:=strings.Split(r.URL.Path[1:],"/")
  //assign page,year
  validPage,validYear:=true,true
  var pageString,yearString string
  if len(pathS)==2{
    validYear=false
    pageString=pathS[1]
  }else if len(pathS)==3{
    pageString=pathS[1]
    yearString=pathS[2]
  }else{
    validPage,validYear=false,false
    out.Status=404
  }

  //show info
  fmt.Println(time.Now())
  fmt.Println("/getinfo from",r.RemoteAddr)
  fmt.Println("Page:",pageString,"Year:",yearString)
  fmt.Println("====================================")

  //process
  if validPage {
    page,_:=strconv.Atoi(pageString)
    //get data
    f,_:=os.OpenFile("./data/info.csv",os.O_RDONLY,0777)
    defer f.Close()
    r:=csv.NewReader(f)
    for i,times:=(page-1)*20,0 ; i<page*20 ; times++ {
      result,err:=r.Read()
      //check eof
      if err==io.EOF {
        if i==(page-1)*20{
          out.Status=404
          validPage=false
        }
        break
      }

      //check year
      if validYear {
        year:=strings.TrimLeft((strings.Split(result[3],"-")[0])," ")
        if year != yearString {
          times--
          continue
        }
      }

      //valid or not
      if i!=times {
        continue;
      }else{
        //valid -> add to out
        i++;
        data:=IndexInfo{result[0],result[1],result[2],result[3],result[4]}
        out.Data=append(out.Data,data)
      }
    }
  }

  //make json and send
  j,_:=json.Marshal(out)
  w.Header().Set("Content-Type", "application/json")
  w.Write(j)
}
