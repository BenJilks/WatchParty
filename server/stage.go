package main

import (
	log "github.com/sirupsen/logrus"
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

func (stage *Stage) SeatForPlayer(token string) *Seat {
	seat, found := stage.seatsUsed[token]
	if !found {
		return nil
	}

	return &seat
}

func (stage *Stage) PlaceViewer(token string) {
	for {
		row := rand.Intn(len(RowSeatCount)-2) + 2
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
	log.WithFields(log.Fields{
		"token": token,
		"row":   seat.Row,
		"seat":  seat.Column,
	}).Info("Assigned seat")
}

func (stage *Stage) RemovePlayer(token string) {
	if _, has := stage.seatsUsed[token]; !has {
		return
	}

	delete(stage.seatsUsed, token)
}
