// src/components/OwnerForm.js
import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';
import axios from 'axios';

const API_URL = 'http://localhost:9090';

function OwnerForm() {
  const [owner, setOwner] = useState({
    name: '',
    description: '',
  });

  useEffect(() => {
    fetchOwner();
  }, []);

  const fetchOwner = async () => {
    try {
      const response = await axios.get(`${API_URL}/owner`);
      setOwner(response.data);
    } catch (error) {
      console.error('Error fetching owner:', error);
    }
  };

  const handleChange = (e) => {
    setOwner({ ...owner, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post(`${API_URL}/owner`, owner);
      console.log('Owner updated:', response.data);
      fetchOwner();
    } catch (error) {
      console.error('Error updating owner:', error);
    }
  };

  return (
    <Form onSubmit={handleSubmit} className="mb-4">
      <Form.Group>
        <Form.Label>Name</Form.Label>
        <Form.Control
          type="text"
          name="name"
          value={owner.name}
          onChange={handleChange}
          required
        />
      </Form.Group>
      <Form.Group>
        <Form.Label>Description</Form.Label>
        <Form.Control
          as="textarea"
          name="description"
          value={owner.description}
          onChange={handleChange}
          required
        />
      </Form.Group>
      <Button variant="primary" type="submit">
        Update Owner
      </Button>
    </Form>
  );
}

export default OwnerForm;
