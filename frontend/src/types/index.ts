export type ProjectStatus = 'published' | 'draft' | 'archived';

export interface Image {
  url: string;
  is_main: boolean;
}

export interface Product {
  id: number;
  project_category_id?: number;
  name: string;
  slug: string;
  description: string;
  price: number;
  price_old?: number;
  images: Image[];
  specs: Record<string, string>;
  ai_tags: string;
  status: ProjectStatus;
  views_count: number;
  orders_count: number;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  icon: string;
  sort_order: number;
}

export type OrderStatus = 'new' | 'processing' | 'done' | 'rejected' | 'spam';

export interface Order {
  id: number;
  project_id?: number;
  project_name?: string; // Added for admin UI
  client_name: string;
  client_phone: string;
  client_email: string;
  comment: string;
  status: OrderStatus;
  created_at: string;
  updated_at: string;
}

export interface ReviewImage {
  url: string;
}

export interface ReviewResponse {
  id: number;
  project_id?: number;
  rating: number;
  comment: string;
  images: ReviewImage[];
  status: 'pending' | 'approved' | 'rejected';
  client_name: string;
  project_name: string;
  created_at: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export interface Stats {
  projects_count: number;
  new_orders_today: number;
  total_orders: number;
  success_rate: number;
  top_projects: Array<{ id: number; name: string; count: number }>;
  orders_by_day: Array<{ date: string; count: number }>;
}
