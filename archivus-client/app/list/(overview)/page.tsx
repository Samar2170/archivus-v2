'use client';
import { Suspense, useEffect,useState } from "react";
import { FileMetaData, ListFileMetaData, getFilesByFolder, getFilesByFolderResponse } from "../../api/files";
import { getFilesList } from "@/app/api/files";
import { useSearchParams } from "next/navigation";
import { useAuth } from "@/app/hooks/useAuth";


export default function Page() {
  useAuth();
    const [files, setFiles] = useState<ListFileMetaData[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const searchParams = useSearchParams();


    useEffect(() => {
        async function loadFiles() {
          try {
            setLoading(true);
            const query = searchParams.get('search') || '';
            const data = await getFilesList(query);
            setFiles(data.files || []);
          } catch (err) {
            console.error("Error fetching files", err);
          } finally {
            setLoading(false);
          }
        }
        loadFiles();
    },[searchParams]);
    
    if (loading) return <p>Loading...</p>;    
    return (
    <ul role="list" className="divide-y divide-gray-100 dark:divide-white/5">
      <Suspense fallback={<p>Loading...</p>}>
      {files.map((file) => (
        <li key={file.ID} className="flex justify-between gap-x-6 py-5">
          <div className="flex min-w-0 gap-x-4">
            <img
              alt=""
            //   src={file.SignedUrl}
              className="size-12 flex-none rounded-full bg-gray-50 dark:bg-gray-800 dark:outline dark:-outline-offset-1 dark:outline-white/10"
            />
            <div className="min-w-0 flex-auto">
              <p className="text-sm/6 font-semibold text-gray-900 dark:text-white">{file.Name}</p>
              <p className="mt-1 truncate text-xs/5 text-green-500 dark:text-gray-400">{file.FilePath}</p>
            </div>
          </div>
          <div className="hidden shrink-0 sm:flex sm:flex-col sm:items-end">
            <p className="text-sm/6 text-white dark:text-white">{file.SizeInMb}</p>
            {
            // person.lastSeen ? (
            //   <p className="mt-1 text-xs/5 text-gray-500 dark:text-gray-400">
            //     Last seen <time dateTime={person.lastSeenDateTime}>{person.lastSeen}</time>
            //   </p>
            // ) : 
            (
              <div className="mt-1 flex items-center gap-x-1.5">
                <div className="flex-none rounded-full bg-emerald-500/20 p-1 dark:bg-emerald-500/30">
                  <div className="size-1.5 rounded-full bg-emerald-500" />
                </div>
                <p className="text-xs/5 text-gray-500 dark:text-gray-400">Online</p>
              </div>
            )
            }
          </div>
        </li>
      ))}
      </Suspense>
    </ul>
        )
}