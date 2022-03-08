package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
	"guarantorplace.com/internal/data"
)

func(app *application) requestForTransaction(ch chan bool, msg string, errs chan error) {
	client := &http.Client{}
	requestLink := fmt.Sprintf("https://toncenter.com/api/v2/getTransactions?address=%s&limit=25&to_lt=0&archival=false&api_key=%s", app.config.Wallet, app.config.ApiKey)
	req, err := http.NewRequest(http.MethodGet, requestLink, nil)
	
	
	if err != nil {
		errs <- err
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		errs <- err
	}

	app.logger.PrintInfo("Sending a request to the ton.center to check the existence of a transaction", nil)

	var res data.Response

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&res)

	var exist bool

	for _, item := range res.Result {

			value, err := strconv.Atoi(item.InMessage.Value)
			if err != nil {
				errs <- err
			}

			if item.InMessage.Message == msg && math.Floor(float64(value)*100)/100000000000 == app.config.Fee {
				app.logger.PrintInfo("Transaction successfully found", nil)
				exist = true
				ch <- true
			}
	}

	if !exist {
		ch <- false
		app.logger.PrintInfo("Transaction not found", nil)
	}

}


func generateUniqueMessage() (string, error) {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))	  
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	bytes := make([]byte, 16)
 
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = chars[seededRand.Intn(len(chars))]
	}
 
	return string(bytes), nil
}


func(app *application) checkTransaction(done chan bool, b *tele.Bot, link string, chat *tele.Chat, message string, errs chan error) {

	_, err := b.Send(chat, app.config.Messages.AfterPaymentResponse, &tele.SendOptions{ParseMode: "MarkdownV2"})
//delete configs

	if err != nil {
		errs <- err
	}

	ch := make(chan bool)

	ticker := time.NewTicker(3 * time.Second)

	ad, err := app.models.Ads.GetByMessage(message)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.logger.PrintInfo("No record in DB", nil)
			errs <- err
		default:
			errs <- err
		}
	}
	
	out:for range ticker.C {
		select {
			case <- done:
				_, err := b.Send(chat, app.config.Messages.FailedPaymentResponse, &tele.SendOptions{ParseMode: "MarkdownV2"})

				if err != nil {
					errs <- err
				}

				ticker.Stop()
				break out
			case status := <-ch: 
				if status && !ad.Paid {
				target, err := b.ChatByID(app.config.ExchangeChannel)

				if err != nil {
					errs <- err
				}


				_, err = b.Send(target, link + "\n\n" + fmt.Sprintf(app.config.Messages.Signature, app.config.RequestLink))
				if err != nil {
					errs <- err
				}
				
	
				input := &data.Ad{
				   	UserId:	ad.UserId,
					Msg: ad.Msg,
					Link: ad.Link,
					Paid: true,
					CreatedAt: ad.CreatedAt,
					ID: ad.ID,
					Version: ad.Version,
				}
			
				err = app.models.Ads.Update(input)

				app.logger.PrintInfo("Update paid status to true in database", nil)
				
				if err != nil {
					switch {
					case errors.Is(err, data.ErrEditConflict):
						errs <- err
					default:
						errs <- err
					}
				}

				_, err = b.Send(chat, app.config.Messages.SuccessPaymentResponse, &tele.SendOptions{ParseMode:"MarkdownV2"})

				if err != nil {
					errs <- err
				}
		
			
				ticker.Stop()
			}else if ad.Paid {
				_, err := b.Send(chat, app.config.Messages.AlreadyPaid, &tele.SendOptions{ParseMode: "MarkdownV2"})

				if err != nil {
					errs <- err
				}

				ticker.Stop()
			}else {
				continue
			}
		case <-ticker.C:
				go app.requestForTransaction(ch, message, errs)
		}
	
	}

}


func (app *application) backToMainMenu(btnBack *tele.Btn, c tele.Context, mainMenu *tele.ReplyMarkup, photo *tele.Photo, text string) {
	app.bot.Handle(btnBack, func(c tele.Context) error {
		c.Delete()
		app.bot.Handle(tele.OnText, func(c tele.Context) error {
			return nil;
		})		

		if photo == nil {
		return c.Send(text, mainMenu)
		}else {
		return c.Send(photo, &tele.SendOptions{ParseMode: "MarkdownV2", ReplyMarkup: mainMenu})
		}
	})
}