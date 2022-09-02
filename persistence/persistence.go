package persistence

import (
	"context"
	"errors"
	"net/http"

	"github.com/ua-parser/uap-go/uaparser"
)

type auth struct {
	user string
	pwd  string
}

type server struct {
	host string
	port int16
	db   string
	ctx  context.Context
}

type opts struct {
	prefix string
}

type persistence struct {
	server server
	conn   connection
	uaP    *uaparser.Parser
	opts   *opts
}

func New(s server, a auth, opts *opts, uaP *uaparser.Parser) (*persistence, error) {
	if conn := s.open(&a); conn != nil {
		return &persistence{s, connection{conn, a}, uaP, opts}, nil
	}

	return nil, errors.New("connection to the server failed")
}

func Server(host string, port int16, db string) server {
	return server{host, port, db, context.Background()}
}

func Auth(user string, pwd string) auth {
	return auth{user, pwd}
}

func Opts(prefix string) *opts {
	return &opts{prefix}
}

func (p *persistence) Save(req *http.Request, table string) {
	tPrefix := &p.opts.prefix

	if table != "" {
		table = *tPrefix + "_" + table
	}

	p.server.save(&p.conn, eventPayload(req, p.uaP), table)
}
