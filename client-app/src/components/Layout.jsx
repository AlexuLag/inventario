import React from 'react'
import { Box, Container } from '@mui/material'
import Navbar from './Navbar'

const Layout = ({ children }) => {
  return (
    <Box sx={{ 
      display: 'flex', 
      flexDirection: 'column', 
      minHeight: '100vh',
      width: '100vw',
      overflow: 'hidden'
    }}>
      <Navbar />
      <Box component="main" sx={{ 
        flex: 1,
        width: '100%',
        height: 'calc(100vh - 64px)', // 64px es la altura del AppBar
        overflow: 'auto',
        p: 3
      }}>
        {children}
      </Box>
    </Box>
  )
}

export default Layout 