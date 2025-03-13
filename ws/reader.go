package ws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func NewWsReader(url string) *Reader {
	return &Reader{
		url:       url,
		retryTime: time.Second,
	}
}

type Reader struct {
	conn      net.Conn
	url       string
	retryTime time.Duration
}

func (p *Reader) Stream(ctx context.Context) (chan []byte, error) {
	ch := make(chan []byte, 1)
	err := p.connect(ctx)
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Break the loop")
				p.conn.Close()
				return
			default:
				msg, err := p.next()
				if err != nil {
					log.Println(err)
					p.reconnect(ctx)
					continue
				}
				if msg == nil {
					continue
				}

				go func() {
					ch <- msg
				}()
			}
		}
	}()

	return ch, nil
}

func (p *Reader) next() ([]byte, error) {
	if p.conn == nil {
		return nil, errors.New("no connection")
	}
	msg, op, err := wsutil.ReadServerData(p.conn)

	if err != nil {
		return nil, err
	}
	if op == ws.OpText || op == ws.OpBinary {
		return msg, err
	}
	return nil, nil
}

func (p *Reader) connect(ctx context.Context) error {
	conn, _, _, err := ws.DefaultDialer.Dial(ctx, p.url)
	if err != nil {
		return err
	}

	p.conn = conn
	p.retryTime = time.Second // Reset retry time on successful connection
	log.Println("Connected successfully")
	return nil
}

func (p *Reader) reconnect(ctx context.Context) {
	t := time.NewTicker(p.retryTime)
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			log.Println("Attempting to reconnect...")

			err := p.connect(ctx)
			if err != nil {
				log.Printf("Failed to reconnect: %v. Retrying in %v...\n", err, p.retryTime)
				p.retryTime *= 2 // Exponential backoff
				if p.retryTime > 30*time.Second {
					p.retryTime = 30 * time.Second // Cap the max retry time
				}
				t.Reset(p.retryTime)
			} else {
				return
			}
		}
	}
}
