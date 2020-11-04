package main

import (
    "encoding/xml"
)

type Database struct {
    XMLName                 xml.Name            `xml:"db"`
    Version                 Version             `xml:"version"`
    Satellites              []Satellite         `xml:"satellite"`
    Channels                []Channel           `xml:"channel"`
}

type Version struct {
    XMLName                 xml.Name            `xml:"version"`
    Winsat_editer           string              `xml:"winsat_editer,attr"`
}

type Satellite struct {
    XMLName                 xml.Name            `xml:"satellite"`
    Name                    string              `xml:"name,attr"`
    Longitude               int                 `xml:"longitude,attr"`
    Lof_hi                  int                 `xml:"lof_hi,attr"`
    Lof_lo                  int                 `xml:"lof_lo,attr"`
    Lof_threshold           int                 `xml:"lof_threshold,attr"`
    Lnb_power               string              `xml:"lnb_power,attr"`
    Signal_22khz            string              `xml:"signal_22khz,attr"`
    Toneburst               string              `xml:"toneburst,attr"`
    Diseqc1_0               string              `xml:"diseqc1_0,attr"`
    Diseqc1_1               string              `xml:"diseqc1_1,attr"`
    Motor                   string              `xml:"motor,attr"`

    Transponders            []Transponder       `xml:"transponder"`
}

type Transponder struct {
    XMLName                 xml.Name            `xml:"transponder"`
    Original_network_id     int                 `xml:"original_network_id,attr"`
    Ts_id                   int                 `xml:"ts_id,attr"`
    Frequency               int                 `xml:"frequency,attr"`
    Symbol_rate             int                 `xml:"symbol_rate,attr"`
    Polarisation            string              `xml:"polarisation,attr"`

    Programs                []Program           `xml:"program"`
}

type Program struct {
    XMLName                 xml.Name            `xml:"program"`
    Name                    string              `xml:"name,attr"`
    Service_id              int                 `xml:"service_id,attr"`
    Channel_number          int                 `xml:"channel_number,attr"`
    Type                    int                 `xml:"type,attr"`
    Scrambled               string              `xml:"scrambled,attr"`
    Parental_lock           string              `xml:"parental_lock,attr"`
    Skip                    string              `xml:"skip,attr"`
    Id                      int                 `xml:"id,attr"`
    Plp_id                  int                 `xml:"plp_id,attr"`
    Sdt_version             int                 `xml:"sdt_version,attr"`
    Default_channnel_num    int                 `xml:"default_channnel_num,attr"`
    Lcn                     int                 `xml:"lcn,attr"`
    Fav                     int                 `xml:"fav,attr"`

    Videos                  []Program_video     `xml:"video"`
    Pcrs                    []Program_pcr       `xml:"pcr"`
    Audios                  []Program_audio     `xml:"audio"`
}

type Program_video struct {
    XMLName                 xml.Name            `xml:"video"`
    Pid                     int                 `xml:"pid,attr"`
    Format                  string              `xml:"format,attr"`
}

type Program_pcr struct {
    XMLName                 xml.Name            `xml:"pcr"`
    Pid                     int                 `xml:"pid,attr"`
}

type Program_audio struct {
    XMLName                 xml.Name            `xml:"audio"`
    Pid                     int                 `xml:"pid,attr"`
    Format                  string              `xml:"format,attr"`
    Language                string              `xml:"language,attr"`
}

type Channel struct {
    XMLName                 xml.Name            `xml:"channel"`
    Fe_type                 string              `xml:"fe_type,attr"`
    Frequency               int                 `xml:"frequency,attr"`
    Original_network_id     int                 `xml:"original_network_id,attr"`
    Ts_id                   int                 `xml:"ts_id,attr"`
    Bandwidth               int                 `xml:"bandwidth,attr"`
    Ofdm_mode               int                 `xml:"ofdm_mode,attr"`

    Programs                []Program           `xml:"program"`
}

type ChListElement struct {
    DisplayName             string
    SearchName              string
}
