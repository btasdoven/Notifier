package me

import(
	"os"
	"io/ioutil"
	"encoding/json"
)

type UserModel struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	DeviceId string `json:"device_id"`
	NotifIds []string `json:"notif_ids"`
	CardIds []string `json:"card_ids"`
	RegId string `json:"reg_id"`
}

type NotifModel struct {
	Id string `json:"id"`
	OwnerId string `json:"owner_id"`	
	ClientId string `json:"client_id"`	
	Name string `json:"name"`
	StartTS []int `json:"start_ts"`
	Period int `json:"period"`		
	NotifPeriod int `json:"notif_period"`
}

type CardModel struct {
	Id string `json:"id"`
	NotifId string `json:"notif_id"`
	Timestamp int `json:"timestamp"`
	Completed bool `json:"completed"`
}

type DBModel struct {
	UserCounter int
	CardCounter int
	NotifCounter int
	User map[string]UserModel
	Card map[string]CardModel	
	Notif map[string]NotifModel
}

func (db DBModel) Save(file string) {
	f, err := os.Create(file);
	if err != nil {
		panic(err)
	}
	
	defer f.Close()
	data, _ := json.Marshal(db)
	f.Write(data)
}

func FromFile(file string) (db DBModel) {
	cnt, _ := ioutil.ReadFile(file)
	if err := json.Unmarshal(cnt, &db); err != nil {
        panic(err)
    }
    
    return db
}

func NewDB() (db DBModel) {
	db.User = make(map[string]UserModel)
	db.User["100001"] = UserModel{"100001", "Batuhan Tasdoven", "btasdoven@gmail.com", "123123123", []string{"200001", "200002", "200003"}, []string{}, ""}
	db.User["100002"] = UserModel{"100002", "Cagla Istanbulluoglu", "cagla.istanbulluoglu@gmail.com", "321321321", []string{"200002"}, []string{}, ""}	
	db.UserCounter = 100002
	
	db.Notif = make(map[string]NotifModel)
	db.Notif["200001"] = NotifModel{"200001", "100001", "100002", "Ilacini ictin mi askim? :)", []int{1443805200, 1443855600}, 86400, 300}
	db.Notif["200002"] = NotifModel{"200002", "100001", "100002", "Bugun 50mg askim :)", []int{1445486400}, 3600*48, 7200}
	db.Notif["200003"] = NotifModel{"200003", "100001", "100002", "Bugun 75mg askim :)", []int{1445572800}, 3600*48, 7200}
	db.Notif["200004"] = NotifModel{"200004", "100002", "100001", "Caglayi sev :)", []int{1443805200}, 37200, 3600}	
	db.NotifCounter = 200004
	
	db.Card = make(map[string]CardModel)	
	db.CardCounter = 300000
	
	return db
}

