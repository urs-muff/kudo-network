// src/components/ConceptList.js
import React, { useState, useEffect } from 'react';
import { ListGroup } from 'react-bootstrap';

const WS_URL = 'ws://localhost:9090/ws';

function ConceptList() {
  const [concepts, setConcepts] = useState([]);

  useEffect(() => {
    const ws = new WebSocket(WS_URL);

    ws.onmessage = (event) => {
      const concept = JSON.parse(event.data);
      setConcepts((prevConcepts) => {
        const index = prevConcepts.findIndex((c) => c.Guid === concept.Guid);
        if (index !== -1) {
          const updatedConcepts = [...prevConcepts];
          updatedConcepts[index] = concept;
          return updatedConcepts;
        } else {
          return [...prevConcepts, concept];
        }
      });
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div>
      <h2>Concepts</h2>
      <ListGroup>
        {concepts.map((concept) => (
          <ListGroup.Item key={concept.Guid}>
            <h5>{concept.Name}</h5>
            <p><strong>Type:</strong> {concept.Type}</p>
            <p><strong>Description:</strong> {concept.Description}</p>
            <p><strong>CID:</strong> {concept.Cid}</p>
            <p><strong>Content:</strong> {concept.Content}</p>
            <p><strong>Timestamp:</strong> {new Date(concept.Timestamp).toLocaleString()}</p>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  );
}

export default ConceptList;
