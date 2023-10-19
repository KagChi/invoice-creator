package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/fogleman/gg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang/freetype/truetype"
	"github.com/kagchi/invoice-creator/structs"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"github.com/gofiber/template/html/v2"
)

const (
	width  = 800  // Width of A4 paper in points (21 cm)
	height = 841  // Height of A4 paper in points (29.7 cm)
	margin = 20.0 // Margin in points
)

func drawString(dc *gg.Context, fontSize float64, s string, x float64, y float64, ax float64, ay float64, width float64, align gg.Align) {
	font, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(font, &truetype.Options{
		Size: fontSize,
	})

	dc.SetFontFace(face)
	dc.DrawStringWrapped(s, x, y, ax, ay, width, fontSize, align)
}

func formatNumberToLocaleString(num int, locale string) string {
	tags, _, err := language.ParseAcceptLanguage(locale)
	if err != nil {
		return fmt.Sprintf("%d", num)
	}

	var tag language.Tag

	if len(tags) > 0 {
		tag = tags[0]
	} else {
		tag = language.English
	}

	p := message.NewPrinter(tag)

	return p.Sprintf("%d", num)
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} [${latency}] ${status} - ${method} ${path}\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("home", fiber.Map{
            "Title": "Hello, World!",
        })
    })

	app.Post("/", func(c *fiber.Ctx) error {
		b := new(structs.Payload)

		if err := c.BodyParser(b); err != nil {
			return err
		}

		dc := gg.NewContext(width, height)
		dc.SetRGB(1, 1, 1)
		dc.Clear()

		dc.SetRGB(0, 0, 0)

		drawString(dc, 28.0, "INVOICE", margin, margin, 0, 0, float64(width)-4*margin, gg.AlignRight)                  // Title
		drawString(dc, 16.0, b.Company.Name, margin, 4*margin, 0, 0, float64(width)-4*margin, gg.AlignRight)           // Company Name
		drawString(dc, 16.0, b.Company.Address.Detail, margin, 5*margin, 0, 0, float64(width)-4*margin, gg.AlignRight) //  Address
		drawString(dc, 16.0, b.Company.Address.Zip, margin, 6*margin, 0, 0, float64(width)-4*margin, gg.AlignRight)    // Zip Address
		drawString(dc, 16.0, b.Company.Address.City, margin, 7*margin, 0, 0, float64(width)-4*margin, gg.AlignRight)   // Country

		dc.SetRGB(0, 0, 0)
		lineY := margin + 20.0 + 10
		dc.DrawLine(margin+20, lineY*4, width-(margin+20), lineY*4)
		dc.SetLineWidth(1.0)
		dc.Stroke()

		drawString(dc, 16.0, b.Client.Name, margin + 20, 11*margin, 0, 0, float64(width)-4*margin, gg.AlignLeft)            // Client Name
		drawString(dc, 16.0, b.Client.Address.Detail, margin + 20, 12 * margin, 0, 0, float64(width)-4*margin, gg.AlignLeft)  //  Client Address
		drawString(dc, 16.0, b.Client.Address.Zip, margin + 20, 13 * margin, 0, 0, float64(width)-4*margin, gg.AlignLeft)     // Client Zip Address
		drawString(dc, 16.0, b.Client.Address.Country, margin + 20, 14 * margin, 0, 0, float64(width)-4*margin, gg.AlignLeft) // Client Country

		drawString(dc, 16.0, fmt.Sprintf("Number : %s", b.Invoice.Number), margin + 20, 11 * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight)   // Invoice Number
		drawString(dc, 16.0, fmt.Sprintf("Date  :  %s", b.Invoice.Date), (margin + 20) -87, 12 * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight)  // Invoice Date
		drawString(dc, 16.0, fmt.Sprintf("Due Date :  %s", b.Invoice.Due), (margin + 20)-87, 13 * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight) // Invoice Due Date

		drawString(dc, 16.0, "Products", margin + 20, 18.5 * margin, 0, 0, float64(width) -4 * margin, gg.AlignLeft)
		drawString(dc, 16.0, "Quantity", margin - 210, 18.5 * margin, 0, 0, float64(width) -4 * margin, gg.AlignRight)
		drawString(dc, 16.0, "Price", margin - 105, 18.5 * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight)
		drawString(dc, 16.0, "Total", margin + 8, 18.5 * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight)

		dc.DrawLine(margin+20, lineY*8, width-(margin+20), lineY*8)
		dc.SetLineWidth(1.0)
		dc.Stroke()

		// TODO: use for const of

		currentMargin := 21.0
		currentPrice := 0
		for _, product := range b.Products {
			drawString(dc, 16.0, product.Name, margin+20, currentMargin * margin, 0, 0, float64(width)-4*margin, gg.AlignLeft)

			quantityStr := fmt.Sprintf("%d", product.Quantity)
			quantityLen := 0.0

			if len(quantityStr) <= 3 {
				quantityLen += float64(len(quantityStr)) * 2.7
			} else {
				quantityLen += float64(len(quantityStr)) * 3.5
			}

			drawString(dc, 16.0, quantityStr, margin-(237.5 - quantityLen), currentMargin * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight) // Quantity

			priceStr := fmt.Sprintf("%d", product.Price)
			priceLen := 0.0

			if len(priceStr) <= 1 {
				priceLen += float64(len(quantityStr)) * 14.0
			} else if len(priceStr) > 1 {
				priceLen -= float64(len(quantityStr)) * 10
				priceLen -= float64(len(quantityStr)) * (float64(len(quantityStr)) * 10)
			}

			drawString(dc, 16.0, fmt.Sprintf("Rp %s", formatNumberToLocaleString(product.Price, b.Setting.Locale)), margin-(107 + priceLen), currentMargin * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight) // Quantity
			drawString(dc, 16.0, fmt.Sprintf("Rp %s", formatNumberToLocaleString(product.Price * product.Quantity, b.Setting.Locale)), margin + 20, currentMargin * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight) // Price
			currentPrice += product.Price * product.Quantity
			currentMargin += 1.0
		}

		lineY = currentMargin + 20.0 + 10.0 + float64(len(b.Products)) - 2.0
		dc.DrawLine(margin + 20, lineY * 9.5, width - (margin + 20), lineY * 9.5) 
		dc.SetLineWidth(1.0)
		dc.Stroke()

		drawString(dc, 16.0, fmt.Sprintf("Sub Total : %s %s", b.Setting.Currency, formatNumberToLocaleString(currentPrice, b.Setting.Locale)), margin+20, (currentMargin + 6.5) * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight)                        // Subtotal\
		drawString(dc, 16.0, fmt.Sprintf("Vat %.2f%% :     %s %s", b.Setting.Vat, b.Setting.Currency, formatNumberToLocaleString(int(float64((currentPrice)) * (b.Setting.Vat / 100)), b.Setting.Locale)), margin + 20, (currentMargin + 7.5) * margin, 0, 0, float64(width)-4*margin, gg.AlignRight) // Subtotal
		
		currentMargin += 1.5
		lineY = currentMargin + 20.0 + 10.0 + (float64(len(b.Products)) - 2.0) + 12.5
		dc.DrawLine(margin + 440, lineY * 9.5, width - (margin + 20), lineY * 9.5)
		dc.SetLineWidth(1.0)
		dc.Stroke()


		drawString(dc, 16.0, fmt.Sprintf("Total :   %s %s", b.Setting.Currency, formatNumberToLocaleString(int(float64(currentPrice) * (b.Setting.Vat / 100)) + currentPrice, b.Setting.Locale)), margin + 20, (currentMargin + 7.5) * margin, 0, 0, float64(width) - 4 * margin, gg.AlignRight) // Subtotal

		var buf bytes.Buffer
		_ = dc.EncodePNG(&buf)

		c.Set("Content-Type", "image/png")

		return c.Send(buf.Bytes())
	})

	log.Fatal(app.Listen(":3000"))
}
