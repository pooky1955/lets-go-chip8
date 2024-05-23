package main

type Speaker interface {
	Play(frequency int)
	Stop()
}

type CLISpeaker struct {
}

func (*CLISpeaker) Play(frequency int) {

}
func (*CLISpeaker) Stop() {

}
