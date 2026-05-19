import type { Metadata } from "next";
import "./globals.scss";

export const metadata: Metadata = {
  title: "Harness Chat",
  description: "Harness 聊天前端（NextChat 迁移）",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  );
}
