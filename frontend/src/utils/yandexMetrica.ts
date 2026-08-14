import { MESSENGER_LINKS, PHONE_HREF } from '@/constants/contacts';

const METRIKA_TAG_URL = 'https://mc.yandex.ru/metrika/tag.js';
const DEFAULT_COUNTER_ID = '111598091';
const rawCounterId = String(__YANDEX_METRIKA_ID__ || DEFAULT_COUNTER_ID).trim();
const parsedCounterId = Number.parseInt(rawCounterId, 10);
const counterId = Number.isFinite(parsedCounterId) && parsedCounterId > 0 ? parsedCounterId : null;

type PrimitiveMetricValue = string | number | boolean | null | undefined;
type MetricParams = Record<string, PrimitiveMetricValue>;
type MetricChannel = 'phone' | 'whatsapp' | 'telegram' | 'max';
type LeadSource = 'quote_quiz' | 'order_form';

type YandexMetricaQueue = ((...args: unknown[]) => void) & {
  a?: unknown[][];
  l?: number;
};

declare global {
  interface Window {
    ym?: YandexMetricaQueue;
    dataLayer?: unknown[];
  }
}

let isInitialized = false;
let lastTrackedUrl = '';
let isContactClickTrackingBound = false;

const canUseBrowser = () => typeof window !== 'undefined' && typeof document !== 'undefined';

const sanitizeParams = (params?: MetricParams) => {
  if (!params) return undefined;

  const entries = Object.entries(params).filter(([, value]) => value !== undefined && value !== null && value !== '');
  return entries.length ? Object.fromEntries(entries) : undefined;
};

const ensureYmQueue = () => {
  if (!canUseBrowser()) return undefined;

  if (typeof window.ym === 'function') {
    return window.ym;
  }

  const ym = ((...args: unknown[]) => {
    (ym.a = ym.a || []).push(args);
  }) as YandexMetricaQueue;

  ym.a = [];
  ym.l = Date.now();
  window.ym = ym;

  return ym;
};

const ensureScriptTag = () => {
  if (!canUseBrowser()) return;
  if (document.querySelector(`script[src="${METRIKA_TAG_URL}"]`)) return;

  const script = document.createElement('script');
  script.async = true;
  script.src = METRIKA_TAG_URL;
  document.head.appendChild(script);
};

const normalizeTrackedPath = (path?: string) => {
  const currentPath = path || `${window.location.pathname}${window.location.search}`;
  const [withoutHash] = currentPath.split('#');

  if (!withoutHash) {
    return `${window.location.pathname}${window.location.search}`;
  }

  return withoutHash.startsWith('/') ? withoutHash : `/${withoutHash}`;
};

const CONTACT_CHANNEL_MATCHERS: Array<{ channel: MetricChannel; matches: (href: string) => boolean }> = [
  { channel: 'phone', matches: (href) => href.startsWith(PHONE_HREF) },
  { channel: 'whatsapp', matches: (href) => href.startsWith(MESSENGER_LINKS.whatsapp) },
  { channel: 'telegram', matches: (href) => href.startsWith(MESSENGER_LINKS.telegram) },
  { channel: 'max', matches: (href) => href.startsWith(MESSENGER_LINKS.max) },
];

const bindContactClickTracking = () => {
  if (!canUseBrowser() || isContactClickTrackingBound) return;

  document.addEventListener('click', (event) => {
    const target = event.target instanceof Element ? event.target : null;
    const link = target?.closest('a[href]') as HTMLAnchorElement | null;
    if (!link) return;

    const href = link.getAttribute('href')?.trim() || '';
    const matchedChannel = CONTACT_CHANNEL_MATCHERS.find((entry) => entry.matches(href));
    if (!matchedChannel) return;

    trackContactClick(matchedChannel.channel, normalizeTrackedPath());
  });

  isContactClickTrackingBound = true;
};

export const initYandexMetrica = () => {
  if (!canUseBrowser() || !counterId || isInitialized) return;

  const ym = ensureYmQueue();
  if (!ym) return;

  ensureScriptTag();
  ym(counterId, 'init', {
    clickmap: true,
    trackLinks: true,
    accurateTrackBounce: true,
    webvisor: true,
    defer: true,
  });
  bindContactClickTracking();

  isInitialized = true;
};

export const trackYandexPageView = (path?: string, title?: string) => {
  if (!canUseBrowser() || !counterId) return;

  const ym = ensureYmQueue();
  if (!ym) return;

  const absoluteUrl = new URL(normalizeTrackedPath(path), window.location.origin).toString();
  if (lastTrackedUrl === absoluteUrl) return;

  const params = sanitizeParams({
    title: title || document.title || undefined,
    referer: lastTrackedUrl || document.referrer || undefined,
  });

  ym(counterId, 'hit', absoluteUrl, params);
  lastTrackedUrl = absoluteUrl;
};

export const trackYandexGoal = (target: string, params?: MetricParams) => {
  if (!canUseBrowser() || !counterId || !target.trim()) return;

  const ym = ensureYmQueue();
  if (!ym) return;

  ym(counterId, 'reachGoal', target, sanitizeParams(params));
};

export const trackContactClick = (channel: MetricChannel, placement: string) => {
  trackYandexGoal('contact_click', {
    channel,
    placement,
  });
};

export const trackLeadSubmitted = (
  source: LeadSource,
  params?: {
    budgetRange?: string;
    contactMethod?: string;
    projectId?: number;
    projectType?: string;
  },
) => {
  trackYandexGoal('lead_submit', {
    source,
    budget_range: params?.budgetRange,
    contact_method: params?.contactMethod,
    project_id: params?.projectId,
    project_type: params?.projectType,
  });
};

export const trackAISearch = (query: string, resultsCount: number, successful: boolean) => {
  trackYandexGoal('ai_search', {
    query_length: query.trim().length,
    results_count: resultsCount,
    successful,
  });
};
