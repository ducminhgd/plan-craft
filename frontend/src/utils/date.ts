import dayjs from 'dayjs';

// Global date format used across the application
export const DATE_FORMAT = 'DD-MMM-YYYY';

// Format a date value to display string
export const formatDate = (date: any): string => {
  if (!date) return '-';
  const parsed = dayjs(date);
  if (!parsed.isValid()) return '-';
  return parsed.format(DATE_FORMAT);
};

// Parse a date string to dayjs object
export const parseDate = (date: any): dayjs.Dayjs | null => {
  if (!date) return null;
  const parsed = dayjs(date);
  if (!parsed.isValid()) return null;
  return parsed;
};

// Convert dayjs to ISO string for backend
export const toISOString = (date: dayjs.Dayjs | null): string | null => {
  if (!date) return null;
  if (!date.isValid()) return null;
  return date.toISOString();
};
