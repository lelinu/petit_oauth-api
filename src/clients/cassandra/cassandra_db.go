package cassandra

import (
	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
)

func init(){
	// Connect to Cassandra cluster:
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	// Set session as a local variable
	var err error
	session, err = cluster.CreateSession()
	if err != nil{
		panic(err)
	}
}

func GetSession() *gocql.Session{
	return session
}
