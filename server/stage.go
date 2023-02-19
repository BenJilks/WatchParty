package main

import (
    "fmt"
    "math/rand"
)

type Seat struct {
    Row    int `json:"row"`
    Column int `json:"column"`
}

type Stage struct {
    seatsUsed map[string]Seat
}

type StageUpdateMessage struct {
    SeatsNotFree []Seat  `json:"seats_not_free"`
    YourToken    *string `json:"your_token"`
    YourSeat     Seat    `json:"your_seat"`
}

func (stage *Stage) UpdateMessage(yourToken *string) StageUpdateMessage {
    var seatsNotFree []Seat
    var yourSeat Seat
    for token, seat := range stage.seatsUsed {
        seatsNotFree = append(seatsNotFree, seat)
        if yourToken != nil && token == *yourToken {
            yourSeat = seat
        }
    }

    return StageUpdateMessage{
        SeatsNotFree: seatsNotFree,
        YourToken:    yourToken,
        YourSeat:     yourSeat,
    }
}

func (stage *Stage) PlayerInSeat(seat Seat) *string {
    for player, playerSeat := range stage.seatsUsed {
        if playerSeat.Row == seat.Row && playerSeat.Column == seat.Column {
            return &player
        }
    }

    return nil
}

func (stage *Stage) PlaceViewer(token string) {
    for {
        row := rand.Intn(len(RowSeatCount))
        column := rand.Intn(RowSeatCount[row])
        seat := Seat{
            Row:    row,
            Column: column,
        }

        existingPlayer := stage.PlayerInSeat(seat)
        if existingPlayer == nil {
            stage.seatsUsed[token] = seat
            break
        }
    }

    seat := stage.seatsUsed[token]
    fmt.Printf("Placed '%s' in row %d, seat %d\n",
        token, seat.Row, seat.Column)
}

func (stage *Stage) RemovePlayer(token string) {
    if _, has := stage.seatsUsed[token]; !has {
        return
    }

    delete(stage.seatsUsed, token)
}
