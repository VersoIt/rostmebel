import api from '@/api/client';

export async function downloadFile(url: string, filename: string) {
  try {
    const response = await api.get(url, {
      responseType: 'blob',
    });
    
    const blob = new Blob([response.data], { 
      type: response.headers['content-type'] 
    });
    
    const downloadUrl = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.setAttribute('download', filename);
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(downloadUrl);
  } catch (error) {
    console.error('Download failed:', error);
    alert('Ошибка при скачивании файла');
  }
}
