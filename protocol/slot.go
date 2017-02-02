package protocol

type Slot struct {
	ID     int16
	Count  int8   `if:"ID,!=,-1"`
	Damage int16  `if:"ID,!=,-1"`
	Tag    []byte `if:"ID,!=,-1" nil:"-1" ltype:"int16"`
}
