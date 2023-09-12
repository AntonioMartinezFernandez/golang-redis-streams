package utils

import "math/rand"

func RandomHello() string {
	greeting := []string{
		"Hola redis streams!",
		"Kaixo redis streams!",
		"Hi redis streams!",
		"Konnichiwa redis streams!",
		"Sveika redis streams!",
		"Salut redis streams!",
		"Hallo redis streams!",
		"Ciao redis streams!",
		"Cześć redis streams!",
		"Nǐ hǎo redis streams!",
	}[rand.Intn(10)]

	return greeting
}
