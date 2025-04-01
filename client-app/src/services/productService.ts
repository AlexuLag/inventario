import { Product, CreateProductRequest } from '../types/Product';

const API_URL = 'http://localhost:8080/api/products';

export const productService = {
    async getAll(): Promise<Product[]> {
        const response = await fetch(API_URL);
        if (!response.ok) {
            throw new Error('Failed to fetch products');
        }
        return response.json();
    },

    async getById(id: string): Promise<Product> {
        const response = await fetch(`${API_URL}?id=${id}`);
        if (!response.ok) {
            throw new Error('Failed to fetch product');
        }
        return response.json();
    },

    async create(product: CreateProductRequest): Promise<Product> {
        const response = await fetch(API_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(product),
        });
        if (!response.ok) {
            throw new Error('Failed to create product');
        }
        return response.json();
    },

    async update(product: Product): Promise<Product> {
        const response = await fetch(API_URL, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(product),
        });
        if (!response.ok) {
            throw new Error('Failed to update product');
        }
        return response.json();
    },

    async delete(id: string): Promise<void> {
        const response = await fetch(`${API_URL}?id=${id}`, {
            method: 'DELETE',
        });
        if (!response.ok) {
            throw new Error('Failed to delete product');
        }
    },
}; 