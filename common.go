package main

const (
	TelegraphLink = "https://telegra.ph"
)

func Url(link string) string {
	return TelegraphLink + link
}

func checkError(err error) {
	if err != nil {
		logger.Panic(err.Error())
	}
}
