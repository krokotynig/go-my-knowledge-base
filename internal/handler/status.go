package handler

import (
	"fmt"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "✅ MyKnowledgeBase Server работает!\n")
	fmt.Fprintf(w, "✅ DB подключена!") // Fprintf пишет в любой io Writer, типа файла , браузера
	// fmt.Fprintf(w, "✅ API готов!")
}
