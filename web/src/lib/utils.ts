import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// T is expected to be the specific HTML attributes interface (e.g., HTMLInputAttributes)
// Use intersection with HTMLAttributes to ensure standard attributes are available if needed,
// but primarily we want to extend T with a ref property.
export type WithElementRef<T> = {
	ref?: HTMLElement | null;
} & T;

export type WithoutChild<T> = Omit<T, 'children'>;

export type WithoutChildrenOrChild<T> = WithoutChild<T> & { children?: any };
