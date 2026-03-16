export const metadata = { title: 'SaaS Ideas Dashboard' };

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body style={{ fontFamily: 'Arial, sans-serif', margin: 0, background: '#f7f7f9' }}>{children}</body>
    </html>
  );
}
