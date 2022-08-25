package controller

import (
        "fmt"
        "strings"
          "strconv"
//	   "github.com/jung-kurt/gofpdf"
  )

      const(
              limit       = 15
              margenlat   = 3
	      layout      = "2006-01-02"
              timeLayout = "15:04:05"
            )
      var (
            TotalCount  int
            offset      int
            posact      int
           )

//-------------------------------------------------------------
func roundU(val float64) int {
    if val > 0 { return int(val+1.0) }
    return int(val)
}
// ---------------------------------------------------
  func atoi32( str string) (nr uint32,err error){
        i, errn := strconv.Atoi(str)
        nr  = uint32(i)
        err =  errn
        return
    }
// ---------------------------------------------------
func rpad(s string,pad string, plength int)string{
    for i:=len(s);i<plength;i++{
        s= s + pad
    }
    return s
}

 func  commify(st string) (str string){
       parts :=   strings.Split(st, ".")
//       fmt.Println(st, parts)
       if len(parts) == 1 {
            parts = append(parts, "00")
        }else  {
            parts[1] = rpad(parts[1], "0",2)
        }
//       fmt.Println(st, parts)
       str    =  strings.Join( parts, "" )
       str    =  strings.Join( strings.Split(str, ","), "" )
       return
 }

 func  money2int64(st string)(val int64, err error){
       str      := commify(st)
       val, err  = strconv.ParseInt(str, 10, 64)
       return
    }

// ---------------------------------------------------
func getNumberOfButtonsForPagination(TotalCount int, limit int) int {
    num := (int)(TotalCount / limit)
    if (TotalCount%limit > 0) {
        num++
    }
    return num
}
// ---------------------------------------------------
func createSliceForBtns(number int, posact int) []int {
    var sliceOfBtn []int
    lffin := margenlat
    rtini := number   -  margenlat  + 1
    inilf := posact   -  margenlat
    finrt := posact   +  margenlat
    if inilf < 1 {
       inilf = 1
      }
    if finrt > number  {
       finrt =  number
      }
    if lffin  > inilf  {
       lffin  = inilf - 1
    }
    if rtini  < finrt  {
        rtini = finrt  + 1
    }
    for i := 1; i <= lffin; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    for i := inilf; i <= finrt; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    for i := rtini; i <= number; i++ {
        sliceOfBtn = append(sliceOfBtn, i)
    }
    return sliceOfBtn
}
// ---------------------------------------------------

 func ConcatNames(s1,s2, sep string) string {
        st1 := strings.Trim(s1, " ")
        st2 := strings.Trim(s2, " ")
        s := []string{st1,st2}
        st :=  fmt.Sprintf(strings.Join(s, sep))
        return st
      }

//--------------------------------------------------------


