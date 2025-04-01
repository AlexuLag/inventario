const API_URL = 'http://localhost:8080/api'

export const getAllProducts = async () => {
  try {
    const response = await fetch(`${API_URL}/products`)
    if (!response.ok) {
      throw new Error('Failed to fetch products')
    }
    return await response.json()
  } catch (error) {
    throw new Error(error.message || 'Error fetching products')
  }
}

export const getProductById = async (id) => {
  try {
    const response = await fetch(`${API_URL}/products/${id}`)
    if (!response.ok) {
      throw new Error('Failed to fetch product')
    }
    return await response.json()
  } catch (error) {
    throw new Error(error.message || 'Error fetching product')
  }
}

export const createProduct = async (productData) => {
  try {
    const response = await fetch(`${API_URL}/products`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(productData),
    })
    if (!response.ok) {
      throw new Error('Failed to create product')
    }
    return await response.json()
  } catch (error) {
    throw new Error(error.message || 'Error creating product')
  }
}

export const updateProduct = async (id, productData) => {
  try {
    const response = await fetch(`${API_URL}/products/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(productData),
    })
    if (!response.ok) {
      throw new Error('Failed to update product')
    }
    return await response.json()
  } catch (error) {
    throw new Error(error.message || 'Error updating product')
  }
}

export const deleteProduct = async (id) => {
  try {
    const response = await fetch(`${API_URL}/products/${id}`, {
      method: 'DELETE',
    })
    if (!response.ok) {
      throw new Error('Failed to delete product')
    }
    return true
  } catch (error) {
    throw new Error(error.message || 'Error deleting product')
  }
} 