import { ref, onMounted } from 'vue';
import type { Product } from '@/types';

const favorites = ref<Product[]>(JSON.parse(localStorage.getItem('favorites') || '[]'));

export function useFavorites() {
  const toggleFavorite = (product: Product) => {
    const index = favorites.value.findIndex(p => p.id === product.id);
    if (index === -1) {
      favorites.value.push(product);
    } else {
      favorites.value.splice(index, 1);
    }
    localStorage.setItem('favorites', JSON.stringify(favorites.value));
  };

  const isFavorite = (id: number) => {
    return favorites.value.some(p => p.id === id);
  };

  return {
    favorites,
    toggleFavorite,
    isFavorite,
  };
}
