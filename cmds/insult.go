package cmds

import (
    "bufio"
    "os"
    "strings"
    "time"
    "math/rand"

    "github.com/bwmarrin/discordgo"
)

type Pair struct {
    wrd string
    cnt int
}

// TODO :: Rework these functions for resuability for compliment markov.
var start []string
var Markov map[string][]Pair

// Takes all possible states and returns a transition matrix.
func transMatrix(pState []Pair) ([]float64) {
    var x float64 = 0.0
    tMat := make([]float64, len(pState))
    for i := range pState {
        x = float64(pState[i].cnt) + x
    }
    for i := range pState {
        tMat[i] = float64(pState[i].cnt) / x
        if i >= 1 {
            var y float64 = 0.0
            for j := i; j >= 0; j-- {
                y = float64(pState[j].cnt) + y
            }
            tMat[i] = y / x
        }
    }
    return tMat
}


// Reads the insult training data and initializes the Markov data structure.
func ReadInsults() {
    file, err := os.Open("Insults.txt")
    Markov = make(map[string][]Pair)
    //var start []string
    if err != nil {
        // Figure out how to do an error handle.
    }
    s := bufio.NewScanner(file)
    for s.Scan() {
        wrds := strings.Fields(s.Text())
        start = append(start, wrds[0])
        for i := range wrds {
            if i + 1 < len(wrds) {
                if Markov[wrds[i]] != nil {
                    exists := false
                    for j := range Markov[wrds[i]] {
                        if Markov[wrds[i]][j].wrd == wrds[i+1] {
                            exists = true
                            Markov[wrds[i]][j].cnt++
                            break
                        }
                    }
                    if exists == false {
                        Markov[wrds[i]] = append(Markov[wrds[i]], Pair{wrds[i+1], 1})
                    }
                } else {
                    Markov[wrds[i]] = append(Markov[wrds[i]], Pair{wrds[i+1], 1})
                }
            }
        }
    }
    return
}

// Loops through Markov training data to craft a uniquely new insult.
func Insult(s *discordgo.Session, m *discordgo.MessageCreate) {
    if len(m.Mentions) == 0 {
        s.ChannelMessageSend(m.ChannelID, ">>> Please include your target(s).\n**Usage:** !insult @victim.")
        return
    }
    rand.Seed(time.Now().UnixNano())
    word := start[rand.Intn(len(start))]
    state := Markov[word]
    insult := word
    for {
        if state == nil {
            break
        }
        tMat := transMatrix(state)
        r := rand.Float64()
        for i := range state {
            if r < tMat[i]{
                word = state[i].wrd
                insult = insult + " " + word
                break
            }
        }
        state = Markov[word]
    }
    s.ChannelMessageSend(m.ChannelID, ("> **" + m.Mentions[0].Username + "** " + insult))
}
