package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type EnvConfig struct {
	FieldWidth  uint
	FieldHeight uint

	FoodStatic uint
	FoodProb   uint

	Dt        uint
	ChunkSize uint
}

func parseUint(varName string) (uint, error) {
	value, err := strconv.Atoi(os.Getenv(varName))
	if err != nil {
		return 0, err
	}
	if value <= 0 {
		return 0, errors.New(varName + "<=0")
	}
	return uint(value), nil

}

func EnvParse() EnvConfig {
	var env EnvConfig
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	// width
	env.FieldWidth, err = parseUint("FIELD_WIDTH")
	if err != nil {
		logrus.Fatal("Error laoding .env file: ", err.Error())
	}
	// height
	env.FieldHeight, err = parseUint("FIELD_HEIGHT")
	if err != nil {
		logrus.Fatal("Error laoding .env file: ", err.Error())
	}
	// dt
	env.Dt, err = parseUint("DT")
	if err != nil {
		logrus.Fatal("Error laoding .env file: ", err.Error())
	}
	// ChunkSize
	env.ChunkSize, err = parseUint("CHUNK_SIZE")
	if err != nil {
		logrus.Fatal("Error laoding .env file: ", err.Error())
	}
	// food static
	env.FoodStatic, err = parseUint("FOOD_STATIC")
	if err != nil {
		logrus.Fatal("Error laoding .env file: ", err.Error())
	}
	// food prob
	env.FoodProb, err = parseUint("FOOD_PROB")
	if err != nil {
		logrus.Fatal("Error laoding .env file: ", err.Error())
	}
	return env
}
