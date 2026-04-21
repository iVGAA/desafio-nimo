// Helpers para tratar data como "YYYY-MM-DD" sem timezone.

export function extractDateOnly(value: string): string {
  if (!value) return "";
  return value.slice(0, 10);
}

export function formatDateBRFromISODate(isoDate: string): string {
  const dateOnly = extractDateOnly(isoDate);
  if (!/^\d{4}-\d{2}-\d{2}$/.test(dateOnly)) return isoDate;

  const [year, month, day] = dateOnly.split("-");
  return `${day}/${month}/${year}`;
}

export function getTodayISODateLocal(): string {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}
