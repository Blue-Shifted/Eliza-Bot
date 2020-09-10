package seen

import (
    "math"
    "time"
    "strconv"
    "strings"
    "regexp"

    "database/sql"
    //"github.com/lloyd/wnram"
    "github.com/bwmarrin/discordgo"
    _ "github.com/mattn/go-sqlite3"
)

// For now a simple lexicon sentiment analysis. Maybe in furture iterations
// I will transform this into a more advanced sentiment analysis.

//func lemma() { }

// Loops through each word looking it up in sentiment database then derives a score.
func getScore(personality float64, words []string) (float64) {
    var neg float64
    var pos float64
    var point float64
    sdb, _ := sql.Open("sqlite3", "./data/SA.db")
    stmt := `SELECT point FROM lexicon WHERE word = $1`
    for i := range(words) {
        row := sdb.QueryRow(stmt, words[i])
        err := row.Scan(&point)
        if err == nil {
            if point < 0 {
                neg = neg - point
            } else if point > 0 { pos = pos + point }
        }
    }
    sdb.Close()
    return personality + (math.Log(pos + 0.5) - math.Log(neg + 0.5))
}

func isOpt(conn *sql.DB, id string) (bool) {
    var opt int
    stmt := `SELECT opt FROM users WHERE id = $1`
    row := conn.QueryRow(stmt, id)
    row.Scan(&opt)
    if opt == 1 {
        return false
    } else {
        return true
    }
}

func OptIn(s *discordgo.Session, m *discordgo.MessageCreate) {
    conn, _ := sql.Open("sqlite3", "./data/user.db")
    ustmt, _ := conn.Prepare("UPDATE users SET opt = 0 WHERE id = ?")
    ustmt.Exec(m.Author.ID)
    s.ChannelMessageSend(m.ChannelID, ("> **" + m.Author.Username + "** is now opted in to data mining services."))
    conn.Close()
}

func OptOut(s *discordgo.Session, m *discordgo.MessageCreate) {
    conn, _ := sql.Open("sqlite3", "./data/user.db")
    ustmt, _ := conn.Prepare("UPDATE users SET opt = 1 WHERE id = ?")
    ustmt.Exec(m.Author.ID)
    s.ChannelMessageSend(m.ChannelID, ("> **" + m.Author.Username + "** is now opted out of data mining services."))
    conn.Close()
}

// Updates the user database with new sentiment score.
func SentAnalysis(s *discordgo.Session, m *discordgo.MessageCreate) {
    var personality float64
    conn, _ := sql.Open("sqlite3", "./data/user.db")

    if !isOpt(conn, m.Author.ID) {
        conn.Close()
        return
    }
    words := strings.Fields(strings.ToLower(m.Content))
    reg, _ := regexp.Compile("[^a-z ]+")
    for i := range words {
        words[i] = reg.ReplaceAllString(words[i], "")
    }
    stmt := `SELECT personality FROM users WHERE id = $1`
    row := conn.QueryRow(stmt, m.Author.ID)
    row.Scan(&personality)
    personality = math.Round(getScore(personality, words) * 100) / 100

    msg := strconv.FormatInt(time.Now().Unix(), 10) + ":" + m.Content
    ustmt, _ := conn.Prepare("UPDATE users SET personality = ?, last_message = ? WHERE id = ?")
    ustmt.Exec(personality, msg, m.Author.ID)
    conn.Close()
}


