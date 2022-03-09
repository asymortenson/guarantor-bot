package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"guarantorplace.com/internal/data"
)



func (app *application) handleApproveChannel() error {

	var (
		selector = &tele.ReplyMarkup{}
	
		btnApprove = selector.Data(app.config.Buttons.ApprovePublicPage, "approve")
		btnDecline = selector.Data(app.config.Buttons.DeclinePublicPage, "decline")
	
	)	



	
	adminBoard, err := app.bot.ChatByID(app.config.AdminChannel)


	selector.Inline(
		selector.Row(btnApprove, btnDecline),
	)		
	
	var publicChannelForApprove string


	
	app.bot.Handle(tele.OnText, func(c tele.Context) error {

		if c.Message().Sender.Username  == ""	{
			publicChannelForApprove = fmt.Sprintf("%d\n\n%s", c.Message().Sender.ID, c.Text())
		_, err = app.bot.Send(adminBoard, publicChannelForApprove, selector)

		if err != nil {
			return err
		}

		}else {
			publicChannelForApprove = fmt.Sprintf("%d\n%s\n\n%s", c.Message().Sender.ID, "@" + c.Message().Sender.Username, c.Text())
			_, err = app.bot.Send(adminBoard, publicChannelForApprove, selector)
			if err != nil {
				return err
			}
	
		}

	
		_, err := app.bot.Send(c.Sender(), app.config.Messages.AfterSubmittingPublicPage, &tele.SendOptions{ParseMode:"MarkdownV2"})

		if err != nil {
			return err
		}

		app.bot.Handle(tele.OnText, func(c tele.Context) error {
			return nil
		})
		return nil
	})

	
	app.bot.Handle(&btnApprove, func(c tele.Context) error {
		var (
			approvedRow = &tele.ReplyMarkup{}
			btnPaid = selector.Data(app.config.Buttons.Paid, "paid")
			btnDeclinePaid = selector.Data(app.config.Buttons.DeclinePaid, "decline_paid")
			btnApproved = selector.Data(app.config.Buttons.Approved, "approved")
		)

		selector.Inline(
			selector.Row(btnPaid),
			selector.Row(btnDeclinePaid),
		)

		approvedRow.Inline(
			approvedRow.Row(btnApproved),
		)

		c.Edit(publicChannelForApprove, approvedRow)


		arrayOfString := strings.Split(c.Text(), "\n")
		

		id, err := strconv.ParseInt(arrayOfString[0], 10, 64)

		if err != nil {
			return err
		}

		textForSend := strings.Split(c.Text(), "\n\n")



		chat, err := app.bot.ChatByID(id)


		if err != nil {
			return err
		}

		message, err := generateUniqueMessage()

		if err != nil {
			return err
		}

		ad := &data.Ad{
				UserId:   id,
				Link: strings.Join(textForSend[1:], ""),
				Msg: message,
		}


		err = app.models.Ads.Insert(ad)
		if err != nil {
			return err
		}



		paymentMessage := fmt.Sprintf(app.config.Messages.PaymentMessage, app.config.Buttons.Paid, "0\\.\\2",  app.config.Wallet, message)
		
		_, err = app.bot.Send(chat, paymentMessage, &tele.SendOptions{ParseMode: "MarkdownV2", ReplyMarkup: selector})
		
		if err != nil {
			return err
		}

		done := make(chan bool) 


		app.bot.Handle(&btnPaid, func(c tele.Context) error {

			var (
				paidRow = &tele.ReplyMarkup{}
				btnDeclineAfterPayment = selector.Data(app.config.Buttons.DeclinePaid, "decline_after_payment")
			)
	
			paidRow.Inline(
				paidRow.Row(btnDeclinePaid),
			)
	
			c.Edit(paymentMessage, &tele.SendOptions{ParseMode: "MarkdownV2", ReplyMarkup: paidRow})
	


			ad, err := app.models.Ads.Get(c.Chat().ID)


			if err != nil {
				return err
			}
			

			errs := make(chan error, 1)

			go app.checkTransaction(done, app.bot, ad.Link, c.Chat(), ad.Msg, errs)
			
			time.Sleep(time.Minute * 5)
			close(done)

			app.bot.Handle(&btnDeclineAfterPayment, func(c tele.Context) error {	
				close(done)
	
				return nil
			})
			
	

			if err := <-errs; err != nil {
				return err
			}
			
			return nil

		})


		

		app.bot.Handle(&btnDeclinePaid, func(c tele.Context) error {	

			var (
				declineRow = &tele.ReplyMarkup{}

				btnDeclinedResponse = selector.Data(app.config.Buttons.DeclinePaid, "decline_after_payment"))

				declineRow.Inline(
					declineRow.Row(btnDeclinedResponse),
				)
		
			c.Edit(paymentMessage, &tele.SendOptions{ParseMode: "MarkdownV2", ReplyMarkup: declineRow})
			return c.Send(app.config.Messages.Responses.FailedPaymentResponse, &tele.SendOptions{ParseMode: "MarkdownV2"})
		})



		

		return nil
	})

	app.bot.Handle(&btnDecline, func(c tele.Context) error {

		var (
			rejectedRow = &tele.ReplyMarkup{}
			btnRejected = rejectedRow.Data(app.config.Buttons.Rejected, "rejected")
		)
		rejectedRow.Inline(
			rejectedRow.Row(btnRejected),
		)


		c.Edit(publicChannelForApprove, rejectedRow)

		id, err := strconv.ParseInt(strings.Split(c.Text(), "\n")[0], 10, 64)

		if err != nil {
			return err
		}

		chat, err := app.bot.ChatByID(id)


		if err != nil {
			return err
		}
  
		_, err = app.bot.Send(chat, app.config.Messages.RejectPublicPage, &tele.SendOptions{
			ParseMode: "MarkdownV2",
		})

		if err != nil {
			return err
		}

		return nil
	})

	return nil

}

