


import (
	"os"
	"log"
	"net"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"bytes"
	kanlog "github.com/kmicheal/docker/daemon/logger/kanlog"
)

type kanlogger struct {
	tag		string
	containerID 	string
	containerName	string
	writer		*kanlog.writer
	extra		map[string]string
}


const (
	name = "kanlog"
	defaultHostname = "localhost"
	defaultPort = "20319"
	defaultTagPrefix = "Prefix"
	msg = "we are now logging...."
)

func init (){
	if err := logger.RegisterLogDriver(name, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(name, ValidateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

// Here, New() creates a kanlog logger using the configuration passed in on the context
// supported configuarations are kanlog address and a tag

func New(ctx logger.Context) (logger.Logger, error) {

	// collect info data for the logrdriver's message
	hostname, err = ctx.Hostname()
	if err != nil {	return nil, fmt.Errorf("kanlog: cannot access hostname to set source field") }

	tag, err := loggerutils.ParseLogTag(ctx, "docker.{{.ID}}")
	if err != nil {	return nil, err }
	extra := ctx.ExtraAttributes(nil)
	logrus.Debugf("logging driver kanlog configured for container:%s, host:%s, port:%d, tag:%s, extra:%v.", ctx.ContainerID, hostname, port, tag, extra)



	//create 'writer'.....write where?
	//writer is a daemon smwhere (on the host)
	//client should pass the log content to the server and have the server create the dir_path while
	//connection should remain open until the container pid closes, during which time it continue to send the log to the file 

	address:="tcp://localhost:20319"
	kanwriter, err := *kanlogger.NewWriter(address)
	if err != nil {
		return nil, fmt.Errorf("kanlogger can not connect to kanlogger daemon at %s: %v", address, err)
	}

	return &kanlogger{
		tag:           tag,
		containerID:   ctx.ContainerID,
		containerName: ctx.ContainerName,
		writer:        kanwriter,
		extra:         extra,
	}, nil
}



//other house-keeping functions 

func (s *kanlogger) Log(msg *logger.Message) error {
	if err := s.writer.WriteMessage(msg); err != nil {
		return fmt.Errorf("kanlog: failed to send kanlog message: %v", err)
	}
	return nil
}

func (s *kanlogger) Close() error {
	return s.writer.Close()
}

func (s *kanlogger) Name() string {
	return name
}






// ValidateLogOpt looks for specific log options kanlog-address & kanlog-tag.
// This evaluates the driver log 'Opts' to check if we know them 
func ValidateLogOpt(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case "kanlog-address":
		case "kanlog-tag":
		case "tag":
		case "labels":
		case "env":
		default:
			return fmt.Errorf("unknown log opt '%s' for kanlog log driver", key)
		}
	}

	//you should verify stderror-log address here

	return nil
}









/*

//not currently needed as we are just writing to files
func parseAddress(address string) (string, error) {
	if address == "" { return "", nil }
	if !urlutil.IsTransportURL(address) { return "", fmt.Errorf("gelf-address should be in form proto://address, got %v", address) }
	url, err := url.Parse(address)
	if err != nil {	return "", err }

	// we support only udp - this GELF's
	if url.Scheme != "udp" { return "", fmt.Errorf("gelf: endpoint needs to be UDP") }

	// get host and port
	if _, _, err = net.SplitHostPort(url.Host); err != nil { return "", fmt.Errorf("gelf: please provide gelf-address as udp://host:port") }
	return url.Host, nil
}

func (s *kanlogger) Log(msg *logger.Message) error {
	data := map[string]string{
		"container_id":   f.containerID,
		"container_name": f.containerName,
		"source":         msg.Source,
		"log":            string(msg.Line),
	}
	for k, v := range f.extra {
		data[k] = v
	}

	if err := s.writer.WriteMessage(&data); err != nil {
		return fmt.Errorf("kanlog: failed to send kanlog message: %v", err)
	}
	return nil

	if msg.Source == "stderr" {
		return s.writer.Err(string(msg.Line))
	}
	return s.writer.Info(string(msg.Line))
}*/

