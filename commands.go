package main

/* Command deployment and initialization goes on here */

import (
    "fmt"

    "./cmds"
    "./cmds/seen"
    "github.com/bwmarrin/discordgo"
)

type cmdfn func(*discordgo.Session, *discordgo.MessageCreate)

var Cmd map[string]cmdfn

// Maybe as needed provide a more detailed Man page for each command?
func man(s *discordgo.Session, m *discordgo.MessageCreate) {
    var manual string
    manual = `>>> 
**!seen: ** Provides information on user specified.
**Usage:** *!seen @pentashift*.

**!insult: ** Insults target provided using Markov chains to generate a uniquely crafted and sometimes completely incoherent insult.
**Usage:** *!insult @Ash Bailey*

**!opt-in: ** Opts into Sentiment Analysis run by this bot along with mining of recent messages to be used within the !seen command.
**Usage:** *!opt-in*

**!opt-out: ** Opts user out of Sentiment Analysis run by this bot along with the mining of recent messages to be used within the *!seen* command. Date joined will still show up but consequently *Personality* and *Last Message* will both be substituted with "Opted out." This is the default upon joining the server.
**Usage:** *!opt-out*
`
    s.ChannelMessageSend(m.ChannelID, manual)
}

func command_init() {
    cmds.ReadInsults()
    fmt.Println("Insult Markov Initialized")

    Cmd = map[string]cmdfn {
        "!insult": cmds.Insult,
        "!seen": seen.Seen,
        "!opt-in": seen.OptIn,
        "!opt-out": seen.OptOut,
        "!man": man,
    }
    fmt.Println("Command Hash Initialized")
}
