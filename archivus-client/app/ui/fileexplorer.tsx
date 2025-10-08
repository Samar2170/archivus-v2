import { useEffect,useState, useRef } from "react";
import { FileMetaData, getFilesByFolder, getFilesByFolderResponse } from "../api/files";
import FileCard from "./components/filecard";
import {Button} from "@/components/ui/button";
import { PlusIcon } from "lucide-react";
import {addFolder } from '@/app/api/folder';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogDescription
} from "@/components/ui/dialog";
import { DndContext, closestCenter,DragOverlay } from '@dnd-kit/core';
import { useSortable } from '@dnd-kit/sortable';
import { moveFile } from "@/app/api/files";
import { uploadFiles } from "@/app/api/files";
import { CSS } from '@dnd-kit/utilities';

const generateRandomId = () => {
    return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
}

function Draggable({id, file, dragDisabled}: {id: number, file: FileMetaData, dragDisabled: boolean}) {
  const [isDragging, setIsDragging] = useState(false);
  const {attributes, listeners, setNodeRef, transform, transition} = useSortable({id});
  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition: transform ? 'transform 200ms cubic-bezier(0.22, 1, 0.36, 1)' : undefined,
    pointerEvents: dragDisabled ? 'none' : 'auto'  
  };

  function handleClick(event: React.MouseEvent<HTMLDivElement>) {
    console.log("Clicked file:", file);
    if (file.IsDir) {
      event.stopPropagation();
      window.location.href = `/?folder=${file.Path}`;
    }
  }
  function handleDragStart() {
    if (dragDisabled) return;
    setIsDragging(true);
  }

  function handleDragEnd() {
    setTimeout(() => setIsDragging(false), 100);
  }


    return (
    <div ref={setNodeRef} style={style} {...listeners} {...attributes}
      onDragStart={handleDragStart} onDragEnd={handleDragEnd} onClick={handleClick}
    className='p-3 bg-white border rounded-md shadow cursor-move'>
      <FileCard key={id} file={file} />
    </div>
  )
}

export default function FileExplorer({folder}: {folder:string}) {
    const [files, setFiles] = useState<FileMetaData[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [size, setSize] = useState<number>(0);
    const [dialogOpen, setDialogOpen] = useState(false);
    const [modalType, setModalType] = useState<"files" | "folder" | null>(null);
    
    const [newFolderName, setNewFolderName] = useState<string>('');

    const [newFiles, setNewFiles] = useState<File[]>([]);
    const [uploading, setUploading] = useState(false);
    const fileInputRef = useRef<HTMLInputElement | null>(null);

    const [activeId, setActiveId] = useState<number | null>(null);
    const [dragDisabled, setDragDisabled] = useState(false);
    const dragTimeout = useRef<NodeJS.Timeout | null>(null);

    const createFolder = async() => {
        try {
          if (folder!=='') {
            await addFolder(folder + '/' + newFolderName);
          } else {
            await addFolder(newFolderName);
          }
            alert("Folder created successfully");
        }  catch (err) {
            console.error("Error creating folder:", err);
            alert("Error creating folder");
        }
    }

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

    const handleButtonClick = () => {
      fileInputRef.current?.click(); 
    };

    // const handleDragStart = (event: any) => setActiveId(event.active.id);

    
    // const handleDragEnd = (event: any) => {
    //   const {active,over}=event;
    //   if (!over || active.id == over.id) return;
    //   const overFile = files[active.id];
    //   if (overFile.IsDir) {
    //     alert("Cannot move a folder");
    //     setActiveId(null);
    //     return;
    //   };


  //     const endFile = files[over.id];
  //     if (endFile.IsDir) {
  //       moveFile(overFile.Path, endFile.Path).then(res => {
  //         if (res) {
  //           alert("File moved successfully");
  //           window.location.reload();
  //         }
  //         setDragDisabled(true);
  //       if (dragTimeout.current) clearTimeout(dragTimeout.current);
  //       dragTimeout.current = setTimeout(() => setDragDisabled(false), 150);
  //     }).catch(err => {
  //       console.error("Error moving file", err);
  //       alert("Error moving file");
  //     });
  //     setActiveId(null);
  //   }
  // }
  const handleDragCancel = () => setActiveId(null);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    
        if (e.target.files) {
          setNewFiles(Array.from(e.target.files));
        }
      };

    const handleUpload = async (e: React.FormEvent) => {
        e.preventDefault();
        if (newFiles.length === 0) return;
        setUploading(true);
        try {
            const res = await uploadFiles(folder,newFiles);
          if (!res) throw new Error("Upload failed");
          alert("Files uploaded successfully!");
          setNewFiles([]);
        } catch (err) {
          console.error(err);
          alert("Error uploading files");
        } finally {
          setUploading(false);
          setModalType(null);
          window.location.reload();
        }
      };
    

    if (loading) return <p>Loading...</p>;

    return (
      <div className="p-4">
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
              <Button type="submit" onClick={(e) => handleUpload(e)}>Submit</Button>
            </form>
          </DialogContent>
        </Dialog>

          <Dialog open={modalType === "folder"} onOpenChange={() => setModalType(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add Folder</DialogTitle>
            </DialogHeader>
            <form className="flex flex-col gap-4">
              <input type="text" multiple value={newFolderName} onChange={(e) => setNewFolderName(e.target.value)}  className="border rounded p-2" />
              <Button type="submit" onClick={() => createFolder()}>Submit</Button>
            </form>
          </DialogContent>
        </Dialog>

      <DndContext collisionDetection={closestCenter} 
        onDragCancel={handleDragCancel} >
        <ul role="list" className="grid grid-cols-6 gap-6 sm:grid-cols-4 lg:grid-cols-7">
      {files.map((file,indx) => (
        
        <Draggable key={indx} id={indx} file={file} dragDisabled={dragDisabled} />
        
      ))}
        </ul>      
        <DragOverlay>
        {activeId ? (
          <div className="p-3 bg-white border rounded-md shadow-lg scale-105 opacity-90">
            <Draggable id={activeId} file={files[activeId]} dragDisabled={true}/>
          </div>
        ) : null}
      </DragOverlay>
      </DndContext>
      </div>
    );

}