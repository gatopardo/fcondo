package plugin

import (
	"html/template"
	"time"
        "fmt"
        "strings"
	"github.com/gatopardo/acad/ch14/fcondo/app/controller"

)

func commas(s string) string {
    lon := len(s)
    if lon <= 3 {
        return s
    } else {
        return commas(s[0:len(s)-3]) + "," + s[len(s)-3:]
    }
}

func Format64() template.FuncMap {
      f := make(template.FuncMap)

      f["FORMAT64"] = func(dat int64)(str string) {
          str  = fmt.Sprintf("%d", dat)
          lon  := len(str)
	  if  dat < 0 {
	          str = str[1:lon]
	  }
          lon  = len(str)  - 2
	  if lon > 0 {
              sini := str[:lon]
              sfin := str[lon:]
              scom := commas(sini)
              str  =  scom + "." + sfin
          }else{
	      pre := "0."
              if lon < 0 {
                  pre = "0.0"
              }
               str = pre + str
	  }
	  if  dat < 0 {
	       str = "-"+str
	  }
          return
     }

     return f
   }

func toString(f float64) string {
    parts := strings.Split(fmt.Sprintf("%.2f", f), ".")
    if parts[0][0] == '-' {
        return "-" + commas(parts[0][1:]) + "." + parts[1]
    }
    return commas(parts[0]) + "." + parts[1]
}

// PrettyTime returns a template.FuncMap
// * PRETTYTIME outputs a nice time format
func PrettyTime() template.FuncMap {
	f := make(template.FuncMap)

	f["PRETTYTIME"] = func(t time.Time) string {
		return t.Format("2006/01/02")
	}

	return f
}

// * DATEFORMAT outputs a nice time format
func DateFormat() template.FuncMap {
	f := make(template.FuncMap)
        layout := "2006/01/02"
	f["DATEFORMAT"] = func(t time.Time) string {
		return t.Format(layout)
	}

	return f
 }


 func ConcatStr() template.FuncMap {
      f := make(template.FuncMap)

      f["CONCATSTR"] = func(s1,s2 string) string {
        st :=  controller.ConcatNames(s1,s2," ")
        return st
      }
     return f
  }


