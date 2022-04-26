package player

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/notify"
)

type Player struct {
	Path string
}

func New(path string) notify.Notifier {
	return &Player{
		Path: path,
	}
}

func Default() notify.Notifier {
	return New("./audio.mp3")
}

func (p *Player) Name() string {
	return "Mp3 Player"
}

func (p *Player) Send() error {
	audioFile, err := os.Open(p.Path)
	if err != nil {
		return errs.Wrap(code.Unexpected, err)
	}
	defer func() {
		_ = audioFile.Close()
	}()

	// 对文件进行解码
	audioStreamer, format, err := mp3.Decode(audioFile)
	if err != nil {
		return errs.Wrap(code.ParseFailed, err)
	}

	defer func() {
		_ = audioStreamer.Close()
	}()
	done := make(chan bool)
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return errs.Wrap(code.Unexpected, err)
	}
	speaker.Play(beep.Seq(audioStreamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	return nil
}
