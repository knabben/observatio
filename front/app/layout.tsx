import "./globals.css";
import '@mantine/core/styles.css';

import { ColorSchemeScript, MantineProvider, mantineHtmlProps, createTheme } from '@mantine/core';

export const metadata = {
  title: 'Observatio',
  description: '',
};

const theme = createTheme({
  fontFamily: 'Roboto, sans-serif',
  fontFamilyMonospace: 'Monaco, Courier, monospace',
  headings: { fontFamily: 'sans-serif' },
});

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" {...mantineHtmlProps}>
    <head>
      <ColorSchemeScript />
    </head>
    <body>
      <MantineProvider theme={theme}>{children}</MantineProvider>
    </body>
    </html>
);
}
