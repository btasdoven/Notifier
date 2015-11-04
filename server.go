package main

import (	
	"me"
	"fmt"
	"sync"
	"time"
    "strconv"
    "net/http"
    
	"github.com/gin-gonic/gin"
    "github.com/alexjlockwood/gcm"
)

func HeapRun() {

	tsHeap := me.NewHeap("notifs")
	now := int(time.Now().Unix())
	
	for _, notif := range db.Notif {
		for _, ts := range notif.StartTS {
			for (ts < now) {
				ts += notif.Period
			}
			tsHeap.PushVal(ts, notif.Id)
		}
	}
		
	for {
		ts := int(time.Now().Unix())
		for ( tsHeap.Min() > -1 && ts >= tsHeap.Min() ) {
			notifTs, notifId := tsHeap.PopVal()
			fmt.Printf("[%v] - Popped: (%d %s)\n", ts, notifTs, notifId)
			
			if notif, ok := db.Notif[notifId]; ok {
				for (ts >= notifTs) {
					notifTs += notif.Period
				}
				
				tsHeap.PushVal(notifTs, notif.Id)
				cardid := CreateCard(notif, ts)
				if (cardid != "") {
					tsHeap.PushVal(notifTs - notif.Period + notif.NotifPeriod, cardid)
				}
			} else if card, okk := db.Card[notifId]; okk {
				if (card.Completed != true) {
					if notif, okkk := db.Notif[card.NotifId]; okkk {
						CreateNotif(card, ts)
						tsHeap.PushVal(ts + notif.NotifPeriod, card.Id)
					}
				}
			}
		} 
		
		if (tsHeap.Min() > -1) {
			fmt.Printf("[%v] - Sleeping for %d seconds\n", ts, tsHeap.Min() - ts)
			time.Sleep( time.Duration(tsHeap.Min() - ts) * time.Second )
		} else {
			fmt.Printf("[%v] - Sleeping for 5 seconds\n", ts)	
			time.Sleep( 5 * time.Second )	
		}
		
		db.Save("dumb.json")
	}
}

var db me.DBModel

func main() {

	var wg sync.WaitGroup
	defer wg.Wait()
	
//	db = me.NewDB()
	db = me.FromFile("dumb.json")
	defer db.Save()
//	gin.SetMode(gin.ReleaseMode)
	
	r := gin.Default()    
	r.GET("/cards/:user_id", Cards)
	r.GET("/notifs/:user_id", Notifs)   
	r.GET("/register/:user_id/:reg_id", Register) 
	r.GET("/done/:card_id", Done) 
	r.GET("/users", Users) 
	r.GET("/heap", Heap)	
	
	wg.Add(1)	 	  	      
	go r.Run("144.122.71.77:8080")
    
	HeapRun()
}


func Heap(c *gin.Context) {
	if h, ok := me.GetHeap("notifs"); ok {
		
		type HeapItem struct {
			Ts time.Time `json:"ts"`
			Id string `json:"id"`
		}
		
		var heap_items []HeapItem
		
		for i := 0; i < h.Len(); i++ {
			ts, id := h.Get(i)
			heap_items = append(heap_items, HeapItem{time.Unix(int64(ts), 0), id}) 		
		}
		c.JSON(http.StatusOK, heap_items)		
	} else {
		c.String(http.StatusNotFound, "")
	}
}

func Users(c *gin.Context) {
	c.JSON(http.StatusOK, db.User)
}

func Cards(c *gin.Context) {
	userId := c.Param("user_id")
	
	type CardNotif struct {
		Id string `json:"id"`
		me.CardModel
		me.NotifModel
	}
	
	var card_notifs []CardNotif
	
	if user, ok := db.User[userId]; ok {
		for _, cardId := range user.CardIds {
			if card, okk := db.Card[cardId]; okk {
				if notif, okkk := db.Notif[card.NotifId]; okkk {
					card_notifs = append(card_notifs, CardNotif{card.Id, card, notif})
				}
			}
		}		
		c.JSON(http.StatusOK, card_notifs)
	} else {  
		c.String(http.StatusNotFound, "")
	}
}

