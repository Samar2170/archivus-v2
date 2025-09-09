'use client';
import { useSearchParams } from 'next/navigation';
import { Suspense } from 'react';

function NotFoundContent() {
  const searchParams = useSearchParams();
  const ref = searchParams.get('ref');
  return <p>Page not found. Ref: {ref}</p>;
}

export default function Page() {
  return (
    <Suspense fallback={<p>Loading...</p>}>
      <NotFoundContent />
    </Suspense>
  );
}