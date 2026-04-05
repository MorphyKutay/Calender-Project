package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID    string `json:"id"`
	Date  string `json:"date"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type App struct {
	ctx      context.Context
	mu       sync.Mutex
	events   []Event
	dataPath string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	cfg, err := os.UserConfigDir()
	if err != nil {
		cfg = "."
	}
	a.dataPath = filepath.Join(cfg, "calender", "events.json")
	if err := a.load(); err != nil {
		println("events yüklenemedi:", err.Error())
	}
}

func (a *App) load() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	data, err := os.ReadFile(a.dataPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.events = nil
			return nil
		}
		return err
	}
	var list []Event
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	a.events = list
	return nil
}

func (a *App) saveLocked() error {
	dir := filepath.Dir(a.dataPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(a.events, "", "  ")
	if err != nil {
		return err
	}
	tmp := a.dataPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, a.dataPath)
}

func (a *App) GetEvents() ([]Event, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	out := make([]Event, len(a.events))
	copy(out, a.events)
	return out, nil
}

func (a *App) AddEvent(date, title, color string) (Event, error) {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return Event{}, errors.New("tarih YYYY-MM-DD olmalı")
	}
	t := title
	if t == "" {
		return Event{}, errors.New("başlık boş olamaz")
	}
	c := color
	if c == "" {
		c = "#facc15"
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	e := Event{
		ID:    uuid.NewString(),
		Date:  date,
		Title: t,
		Color: c,
	}
	a.events = append(a.events, e)
	if err := a.saveLocked(); err != nil {
		a.events = a.events[:len(a.events)-1]
		return Event{}, err
	}
	return e, nil
}

func (a *App) DeleteEvent(id string) error {
	if id == "" {
		return errors.New("id gerekli")
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	backup := append([]Event(nil), a.events...)
	idx := -1
	for i := range a.events {
		if a.events[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return errors.New("etkinlik bulunamadı")
	}
	a.events = append(a.events[:idx], a.events[idx+1:]...)
	if err := a.saveLocked(); err != nil {
		a.events = backup
		return err
	}
	return nil
}
