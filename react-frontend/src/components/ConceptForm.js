// src/components/ConceptForm.js
import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import axios from 'axios';

const API_URL = 'http://localhost:9090';

function ConceptForm() {
  const [concept, setConcept] = useState({
    name: '',
    description: '',
    type: '',
    content: '',
  });

  const handleChange = (e) => {
    setConcept({ ...concept, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post(`${API_URL}/concept`, concept);
      console.log('Concept added:', response.data);
      setConcept({ name: '', description: '', type: '', content: '' });
    } catch (error) {
      console.error('Error adding concept:', error);
    }
  };

  return (
    <Form onSubmit={handleSubmit} className="mb-4">
      <Form.Group>
        <Form.Label>Name</Form.Label>
        <Form.Control
          type="text"
          name="name"
          value={concept.name}
          onChange={handleChange}
          required
        />
      </Form.Group>
      <Form.Group>
        <Form.Label>Description</Form.Label>
        <Form.Control
          as="textarea"
          name="description"
          value={concept.description}
          onChange={handleChange}
          required
        />
      </Form.Group>
      <Form.Group>
        <Form.Label>Type</Form.Label>
        <Form.Control
          type="text"
          name="type"
          value={concept.type}
          onChange={handleChange}
          required
        />
      </Form.Group>
      <Form.Group>
        <Form.Label>Content</Form.Label>
        <Form.Control
          as="textarea"
          name="content"
          value={concept.content}
          onChange={handleChange}
          required
        />
      </Form.Group>
      <Button variant="primary" type="submit">
        Add Concept
      </Button>
    </Form>
  );
}

export default ConceptForm;