func Notifs(c *gin.Context) {
	userId := c.Param("user_id")
	
	if user, ok := db.User[userId]; ok {
		notifs := []me.NotifModel{}
		for _, notifId := range user.NotifIds {
			if notif, okk := db.Notif[notifId]; okk {
				notifs = append(notifs, notif)
			}
		}		
		c.JSON(http.StatusOK, notifs)
	} else {  
		c.String(http.StatusNotFound, "")
	} 
}

func Register(c *gin.Context) {
	userId := c.Param("user_id")
	regId := c.Param("reg_id")
		
	if user, ok := db.User[userId]; ok {
		user.RegId = regId		
		db.User[userId] = user
		c.String(http.StatusOK, "")
	} else {  
		c.String(http.StatusNotFound, "")
	} 
}


/* Mark the card as completed and
   Send a notification to the owner of the card
*/
func Done(c *gin.Context) {
	cardId := c.Param("card_id")
	
	if card, ok := db.Card[cardId]; ok {
		card.Completed = true
		db.Card[cardId] = card
		if notif, okk := db.Notif[card.NotifId]; okk {
			if user, okkk := db.User[notif.OwnerId]; okkk {
				data := map[string]interface{}{"title": notif.Name + " is done.", "id": cardId}
				regIDs := []string{user.RegId}
				msg := gcm.NewMessage(data, regIDs...)

				// Create a Sender to send the message.
				sender := &gcm.Sender{ApiKey: "AIzaSyDhdyFnigm2EfKj4LgccjytRYcvUWl6aLA"}
		
				// Send the message and receive the response after at most two retries.
				_, err := sender.Send(msg, 2)
				if err != nil {
					c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to send message: %v", err))
				} else {		
					c.String(http.StatusOK, "Message is sent.")
				}
			}
		}
	} else {  
		c.String(http.StatusNotFound, "")
	} 
}

func CreateCard(notif me.NotifModel, ts int) string {
	db.CardCounter += 1
	cardId := strconv.Itoa(db.CardCounter)
	db.Card[cardId] = me.CardModel{cardId, notif.Id, ts, false}	
	if owner, ok := db.User[notif.OwnerId]; ok {
		if client, okk := db.User[notif.ClientId]; okk {
			owner.CardIds = append(owner.CardIds, cardId)
			db.User[notif.OwnerId] = owner
			
			data := map[string]interface{}{"title": notif.Name, "id": cardId}
			regIDs := []string{client.RegId}
			msg := gcm.NewMessage(data, regIDs...)

			// Create a Sender to send the message.
			sender := &gcm.Sender{ApiKey: "AIzaSyDhdyFnigm2EfKj4LgccjytRYcvUWl6aLA"}
	
			// Send the message and receive the response after at most two retries.
			_, err := sender.Send(msg, 2)
			if err != nil {
				fmt.Println("Failed to send message: %v", err)
			} else {		
				fmt.Println("Message is sent.")
			}
			
			return cardId
		}
	}
	return ""
}

func CreateNotif(card me.CardModel, ts int) {
	if notif, ok := db.Notif[card.NotifId]; ok {
		if client, okk := db.User[notif.ClientId]; okk {
			data := map[string]interface{}{"title": notif.Name, "id": card.Id}
			regIDs := []string{client.RegId}
			msg := gcm.NewMessage(data, regIDs...)

			// Create a Sender to send the message.
			sender := &gcm.Sender{ApiKey: "AIzaSyDhdyFnigm2EfKj4LgccjytRYcvUWl6aLA"}

			// Send the message and receive the response after at most two retries.
			_, err := sender.Send(msg, 2)
			if err != nil {
				fmt.Println("Failed to send message: %v", err)
			} else {		
				fmt.Println("Message is sent.")
			}
		}
	}
}

