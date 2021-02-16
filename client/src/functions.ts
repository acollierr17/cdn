import type { ImageResult } from './api';

export const formatFileSize = (size: number): string => {
  if (size < 10240) return `${size}B`;
  else if (size < 1048576) return `${(size / 1024).toFixed(2)}KiB`;
  else if (size < 1073741824) return `${(size / 1048576).toFixed(2)}MiB`;
  else return `${(size / 1073741824).toFixed(2)}GiB`;
};

const getOrdinalString = (n: number): string => {
  return (
    n +
    (n > 0
      ? ['th', 'st', 'nd', 'rd'][(n > 3 && n < 21) || n % 10 > 3 ? 0 : n % 10]
      : '')
  );
};

const months = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'May',
  'Jun',
  'Jul',
  'Aug',
  'Sep',
  'Oct',
  'Nov',
  'Dec',
];

export const formatDate = (date: ImageResult['last_modified']): string => {
  const d = new Date(date);
  return `${months[d.getMonth()]} ${getOrdinalString(
    d.getDate(),
  )}, ${d.getFullYear()}`;
};
