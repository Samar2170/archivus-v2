'use client'
import { useState, useRef } from "react";
import { Button } from "@/components/ui/button";
import { PlusIcon } from "lucide-react";
import { addFolder } from '@/app/api/folder';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription
} from "@/components/ui/dialog";
import { uploadFiles, uploadFilesWithProgress } from "@/app/api/files";
import { Progress } from "@/components/ui/progress";


export default function FileFolderDialog({ folder }: { folder: string }) {
  const [newFiles, setNewFiles] = useState<File[]>([]);
  const [modalType, setModalType] = useState<"files" | "folder" | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [newFolderName, setNewFolderName] = useState<string>("");
  const [uploading, setUploading] = useState<boolean>(false);
  const [progress, setProgress] = useState<number>(0);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const createFolder = async () => {
    try {
      if (folder !== '') {
        await addFolder(folder + '/' + newFolderName);
      } else {
        await addFolder(newFolderName);
      }
      alert("Folder created successfully");
    } catch (err) {
      console.error("Error creating folder:", err);
      alert("Error creating folder");
    }
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {

    if (e.target.files) {
      setNewFiles(Array.from(e.target.files));
    }
  };

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newFiles.length === 0) return;
    setUploading(true);
    setProgress(0);
    try {
      const res = await uploadFilesWithProgress(folder, newFiles, (p) => setProgress(p));
      if (!res) throw new Error("Upload failed");
      alert("Files uploaded successfully!");
      setNewFiles([]);
    } catch (err) {
      console.error(err);
      alert("Error uploading files");
    } finally {
      setUploading(false);
      setModalType(null);
      setProgress(0);
      window.location.reload();
    }
  };

  return (
    <>
      <div className="fixed bottom-6 right-6">
        <Button
          onClick={() => setDialogOpen(true)}
          className="rounded-full h-14 w-14 p-0 shadow-lg"
        >
          <PlusIcon className="h-6 w-6" />
        </Button>
      </div>
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Select an option</DialogTitle>
            <DialogDescription>
              Choose what you want to upload.
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-3 mt-4">
            <Button
              onClick={() => {
                setDialogOpen(false);
                setModalType("files");
              }}
              variant="outline"
            >
              Upload Files
            </Button>
            <Button
              onClick={() => {
                setDialogOpen(false);
                setModalType("folder");
              }}
              variant="outline"
            >
              Add Folder
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={modalType === "files"} onOpenChange={() => setModalType(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Upload Files</DialogTitle>
          </DialogHeader>
          <form className="flex flex-col gap-4">
            <input type="file" multiple ref={fileInputRef} onChange={handleFileChange} className="border rounded p-2" />
            {uploading && (
              <div className="flex flex-col gap-2">
                <Progress value={progress} className="w-full" />
                <p className="text-sm text-gray-500 text-center">{Math.round(progress)}%</p>
              </div>
            )}

            <Button type="submit" disabled={uploading} onClick={(e) => handleUpload(e)}>
              {uploading ? "Uploading..." : "Submit"}
            </Button>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={modalType === "folder"} onOpenChange={() => setModalType(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add Folder</DialogTitle>
          </DialogHeader>
          <form className="flex flex-col gap-4">
            <input type="text" multiple value={newFolderName} onChange={(e) => setNewFolderName(e.target.value)} className="border rounded p-2" />
            <Button type="submit" onClick={() => createFolder()}>Submit</Button>
          </form>
        </DialogContent>
      </Dialog>
    </>
  )
}
