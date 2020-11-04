package main

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "os"
    "bufio"
    "strings"
    "time"
)

//read the exported XML
//post: a database with the contents of input.xml
func getDatabase() (int, *Database) {
    var database Database
    xmlFile, err := os.Open("input.xml")
    if err != nil {
        return -1,&database
    }
    defer xmlFile.Close()
    byteValue, _ := ioutil.ReadAll(xmlFile)
    xml.Unmarshal(byteValue, &database)
    //generate new transponder at first position.
    return 1,&database
}

//read the desired channel list
//post: an array of ChList elements
func getDesiredChList() (int, []ChListElement) {
    //initialize the ret value
    var retSlice []ChListElement

    //initialize aux variables
    var current_line string

    //READING THE FILE
    file, err := os.Open("desired.txt")
    if err != nil {
        return -1,retSlice
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      current_line = scanner.Text() //new line
      var newEl ChListElement
      if (len(current_line)>0) { //check if the line is not empty
          if ((len(current_line)>1) && current_line[0]=='/' && current_line[1]=='/') { continue }
          var splitted_string []string = strings.Split(current_line, "|")
          newEl.DisplayName = splitted_string[0]
          if ((len(splitted_string)>1) && (len(splitted_string[1]))>0) {
            newEl.SearchName = splitted_string[1]
          } else {
            newEl.SearchName = newEl.DisplayName //if no searchname is given, use displayname
          }

      } else {
          newEl.DisplayName = "%empty%"
          newEl.SearchName = newEl.DisplayName
      }

      retSlice = append(retSlice,newEl)
    }
    return 1,retSlice;
}

//Returns an array of pointers
//Pre: a database with Channels
//Post: an array of pointers to channels
func program_list(db *Database) []*Program {
    var ptrs []*Program

    //add to the pointer array the SAT programs
    for i := 0; i < len(db.Satellites[0].Transponders); i++ {
        for j := 0; j < len(db.Satellites[0].Transponders[i].Programs); j++ {
            ptrs = append(ptrs, &(db.Satellites[0].Transponders[i].Programs[j]))
        }
    }

    //add to the pointer array the DVB programs
    for i := 0; i < len(db.Channels); i++ {
        for j := 0; j < len(db.Channels[i].Programs); j++ {
            ptrs = append(ptrs, &(db.Channels[i].Programs[j]))
        }
    }
    return ptrs;
}

//Pre: a database
//Post: a satellites_YYYY_MM_DD_HH_mm_SS.xml is generated with the db's content
func writeFile(db *Database) string {
    filename := "satellites_"
    t := time.Now()
    filename = filename + fmt.Sprintf(t.Format("2006_01_02_15_04_05"))

    filename = filename + ".xml"

    file, _ := xml.MarshalIndent(db, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)

    return filename

}
