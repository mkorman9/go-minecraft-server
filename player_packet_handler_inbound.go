package main

import (
	"crypto/x509"
	"fmt"
	"github.com/mkorman9/go-minecraft-server/packets"
	"io"
	"log"
	"time"
)

func (pph *PlayerPacketHandler) HandlePacket(packetDelivery *packets.PacketDelivery) (err error) {
	switch pph.state {
	case PlayerStateBeforeHandshake:
		err = pph.OnBeforeHandshakePacket(packetDelivery.PacketID, packetDelivery.Reader)
	case PlayerStateLogin:
		err = pph.OnLoginPacket(packetDelivery.PacketID, packetDelivery.Reader)
	case PlayerStateEncryption:
		err = pph.OnEncryptionPacket(packetDelivery.PacketID, packetDelivery.Reader)
	case PlayerStatePlay:
		err = pph.OnPlayPacket(packetDelivery.PacketID, packetDelivery.Reader)
	}

	return
}

func (pph *PlayerPacketHandler) OnBeforeHandshakePacket(packetId int, packetReader io.Reader) error {
	switch packetId {
	case 0x00:
		return pph.OnHandshakeRequest(packetReader)
	case 0x01:
		return pph.OnPing(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in before handshake state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnLoginPacket(packetId int, packetReader io.Reader) error {
	switch packetId {
	case 0x00:
		return pph.OnLoginStartRequest(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in login state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnEncryptionPacket(packetId int, packetReader io.Reader) error {
	switch packetId {
	case 0x01:
		return pph.OnEncryptionResponse(packetReader)
	default:
		return fmt.Errorf("unrecognized packet id: 0x%x in encryption state", packetId)
	}
}

func (pph *PlayerPacketHandler) OnPlayPacket(packetId int, packetReader io.Reader) error {
	switch packetId {
	case 0x00:
		return pph.OnTeleportConfirm(packetReader)
	case 0x03:
		return pph.OnChatCommand(packetReader)
	case 0x04:
		return pph.OnChatMessage(packetReader)
	case 0x07:
		return pph.OnSettings(packetReader)
	case 0x0c:
		return pph.OnCustomPayload(packetReader)
	case 0x11:
		return pph.OnKeepAliveResponse(packetReader)
	case 0x13:
		return pph.OnPosition(packetReader)
	case 0x14:
		return pph.OnPositionLook(packetReader)
	case 0x15:
		return pph.OnLook(packetReader)
	case 0x1b:
		return pph.OnAbilities(packetReader)
	case 0x1d:
		return pph.OnEntityAction(packetReader)
	case 0x2a:
		return pph.OnSetCreativeSlot(packetReader)
	case 0x2e:
		return pph.OnArmAnimation(packetReader)
	default:
		log.Printf("unrecognized packet id: 0x%x in play state\n", packetId)
		return nil
	}
}

func (pph *PlayerPacketHandler) OnHandshakeRequest(packetReader io.Reader) error {
	log.Println("received HandshakeRequest")

	handshakeRequest, err := HandshakeRequest.Read(packetReader)
	if err != nil {
		return err
	}

	if handshakeRequest.VarInt("protocolVersion") == 0 {
		return nil // ignore status request
	}

	switch handshakeRequest.VarInt("nextState") {
	case HandshakeTypeStatus:
		return pph.sendHandshakeStatusResponse()
	case HandshakeTypeLogin:
		pph.state = PlayerStateLogin
	}

	return nil
}

func (pph *PlayerPacketHandler) OnStatusRequest(_ io.Reader) error {
	log.Println("received StatusRequest")

	// ignore

	return nil
}

func (pph *PlayerPacketHandler) OnPing(packetReader io.Reader) error {
	log.Println("received PingRequest")

	pingRequest, err := PingRequest.Read(packetReader)
	if err != nil {
		return err
	}

	return pph.sendPongResponse(pingRequest.Int64("payload"))
}

func (pph *PlayerPacketHandler) OnLoginStartRequest(packetReader io.Reader) error {
	log.Println("received LoginStartRequest")

	loginStartRequest, err := LoginStartRequest.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.Name = loginStartRequest.String("name")
	pph.player.DisplayName = NewChatMessage(loginStartRequest.String("name"))
	pph.verifyToken, _ = getSecureRandomString(VerifyTokenLength)

	if loginStartRequest.Bool("hasSigData") {
		pph.player.Signature = loginStartRequest.String("signature")

		if loginStartRequest.ByteArray("publicKey") != nil {
			publicKey, err := loadPublicKey(loginStartRequest.ByteArray("publicKey"))
			if err != nil {
				log.Printf("%v\n", err)
				return NewPacketHandlingError(err, NewChatMessage("Malformed Public Key"))
			}

			publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
			if err != nil {
				return err
			}

			pph.player.PublicKey = publicKey
			pph.player.PublicKeyDER = publicKeyDER
			pph.player.Timestamp = loginStartRequest.Int64("timestamp")
		}
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

func (pph *PlayerPacketHandler) OnEncryptionResponse(packetReader io.Reader) error {
	log.Println("received EncryptionResponse")

	encryptionResponse, err := EncryptionResponse.Read(packetReader)
	if err != nil {
		return err
	}

	sharedSecret, err := pph.world.Server().DecryptMessage(encryptionResponse.ByteArray("sharedSecret"))
	if err != nil {
		log.Printf("%v\n", err)
		return NewPacketHandlingError(err, NewChatMessage("Malformed Shared Secret"))
	}

	pph.sharedSecret = sharedSecret
	pph.serverHash = pph.world.Server().GenerateServerHash(sharedSecret)

	if encryptionResponse.Bool("hasVerifyToken") {
		verifyToken, err := pph.world.Server().DecryptMessage(encryptionResponse.ByteArray("verifyToken"))
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
			encryptionResponse.Int64("salt"),
			encryptionResponse.ByteArray("messageSignature"),
		)
		if err != nil {
			log.Printf("%v\n", err)
			return NewPacketHandlingError(err, NewChatMessage("Signature verification error"))
		}
	}

	verificationResult, err := MojangVerifyPlayer(pph.player.Name, pph.serverHash)
	if err != nil {
		fmt.Println(err)
	}

	if !verificationResult.Verified {
		if err != nil {
			return NewPacketHandlingError(err, NewChatMessage("Username verification failed"))
		}
	}

	pph.player.UUID = *verificationResult.UUID
	pph.player.Textures = verificationResult.Textures
	pph.player.TexturesSignature = verificationResult.TexturesSignature

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

func (pph *PlayerPacketHandler) OnTeleportConfirm(packetReader io.Reader) error {
	log.Println("received TeleportConfirm")

	_, err := TeleportConfirmPacket.Read(packetReader)
	if err != nil {
		return err
	}

	// nop

	return nil
}

func (pph *PlayerPacketHandler) OnSettings(packetReader io.Reader) error {
	log.Println("received Settings")

	settingsPacket, err := SettingsPacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnClientSettings(&PlayerClientSettings{
		Locale:              settingsPacket.String("locale"),
		ViewDistance:        settingsPacket.Byte("viewDistance"),
		ChatColors:          settingsPacket.Bool("chatColors"),
		SkinParts:           settingsPacket.Byte("skinParts"),
		MainHand:            settingsPacket.VarInt("mainHand"),
		EnableTextFiltering: settingsPacket.Bool("enableTextFiltering"),
		EnableServerListing: settingsPacket.Bool("enableServerListing"),
	})

	return nil
}

func (pph *PlayerPacketHandler) OnPosition(packerReader io.Reader) error {
	positionPacket, err := PositionPacket.Read(packerReader)
	if err != nil {
		return err
	}

	pph.player.OnPositionUpdate(
		positionPacket.Float64("x"),
		positionPacket.Float64("y"),
		positionPacket.Float64("z"),
	)
	pph.player.OnGroundUpdate(positionPacket.Bool("onGround"))

	return nil
}

func (pph *PlayerPacketHandler) OnPositionLook(packetReader io.Reader) error {
	positionLookPacket, err := PositionLookPacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnPositionUpdate(
		positionLookPacket.Float64("x"),
		positionLookPacket.Float64("y"),
		positionLookPacket.Float64("z"),
	)
	pph.player.OnGroundUpdate(positionLookPacket.Bool("onGround"))
	pph.player.OnLookUpdate(positionLookPacket.Float32("yaw"), positionLookPacket.Float32("pitch"))

	return nil
}

func (pph *PlayerPacketHandler) OnLook(packetReader io.Reader) error {
	lookPacket, err := LookPacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnGroundUpdate(lookPacket.Bool("onGround"))
	pph.player.OnLookUpdate(lookPacket.Float32("yaw"), lookPacket.Float32("pitch"))

	return nil
}

func (pph *PlayerPacketHandler) OnCustomPayload(packetReader io.Reader) error {
	log.Println("received CustomPayload")

	customPayloadPacket, err := CustomPayloadPacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnPluginChannel(
		customPayloadPacket.String("channel"),
		customPayloadPacket.ByteArray("data"),
	)

	return nil
}

func (pph *PlayerPacketHandler) OnArmAnimation(packetReader io.Reader) error {
	log.Println("received ArmAnimation")

	armAnimationPacket, err := ArmAnimationPacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnArmAnimation(armAnimationPacket.VarInt("hand"))

	return nil
}

func (pph *PlayerPacketHandler) OnAbilities(packetReader io.Reader) error {
	log.Println("received Abilities")

	_, err := AbilitiesPacket.Read(packetReader)
	if err != nil {
		return err
	}

	// TODO

	return nil
}

func (pph *PlayerPacketHandler) OnSetCreativeSlot(packetReader io.Reader) error {
	log.Println("received SetCreativeSlot")

	_, err := SetCreativeSlotPacket.Read(packetReader)
	if err != nil {
		return err
	}

	// TODO

	return nil
}

func (pph *PlayerPacketHandler) OnKeepAliveResponse(packetReader io.Reader) error {
	keepAliveResponsePacket, err := KeepAliveResponsePacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnKeepAliveResponse(keepAliveResponsePacket.Int64("keepAliveId"))

	return nil
}

func (pph *PlayerPacketHandler) OnEntityAction(packetReader io.Reader) error {
	log.Println("received EntityAction")

	entityActionPacket, err := EntityActionPacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnAction(
		entityActionPacket.VarInt("entityId"),
		entityActionPacket.VarInt("actionId"),
		entityActionPacket.VarInt("jumpBoost"),
	)

	return nil
}

func (pph *PlayerPacketHandler) OnChatCommand(packetReader io.Reader) error {
	log.Println("received ChatCommand")

	chatCommandPacket, err := ChatCommandPacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnChatCommand(
		chatCommandPacket.String("message"),
		time.UnixMilli(chatCommandPacket.Int64("timestamp")),
	)

	return nil
}

func (pph *PlayerPacketHandler) OnChatMessage(packetReader io.Reader) error {
	log.Println("received ChatMessage")

	chatMessagePacket, err := ChatMessagePacket.Read(packetReader)
	if err != nil {
		return err
	}

	pph.player.OnChatMessage(
		chatMessagePacket.String("message"),
		time.UnixMilli(chatMessagePacket.Int64("timestamp")),
	)

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

	err = pph.sendPlayerListUpdate([]*Player{pph.player})
	if err != nil {
		return err
	}

	pph.state = PlayerStatePlay
	pph.player.OnJoin(GameModeSurvival)

	return nil
}
