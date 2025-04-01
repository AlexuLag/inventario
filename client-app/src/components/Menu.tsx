import React from 'react';
import {
  Box,
  Button,
  Container,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Divider,
} from '@mui/material';

interface MenuProps {
  onLogout: () => void;
}

const Menu: React.FC<MenuProps> = ({ onLogout }) => {
  const handleInventory = () => {
    // TODO: Implementar navegación al inventario
    console.log('Ver inventario');
  };

  const handleInvoice = () => {
    // TODO: Implementar navegación a crear factura
    console.log('Crear factura');
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Paper elevation={3} sx={{ p: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom align="center">
            Menú Principal
          </Typography>

          <List>
            <ListItem disablePadding>
              <ListItemButton onClick={handleInventory}>
                <ListItemText primary="Ver Inventario" />
              </ListItemButton>
            </ListItem>
            <Divider />
            <ListItem disablePadding>
              <ListItemButton onClick={handleInvoice}>
                <ListItemText primary="Crear Factura" />
              </ListItemButton>
            </ListItem>
            <Divider />
            <ListItem disablePadding>
              <ListItemButton onClick={onLogout}>
                <ListItemText primary="Salir" />
              </ListItemButton>
            </ListItem>
          </List>
        </Paper>
      </Box>
    </Container>
  );
};

export default Menu; 