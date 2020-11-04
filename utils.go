package main

import (
    "fmt"
    "strconv"
)

//Pre: a string with more than 3 runs/chars
//Post: the string without the first 3 runs/chars
func trimFirstThirdRuns(s string) string {
    m := 0
    for i := range s {
        if m >= 3 {
            return s[i:]
        }
        m++
    }
    return s[:0]
}

//Pre: two channels
//Post: true if the two channels have the same audio tracks
func sameAudioLanguages(program1 *Program, program2 *Program) bool {
    if (len(program1.Audios) != len(program2.Audios)) { return false }
    for i := 0; i < len (program1.Audios); i++ {
        if (program1.Audios[i].Language != program2.Audios[i].Language) { return false }
    }
    return true
}


func dbClean(db *Database, maxfreq int) {
    //TRANSPONDER CLEANUP
    var cleanTransponder bool = false
    if (maxfreq>0) { cleanTransponder = true }


    for i := 0; i < len(db.Satellites[0].Transponders); i++ {
        if ((cleanTransponder==true) && (db.Satellites[0].Transponders[i]).Frequency>maxfreq) {
            db.Satellites[0].Transponders = append(db.Satellites[0].Transponders[:i],
                db.Satellites[0].Transponders[i+1:]...)
            i--
            continue
        }

        for j := 0; j < len(db.Satellites[0].Transponders[i].Programs); j++ {
            //add "xxx" at the beggining of the channel name
            if (len(db.Satellites[0].Transponders[i].Programs[j].Name)>2) {
                name := []byte(db.Satellites[0].Transponders[i].Programs[j].Name)
                name[0]='x'
                name[1]='x'
                name[2]='x'
                db.Satellites[0].Transponders[i].Programs[j].Name = string(name)
            }
            if ((len(db.Satellites[0].Transponders[i].Programs[j].Videos)==0) ||
            (len(db.Satellites[0].Transponders[i].Programs[j].Name) == 0) ||
            (db.Satellites[0].Transponders[i].Programs[j].Name == "xxxNo Name")) {
                db.Satellites[0].Transponders[i].Programs = append(db.Satellites[0].Transponders[i].Programs[:j],
                    db.Satellites[0].Transponders[i].Programs[j+1:]...)
                j--
            }
        }
        //add to the pointer array the DVB programs
        for i := 0; i < len(db.Channels); i++ {
            for j := 0; j < len(db.Channels[i].Programs); j++ {
                if (len(db.Channels[i].Programs[j].Name)>2) {
                    name := []byte(db.Channels[i].Programs[j].Name)
                    name[0]='x'
                    name[1]='x'
                    name[2]='x'
                    db.Channels[i].Programs[j].Name = string(name)
                }
            }
        }
    }
}

//Pre: a channel which name is xxxChannel and the desired number
//Post: the channel gets the number insted of xxx
func numerateChannel(p *Program, n int, displayName string) int {
    //check if the program has a number
    if (p.Name[0]!='x') {
        var original_number int
        original_number, err := strconv.Atoi(p.Name[0:3])
        if (err != nil) { panic("strange error numerating the channel") }
        return -1*original_number
    } else if (displayName == "%empty%") {
        displayName = "Sense nom"
    }
    number := fmt.Sprintf("%03d", n)
    p.Name=number+displayName
    return 0;
}


func getFakeTransponder(db *Database) *Transponder {
    if (db.Satellites[0].Transponders[len(db.Satellites[0].Transponders)-1].Frequency != 10101010) {
        var newT Transponder
        var prgs []Program

        newT.Ts_id = 1050
        newT.Frequency = 10101010
        newT.Symbol_rate = 21999000
        newT.Polarisation = "V"
        newT.Programs = prgs

        db.Satellites[0].Transponders = append(db.Satellites[0].Transponders,newT)
    }
    return &db.Satellites[0].Transponders[len(db.Satellites[0].Transponders)-1]
}


//Pre:
//Post: a non-existant program used tu fulfill the desired channel list
func generateEmptyChannel(number int, db *Database) *Program {
    var emptyChannel Program
    var videos []Program_video
    var video0 Program_video
    var audios []Program_audio
    var audio0 Program_audio

    aux_str := fmt.Sprintf("%03d", number)
    emptyChannel.Name = aux_str + "Sense nom"
    emptyChannel.Service_id = 30817
    emptyChannel.Channel_number = 1
    emptyChannel.Type = 1
    emptyChannel.Scrambled = "false"
    emptyChannel.Parental_lock = "false"
    emptyChannel.Skip = "false"
    emptyChannel.Id = 2300 + number
    emptyChannel.Plp_id = 255
    emptyChannel.Sdt_version = 0
    emptyChannel.Default_channnel_num = 22
    emptyChannel.Lcn = 2222
    emptyChannel.Fav = 0

    video0.Format = "h264"
    audio0.Format = "aac3"
    audio0.Language = "cat"

    videos = append(videos,video0)
    audios = append(audios,audio0)

    emptyChannel.Videos = videos
    emptyChannel.Audios = audios

    //Add the channel to a Transponder
    fakeT := getFakeTransponder(db)
    fakeT.Programs = append(fakeT.Programs,emptyChannel)

    return &fakeT.Programs[len(fakeT.Programs)-1]
}
