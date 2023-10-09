package parse

/*
func getUdpOriginalDstAddr(conn *net.UDPConn, b []byte) (n int, remoteAddr *net.UDPAddr, dstAddr *net.UDPAddr, err error) {
	oob := u.GetBufferPool(1024)
	defer u.PutBufferPool(oob)
	n, oobn, _, remoteAddr, err := conn.ReadMsgUDP(b, *oob)
	if err != nil {
		return 0, nil, nil, err
	}

	msgs, err := unix.ParseSocketControlMessage((*oob)[:oobn])
	if err != nil {
		return 0, nil, nil, fmt.Errorf("parsing socket control message: %s", err)
	}

	for _, msg := range msgs {
		if msg.Header.Level == unix.SOL_IP && msg.Header.Type == unix.IP_RECVORIGDSTADDR {
			originalDstRaw := &unix.RawSockaddrInet4{}
			if err = binary.Read(bytes.NewReader(msg.Data), binary.LittleEndian, originalDstRaw); err != nil {
				return 0, nil, nil, fmt.Errorf("reading original destination address: %s", err)
			}

			switch originalDstRaw.Family {
			case unix.AF_INET:
				pp := (*unix.RawSockaddrInet4)(unsafe.Pointer(originalDstRaw))
				p := (*[2]byte)(unsafe.Pointer(&pp.Port))
				dstAddr = &net.UDPAddr{
					IP:   net.IPv4(pp.Addr[0], pp.Addr[1], pp.Addr[2], pp.Addr[3]),
					Port: int(p[0])<<8 + int(p[1]),
				}

			case unix.AF_INET6:
				pp := (*unix.RawSockaddrInet6)(unsafe.Pointer(originalDstRaw))
				p := (*[2]byte)(unsafe.Pointer(&pp.Port))
				dstAddr = &net.UDPAddr{
					IP:   net.IP(pp.Addr[:]),
					Port: int(p[0])<<8 + int(p[1]),
					Zone: strconv.Itoa(int(pp.Scope_id)),
				}

			default:
				return 0, nil, nil, fmt.Errorf("original destination is an unsupported network family")
			}
			break
		}
	}

	if dstAddr == nil {
		return 0, nil, nil, fmt.Errorf("unable to obtain original destination: %s", err)
	}

	return n, remoteAddr, dstAddr, nil

}
*/
