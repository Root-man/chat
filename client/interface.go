package client

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/root-man/chat/packets"
)

func StartInterface() {
	var c *Client
	app := tview.NewApplication()
	connectForm := tview.NewForm()

	connectForm.
		AddInputField("username", "", 20, nil, nil).
		AddButton("Connect", func() {
			username := connectForm.GetFormItemByLabel("username").(*tview.InputField).GetText()
			c = New(username)

			msgChan, err := c.Connect("localhost", 4444)
			if err != nil {
				log.Printf("Failed to init the client: %s", err)
				os.Exit(2)
			}

			renderChatView(app, c, msgChan)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	connectForm.SetBorder(true).SetTitle("Connect to the chat server").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(connectForm, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func renderChatView(app *tview.Application, c *Client, msgChan <-chan packets.Message) {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	button := tview.NewButton("Quit").SetSelectedFunc(func() {
		app.Stop()
	})

	headerText := newPrimitive("Welcome to mega chat! you are connected as " + c.name)

	header := tview.NewFlex().
		AddItem(headerText, 0, 1, false).
		AddItem(button, 20, 1, false)

	chatBox := tview.NewTextView().SetTextAlign(tview.AlignLeft).SetDynamicColors(true)

	inputField := tview.NewInputField()
	inputField.SetLabel("Message: ").SetFieldWidth(0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			msgText := inputField.GetText()
			msg, err := c.Send(msgText)
			if err != nil {
				// do something
			}
			chatBox.Write(chatViewOwnMsgFormat(msg))
			inputField.SetText("")
		}
	})

	grid := tview.NewGrid().
		SetRows(1, 0, 2).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(chatBox, 1, 0, 1, 3, 0, 0, false).
		AddItem(inputField, 2, 0, 1, 3, 0, 0, false)

	app.SetRoot(grid, true).SetFocus(inputField).Sync()

	// Goroutine to receive messages
	go func() {
		for msg := range msgChan {
			app.QueueUpdateDraw(func() {
				chatBox.Write(chatViewMsgFormat(&msg))
			})
		}
	}()
}

func chatViewMsgFormat(m *packets.Message) []byte {
	return []byte("[blue]" + m.Timestamp.Format(time.DateTime) + " [yellow]" + m.From + "[white]: " + m.Payload + "\n")
}

func chatViewOwnMsgFormat(m *packets.Message) []byte {
	return []byte("[blue]" + m.Timestamp.Format(time.DateTime) + " [yellow]Me" + "[white]: " + m.Payload + "\n")
}
