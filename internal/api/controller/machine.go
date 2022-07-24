package controller

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func MachineInfo(_ *gin.Context, _ *struct{}) (_ interface{}, err error) {
	machineId, err := machineid.ProtectedID("whatsapp-client")
	if err != nil {
		log.Fatal(err)
	}
	machineId = strings.ToUpper(machineId[:16])

	return gin.H{
		"machineId": machineId,
		"version":   version,
	}, nil
}
