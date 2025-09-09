import { useEffect,useState } from "react";
import { FileMetaData, getFilesByFolder, getFilesByFolderResponse } from "../api/files";
import FileCard from "./components/filecard";

const generateRandomId = () => {
    return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
}

export default function FileExplorer({folder}: {folder:string}) {
    const [files, setFiles] = useState<FileMetaData[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [size, setSize] = useState<number>(0);
    useEffect(() => {
        async function loadFiles() {
          try {
            setLoading(true);
            const data = await getFilesByFolder(folder);
            for (const file of data.files) {
              if (file.IsDir) {
                file.ID = generateRandomId();
              }
            }
            setFiles(data.files || []);
            setSize(data.size);
          } catch (err) {
            console.error("Error fetching files", err);
          } finally {
            setLoading(false);
          }
        }
        loadFiles();
    },[folder]);

    if (loading) return <p>Loading...</p>;

    return (
      <div className="p-4">
        <ul role="list" className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
      {files.map((file) => (
        <FileCard key={file.ID} file={file} />
      ))}
    </ul>      
      </div>
    );

}