"use client";

import Image from "next/image";
import { FileMetaData, getFilesByFolderResponse } from "@/app/api/files";
import { Folder, File as FileIcon } from "lucide-react";
import Link from "next/link";

interface FileCardProps {
  file: getFilesByFolderResponse;
}

const getFileIcon = (extension: string) => {
  switch (extension) {
    case ".pdf":
      return <Image src="/pdft.webp" alt="pdf" width={80} height={80} className="object-cover w-20 h-20" />
    case ".doc":
    case ".docx":
      return <Image src="/doct.webp" alt="word" width={80} height={80} className="object-cover w-20 h-20" />
    case ".xlsx":
    case ".csv":
    case ".xls":
      return <Image src="/excelt.avif" alt="excel" width={80} height={80} className="object-cover w-20 h-20" />
    default:
      return <Image src="/filet.webp" alt="file" width={80} height={80} className="object-cover w-20 h-20" />
  }
}

export default function FileCard({ file }: { file: FileMetaData }) {
  const isDir = file.IsDir;

  return (
    <div className="w-auto flex flex-col items-center gap-3 
  hover:shadow-lg hover:-translate-y-1 transition-all duration-200 ">

      {isDir ? (
        // <Link href={`/?folder=${file.Path}`} className="flex flex-col items-center">
        <Folder className="w-20 h-20 text-blue-500" />
        // </Link> 
      ) : file.SignedUrl && file.Extension.match(/(png|jpg|jpeg|gif|webp)$/i) ? (
        <Link href={file.SignedUrl} className="flex flex-col items-center">
          {file.Thumbnail && file.Thumbnail != "" ?
            <Image
              src={file.Thumbnail}
              alt={file.Name}
              width={80}
              height={80}
              className="object-cover w-20 h-20"
            />
            :
            getFileIcon(file.Extension)
          }
        </Link>
      ) : (
        <Link href={file.SignedUrl} className="flex flex-col items-center">
          {getFileIcon(file.Extension)}
        </Link>
      )}

      <div className="w-full text-center">
        <p className="text-sm font-medium truncate">{file.Name}</p>
        {!isDir && (
          <p className="text-xs text-gray-500 mt-1">
            {(file.Size).toFixed(2)} MB
          </p>
        )}
      </div>
    </div>
  );
}