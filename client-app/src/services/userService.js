const API_URL = 'http://localhost:8080/api'

export const createUser = async (userData) => {
  try {
    const response = await fetch(`${API_URL}/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userData),
    })
    if (!response.ok) {
      throw new Error('Failed to create user')
    }
    return await response.json()
  } catch (error) {
    throw new Error(error.message || 'Error creating user')
  }
}

export const getUserById = async (id) => {
  try {
    const response = await fetch(`${API_URL}/users/${id}`)
    if (!response.ok) {
      throw new Error('Failed to fetch user')
    }
    return await response.json()
  } catch (error) {
    throw new Error(error.message || 'Error fetching user')
  }
}

export const getAllUsers = async () => {
  try {
    const response = await fetch(`${API_URL}/users`)
    if (!response.ok) {
      throw new Error('Failed to fetch users')
    }
    return await response.json()
  } catch (error) {
    throw new Error(error.message || 'Error fetching users')
  }
}

export const updateUser = async (id, userData) => {
  try {
    const response = await fetch(`${API_URL}/users/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userData),
    })
    if (!response.ok) {
      throw new Error('Failed to update user')
    }
    return await response.json()
  } catch (error) {
    throw new Error(error.message || 'Error updating user')
  }
}

export const deleteUser = async (id) => {
  try {
    const response = await fetch(`${API_URL}/users/${id}`, {
      method: 'DELETE',
    })
    if (!response.ok) {
      throw new Error('Failed to delete user')
    }
    return true
  } catch (error) {
    throw new Error(error.message || 'Error deleting user')
  }
} 