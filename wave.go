package wave

import (
	"os"
	"io"
	"bytes"
	"encoding/binary"
)

type WaveFile struct {
	Riff	string	
	ChunkSize uint32
	Wave	string
	Fmt	string
	Subchunk1Size uint32
	AudioFormat uint16
	NumOfChan	uint16
	SampleRate	uint32
	BytesPerSec uint16
	BlockAlign uint16
	BitsPerSample uint16
	Subchunk2ID string
	DataSize uint32
	Data []byte
	
	ReadPtr int

	FilePtr  *os.File

}

func OpenWav(filename string, mode string) (*WaveFile) {
	
	var f *os.File
	var err error
	
	if (mode == "w") {
		f,err = os.Create(filename)
	}
	if (mode == "r") {
		f,err = os.Open(filename)
	}

	
	if err != nil {
   		panic(err)
   	}

	wv := readWav(f) 
	wv.FilePtr = f
	return wv
}


func (wv WaveFile) GetNChannels() int {
	return int(wv.NumOfChan)
}

func (wv WaveFile) SetNChannels(n int) {
	wv.NumOfChan = uint16(n)
}

func (wv WaveFile) SetSampWidth(n int) {
	wv.BitsPerSample = uint16(n) * 8
}


func (wv WaveFile) SetFrameRate(n int) {
	wv.SampleRate = uint32(n)	
}

func (wv WaveFile) GetSampWidth() int {

	return int(wv.BitsPerSample/8)

}

func (wv WaveFile) GetFrameRate() int {

	return int(wv.SampleRate)

}



func (wv WaveFile) WriteFrames(data *[] byte) {
//Write audio frames and make sure nframes is correct.
	if len(*data)/(wv.GetNChannels() * wv.GetSampWidth()) != wv.GetNFrames() {
		panic ("frame size not valid")
	}
	
	wv.FilePtr.Write(*data) 
	
}


func (wv WaveFile) WriteFramesRaw(data *[] byte) {
//Write audio frames with no check of nframes
	
	wv.FilePtr.Write(*data) 
	
}


func (wv WaveFile) GetNFrames() int {

	numSamples  := wv.DataSize/wv.SampleRate
	numChan := uint32( wv.NumOfChan)
	return int(numSamples/numChan)

}

func (wv WaveFile) ReadFrames(n int) []byte {
  // n frame reads is reading length of n * sampwidth * channels 
	bytesToRead := n * wv.GetSampWidth() * wv.GetNChannels()
	defer func() {
		wv.ReadPtr += bytesToRead
	}()
	return wv.Data[wv.ReadPtr:bytesToRead+wv.ReadPtr]

}

func (wv WaveFile) Close() {
	wv.FilePtr.Close()

}


func (wv WaveFile) Rewind() {
	wv.ReadPtr = 0
}

func (wv WaveFile) SetPos(n int) {
	wv.ReadPtr = n
}

func (wv WaveFile) Tell() int {
	return wv.ReadPtr
}


func (wv WaveFile) GetParams() (int,int,int,int) {
	return wv.GetNChannels(),wv.GetSampWidth(),wv.GetFrameRate(),wv.GetNFrames()
}
 
func readWav(f *os.File) (*WaveFile) {
	  //make slice for the 44 byte header
		data := make([]byte, 44)
		_, err:= f.Read(data)
		if err != nil {
		           panic("file error")
		}
		wv := new(WaveFile)
		wv.Riff = string(data[0:4])

		// conversion needed to move the raw binary data to integer value
		// create a buffer of bytes to hold the data
		buf := bytes.NewBuffer(data[4:8])
		// read the data in the buffer into the integer using LittleEndian conversion
		binary.Read(buf, binary.LittleEndian, &wv.ChunkSize)
		wv.Wave = string(data[8:12])
		wv.Fmt = string(data[12:16])
		buf = bytes.NewBuffer(data[16:20])
		binary.Read(buf, binary.LittleEndian, &wv.Subchunk1Size)

		buf = bytes.NewBuffer(data[20:22])
		binary.Read(buf, binary.LittleEndian, &wv.AudioFormat)

		buf = bytes.NewBuffer(data[22:24])
		binary.Read(buf, binary.LittleEndian, &wv.NumOfChan)

		buf = bytes.NewBuffer(data[24:28])
		binary.Read(buf, binary.LittleEndian, &wv.SampleRate)

		buf = bytes.NewBuffer(data[28:32])
		binary.Read(buf, binary.LittleEndian, &wv.BytesPerSec)

		buf = bytes.NewBuffer(data[32:34])
		binary.Read(buf, binary.LittleEndian, &wv.BlockAlign)

		buf = bytes.NewBuffer(data[34:36])
		binary.Read(buf, binary.LittleEndian, &wv.BitsPerSample)

		wv.DataSize = uint32(wv.BitsPerSample/8) * uint32( wv.SampleRate)    

		buf = bytes.NewBuffer(data[36:40])
		binary.Read(buf, binary.LittleEndian, &wv.Subchunk2ID)

		buf = bytes.NewBuffer(data[40:44])
		binary.Read(buf, binary.LittleEndian, &wv.DataSize)

		data2 := make([]byte, wv.DataSize)
		_, err = f.Read(data2)
		if err != nil {
		           if err == io.EOF {
		               panic ("file error")
		           }
		}
		
		wv.Data = data2;

		return wv
}

