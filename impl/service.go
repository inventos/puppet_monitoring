// Service to handle connections from puppet server
package impl

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"
)

type Service struct {
	waitGroup *sync.WaitGroup
	listener  *net.TCPListener
	env       *EnvironmentCollection
}

// Store a pointer to environment collections to call funcs
func (s *Service) SetEnvCollection(envs *EnvironmentCollection) {
	s.env = envs
}

// Wait when all connections will be handled and close listener
func (s *Service) Stop() {
	s.waitGroup.Wait()
	s.listener.Close()
}

// Creates new service instance
func (s Service) NewService() *Service {
	var srv = &Service{
		waitGroup: &sync.WaitGroup{},
	}
	return srv
}

// Handle listener in async way
func (s *Service) HandleListener(listener *net.TCPListener) {
	s.listener = listener
	for {
		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		s.waitGroup.Add(1)
		// process the connection asynchronously
		go s.HandleConnection(conn)
	}
}

const (
	timeout = 5 * time.Second // 5 seconds to process connection
)

// recieving data and send it to Encironment collection to process
func (s *Service) HandleConnection(conn *net.TCPConn) {
	defer conn.Close()
	defer s.waitGroup.Done()
	// set 5 secs timeout (or we will live forever)
	conn.SetDeadline(time.Now().Add(timeout))
	data, _ := ioutil.ReadAll(conn)
	// process report data
	res := s.env.ProcessReport(data)

	// if something goes wrong (not json string, for example), we say to back off
	if !res {
		w := bufio.NewWriter(conn)
		w.WriteString("HTTP/1.1 404 Not Found\r\n\r\n")
		w.Flush()
	}
}
