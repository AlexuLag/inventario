import React from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Box,
  Container,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Divider,
} from '@mui/material'

const Menu = () => {
  const navigate = useNavigate()

  const handleInventory = () => {
    navigate('/products')
  }

  const handleInvoice = () => {
    // TODO: Implementar navegación a crear factura
    console.log('Crear factura')
  }

  const handleLogout = () => {
    // TODO: Implementar limpieza de estado/sesión si es necesario
    navigate('/register')
  }

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
              <ListItemButton onClick={handleLogout}>
                <ListItemText primary="Salir" />
              </ListItemButton>
            </ListItem>
          </List>
        </Paper>
      </Box>
    </Container>
  )
}

export default Menu 