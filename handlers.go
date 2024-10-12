package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/olahol/melody"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("./html/index.html"))
	if err := templ.Execute(w, m.GetFields()); err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}
}

func HandleWS(m *melody.Melody) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if ss, _ := m.Sessions(); len(ss) == 2 {
			w.WriteHeader(403)
		}
		if err := m.HandleRequest(w, r); err != nil {
			log.Println(err)
			w.WriteHeader(500)
		}
	}
}

type WSOutgoingCellMessage struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func (ocm *WSOutgoingCellMessage) WriteButton(id, content string) {
	ocm.Id = id
	ocm.Content = fmt.Sprintf("<button class='cell' id='%s'>%s</button>", id, content)
}

func (ocm *WSOutgoingCellMessage) WriteError(content string) {
	ocm.Id = "error"
	ocm.Content = fmt.Sprintf("<p id=\"error\">%s</p>", content)
}

func (m *TMap) SendCell(mel *melody.Melody, from *melody.Session, row, col int, which string) {
	cell := m.First[row][col]
	if which == "s" {
		cell = m.Second[row][col]
	}

	var ocm WSOutgoingCellMessage
	ocm.WriteButton(fmt.Sprintf("%v-%v-%v", which, row, col), cell)

	bytes, _ := json.Marshal(ocm)
	err := mel.BroadcastOthers(bytes, from)
	if err != nil {
		log.Println(err)
	}
}

type WSIncomingCellMessage struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

func HandleMessage(m *TMap) func(s *melody.Session, msg []byte) {
	return func(s *melody.Session, msg []byte) {
		var icm WSIncomingCellMessage
		var ocm WSOutgoingCellMessage
		err := json.Unmarshal(msg, &icm)
		if err != nil {
			log.Println(err)
			ocm.WriteError("Bad request")
			bytes, _ := json.Marshal(ocm)
			s.Write(bytes)
			return
		}
		fmt.Println(icm)
		if icm.Key == m.FirstKey {

		} else if icm.Key == m.SecondKey {

		} else {
			ocm.WriteError("Nuhai bebru")
			bytes, _ := json.Marshal(ocm)
			s.Write(bytes)
			return
		}
	}
}
