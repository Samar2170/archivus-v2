"use client";

import Image from "next/image";
import { FileMetaData, getFilesByFolderResponse } from "@/app/api/files";
import { Folder, File as FileIcon } from "lucide-react";
import Link from "next/link";

interface FileCardProps {
  file: getFilesByFolderResponse;
}

export default function FileCard({ file }:{file:FileMetaData}) {
    const isDir = file.IsDir;
  
    return (
<div className="w-44 rounded-2xl shadow-sm p-4 flex flex-col items-center gap-3 
  hover:shadow-lg hover:-translate-y-1 transition-all duration-200 ">
  
  {isDir ? (
    <Link href={`/?folder=${file.Path}`} className="flex flex-col items-center">
      <Folder className="w-20 h-20 text-blue-500" />
    </Link>
  ) : file.SignedUrl && file.Extension.match(/(png|jpg|jpeg|gif)$/i) ? (
    <Link href={file.SignedUrl} className="flex flex-col items-center">
      <Image
        src={file.SignedUrl}
        alt={file.Name}
        width={80}
        height={80}
        className="rounded-lg object-cover w-20 h-20 shadow-sm"
      />
    </Link>
  ) : (
    <Link href={file.SignedUrl} className="flex flex-col items-center">
      <FileIcon className="w-16 h-16 text-gray-500" />
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