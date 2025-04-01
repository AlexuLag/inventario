export interface Product {
    id: string;
    name: string;
    code: string;
    created_at: string;
    updated_at: string;
    image_url: string;
}

export interface CreateProductRequest {
    name: string;
    code: string;
    image_url: string;
} 