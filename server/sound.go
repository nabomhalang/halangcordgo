package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nabomhalang/halangcordgo/config"
	"github.com/nabomhalang/halangcordgo/queue"
)

func PlaySound(guildID string, el *queue.Elements) (SkipReason, error) {
	var (
		opuslen    int16
		skip       bool
		skipReason SkipReason
		err        error
	)

	_ = SV[guildID].VC.Speaking(true)

	for {
		select {
		case <-SV[guildID].Pause:
			select {
			case <-SV[guildID].Resume:
				el.Segments = SV[guildID].Queue.GetFirst().Segments
			case skipReason = <-SV[guildID].Skip:
				cleanUp(guildID, el.Closer)
				return skipReason, nil
			}
		case skipReason = <-SV[guildID].Skip:
			cleanUp(guildID, el.Closer)
			return skipReason, nil
		default:
			if el.Segments[SV[guildID].Frames] {
				skip = !skip
			}

			err = binary.Read(el.Reader, binary.LittleEndian, &opuslen)

			if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
				if el.Loop {
					if el.Closer != nil {
						_ = el.Closer.Close()
					}

					f, _ := os.Open(fmt.Sprintf("%s%s.dca", config.Get().CachePath, el.ID))
					el.Reader = f
					el.Closer = f
					continue
				} else {
					cleanUp(guildID, el.Closer)
					return Finished, nil
				}
			}

			InBuf := make([]byte, opuslen)
			err = binary.Read(el.Reader, binary.LittleEndian, &InBuf)

			SV[guildID].Frames++

			if skip {
				continue
			}

			if err != nil {
				cleanUp(guildID, el.Closer)
				_ = os.Remove(fmt.Sprintf("%s%s.dca", config.Get().CachePath, el.ID))
				return Error, err
			}

			SV[guildID].VC.OpusSend <- InBuf
		}
	}
}

func cleanUp(guildID string, closer io.Closer) {
	_ = SV[guildID].VC.Speaking(false)
	SV[guildID].Frames = 0

	if closer != nil {
		_ = closer.Close()
	}
}
