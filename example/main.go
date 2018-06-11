package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/y4v8/errors"
	"github.com/y4v8/train"
)

func main() {
	args, err := getArgs()
	if err != nil {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}
	monitor(args.fromName, args.toName, args.trainNumber, args.wagonTypeID, args.date, args.intPlaces, 60*time.Second)
}

type Args struct {
	fromName    string
	toName      string
	trainNumber string
	wagonTypeID string
	date        string
	strPlaces   string
	intPlaces   []int
}

func getArgs() (*Args, error) {
	var args Args

	flag.StringVar(&args.fromName, "from", "", "Марганец")
	flag.StringVar(&args.toName, "to", "", "Харьков")
	flag.StringVar(&args.trainNumber, "train", "", "074О")
	flag.StringVar(&args.wagonTypeID, "type", "", "П")
	flag.StringVar(&args.date, "date", "", "2018-06-21")
	flag.StringVar(&args.strPlaces, "places", "", "17,19 [not required]")
	flag.Parse()

	if args.fromName == "" || args.toName == "" || args.trainNumber == "" || args.wagonTypeID == "" || args.date == "" {
		return nil, errors.New("Bad arguments")
	}

	strPlaces := strings.Split(args.strPlaces, ",")
	args.intPlaces = make([]int, 0, len(strPlaces))
	for _, s := range strPlaces {
		if s == "" {
			continue
		}

		intPlace, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.New("Bad arguments")
		}

		args.intPlaces = append(args.intPlaces, intPlace)
	}

	return &args, nil
}

func monitor(fromName, toName, trainNumber, wagonTypeID, date string, places []int, repeatTime time.Duration) {
	log.Printf("[%v] %v->%v %v %v %v\n", trainNumber, fromName, toName, wagonTypeID, date, places)

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

		freePlaces := "[Free] "
		for _, wagonType := range dataTrainWagons.Types {
			freePlaces += wagonType.TypeID + ":" + strconv.Itoa(wagonType.Free) + " "
		}
		log.Println(freePlaces)

		for _, wagon := range dataTrainWagons.Wagons {
			if wagon.TypeID == paramsTrainWagons.WagonTypeID && wagon.Free > 0 {
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

				if len(places) == 0 {
					exists = true
					freePlaces = append(freePlaces, place)
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
