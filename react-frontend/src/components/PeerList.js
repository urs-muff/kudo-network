// src/components/PeerList.js
import React, { useState, useEffect } from 'react';
import { ListGroup } from 'react-bootstrap';

const WS_URL = 'ws://localhost:9090/ws/peers';

function PeerList() {
  const [peers, setPeers] = useState({});

  useEffect(() => {
    const ws = new WebSocket(WS_URL);

    ws.onmessage = (event) => {
      const peerData = JSON.parse(event.data);
      setPeers(peerData ?? {});
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div>
      <h2>Peers</h2>
      <ListGroup>
        {Object.entries(peers).map(([peerId, peerInfo]) => (
          <ListGroup.Item key={peerId}>
            <h5>Peer ID: {peerId}</h5>
            <p><strong>Last Updated:</strong> {new Date(peerInfo.Timestamp).toLocaleString()}</p>
            <p><strong>CIDs:</strong></p>
            <ul>
              {peerInfo.CIDs.map((cid) => (
                <li key={cid}>{cid}</li>
              ))}
            </ul>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  );
}

export default PeerList;