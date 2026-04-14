import type { AxiosError } from 'axios';

export interface ApiErrorPayload {
  code: string;
  message: string;
  meta?: Record<string, unknown>;
}

const fieldLabels: Record<string, string> = {
  client_name: 'имя',
  client_phone: 'телефон',
  client_email: 'email',
  comment: 'комментарий',
  project_type: 'тип проекта',
  budget_range: 'бюджет',
  city: 'город',
  contact_method: 'удобный способ связи',
  username: 'логин',
  password: 'пароль',
  refresh_token: 'сессия',
  name: 'название',
  slug: 'URL',
  price: 'стоимость',
  query: 'запрос',
  rating: 'оценка',
};

const messages: Record<string, string> = {
  INVALID_JSON: 'Не получилось прочитать данные. Проверьте форму и попробуйте еще раз.',
  INVALID_ID: 'Некорректный идентификатор записи.',
  INVALID_QUERY: 'Некорректные параметры фильтра.',
  UNAUTHORIZED: 'Сессия закончилась. Войдите заново.',
  AUTH_INVALID_CREDENTIALS: 'Неверный логин или пароль.',
  AUTH_INVALID_REFRESH_TOKEN: 'Сессия закончилась. Войдите заново.',
  PROJECT_NOT_FOUND: 'Проект не найден или был снят с публикации.',
  UPLOAD_FILE_MISSING: 'Выберите изображение для загрузки.',
  UPLOAD_INVALID_TYPE: 'Загрузите изображение в формате JPG, PNG, WebP или GIF.',
  ORDER_IP_BLOCKED: 'Заявки с этого адреса временно ограничены. Напишите нам напрямую в WhatsApp, Telegram или MAX.',
  ORDER_RATE_LIMITED: 'Слишком много заявок за сутки. Напишите нам напрямую или попробуйте позже.',
  ORDER_NOT_FOUND: 'Заявка не найдена. Обновите список и попробуйте еще раз.',
  REVIEW_INVALID_PHONE: 'Укажите корректный номер телефона.',
  REVIEW_NOT_ALLOWED: 'Отзыв можно оставить только по номеру завершенного заказа.',
  INTERNAL_ERROR: 'Сервис временно недоступен. Мы уже разбираемся.',
  NETWORK_ERROR: 'Нет связи с сервером. Проверьте интернет и попробуйте еще раз.',
};

export function parseApiError(error: unknown): ApiErrorPayload {
  const axiosError = error as AxiosError<{ error?: ApiErrorPayload | string }>;
  const payload = axiosError.response?.data?.error;

  if (payload && typeof payload === 'object' && 'code' in payload) {
    return payload;
  }

  if (typeof payload === 'string') {
    return {
      code: 'LEGACY_ERROR',
      message: payload,
    };
  }

  if (axiosError.request && !axiosError.response) {
    return {
      code: 'NETWORK_ERROR',
      message: 'Network error',
    };
  }

  return {
    code: 'INTERNAL_ERROR',
    message: 'Unexpected error',
  };
}

export function getApiErrorMessage(error: unknown): string {
  const apiError = parseApiError(error);

  if (apiError.code === 'VALIDATION_FAILED') {
    return validationMessage(apiError);
  }

  if (apiError.code === 'UPLOAD_FILE_TOO_LARGE') {
    const maxBytes = Number(apiError.meta?.max_bytes || 0);
    const maxMb = maxBytes ? Math.round(maxBytes / 1024 / 1024) : 10;
    return `Файл слишком большой. Максимум ${maxMb} МБ.`;
  }

  return messages[apiError.code] || apiError.message || messages.INTERNAL_ERROR;
}

function validationMessage(error: ApiErrorPayload): string {
  const fields = error.meta?.fields;
  if (!Array.isArray(fields) || fields.length === 0) {
    return 'Проверьте обязательные поля формы.';
  }

  const labels = fields
    .map((item) => {
      if (!item || typeof item !== 'object' || !('field' in item)) return '';
      const field = String(item.field);
      return fieldLabels[field] || field;
    })
    .filter(Boolean);

  if (labels.length === 0) {
    return 'Проверьте обязательные поля формы.';
  }

  return `Заполните: ${Array.from(new Set(labels)).join(', ')}.`;
}
