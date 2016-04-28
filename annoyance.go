/* Annoyance Discord Chat Bot by IMcPwn.
 * Copyright 2016 Carleton Stuberg

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.

 * For the latest code and contact information visit: http://imcpwn.com
 */

package main

import (
    "fmt"
    "io/ioutil"
    "flag"
    "math/rand"

    "github.com/bwmarrin/discordgo"
    "github.com/bwmarrin/dgvoice"
)

var FOLDER *string

func main() {
    TOKEN := flag.String("t", "", "Discord authentication token")
    FOLDER = flag.String("f", "", "Folder that contains the mp3s to play")
    flag.Parse()

    dg, err := discordgo.New(*TOKEN)
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        return
    }

    dg.AddHandler(VoiceStateUpdate)

    // Open the websocket and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        return
    }

    // Make sure we're logged in successfully
    prefix, err := dg.User("@me")
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        return
    }
    fmt.Println("Logged in as " + prefix.Username)

    // Set status to "away"
    dg.UpdateStatus(1, "")

    fmt.Println("Welcome to Annoyance! Press enter to quit.")
    var input string
    fmt.Scanln(&input)
    return
}

// This function is called whenever there is a voice state update
// i.e mute/unmute, channel join/leave, etc.
func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
    fmt.Println("[*] Called")
    if v.ChannelID == "" {
        fmt.Println("[X] Invalid channel")
        return
    }
    if len(s.VoiceConnections) != 0 {
        fmt.Println("[X] Already speaking")
        return
    }
    // 1/20 chance of not ignoring the call 
    // This is so the bot is not triggered every time there is a voice update.
    if rand.Intn(20) != 1 {
        fmt.Println("[X] Ignoring call")
        return
    }
    fmt.Println("[*] Responding to call")

    fmt.Println("[*] Joining Channel ID #" + v.ChannelID)
    // Join the server unmuted and deafened
    dgv, err := s.ChannelVoiceJoin(v.GuildID, v.ChannelID, false, true)
    if err != nil {
        fmt.Println(err)
        return
    }
    // Start loop and attempt to play all files in the given folder
    fmt.Println("[*] Reading Folder: ", *FOLDER)
    files, _ := ioutil.ReadDir(*FOLDER)
    for _, f := range files {
        fmt.Println("[*] PlayAudioFile:", f.Name())
        // Say we're "playing" the name of the audio file
        s.UpdateStatus(0, f.Name())
        dgvoice.PlayAudioFile(dgv, fmt.Sprintf("%s/%s", *FOLDER, f.Name()))
    }
    // Set status to away
    s.UpdateStatus(1, "")
    defer dgv.Disconnect()
    defer dgv.Close()
}