import React, { useState } from 'react';
import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  Paper,
  SelectChangeEvent,
} from '@mui/material';
import { createUser } from '../services/userService';
import Menu from './Menu';

interface RegistrationFormData {
  name: string;
  email: string;
  password: string;
  role: string;
}

const UserRegistration: React.FC = () => {
  const [formData, setFormData] = useState<RegistrationFormData>({
    name: '',
    email: '',
    password: '',
    role: 'user',
  });

  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<boolean>(false);
  const [isRegistered, setIsRegistered] = useState<boolean>(false);

  const handleTextChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSelectChange = (e: SelectChangeEvent) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess(false);

    try {
      await createUser(formData);
      setSuccess(true);
      setIsRegistered(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred during registration');
    }
  };

  const handleLogout = () => {
    setIsRegistered(false);
    setFormData({
      name: '',
      email: '',
      password: '',
      role: 'user',
    });
  };

  if (isRegistered) {
    return <Menu onLogout={handleLogout} />;
  }

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Paper elevation={3} sx={{ p: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom align="center">
            User Registration
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          {success && (
            <Alert severity="success" sx={{ mb: 2 }}>
              Registration successful!
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
            <TextField
              fullWidth
              required
              label="Name"
              name="name"
              value={formData.name}
              onChange={handleTextChange}
              margin="normal"
            />

            <TextField
              fullWidth
              required
              type="email"
              label="Email"
              name="email"
              value={formData.email}
              onChange={handleTextChange}
              margin="normal"
            />

            <TextField
              fullWidth
              required
              type="password"
              label="Password"
              name="password"
              value={formData.password}
              onChange={handleTextChange}
              margin="normal"
            />

            <FormControl fullWidth margin="normal">
              <InputLabel id="role-label">Role</InputLabel>
              <Select
                labelId="role-label"
                name="role"
                value={formData.role}
                label="Role"
                onChange={handleSelectChange}
              >
                <MenuItem value="user">User</MenuItem>
                <MenuItem value="admin">Admin</MenuItem>
              </Select>
            </FormControl>

            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              Register
            </Button>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default UserRegistration; 