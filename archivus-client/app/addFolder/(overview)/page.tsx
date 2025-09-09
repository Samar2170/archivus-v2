'use client';
import { useEffect, useState } from 'react';
import { FolderResponse, listFolders,addFolder } from '@/app/api/folder';
import { time } from 'console';
import { useAuth } from '@/app/hooks/useAuth';

export default function AddFolderForm() {
    useAuth();
    const [folders, setFolders] = useState<FolderResponse[]>([]);
    const [newFolderName, setNewFolderName] = useState('');
    const [selectedFolder, setSelectedFolder] = useState<string>('');
    useEffect(() => {
        const delay = (ms:number) => new Promise(res => setTimeout(res, ms));
        const fetchFolders = async () => {
            try {
                const data = await listFolders();
                setFolders(data);
              } catch (err) {
                console.error("Error fetching folders:", err);
              }
            };
            fetchFolders();    
    },[]);

    const createFolder = async() => {
        let finalFolderName = '';
        if (selectedFolder !== ''){
            finalFolderName = `${selectedFolder}/${newFolderName}`;
        } else {
            finalFolderName = newFolderName;
        }
        try {
            await addFolder(selectedFolder + '/' + newFolderName);
            alert("Folder created successfully");
            setNewFolderName('');
        } catch (err) {
            console.error("Error creating folder:", err);
            alert("Error creating folder");
        } 
    }

    return (

    <div className="flex justify-center items-center h-24">
    <div className="p-4 border rounded-lg flex flex-col gap-4 w-96 mt-32">
    <select className="border p-2 rounded"
        value={selectedFolder}
        onChange={(e) => {
            setSelectedFolder(e.target.value)
        }}
      >
        <option value="" disabled>Select a folder</option>
        {folders?.map((folder) => (
            <option key={folder.Path} value={folder.Path}>{folder.Path}</option>
        ))}
      </select>
      <input type="text"
        className="border p-2 rounded"
        placeholder="New Folder Name"
        value={newFolderName}
        onChange={(e) => setNewFolderName(e.target.value)}
        />
        <button
        onClick={createFolder}
        className="bg-blue-500 text-white p-2 rounded hover:bg-blue-600"
        > Add Folder</button>

    </div>
    </div>
    );
}