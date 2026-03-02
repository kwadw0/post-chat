package main

import "log"

func main() {
	cfg := config {
		addr: ":8000",
		db: dbConfig{},
	}
	api := application{
		config: cfg,
	}
  	api.run(api.mount())

	if err := api.run((api.mount())); err != nil {
		log.Printf("Application failed to start, err: %s", err)
	}
}






