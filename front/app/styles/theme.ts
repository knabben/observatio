import {createTheme, MantineColorsTuple} from '@mantine/core';
import {openSans} from '@/app/styles/fonts';

/**
 * Single accent scale consolidating the green hex literals previously duplicated
 * across nav-links, tabs, and detail headers (#aaf16a, #48654a, ...). Index 4 is
 * the bright CTA/active-nav lime; index 8 is the muted dark-green heading/tab accent.
 */
const brand: MantineColorsTuple = [
  '#f2fbe6',
  '#e3f5c4',
  '#d0f0a0',
  '#bcec7d',
  '#aaf16a',
  '#8fdb52',
  '#6ec93f',
  '#54a838',
  '#48654a',
  '#2f4a34',
];

/**
 * Tri-state status semantics used by `StatusIndicator`. Mantine's built-in
 * green/red/gray already carry the right contrast in both color schemes.
 */
export const STATUS_COLORS = {
  healthy: 'green',
  degraded: 'orange',
  notready: 'red',
  unknown: 'gray',
} as const;

export const theme = createTheme({
  primaryColor: 'brand',
  colors: {brand},
  fontFamily: openSans.style.fontFamily,
  headings: {fontFamily: openSans.style.fontFamily},
});
