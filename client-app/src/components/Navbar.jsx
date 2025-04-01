import React from 'react'
import { useNavigate } from 'react-router-dom'
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  IconButton,
} from '@mui/material'
import InventoryIcon from '@mui/icons-material/Inventory'
import ReceiptIcon from '@mui/icons-material/Receipt'
import LogoutIcon from '@mui/icons-material/Logout'

const Navbar = () => {
  const navigate = useNavigate()

  const handleInventory = () => {
    navigate('/products')
  }

  const handleInvoice = () => {
    // TODO: Implementar navegación a crear factura
    console.log('Crear factura')
  }

  const handleLogout = () => {
    navigate('/register')
  }

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          Inventory System
        </Typography>
        <Box>
          <IconButton color="inherit" onClick={handleInventory}>
            <InventoryIcon />
          </IconButton>
          <IconButton color="inherit" onClick={handleInvoice}>
            <ReceiptIcon />
          </IconButton>
          <IconButton color="inherit" onClick={handleLogout}>
            <LogoutIcon />
          </IconButton>
        </Box>
      </Toolbar>
    </AppBar>
  )
}

export default Navbar 