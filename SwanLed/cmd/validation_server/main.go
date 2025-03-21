package validation_server

//import (
//	"fmt"
//	"github.com/sirupsen/logrus"
//	"io"
//	"net"
//	"os"
//)
//
//var log = logrus.New()
//
//func SendFile(conn net.Conn, filePath string) error {
//	file, err := os.Open(filePath)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	defer func(file *os.File) {
//		err := file.Close()
//		if err != nil {
//			log.Error(err)
//		}
//	}(file)
//
//	info, err := file.Stat()
//	if err != nil {
//		log.Error("Error getting file stats:", err)
//		return err
//	}
//	if info.Size() == 0 {
//		log.Error("File is empty.")
//		return fmt.Errorf("file is empty")
//	}
//
//	_, err = io.Copy(conn, file)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	return nil
//}
//
//func ListenTheValidatorConnection(validatorAddress string) {
//	listener, err := net.Listen("tcp", ":"+validatorAddress)
//	if err != nil {
//		log.Error(err)
//	}
//
//	log.Info("Listening on :9008. Validating server...")
//	defer func(listener net.Listener) {
//		err := listener.Close()
//		if err != nil {
//			log.Error(err)
//		}
//	}(listener)
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			log.Error(err)
//		}
//		remoteAddr := conn.RemoteAddr().String()
//		log.Info("New connection from : ", remoteAddr)
//		defer func(conn net.Conn) {
//			err := conn.Close()
//			if err != nil {
//				log.Error(err)
//			}
//		}(conn)
//		go func(conn net.Conn) {
//			err := SendFile(conn, "./dag_data.bin")
//			if err != nil {
//				log.Error(err)
//			}
//		}(conn)
//		go func(conn net.Conn) {
//
//		}(conn)
//	}
//}
//
//func Start(validatorsPort string) {
//	log.Info("Validator server start at localhost:0")
//	ListenTheValidatorConnection(validatorsPort)
//}
