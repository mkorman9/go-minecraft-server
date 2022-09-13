package main

import (
	"fmt"
	"log"
	"time"
)

func (pph *PlayerPacketHandler) HandlePacket(packet []byte) (err error) {
	packetReader := NewPacketReaderContext(packet)
	packetId := packetReader.FetchVarInt()

	switch pph.state {
	case PlayerStateBeforeHandshake:
		err = pph.OnBeforeHandshakePacket(packetId, packetReader)
	case PlayerStateLogin:
		err = pph.OnLoginPacket(packetId, packetReader)
	case PlayerStateEncryption:
		err = pph.OnEncryptionPacket(packetId, packetReader)
	case PlayerStatePlay:
		err = pph.OnPlayPacket(packetId, packetReader)
	}

	return
}

func (pph *PlayerPacketHandler) OnBeforeHandshakePacket(packetId int, packetReader *PacketReaderContext) error {
	switch packetId {
	case 0x00:
		if packetReader.BytesLeft() > 0 {
			return pph.OnHandshakeRequest(packetReader)
		} else {
			return pph.OnStatusRequest(packetReader)
		}
	case 0x01:
		return pph.OnPing(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in before handshake state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnLoginPacket(packetId int, packetReader *PacketReaderContext) error {
	switch packetId {
	case 0x00:
		return pph.OnLoginStartRequest(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in login state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnEncryptionPacket(packetId int, packetReader *PacketReaderContext) error {
	switch packetId {
	case 0x01:
		return pph.OnEncryptionResponse(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in encryption state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnPlayPacket(packetId int, packetReader *PacketReaderContext) error {
	switch packetId {
	case 0x00:
		return pph.OnTeleportConfirm(packetReader)
	case 0x07:
		return pph.OnSettings(packetReader)
	case 0x0c:
		return pph.OnCustomPayload(packetReader)
	case 0x13:
		return pph.OnPosition(packetReader)
	case 0x14:
		return pph.OnPositionLook(packetReader)
	case 0x15:
		return pph.OnLook(packetReader)
	case 0x2e:
		return pph.OnArmAnimation(packetReader)
	case 0x1b:
		return pph.OnAbilities(packetReader)
	case 0x2a:
		return pph.OnSetCreativeSlot(packetReader)
	case 0x03:
		return pph.OnChatCommand(packetReader)
	case 0x04:
		return pph.OnChatMessage(packetReader)
	case 0x11:
		return pph.OnKeepAliveResponse(packetReader)
	case 0x1d:
		return pph.OnEntityAction(packetReader)
	default:
		log.Printf("unrecognized packet id: 0x%x in play state\n", packetId)
		return nil
	}
}

func (pph *PlayerPacketHandler) OnHandshakeRequest(packetReader *PacketReaderContext) error {
	log.Println("received HandshakeRequest")

	var request HandshakeRequest
	err := request.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	switch request.NextState {
	case HandshakeTypeStatus:
		return pph.sendHandshakeStatusResponse()
	case HandshakeTypeLogin:
		pph.state = PlayerStateLogin
	}

	return nil
}

func (pph *PlayerPacketHandler) OnStatusRequest(_ *PacketReaderContext) error {
	log.Println("received StatusRequest")

	// ignore

	return nil
}

func (pph *PlayerPacketHandler) OnPing(packetReader *PacketReaderContext) error {
	log.Println("received PingRequest")

	var request PingRequest
	err := request.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	return pph.sendPongResponse(request.Payload)
}

func (pph *PlayerPacketHandler) OnLoginStartRequest(packetReader *PacketReaderContext) error {
	log.Println("received LoginStartRequest")

	var request LoginStartRequest
	err := request.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.Name = request.Name
	pph.player.Signature = request.Signature
	pph.verifyToken, _ = getSecureRandomString(VerifyTokenLength)

	if request.PublicKey != nil {
		publicKey, err := loadPublicKey(request.PublicKey)
		if err != nil {
			log.Printf("%v\n", err)
			return NewPacketHandlingError(err, NewChatMessage("Malformed Public Key"))
		}

		pph.player.PublicKey = publicKey
	}

	if pph.world.Settings().OnlineMode {
		pph.state = PlayerStateEncryption
		return pph.sendEncryptionRequest()
	} else {
		err := pph.setupCompression()
		if err != nil {
			return err
		}

		err = pph.sendLoginSuccessResponse()
		if err != nil {
			return err
		}

		return pph.OnJoin()
	}
}

func (pph *PlayerPacketHandler) OnEncryptionResponse(packetReader *PacketReaderContext) error {
	log.Println("received EncryptionResponse")

	var response EncryptionResponse
	err := response.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	sharedSecret, err := pph.world.Server().DecryptMessage(response.SharedSecret)
	if err != nil {
		log.Printf("%v\n", err)
		return NewPacketHandlingError(err, NewChatMessage("Malformed Shared Secret"))
	}

	pph.sharedSecret = sharedSecret
	pph.serverHash = pph.world.Server().GenerateServerHash(sharedSecret)

	if response.VerifyToken != nil {
		verifyToken, err := pph.world.Server().DecryptMessage(response.VerifyToken)
		if err != nil {
			log.Printf("%v\n", err)
			return NewPacketHandlingError(err, NewChatMessage("Malformed Verify Token"))
		}

		if string(verifyToken) != pph.verifyToken {
			log.Printf("token mismatch\n")
			return NewPacketHandlingError(err, NewChatMessage("Token mismatch"))
		}
	} else {
		err = verifyRsaSignature(
			pph.player.PublicKey,
			pph.verifyToken,
			response.Salt,
			response.MessageSignature,
		)
		if err != nil {
			log.Printf("%v\n", err)
			return NewPacketHandlingError(err, NewChatMessage("Signature verification error"))
		}
	}

	err = pph.setupEncryption()
	if err != nil {
		return err
	}

	err = pph.setupCompression()
	if err != nil {
		return err
	}

	err = pph.sendLoginSuccessResponse()
	if err != nil {
		return err
	}

	return pph.OnJoin()
}

func (pph *PlayerPacketHandler) OnTeleportConfirm(packetReader *PacketReaderContext) error {
	log.Println("received TeleportConfirm")

	var packet TeleportConfirmPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	// nop

	return nil
}

func (pph *PlayerPacketHandler) OnSettings(packetReader *PacketReaderContext) error {
	log.Println("received Settings")

	var packet SettingsPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnClientSettings(&PlayerClientSettings{
		Locale:              packet.Locale,
		ViewDistance:        packet.ViewDistance,
		ChatColors:          packet.ChatColors,
		SkinParts:           packet.SkinParts,
		MainHand:            packet.MainHand,
		EnableTextFiltering: packet.EnableTextFiltering,
		EnableServerListing: packet.EnableServerListing,
	})

	return nil
}

func (pph *PlayerPacketHandler) OnPosition(packetReader *PacketReaderContext) error {
	var packet PositionPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnPositionUpdate(packet.X, packet.Y, packet.Z)
	pph.player.OnGroundUpdate(packet.OnGround)

	return nil
}

func (pph *PlayerPacketHandler) OnPositionLook(packetReader *PacketReaderContext) error {
	var packet PositionLookPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnPositionUpdate(packet.X, packet.Y, packet.Z)
	pph.player.OnGroundUpdate(packet.OnGround)
	pph.player.OnLookUpdate(packet.Yaw, packet.Pitch)

	return nil
}

func (pph *PlayerPacketHandler) OnLook(packetReader *PacketReaderContext) error {
	var packet LookPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnGroundUpdate(packet.OnGround)
	pph.player.OnLookUpdate(packet.Yaw, packet.Pitch)

	return nil
}

func (pph *PlayerPacketHandler) OnCustomPayload(packetReader *PacketReaderContext) error {
	log.Println("received CustomPayload")

	var packet CustomPayloadPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnPluginChannel(packet.Channel, packet.Data)

	return nil
}

func (pph *PlayerPacketHandler) OnArmAnimation(packetReader *PacketReaderContext) error {
	log.Println("received ArmAnimation")

	var packet ArmAnimationPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnArmAnimation(packet.Hand)

	return nil
}

func (pph *PlayerPacketHandler) OnAbilities(packetReader *PacketReaderContext) error {
	log.Println("received Abilities")

	var packet AbilitiesPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	// TODO

	return nil
}

func (pph *PlayerPacketHandler) OnSetCreativeSlot(packetReader *PacketReaderContext) error {
	log.Println("received SetCreativeSlot")

	var packet SetCreativeSlotPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	// TODO

	return nil
}

func (pph *PlayerPacketHandler) OnKeepAliveResponse(packetReader *PacketReaderContext) error {
	var packet KeepAliveResponsePacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnKeepAliveResponse(packet.KeepAliveID)

	return nil
}

func (pph *PlayerPacketHandler) OnEntityAction(packetReader *PacketReaderContext) error {
	log.Println("received EntityAction")

	var packet EntityActionPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnAction(packet.EntityID, packet.ActionID, packet.JumpBoost)

	return nil
}

func (pph *PlayerPacketHandler) OnChatCommand(packetReader *PacketReaderContext) error {
	log.Println("received ChatCommand")

	var packet ChatCommandPacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnChatCommand(packet.Message, time.UnixMilli(packet.Timestamp))

	return nil
}

func (pph *PlayerPacketHandler) OnChatMessage(packetReader *PacketReaderContext) error {
	log.Println("received ChatMessage")

	var packet ChatMessagePacket
	err := packet.Unmarshal(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnChatMessage(packet.Message, time.UnixMilli(packet.Timestamp))

	return nil
}

func (pph *PlayerPacketHandler) OnJoin() error {
	pph.player.EntityID = pph.world.GenerateEntityID()

	err := pph.sendPlayPacket(pph.player.EntityID)
	if err != nil {
		return err
	}

	err = pph.sendSpawnPosition()
	if err != nil {
		return err
	}

	pph.state = PlayerStatePlay
	pph.player.OnJoin(GameModeSurvival)

	return nil
}