func (app *application) handleStartCommand() error {

	var (
		mainMenu = &tele.ReplyMarkup{}
		backToMenu = &tele.ReplyMarkup{}
	    mainMenuPhoto = &tele.Photo{File: tele.FromURL("https://ibb.co/D8h9XKN")}

		btnBack = backToMenu.Data(app.config.Buttons.BackToPrevious, "back")
		btnChoosePlace = mainMenu.Data(app.config.Buttons.ChoosePublicPage, "place")
		btnCreateRequest = mainMenu.Data(app.config.Buttons.CreateRequest, "create_request")
	)

	backToMenu.Inline(
		backToMenu.Row(btnBack),
	)
	
	mainMenu.Inline(
		mainMenu.Row(btnCreateRequest),
		mainMenu.Row(btnChoosePlace),
	)
	
	
	app.bot.Handle("/start", func(c tele.Context) error {
		return c.Send(mainMenuPhoto, mainMenu)
	})

	app.bot.Handle(&btnChoosePlace, func(c tele.Context) error {
		var (
			publicPagesMenu = &tele.ReplyMarkup{}
			btnPublicPageProgrammer = publicPagesMenu.Data(app.config.PublicPages.Programmer, "programmer")
			btnPublicPageAboutTon = publicPagesMenu.Data(app.config.PublicPages.AboutTon, "about_ton")
			btnBackToMainMenu = publicPagesMenu.Data(app.config.Buttons.BackToPrevious, "btn_back_to_main_menu")
			choosePublicPagePhoto = &tele.Photo{File: tele.FromURL("https://i.ibb.co/3MWWymj/Screenshot-2022-03-06-at-11-58-49.png"), Caption: "📣 *Выберите сообщество для покупки рекламы:*"}
		)

		c.Delete()

		app.backToMainMenu(&btnBackToMainMenu, c, mainMenu, mainMenuPhoto, "")		

		publicPagesMenu.Inline(
			publicPagesMenu.Row(btnPublicPageProgrammer),
			publicPagesMenu.Row(btnPublicPageAboutTon),
			publicPagesMenu.Row(btnBackToMainMenu),
		)

		app.bot.Handle(&btnPublicPageProgrammer, func(c tele.Context) error {
			var (
				publicPageInfoMenu = &tele.ReplyMarkup{}
				btnBuy = publicPageInfoMenu.URL(app.config.Buttons.BuyAd, "https://telegra.ph/Prajs-list-dlya-soobshchestva-Programmist-03-06")
				btnBackToPublicPages = publicPageInfoMenu.Data(app.config.Buttons.BackToPrevious, "btn_back_to_public_pages")
				programmerText = "👥 *Сообщество*: \n• ton\\.\\place/group10316 \n\n☎️ *Контакты*: \n• ton\\.\\place/id39469\n• @goreactdev"
				choosePublicPagePhoto = &tele.Photo{File: tele.FromURL("https://i.ibb.co/3MWWymj/Screenshot-2022-03-06-at-11-58-49.png"), Caption: "📣 *Выберите сообщество для покупки рекламы:*"}

				programmerMenuPhoto = &tele.Photo{File: tele.FromURL("https://ibb.co/4SsZPmv"), Caption: programmerText}
			)

			publicPageInfoMenu.Inline(
				publicPageInfoMenu.Row(btnBuy),
				publicPageInfoMenu.Row(btnBackToPublicPages),
			)

			app.backToMainMenu(&btnBackToPublicPages, c, publicPagesMenu, choosePublicPagePhoto, "")
			
			c.Delete()

			return c.Send(programmerMenuPhoto, &tele.SendOptions{ReplyMarkup: publicPageInfoMenu, ParseMode: "MarkdownV2"})
		})

		app.bot.Handle(&btnPublicPageAboutTon, func(c tele.Context) error {
			var (
				publicPageInfoMenu = &tele.ReplyMarkup{}
				btnBuy = publicPageInfoMenu.URL(app.config.Buttons.BuyAd, "https://telegra.ph/price-03-06-2")
				btnBackToPublicPages = publicPageInfoMenu.Data(app.config.Buttons.BackToPrevious, "btn_back_to_public_pages")
				programmerText = "👥 *Сообщество*: \n• ton\\.\\place/group15 \n\n☎️ *Контакты*: \n• ton\\.\\place/math\\_is\n• @math\\_is"
				choosePublicPagePhoto = &tele.Photo{File: tele.FromURL("https://i.ibb.co/3MWWymj/Screenshot-2022-03-06-at-11-58-49.png"), Caption: "📣 *Выберите сообщество для покупки рекламы:*"}
				programmerMenuPhoto = &tele.Photo{File: tele.FromURL("https://ibb.co/8jDHqH1"), Caption: programmerText}
			)

			publicPageInfoMenu.Inline(
				publicPageInfoMenu.Row(btnBuy),
				publicPageInfoMenu.Row(btnBackToPublicPages),
			)

			app.backToMainMenu(&btnBackToPublicPages, c, publicPagesMenu, choosePublicPagePhoto, "")
			
			c.Delete()

			return c.Send(programmerMenuPhoto, &tele.SendOptions{ReplyMarkup: publicPageInfoMenu, ParseMode: "MarkdownV2"})
		})

		
		return c.Send(choosePublicPagePhoto, &tele.SendOptions{ReplyMarkup: publicPagesMenu, ParseMode: "MarkdownV2"})
	})

	app.bot.Handle(&btnCreateRequest, func(c tele.Context) error {

		var (
			putPublicPagePhoto = &tele.Photo{File: tele.FromURL("https://i.ibb.co/CMvLWDL/Screenshot-2022-03-06-at-11-57-43.png"), Caption: app.config.Messages.PutPublicPage}
		)

		c.Delete()

		_, err := app.bot.Send(c.Sender(), putPublicPagePhoto, &tele.SendOptions{ParseMode:"MarkdownV2", ReplyMarkup: backToMenu})
		if err != nil {
			return err
		}
		
		app.backToMainMenu(&btnBack, c, mainMenu, mainMenuPhoto , "")		

			
		err = app.handleApproveChannel()

		if err != nil {
			return err
		}

		return nil

	})
	return nil


}

