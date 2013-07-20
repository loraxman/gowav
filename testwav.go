package main

import (
	"flag"
	"fmt"
	"wave"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "fileName", "test.wav", "pass the filename for the wav")
	flag.Parse()
	wv := wave.OpenWav(fileName, "r")
	fmt.Printf("%s\t%d\t%s\t%s\t%d\t%d\t%d\t%d\t%d\n", wv.Riff, wv.ChunkSize, wv.Wave, wv.Fmt, wv.AudioFormat, wv.NumOfChan,
		wv.SampleRate, wv.DataSize, wv.BitsPerSample)

	fmt.Printf("%d\t%d\n", wv.GetNChannels(), wv.GetNFrames())
	data := wv.ReadFrames(30)
	fmt.Printf("%s\n", data)
	a, b, c, d := wv.GetParams()
	fmt.Printf("%d\t%d\t%d\t%d\n", a, b, c, d)

	wv2 := wave.OpenWav("test2.wav", "w")
	wv2.SetNChannels(2)
	wv2.SetFrameRate(48000)
	wv2.SetSampWidth(2)
	wv2.WriteFrames(&data)
	wv2.Close()
}
