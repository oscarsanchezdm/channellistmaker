package main

import (
    "strings"
    "github.com/agext/levenshtein" //run  go get github.com/agext/levenshtein
)

//Pre: a desired channel and the list of all the exported channels
//Post: an array of channels that have similar name as the desired
func getSimilarChannels(search ChListElement, programs []*Program, imp_param float64) []*Program  {
    var SimilarPrograms []*Program
    var sim_value float64

    var param float64 = imp_param

    for ok := true; ok; ok = ((len(SimilarPrograms)<1) && (param > 0.4)) {
        for i := 0; i < len (programs); i++ {
            string1 := strings.ToUpper(search.SearchName)
            string1 = strings.Replace(string1, " ", "", -1)
            string1 = strings.Split(string1, "HD")[0]

            string2 := trimFirstThirdRuns(strings.ToUpper(programs[i].Name))
            string2 = strings.Replace(string2, " ", "", -1)
            string2 = strings.Split(string2, "HD")[0]

            sim_value = levenshtein.Similarity(string1, string2, levenshtein.NewParams())
            if (sim_value > param) { SimilarPrograms = append(SimilarPrograms, programs[i])}
        }
        param = param*0.95
    }
    return SimilarPrograms;
}

//@Pre: a list of channels that are theorically the same

//@Post: the list of channels after a deletion of channels that have the same name
//with the same or prior characteritics (less quality, scrambled, etc)
func deleteCopies(similarChannels []*Program) []*Program {
    for i := 0; i < len (similarChannels); i++ {
        //first step: check if theare are other channels with the SAME name
        string1 := strings.ToUpper(strings.Replace(similarChannels[i].Name, " ", "", -1))
        string1 = strings.Split(string1, "HD")[0]
        for j := 0; ((j < len (similarChannels)) && (i!=j) && (i < len (similarChannels))); j++ {
            string2 := strings.ToUpper(strings.Replace(similarChannels[j].Name, " ", "", -1))
            string2 = strings.Split(string2, "HD")[0]

            //check if the name is the same
            if (string1 == string2) {
                  //check which has more video quality. mpeg for SD, h264 for HD, hevc for 4K (not supported for most countries)
                  if (similarChannels[i].Videos[0].Format == similarChannels[j].Videos[0].Format) {
                      //check if one of them is not scrambled
                      if (similarChannels[i].Scrambled != similarChannels[j].Scrambled) {
                          if (similarChannels[i].Scrambled == "true") {
                              similarChannels = append(similarChannels[:i], similarChannels[i+1:]...)
                              if (i != 0) { i-- }
                          } else {
                              similarChannels = append(similarChannels[:j], similarChannels[j+1:]...)
                              j--
                          }
                      } else {
                          //check Audios
                          if (sameAudioLanguages(similarChannels[i],similarChannels[j])) {
                              similarChannels = append(similarChannels[:i], similarChannels[i+1:]...)
                              if (i != 0) { i-- }
                          } else { continue }
                      }
                  } else if (similarChannels[i].Videos[0].Format == "h264") {
                      similarChannels = append(similarChannels[:j], similarChannels[j+1:]...)
                      j--
                  } else if (similarChannels[j].Videos[0].Format == "h264") {
                      similarChannels = append(similarChannels[:i], similarChannels[i+1:]...)
                      if (i != 0) { i-- }
                  }
             }
        }
    }
    return similarChannels
}


//Pre:
//Post: it returns a matrix. Every [i] entry is the desiredch, and the [j]
//are the matches for the desired channelm
func makeMatchMatrix(programs []*Program, desiredChList []ChListElement, db *Database) [][]*Program {
    var returntable [][]*Program

    for i := 0; i < len (desiredChList); i++ {
        if (desiredChList[i].SearchName != "%empty%") {
            var similarChannels []*Program = getSimilarChannels(desiredChList[i], programs, 1.0)
            similarChannels = deleteCopies(similarChannels)
            returntable = append(returntable, similarChannels)
      } else {
          var similarChannels []*Program
          similarChannels = append(similarChannels,generateEmptyChannel(i+1,db))
          returntable = append(returntable, similarChannels)
      }
    }
    return returntable
}


//Pre: a database. a boolean asking for keeping the radio channels
//Post: the database with the channels starting by 'x' removed
func cleanup(db *Database, keepRadio bool) {
    //cleanup of SAT programs
    for i := 0; i < len(db.Satellites[0].Transponders); i++ {
        for j := 0; j < len(db.Satellites[0].Transponders[i].Programs); j++ {
            if ((len(db.Satellites[0].Transponders[i].Programs[j].Name)==0) ||
             (db.Satellites[0].Transponders[i].Programs[j].Name[0]=='x')) {
                 if ((len(db.Satellites[0].Transponders[i].Programs[j].Videos)==0) &&
                 (len(db.Satellites[0].Transponders[i].Programs[j].Audios)>0) && keepRadio) {
                     //radio channel
                 } else {
                    db.Satellites[0].Transponders[i].Programs =
                    append(db.Satellites[0].Transponders[i].Programs[:j],
                        db.Satellites[0].Transponders[i].Programs[j+1:]...)
                    j--
                }
            }
        }
    }

    //cleanup of DVB programs
    for i := 0; i < len(db.Channels); i++ {
        for j := 0; j < len(db.Channels[i].Programs); j++ {
            if ((len(db.Channels[i].Programs[j].Name)==0) ||
             (db.Channels[i].Programs[j].Name[0]=='x')) {
                 if ((len(db.Channels[i].Programs[j].Videos)==0) &&
                 (len(db.Channels[i].Programs[j].Audios)>0) && keepRadio) {
                     //radio channel
                 } else {
                    db.Channels[i].Programs =append(db.Channels[i].Programs[:j],
                        db.Channels[i].Programs[j+1:]...)
                    j--
                }
            }
        }
    }
}
