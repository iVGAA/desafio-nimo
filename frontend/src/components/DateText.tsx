import { formatDateBRFromISODate } from "@/lib/date";

type DateTextProps = {
  value: string;
};

export default function DateText({ value }: DateTextProps) {
  return <>{formatDateBRFromISODate(value)}</>;
}
