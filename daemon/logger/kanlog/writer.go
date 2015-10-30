// 1. establish the connection, and get the writer_handler
// 2. connect using the established connection
// 3. send log message to the daemon...with priority-tag or not
// 4. close the connection



type Writer struct {
	mu               sync.Mutex
	conn             net.Conn
	hostname         string
}

type Message struct {
	Short    string
}


addr :="localhost:20319"

// New returns a new kanlog Writer.  This writer can be used to send the
// output of the standard Go log functions to a central kanlog server by
// passing it to log.SetOutput()
func NewWriter(addr string) (*Writer, error) {
	var err error
	w := new(Writer)

	if w.conn, err = net.Dial("tcp", addr); err != nil {
		return nil, err
	}
	if w.hostname, err = os.Hostname(); err != nil {
		return nil, err
	}

	return w, nil
}


// WriteMessage sends the specified message to the kanlog server
// specified in the call to New().  It assumes all the fields are
// filled out appropriately.  So you may want to come up with 
//something that does some type-/field-checking.
func (w *Writer) WriteMessage(m *Message) (err error) {

	n, err := fmt.Fprintf(w.conn, m + "\n")
	if err != nil {
		return
	}
	return nil
}




