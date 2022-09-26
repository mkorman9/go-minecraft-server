package packets

import (
	"bytes"
	"compress/zlib"
	"github.com/mkorman9/go-minecraft-server/types"
	"io"
)

type PacketWriter struct {
	writer               io.Writer
	compressionThreshold int
}

func NewPacketWriter(writer io.Writer) *PacketWriter {
	return &PacketWriter{
		writer:               writer,
		compressionThreshold: -1,
	}
}

func (pw *PacketWriter) SetCompression(threshold int) {
	pw.compressionThreshold = threshold
}

func (pw *PacketWriter) SetEncryption(cipherStream *CipherStream) {
	pw.writer = cipherStream.WrapWriter(pw.writer)
}

func (pw *PacketWriter) Write(packet *PacketData) error {
	var packetData bytes.Buffer

	_, err := packet.WriteTo(&packetData)
	if err != nil {
		return err
	}

	switch pw.compressionThreshold {
	case -1:
		// no compression
		err = types.WriteVarInt(pw.writer, packetData.Len())
		if err != nil {
			return err
		}

		_, err = pw.writer.Write(packetData.Bytes())
		if err != nil {
			return err
		}
	default:
		// compression
		if packetData.Len() >= pw.compressionThreshold {
			var zlibBuffer bytes.Buffer

			zlibWriter := zlib.NewWriter(&zlibBuffer)
			_, err = zlibWriter.Write(packetData.Bytes())
			if err != nil {
				return err
			}

			err = zlibWriter.Close()
			if err != nil {
				return err
			}

			err = types.WriteVarInt(pw.writer, types.GetVarIntSize(packetData.Len())+zlibBuffer.Len())
			if err != nil {
				return err
			}

			err = types.WriteVarInt(pw.writer, packetData.Len())
			if err != nil {
				return err
			}

			_, err := pw.writer.Write(zlibBuffer.Bytes())
			if err != nil {
				return err
			}
		} else {
			err = types.WriteVarInt(pw.writer, packetData.Len()+1)
			if err != nil {
				return err
			}

			err = types.WriteVarInt(pw.writer, 0)
			if err != nil {
				return err
			}

			_, err := pw.writer.Write(packetData.Bytes())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
