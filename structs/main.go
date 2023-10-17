package structs

type AddressPayload struct {
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Country string `json:"country"`
	Detail  string `json:"detail"`
}

type Profile struct {
	Name    string `json:"name"`
	Address AddressPayload `json:"address"`
}

type Invoice struct {
	Number string `json:"number"`
	Date   string `json:"date"`
	Due    string `json:"due"`
}

type Setting struct {
	Currency string  `json:"currency"`
	Vat 	 float64 `json:"vat"`
	Locale   string  `json:"locale"`
}

type Product struct {
	Name     string `json:"name"`
	Price    int `json:"price"`
	Quantity int `json:"quantity"`
}

type Payload struct {
	Invoice  Invoice  `json:"Invoice"`
	Company  Profile `json:"company"`
	Client   Profile `json:"client"`
	Setting	 Setting `json:"setting"`
	Products []Product  `json:"products"`
}