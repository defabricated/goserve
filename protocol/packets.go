package protocol

import "reflect"

type Packet interface{}

var packets = [4][2][]reflect.Type{
	Handshaking: [2][]reflect.Type{
		Clientbound: []reflect.Type{},
		Serverbound: []reflect.Type{
			reflect.TypeOf((*Handshake)(nil)).Elem(),
		},
	},
	Play: [2][]reflect.Type{
		Clientbound: []reflect.Type{},
		Serverbound: []reflect.Type{},
	},
	Login: [2][]reflect.Type{
		Clientbound: []reflect.Type{},
		Serverbound: []reflect.Type{},
	},
	Status: [2][]reflect.Type{
		Clientbound: []reflect.Type{
			reflect.TypeOf((*StatusResponse)(nil)).Elem(),
			reflect.TypeOf((*StatusPing)(nil)).Elem(),
		},
		Serverbound: []reflect.Type{
			reflect.TypeOf((*StatusGet)(nil)).Elem(),
			reflect.TypeOf((*ClientStatusPing)(nil)).Elem(),
		},
	},
}

var packetIDs = [2]map[reflect.Type]int{
	Clientbound: map[reflect.Type]int{},
	Serverbound: map[reflect.Type]int{},
}

func init() {
	for _, st := range packets {
		for d, dir := range st {
			for i, p := range dir {
				if _, ok := packetIDs[d][p]; ok {
					panic("Duplicate packet " + p.Name())
				}
				packetIDs[d][p] = i
			}
		}
	}
}