func (app *application) handleMailingCommand() error {
	adminOnly := app.bot.Group()

	chatIds := []int64{2086435608, 359499553}

	adminOnly.Use(middleware.Whitelist(chatIds...))

	adminOnly.Handle("/mailing", func(c tele.Context) error {
		return c.Send("Введите текст для массовой рассылки пользователям:")
	})

	adminOnly.Handle(tele.OnText, func(c tele.Context) error {

		users, err := app.models.Ads.GetAll()

		if err != nil {
			return err
		}

		for index, user := range users {
			if index % 25 == 0 {
				time.Sleep(time.Second * 5)
			}

			chat, err := app.bot.ChatByID(user.UserId)

			if err != nil {
				return err
			}
			_, err = app.bot.Send(chat, c.Text())

			if err != nil {
				return err
			}
		}
		return nil
	})

	return nil

}


func (app *application) handleUpdates() error {
	app.bot.Handle(tele.OnMyChatMember, func(c tele.Context) error {
		if c.ChatMember().NewChatMember.Role == "member" {
		user := data.User{
			ID: c.ChatMember().NewChatMember.User.ID,
		}
		err := app.models.Users.Insert(&user)
		if err != nil {
			return err
		}
		return nil
		}else {
			err := app.models.Users.Delete(c.ChatMember().OldChatMember.User.ID)
			if err != nil {
				return err
			}
			return nil
		}
	})


	err := app.handleMailingCommand()

	if err != nil {
		return err
	}

	err = app.handleStartCommand()
	

	if err != nil {
		return err
	}
	return nil
}



// func (app *application) publicPages(btn *tele.Btn, publicPagesMenu *tele.ReplyMarkup, choosePublicPagePhoto *tele.Photo, public data.Public) {
// 		app.bot.Handle(&btn, func(c tele.Context) error {
// 			var (
// 				publicPageInfoMenu = &tele.ReplyMarkup{}
// 				btnBuy = publicPageInfoMenu.URL(app.config.Buttons.BuyAd, public.TelegraphLink)
// 				btnBackToPublicPages = publicPageInfoMenu.Data(app.config.Buttons.BackToPrevious, "btn_back_to_public_pages")
// 				programmerText = "👥 *Сообщество*: \n• ton\\.\\place/group10316 \n\n☎️ *Контакты*: \n• ton\\.\\place/id39469\n• @goreactdev"

// 				programmerMenuPhoto = &tele.Photo{File: tele.FromURL("https://ibb.co/4SsZPmv"), Caption: programmerText}
// 			)

// 			publicPageInfoMenu.Inline(
// 				publicPageInfoMenu.Row(btnBuy),
// 				publicPageInfoMenu.Row(btnBackToPublicPages),
// 			)

// 			app.backToMainMenu(&btnBackToPublicPages, c, publicPagesMenu, choosePublicPagePhoto, "")
			
// 			c.Delete()

// 			return c.Send(programmerMenuPhoto, &tele.SendOptions{ReplyMarkup: publicPageInfoMenu, ParseMode: "MarkdownV2"})
// 		})

// }
