package client

import (
	"context"
	"log"
	"math/rand"
	"sync"

	"github.com/JUSSIAR/Golang-gRPC-blog/service/api"
)

const (
	COUNT_CREATES             = 100
	COUNT_READS               = 1000
	MESSAGE_COUNT_CREATES_OK  = "Create OK reqs"
	MESSAGE_COUNT_CREATES_BAD = "Create BAD reqs"
	MESSAGE_COUNT_READS_OK    = "Read OK reqs"
	MESSAGE_COUNT_READS_BAD   = "Read BAD reqs"
)

type BlogServerClient struct {
	api.BlogServerClient
}

func generateRandomBool() bool {
	return rand.Int()%2 == 0
}

func generateRandomString() string {
	len := int(rand.Uint32() % 256)
	line := ""
	for i := 0; i < len; i++ {
		line += string(byte('a' + (rand.Uint32() % 26)))
	}
	return line
}

func generateComment(ctx context.Context, client *BlogServerClient, wg *sync.WaitGroup, ch chan<- bool) {
	defer wg.Done()

	_, err := client.CreateComment(ctx, &api.CommentRequestCreate{
		AuthorLogin: "jussiar",
		Text:        generateRandomString(),
		PostId:      rand.Uint64(),
	})

	if err == nil {
		ch <- true
	} else {
		ch <- false
	}
}

func generatePost(ctx context.Context, client *BlogServerClient, wg *sync.WaitGroup, ch chan<- bool) {
	defer wg.Done()

	_, err := client.CreatePost(ctx, &api.PostRequestCreate{
		AuthorLogin: "jussiar",
		Text:        generateRandomString(),
		Title:       generateRandomString(),
	})

	if err == nil {
		ch <- true
	} else {
		ch <- false
	}
}

func readComment(ctx context.Context, client *BlogServerClient, wg *sync.WaitGroup, ch chan<- bool) {
	defer wg.Done()

	_, err := client.ReadComment(ctx, &api.CommentRequestReadDelete{
		Id: rand.Uint64(),
	})

	if err == nil {
		ch <- true
	} else {
		ch <- false
	}
}

func readPost(ctx context.Context, client *BlogServerClient, wg *sync.WaitGroup, ch chan<- bool) {
	defer wg.Done()

	_, err := client.ReadPost(ctx, &api.PostRequestReadDelete{
		Id: rand.Uint64(),
	})

	if err == nil {
		ch <- true
	} else {
		ch <- false
	}
}

func logStats(ch chan bool, messageOk string, messageBad string) {
	countOk := 0
	countBad := 0
	for flag := range ch {
		if flag {
			countOk++
		} else {
			countBad++
		}
	}

	log.Println(messageOk, countOk)
	log.Println(messageBad, countBad)
	close(ch)
}

func main() {
	ctx := context.Background()
	client := &BlogServerClient{}

	chanWriter := make(chan bool, COUNT_CREATES)
	var wgCreator sync.WaitGroup
	for i := 0; i < COUNT_CREATES; i++ {
		wgCreator.Add(1)
		if generateRandomBool() {
			go generateComment(ctx, client, &wgCreator, chanWriter)
		} else {
			go generatePost(ctx, client, &wgCreator, chanWriter)
		}
	}
	wgCreator.Wait()

	chanReader := make(chan bool, COUNT_READS)
	var wgReader sync.WaitGroup
	for i := 0; i < COUNT_CREATES; i++ {
		wgCreator.Add(1)
		if generateRandomBool() {
			go readComment(ctx, client, &wgReader, chanReader)
		} else {
			go readPost(ctx, client, &wgReader, chanReader)
		}
	}
	wgReader.Wait()

	go logStats(chanWriter, MESSAGE_COUNT_CREATES_OK, MESSAGE_COUNT_CREATES_BAD)
	go logStats(chanReader, MESSAGE_COUNT_READS_OK, MESSAGE_COUNT_READS_BAD)
}
