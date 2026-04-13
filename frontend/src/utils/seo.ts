const DEFAULT_SITE_URL = 'https://rostmebel.shop';
const DEFAULT_TITLE = 'РОСТ Мебель — кухни и корпусная мебель по размеру';
const DEFAULT_DESCRIPTION =
  'Проектируем, производим и устанавливаем кухни, шкафы и системы хранения по Крыму. Смета, договор, доставка и монтаж под ключ.';
const DEFAULT_IMAGE = '/assets/images/hero-1.jpg';

const publicSiteURL = (import.meta.env.VITE_PUBLIC_SITE_URL || DEFAULT_SITE_URL).replace(/\/+$/, '');

type PageSeo = {
  title?: string;
  description?: string;
  path?: string;
  image?: string;
  imageAlt?: string;
  type?: string;
  robots?: string;
};

const isAbsoluteURL = (value: string) => /^https?:\/\//i.test(value);

export const absoluteUrl = (path = '/') => {
  if (isAbsoluteURL(path)) return path;
  return `${publicSiteURL}${path.startsWith('/') ? path : `/${path}`}`;
};

export const compactDescription = (value: string, maxLength = 155) => {
  const normalized = value.replace(/\s+/g, ' ').trim();
  if (normalized.length <= maxLength) return normalized;
  return `${normalized.slice(0, maxLength - 1).trimEnd()}…`;
};

const upsertMeta = (selector: string, attrs: Record<string, string>, content: string) => {
  let meta = document.head.querySelector(selector) as HTMLMetaElement | null;
  if (!meta) {
    meta = document.createElement('meta');
    Object.entries(attrs).forEach(([name, value]) => meta?.setAttribute(name, value));
    document.head.appendChild(meta);
  }
  meta.setAttribute('content', content);
};

const upsertCanonical = (url: string) => {
  let link = document.head.querySelector('link[rel="canonical"]') as HTMLLinkElement | null;
  if (!link) {
    link = document.createElement('link');
    link.rel = 'canonical';
    document.head.appendChild(link);
  }
  link.href = url;
};

export const setPageSeo = ({
  title = DEFAULT_TITLE,
  description = DEFAULT_DESCRIPTION,
  path = '/',
  image = DEFAULT_IMAGE,
  imageAlt = 'Кухни и корпусная мебель РОСТ Мебель',
  type = 'website',
  robots = 'index,follow',
}: PageSeo) => {
  const canonical = absoluteUrl(path);
  const imageURL = absoluteUrl(image);
  const safeDescription = compactDescription(description);

  document.title = title;
  upsertCanonical(canonical);
  upsertMeta('meta[name="description"]', { name: 'description' }, safeDescription);
  upsertMeta('meta[name="robots"]', { name: 'robots' }, robots);
  upsertMeta('meta[property="og:type"]', { property: 'og:type' }, type);
  upsertMeta('meta[property="og:title"]', { property: 'og:title' }, title);
  upsertMeta('meta[property="og:description"]', { property: 'og:description' }, safeDescription);
  upsertMeta('meta[property="og:url"]', { property: 'og:url' }, canonical);
  upsertMeta('meta[property="og:image"]', { property: 'og:image' }, imageURL);
  upsertMeta('meta[property="og:image:alt"]', { property: 'og:image:alt' }, imageAlt);
  upsertMeta('meta[property="og:site_name"]', { property: 'og:site_name' }, 'РОСТ Мебель');
  upsertMeta('meta[property="og:locale"]', { property: 'og:locale' }, 'ru_RU');
  upsertMeta('meta[name="twitter:card"]', { name: 'twitter:card' }, 'summary_large_image');
  upsertMeta('meta[name="twitter:title"]', { name: 'twitter:title' }, title);
  upsertMeta('meta[name="twitter:description"]', { name: 'twitter:description' }, safeDescription);
  upsertMeta('meta[name="twitter:image"]', { name: 'twitter:image' }, imageURL);
};

export const setJsonLd = (id: string, payload: unknown) => {
  let script = document.getElementById(id) as HTMLScriptElement | null;
  if (!script) {
    script = document.createElement('script');
    script.id = id;
    script.type = 'application/ld+json';
    document.head.appendChild(script);
  }
  script.textContent = JSON.stringify(payload);
};

export const removeJsonLd = (id: string) => {
  document.getElementById(id)?.remove();
};

export const buildBusinessSchema = () => ({
  '@context': 'https://schema.org',
  '@type': 'FurnitureStore',
  name: 'РОСТ Мебель',
  url: absoluteUrl('/'),
  logo: absoluteUrl('/assets/logo.svg'),
  image: absoluteUrl(DEFAULT_IMAGE),
  telephone: '+7 978 763-16-03',
  email: 'rost.salon2003@mail.ru',
  priceRange: '₽₽',
  address: {
    '@type': 'PostalAddress',
    addressLocality: 'Симферополь',
    addressRegion: 'Крым',
    addressCountry: 'RU',
  },
  areaServed: ['Симферополь', 'Крым'],
});
