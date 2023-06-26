package controller

import (
	"log"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/weilence/whatsapp-client/internal/api"
)

type MachineInfoRes struct {
	MachineID string `json:"machineId"`
	Version   string `json:"version"`
}

func MachineInfo(_ *api.HttpContext, _ *struct{}) (_ interface{}, err error) {
	machineID, err := machineid.ProtectedID("whatsapp-client")
	if err != nil {
		log.Fatal(err)
	}
	machineID = strings.ToUpper(machineID[:16])

	return MachineInfoRes{
		MachineID: machineID,
		Version:   version,
	}, nil
}
