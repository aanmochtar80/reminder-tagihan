package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.mau.fi/whatsmeow/types"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

var (
	WAClient       *whatsmeow.Client
	CurrentQR      string
	qrMutex        sync.RWMutex
	storeContainer *sqlstore.Container
)

func InitWhatsApp() {
	if storeContainer == nil {
		dbLog := waLog.Stdout("Database", "WARN", true)
		var err error
		storeContainer, err = sqlstore.New(context.Background(), "sqlite", "file:database/whatsapp.db?_pragma=foreign_keys(1)", dbLog)
		if err != nil {
			log.Fatalf("Failed to connect to WhatsApp database: %v", err)
		}
	}
	if err != nil {
		log.Fatalf("Failed to connect to WhatsApp database: %v", err)
	}

	deviceStore, err := storeContainer.GetFirstDevice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get device: %v", err)
	}

	clientLog := waLog.Stdout("Client", "WARN", true)
	if WAClient != nil {
		WAClient.Disconnect()
	}
	
	// Reset CurrentQR before new client
	qrMutex.Lock()
	CurrentQR = ""
	qrMutex.Unlock()

	WAClient = whatsmeow.NewClient(deviceStore, clientLog)
	
	if WAClient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := WAClient.GetQRChannel(context.Background())
		err = WAClient.Connect()
		if err != nil {
			log.Fatalf("Failed to connect to WhatsApp: %v", err)
		}
		
		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					qrMutex.Lock()
					CurrentQR = evt.Code
					qrMutex.Unlock()
				} else {
					fmt.Println("WhatsApp Login event:", evt.Event)
				}
			}
		}()
	} else {
		// Already logged in, just connect
		err = WAClient.Connect()
		if err != nil {
			log.Fatalf("Failed to connect to WhatsApp: %v", err)
		}
		log.Println("WhatsApp already logged in")
	}
}

// SendMessage sends a text message to a specific number (e.g. "62812345678@s.whatsapp.net")
func SendMessage(jid string, message string) error {
	if WAClient == nil || !WAClient.IsLoggedIn() {
		return fmt.Errorf("whatsapp not logged in")
	}

	jidParsed, err := types.ParseJID(jid)
	if err != nil {
		return err
	}
	msg := &waProto.Message{Conversation: proto.String(message)}
	_, err = WAClient.SendMessage(context.Background(), jidParsed, msg)
	return err
}

func GetCurrentQR() string {
	qrMutex.RLock()
	defer qrMutex.RUnlock()
	return CurrentQR
}
