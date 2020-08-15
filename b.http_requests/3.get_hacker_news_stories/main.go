package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/caser/gophernews"
)

var hackerNewsClient *gophernews.Client

func init() {
	hackerNewsClient = gophernews.NewClient()
}

type Story struct {
	title  string
	url    string
	author string
	source string
}

func newHnStories(c chan<- Story) {

	defer close(c)

	changes, err := hackerNewsClient.GetChanges()
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup

	for _, id := range changes.Items {
		wg.Add(1)
		go getHNStoryDetails(id, c, &wg)
	}
	wg.Wait()

}
func getHNStoryDetails(id int, c chan<- Story, wg *sync.WaitGroup) {
	defer wg.Done()

	story, err := hackerNewsClient.GetStory(id)
	if err != nil {
		return
	}

	newStory := Story{
		title:  story.Title,
		url:    story.URL,
		author: story.By,
		source: "HackerNews",
	}
	c <- newStory
}

func main() {
	fromHn := make(chan Story, 8)
	toPrint := make(chan Story, 8)
	toFile := make(chan Story, 8)

	go newHnStories(fromHn)

	file, err := os.Create("stories.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	go outputToConsole(fromHn)
	go outputToFile(toFile, file)

	hnOpen := true
	for hnOpen {
		select {
		case story, open := <-fromHn:
			if open {
				toFile <- story
				toPrint <- story
			} else {
				hnOpen = false
			}
		}
	}
}

func outputToFile(c chan Story, file *os.File) {
	for {
		s := <-c
		fmt.Fprintf(file, "%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

func outputToConsole(c <-chan Story) {
	var s Story
	more := true
	for more {
		select {
		case s, more = <-c:
			if more {
				fmt.Printf("%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
			} else {
				return
			}
		}
	}
}
