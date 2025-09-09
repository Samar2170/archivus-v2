"use client";

import { uploadFiles } from "@/app/api/files";
import { FolderResponse,listFolders } from "@/app/api/folder";
import { useAuth } from "@/app/hooks/useAuth";
import { useState, useRef, useEffect } from "react";

export default function FileUploadForm() {
  useAuth();
  const [files, setFiles] = useState<File[]>([]);
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [folders, setFolders] = useState<FolderResponse[]>([]);
  const [selectedFolder, setSelectedFolder] = useState<string>('');

  useEffect(() => {
    const fetchFolders = async () => {
      try {
        const data = await listFolders();
        setFolders(data);
      } catch (err) {
        console.error("Error fetching folders:", err);
      }
    };
    fetchFolders();

  },[])

  const handleButtonClick = () => {
    fileInputRef.current?.click(); 
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    
    if (e.target.files) {
      setFiles(Array.from(e.target.files));
    }
  };

//   const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
//     if (e.target.files) {
//       setFiles(Array.from(e.target.files));
//     }
//   };

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (files.length === 0) return;
    setUploading(true);
    try {
        const res = await uploadFiles(selectedFolder,files);
      if (!res) throw new Error("Upload failed");
      alert("Files uploaded successfully!");
      setFiles([]);
    } catch (err) {
      console.error(err);
      alert("Error uploading files");
    } finally {
      setUploading(false);
    }
  };

  return (
    <>
    <div className="flex justify-center items-center h-24">
    <div className="p-4 border rounded-lg flex flex-col gap-4 w-96 mt-32">
    <select className="border p-2 rounded"
        value={selectedFolder}
        onChange={(e) => setSelectedFolder(e.target.value)}
      >
        <option value="">Select Folder</option>
        {folders?.map((folder) => (
            <option key={folder.Path} value={folder.Path}>{folder.Path}</option>

        ))}
      </select>
      
      {/* Hidden input */}
      <input
        type="file"
        multiple
        ref={fileInputRef}
        onChange={handleFileChange}
        className="hidden"
      />

      {/* Styled button that triggers input */}
      <button
        type="button"
        onClick={handleButtonClick}
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
      >
        Select Files
      </button>

      {files.length > 0 && (
        <ul className="text-sm text-gray-700">
          {files.map((file) => (
            <li key={file.name}>📄 {file.name}</li>
          ))}
        </ul>
      )}
    
        <button
            onClick={handleUpload}
            disabled={uploading || files.length === 0}
            className={`px-4 py-2 rounded text-white ${uploading || files.length === 0 ? 'bg-gray-400 cursor-not-allowed' : 'bg-green-600 hover:bg-green-700'}`}
        > Upload Files </button>
    </div>
    </div>

    </>
  );
}


