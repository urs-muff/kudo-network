// src/App.js
import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import ConceptForm from './components/ConceptForm';
import ConceptList from './components/ConceptList';
import PeerList from './components/PeerList';
import OwnerForm from './components/OwnerForm';

function App() {
  return (
    <Container className="mt-5">
      <h1 className="text-center mb-4">IPFS Concept Manager</h1>
      <Row>
        <Col md={6}>
          <OwnerForm />
          <ConceptForm />
          <ConceptList />
        </Col>
        <Col md={6}>
          <PeerList />
        </Col>
      </Row>
    </Container>
  );
}

export default App;