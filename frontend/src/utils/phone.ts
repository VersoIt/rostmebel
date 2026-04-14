export const formatRussianPhone = (value: string) => {
  let input = value.replace(/\D/g, '');

  if (input.startsWith('7')) input = input.substring(1);
  if (input.startsWith('8')) input = input.substring(1);

  let formatted = '+7 ';
  if (input.length > 0) formatted += `(${input.substring(0, 3)}`;
  if (input.length >= 4) formatted += `) ${input.substring(3, 6)}`;
  if (input.length >= 7) formatted += `-${input.substring(6, 8)}`;
  if (input.length >= 9) formatted += `-${input.substring(8, 10)}`;

  return formatted.substring(0, 18);
};
