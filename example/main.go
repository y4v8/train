package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/y4v8/errors"
	"github.com/y4v8/train"
)

func main() {
	places := []int{13, 15, 17, 19, 21, 23}
	monitor("Константиновка", "Киев", "126О", "П", "2018-04-28", places, 60*time.Second)
}

func monitor(fromName, toName, trainNumber, wagonTypeID, date string, places []int, repeatTime time.Duration) {
	api := train.NewApi()

	from, err := getStation(api, fromName)
	panicIf(err)

	to, err := getStation(api, toName)
	panicIf(err)

	paramsTrainWagons := train.ParamsTrainWagons{
		From:        from,
		To:          to,
		Date:        date,
		Train:       trainNumber,
		WagonTypeID: wagonTypeID,
	}

	timer := time.NewTicker(repeatTime)

	for {
		if printFreePlaceWagons(api, paramsTrainWagons, places) {
			break
		}

		<-timer.C
	}
}

func printFreePlaceWagons(api *train.Api, paramsTrainWagons train.ParamsTrainWagons, places []int) bool {
	exists := false

	dataTrainWagons, err := api.TrainWagons(paramsTrainWagons)
	if err != nil {
		log.Println(err)
	} else if len(dataTrainWagons.Wagons) > 0 {
		paramsTrainWagon := train.ParamsTrainWagon{
			From:      paramsTrainWagons.From,
			To:        paramsTrainWagons.To,
			Train:     paramsTrainWagons.Train,
			Date:      paramsTrainWagons.Date,
			WagonType: paramsTrainWagons.WagonTypeID,
		}

		for _, wagon := range dataTrainWagons.Wagons {
			if wagon.Free > 0 {
				paramsTrainWagon.WagonNum = wagon.Num
				paramsTrainWagon.WagonClass = wagon.Class
				exists = printFreePlaceWagon(api, paramsTrainWagon, wagon.Free, places)
			}
		}
	}
	return exists
}

func printFreePlaceWagon(api *train.Api, paramsTrainWagon train.ParamsTrainWagon, free int, places []int) bool {
	exists := false

	dataTrainWagon, err := api.TrainWagon(paramsTrainWagon)
	if err != nil {
		log.Printf("wagon:%v, free:%v, error:%v", paramsTrainWagon.WagonNum, free, err)
	} else {
		for _, group := range dataTrainWagon.Places {
			var freePlaces []string
			for _, place := range group {
				freePlace, err := strconv.Atoi(place)
				if err != nil {
					continue
				}
				for _, checkPlace := range places {
					if checkPlace == freePlace {
						exists = true
						freePlaces = append(freePlaces, place)
					}
				}
			}
			log.Printf("wagon:%v, free:%v, places:%v", paramsTrainWagon.WagonNum, free, strings.Join(freePlaces, ","))
		}
	}

	return exists
}

func getStation(api *train.Api, name string) (string, error) {
	paramsStation := train.ParamsStation{
		Term: name,
	}
	dataStations, err := api.Station(paramsStation)
	if err != nil {
		return "", errors.Wrap(err)
	}

	if len(dataStations) < 1 {
		return "", errors.New("station '%v' is not found", name)
	}

	return strconv.Itoa(dataStations[0].Value), nil
}

func panicIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
