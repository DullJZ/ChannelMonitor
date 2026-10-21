package main

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestFetchChannelsDoesNotLogAPIKeys(t *testing.T) {
	previousDB, previousConfig := db, config
	previousOutput, previousFlags, previousPrefix := log.Writer(), log.Flags(), log.Prefix()
	t.Cleanup(func() {
		db, config = previousDB, previousConfig
		log.SetOutput(previousOutput)
		log.SetFlags(previousFlags)
		log.SetPrefix(previousPrefix)
	})

	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(&Channel{}); err != nil {
		t.Fatal(err)
	}
	const apiKey = "fake-channel-api-key-for-regression-test"
	channel := Channel{ID: 7, Type: 1, Name: "test-channel", Key: apiKey, Status: 1}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	config = &Config{}
	var output bytes.Buffer
	log.SetOutput(&output)
	log.SetFlags(0)
	log.SetPrefix("")

	channels, err := fetchChannels()
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 1 || channels[0].Key != apiKey {
		t.Fatal("channel credentials must remain available for requests")
	}
	if strings.Contains(output.String(), apiKey) {
		t.Fatal("channel discovery log exposes the API key")
	}
	if !strings.Contains(output.String(), channel.Name) || !strings.Contains(output.String(), "ID:7") {
		t.Fatal("channel discovery log should retain the channel name and ID")
	}
}
