'use client';
import Image from "next/image";
import { useAuth } from "./hooks/useAuth";
import FileExplorer from "./ui/fileexplorer";
import { useSearchParams } from "next/navigation";
import { Suspense } from "react";

export default function Home() {
  useAuth();
  const searchParams = useSearchParams();
  const folder = searchParams.get('folder') || '';
  return (
    <>
    <Suspense fallback={<p>Loading...</p>}>
    <FileExplorer folder={folder} />
    </Suspense>
    </>
  );
}
