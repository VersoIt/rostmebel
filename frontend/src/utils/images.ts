export interface ResponsiveImageSource {
  src: string;
  srcset?: string;
  sizes?: string;
}

const DEFAULT_VARIANT_QUALITY = 80;
const IMAGE_VARIANT_ENDPOINT = '/api/v1/images/variant';

const normalizeUploadPath = (url: string): string | null => {
  const trimmedUrl = String(url || '').trim();
  if (!trimmedUrl || trimmedUrl.startsWith('data:') || trimmedUrl.startsWith('blob:')) {
    return null;
  }

  try {
    const baseOrigin = typeof window !== 'undefined' ? window.location.origin : 'https://rostmebel.shop';
    const parsedUrl = new URL(trimmedUrl, baseOrigin);
    return parsedUrl.pathname.startsWith('/uploads/') ? parsedUrl.pathname : null;
  } catch {
    return trimmedUrl.startsWith('/uploads/') ? trimmedUrl : null;
  }
};

export const buildImageVariantUrl = (url: string, width: number, quality = DEFAULT_VARIANT_QUALITY) => {
  const normalizedUploadPath = normalizeUploadPath(url);
  if (!normalizedUploadPath) {
    return url;
  }

  const searchParams = new URLSearchParams({
    src: normalizedUploadPath,
    w: String(Math.max(1, Math.round(width))),
    q: String(Math.max(40, Math.min(95, Math.round(quality)))),
  });

  return `${IMAGE_VARIANT_ENDPOINT}?${searchParams.toString()}`;
};

export const buildResponsiveImageSet = (
  url: string,
  widths: number[],
  sizes: string,
  quality = DEFAULT_VARIANT_QUALITY,
): ResponsiveImageSource => {
  const normalizedWidths = Array.from(
    new Set(widths.map((width) => Math.max(1, Math.round(width))).filter(Boolean)),
  ).sort((left, right) => left - right);

  if (!normalizedWidths.length) {
    return { src: url };
  }

  const normalizedUploadPath = normalizeUploadPath(url);
  if (!normalizedUploadPath) {
    return { src: url, sizes };
  }

  return {
    src: buildImageVariantUrl(normalizedUploadPath, normalizedWidths[normalizedWidths.length - 1], quality),
    srcset: normalizedWidths
      .map((width) => `${buildImageVariantUrl(normalizedUploadPath, width, quality)} ${width}w`)
      .join(', '),
    sizes,
  };
};
