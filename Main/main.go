package main

import log "gorm/Log"

func main() {
	log.SetLogLevel(log.WARN)
	log.Infof("infof")
	log.Warnf("warnf")
	log.Errorf("errorf %v", 3.14)
}
