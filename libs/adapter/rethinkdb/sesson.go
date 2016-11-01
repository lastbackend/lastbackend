package rethinkdb

import (
	r "gopkg.in/dancannon/gorethink.v2"

)

// Get RethinkDB session
func Get(conf r.ConnectOpts) (*r.Session, error) {

	session, err := r.Connect(conf)
	if err != nil {
		panic(err.Error())
	}
	return session, nil
}
