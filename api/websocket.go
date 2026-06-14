package api

import (
    "encoding/base64"
    "crypto/sha1"
    "net/http"
    "errors"
    "bufio"
    "net"
    "fmt"
    "mygo/config"
    "mygo/log"
)

func websocket_upgrade(w http.ResponseWriter, r *http.Request) error {
    if !r.ProtoAtLeast(1, 1) {
        return errors.New("need \"HTTP/1.1\"")
    }

    if r.Host != config.SERVER_HOST {
        return errors.New("invalid \"Host\"")
    }

    if r.Header.Get("Sec-WebSocket-Version") != "13" {

        w.WriteHeader(426)

        return errors.New("invalid \"Set-WebSocket-Version\"")
    }

    if r.Header.Get("Connection") != "Upgrade" {
        return errors.New("invalid \"Connection\"")
    }

    if r.Header.Get("Upgrade") != "websocket" {
        return errors.New("invalid \"Upgrade\"")
    }

    key := r.Header.Get("Sec-WebSocket-Key")

    decoded, err := base64.StdEncoding.DecodeString(key)

    if err != nil || len(decoded) != 16 {
        return errors.New("invalid \"Sec-WebSocket-Key\"")
    }

    magic  := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
    hash   := sha1.Sum([]byte(key + magic))
    accept := base64.StdEncoding.EncodeToString(hash[:])

    response := fmt.Sprintf(
        "HTTP/1.1 101 Switching Protocols\r\n" +
        "Upgrade: websocket\r\n" +
        "Connection: Upgrade\r\n" +
        "Sec-WebSocket-Accept: %s\r\n" +
        "\r\n", accept);

    hj, ok := w.(http.Hijacker)
    if !ok {
        w.WriteHeader(500)
        return errors.New("unsupported websocket")
    }

    conn, rw, err := hj.Hijack()
    if err != nil {
        w.WriteHeader(500)
        return errors.New("unsupported websocket")
    }

    rw.WriteString(response)
    rw.Flush()

    go func(conn net.Conn, rw *bufio.ReadWriter) {
        defer conn.Close()

        buffer := make([]byte, 0x100)

        for {
            _, _ = rw.Read(buffer)

            log.Info <- buffer
        }
    }(conn, rw)

    return nil
}

func websocket_make(payload_type string, length uint64) ([]byte, error) {
    if (length >> 63) != 0 {
        return nil, errors.New("error length")
    }

    var fin         byte   = 0x80
    var rsv         byte   = 0x00
    var opcode      byte
    var payload_len byte
    var extend_len  []byte
    var header      []byte

    switch payload_type {
        case "continuation": opcode = 0x0
        case "text":         opcode = 0x1
        case "binary":       opcode = 0x2
        case "close":        opcode = 0x8
        case "ping":         opcode = 0x9
        case "pong":         opcode = 0xa
        default: return nil, errors.New("unknown type")
    }

    if length <= 0x7d {

        payload_len = byte(length)
        extend_len  = []byte{}

    } else if length <= 0xffff {

        payload_len = 0x7e
        extend_len  = []byte{
            byte(length >> 8), 
            byte(length >> 0)}

    } else {

        payload_len = 0x7f
        extend_len  = []byte{
            byte(length >> 0x38),
            byte(length >> 0x30),
            byte(length >> 0x28),
            byte(length >> 0x20),
            byte(length >> 0x18),
            byte(length >> 0x10),
            byte(length >> 0x08),
            byte(length >> 0x00)}
    }

    header = []byte{fin | rsv | opcode}
    header = append(header, payload_len)
    header = append(header, extend_len...)
    
    return header, nil
}

func websocket_type(data []byte) (string, error) {
    if len(data) < 1 {
        return "", errors.New("invalid data")
    }

    var payload_type string

    switch data[0] & 0x0f {
        case 0x0: payload_type = "continuation"
        case 0x1: payload_type = "text"
        case 0x2: payload_type = "binary"
        case 0x8: payload_type = "close"
        case 0x9: payload_type = "ping"
        case 0xa: payload_type = "pong"
        default: return "", errors.New("unknown type")
    }

    return payload_type, nil
}

func websocket_has_mask(data []byte) (bool, error) {
    if len(data) < 2 {
        return false, errors.New("invalid data")
    }

    return (data[1] & 0x80) != 0, nil
}

func websocket_payload_len(data []byte) (uint8, error) {
    if len(data) < 2 {
        return 0, errors.New("invalid data")
    }

    return uint8(data[1] & 0x7f), nil
}

func websocket_extend_len(data []byte) (uint64, error) {
    payload_len, err := websocket_payload_len(data)
    if err != nil {
        return 0, err
    }

    if payload_len <= 0x7d {
        return 0, nil
    }

    data_len := len(data)

    if  (payload_len == 0x7e && data_len < 4) ||
        (payload_len == 0x7f && data_len < 10) {

        return 0, errors.New("invalid data")
    }

    if payload_len == 0x7e {
        return (uint64(data[2]) << 0x8) | 
               (uint64(data[3]) << 0x0), nil
    }

    return (uint64(data[2]) << 0x38) |
           (uint64(data[3]) << 0x30) |
           (uint64(data[4]) << 0x28) |
           (uint64(data[5]) << 0x20) |
           (uint64(data[6]) << 0x18) |
           (uint64(data[7]) << 0x10) |
           (uint64(data[8]) << 0x08) |
           (uint64(data[9]) << 0x00), nil
}

func websocket_get_maskkey(data []byte) ([4]byte, error) {

    var maskkey [4]byte

    data_len := len(data)

    payload_len, err := websocket_payload_len(data)
    if err != nil {
        return maskkey, err
    }

    offset := 2

    if payload_len == 0x7e {
        offset += 2
    }

    if payload_len == 0x7f {
        offset += 8
    }

    has_mask, err := websocket_has_mask(data)
    if err != nil {
        return maskkey, err
    }

    if !has_mask {
        return maskkey, errors.New("no mask")
    }

    if data_len < offset + 4 {
        return maskkey, errors.New("no mask key")
    }

    maskkey = [4]byte(data[offset : offset + 4])

    return maskkey, nil
}

func websocket_unmask(maskkey [4]byte, index *uint64, 
    payload_len uint64, payload_part []byte) {

    end := (*index) + uint64(len(payload_part))

    if end > payload_len {
        end = payload_len
    }

    end -= (*index)

    for i := uint64(0); i < end; i++ {
        payload_part[i] ^= maskkey[(*index) % 4]
        (*index)++
    }
}
