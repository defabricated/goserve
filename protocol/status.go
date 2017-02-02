package protocol

type StatusResponse struct {
	Data string
}

type StatusPing struct {
	Time int64
}

type StatusGet struct {
}

type ClientStatusPing struct {
	Time int64
}
