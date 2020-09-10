package seen

import (
    "strconv"
    "time"

    "database/sql"
    "github.com/bwmarrin/discordgo"
    _ "github.com/mattn/go-sqlite3"
)

func role_color(conn *sql.DB, id string) int {
    var role string
    row := conn.QueryRow(`SELECT rank FROM users WHERE id = ?`, id)
    err := row.Scan(&role)
    if err == nil {
        switch role {
        case "Owner":
            return 0x00b8ff
        case "Admin":
            return 0xd600ff
        case "Moderator":
            return 0xfd0e35
        case "Game Master":
            return 0x00ff9f
        }
    }
    return 0xf8f8ff
}

func Seen(s *discordgo.Session, m *discordgo.MessageCreate) {
    if len(m.Mentions) == 0 {
        s.ChannelMessageSend(m.ChannelID, ">>> Please include your target.\n**Usage:** !seen @target")
        return
    } else if m.Mentions[0].ID == s.State.User.ID {
        s.ChannelMessageSend(m.ChannelID, "Pathetic.")
        return
    }

    var last_message string
    var server_joined string
    var personality_score string

    conn, _ := sql.Open("sqlite3", "./data/user.db")
    color := role_color(conn, m.Mentions[0].ID)
    stmt := `SELECT last_message, joined, personality FROM users WHERE id = ?`
    row := conn.QueryRow(stmt, m.Mentions[0].ID)
    row.Scan(&last_message, &server_joined, &personality_score)
    if !isOpt(conn, m.Mentions[0].ID) {
        last_message = "Opted out."
        personality_score = "Opted out."
    } else {
         timestamp, _ := strconv.Atoi(last_message[:10])
         ltime := time.Unix(int64(timestamp), 0)
         last_message = ltime.Format(time.ANSIC) + "\n*" + last_message[11:] + "*"
    }
    joinedtime, _ := strconv.Atoi(server_joined)
    joined_timestamp := time.Unix(int64(joinedtime), 0)
    server_joined = joined_timestamp.Format(time.ANSIC)

    avi := discordgo.MessageEmbedThumbnail {
            URL: m.Mentions[0].AvatarURL(""),
            Width: 1000,
            Height: 1000,
    }

    joined := discordgo.MessageEmbedField {
                    Name: "Joined: ",
                    Value: server_joined,
                    Inline: false,
    }

    personality := discordgo.MessageEmbedField {
                    Name: "Personality: ",
                    Value: personality_score,
                    Inline: true,
    }

    lastmessage := discordgo.MessageEmbedField {
                    Name: "Last Message: ",
                    Value: last_message,
                    Inline: false,
    }

    embed := discordgo.MessageEmbed {
        Title: m.Mentions[0].Username,
        Color: color,
        Thumbnail: &avi,
        Fields: []*discordgo.MessageEmbedField {
                    &lastmessage,
                    &personality,
                    &joined,
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func AddUser(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
    if s.State.User.ID == e.User.ID {
        return
    }

    conn, _ := sql.Open("sqlite3", "./data/user.db")

    /* Check if user already exists and if so just return */
    var username string
    stmt := `SELECT username FROM users WHERE id = ?`
    row := conn.QueryRow(stmt, e.User.ID)
    err := row.Scan(&username)
    if err != nil {
        stmt = `INSERT INTO users (id, username, personality, last_message, joined, opt) VALUES (?, ?, ?, ?, ?, ?)`
        ustmt, _ := conn.Prepare(stmt)
        timestamp := strconv.FormatInt(time.Now().Unix(), 10)
        ustmt.Exec(e.User.ID, e.User.Username, "0.0", timestamp + ":  ", timestamp, 1)
    }

    conn.Close()
}
