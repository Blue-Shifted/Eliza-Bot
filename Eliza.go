package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "strings"
    "./cmds/seen"

    "github.com/bwmarrin/discordgo"
)

// TODO :: 
//          Do A LOT more documentation and some more cleaning up of code.
//          Also probably would be smart to throw some error handling in there somewhere.
//          Look into migrating towards a big boy database like PostgreSQL or even MySQL

func main() {
    // Initializing Markov data structures.
    dg, err := discordgo.New("Bot " + "")

    dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages)
    dg.AddHandler(onUserJoin)
    dg.AddHandler(onUserLeave)
    dg.AddHandler(messageCreate)
    dg.AddHandler(messageDelete)    // Catches all these messages as nil, handle at later date?
    dg.AddHandler(messageUpdate)    // Same here.
    dg.AddHandler(onReady)

    command_init()
    fmt.Println("Connecting to Discord.")
    err = dg.Open()
    if err != nil {
        fmt.Println("Error in opening connection.")
        return
    }
    fmt.Println("Connected to Discord.")

    //dg.ChannelVoiceJoin("711031646421254164", "", false, false)

    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
    dg.Close()
    fmt.Println("\nThanks for using Eliza")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    /* Ignore messages from bot self. */
    if m.Author.ID == s.State.User.ID {
        return
    }
    /* Cause apparently pictures will break the bot otherwise. */
    if m.Content == "" {
        return
    }
    words := strings.Fields(m.Content)
    if strings.HasPrefix(words[0], "!") {
        if Cmd[words[0]] == nil {
            s.ChannelMessageSend(m.ChannelID, ("> Command: *" + words[0] + "* does not exist. Stop."))
            return
        } else { Cmd[words[0]](s, m) }
    } else { seen.SentAnalysis(s, m) }
}

/* Placeholder code to see if this even works */
func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
    if m.BeforeDelete != nil {
        s.ChannelMessageSend(m.ChannelID, "Deleted: " + m.BeforeDelete.Content)
    }
}

/* Placeholder code to see if this even works */
func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
    if m.BeforeUpdate != nil {
        s.ChannelMessageSend(m.ChannelID, "Updated: " + m.BeforeUpdate.Content)
    }
}

/* Users are added to the user.db upon entering the server here. */
func onUserJoin(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
    seen.AddUser(s, e)
}

/* May do exit notifications. */
func onUserLeave(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
    fmt.Println(e.User.Username + ": has left.\n")
}

func onReady(s *discordgo.Session, m *discordgo.Ready) {
    s.UpdateStatus(0, "Type !man for more information.")
    /* Probably should be moved to commands.go file for initialization*/
    /* Verfies user.db is up to date. */
    mem, _ := s.GuildMembers("711031646421254164", "0", 100)
    for i := range mem {
        seen.AddUser(s, &discordgo.GuildMemberAdd{mem[i]})
    }
}
